// Package client contains struct for client third parties
package client

import (
	"crypto/tls"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/xerrors"
	"gopkg.in/eapache/go-resiliency.v1/breaker"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// NewAMQPConfig is a simple AMQP Config constructor
func NewAMQPConfig(vhost string, channelMax int, heartbeat time.Duration) amqp.Config {
	return amqp.Config{
		Vhost:      vhost,
		ChannelMax: channelMax,
		Heartbeat:  heartbeat,
	}
}

// AMQPConnectionConfig is a struct to manage AMQP connections
type AMQPConnectionConfig struct {
	host     string
	port     int
	user     string
	password string
	config   *amqp.Config
}

// NewAMQPConnectionConfig is AMQP client constructor
func NewAMQPConnectionConfig(host string, port int, user, password string, config *amqp.Config) *AMQPConnectionConfig {
	return &AMQPConnectionConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		config:   config,
	}
}

// AMQPConnection is a client wrapper interface
type AMQPConnection interface {
	LocalAddr() net.Addr
	ConnectionState() tls.ConnectionState
	NotifyClose(receiver chan *amqp.Error) chan *amqp.Error
	NotifyBlocked(receiver chan amqp.Blocking) chan amqp.Blocking
	Close() error
	IsClosed() bool
	Channel() (AMQPChannel, error)
}

// amqpConnection is a struct to handle AMQP connection
type amqpConnection struct {
	*amqp.Connection
	config *AMQPConnectionConfig
}

// NewAMQPConnection is AMQPConnection constructor
func NewAMQPConnection(config *AMQPConnectionConfig) (AMQPConnection, error) {
	conn := &amqpConnection{
		config: config,
	}
	err := conn.createAMQPConnection()
	if err != nil {
		return nil, err
	}
	go conn.reconnectConn()
	return conn, nil
}

// createAMQPConnection is AMQP connection setter
func (ac *amqpConnection) createAMQPConnection() (err error) {
	var connection *amqp.Connection
	url := fmt.Sprintf(
		"%s://%s:%s@%s:%d/",
		"amqp",
		ac.config.user,
		ac.config.password,
		ac.config.host,
		ac.config.port,
	)
	if ac.config.config != nil {
		connection, err = amqp.DialConfig(
			url,
			*ac.config.config,
		)
	} else {
		connection, err = amqp.Dial(url)
	}
	if err != nil {
		return xerrors.Errorf("unable to start AMQP subsystem: %v", err)
	}
	ac.Connection = connection
	return nil
}

// Private method to handle connection reconnect on error / close / timeout
func (ac *amqpConnection) reconnectConn() {
	connClose := ac.Connection.NotifyClose(make(chan *amqp.Error, 1))

	closeReason, open := <-connClose
	if !open {
		debug("AMQP connection closed with reason: %v", closeReason)
		return
	}
	cb := breaker.New(3, 1, 5*time.Second)
	var hasReconnected = false
	for !hasReconnected {
		result := cb.Run(func() (err error) {
			err = ac.createAMQPConnection()
			if err != nil {
				return err
			}
			return nil
		})

		switch result {
		case nil:
			debug("AMQP connection reconnected")
			hasReconnected = true
		case breaker.ErrBreakerOpen:
		default:
		}
	}
	// New connections set, rerun async reconnect
	go ac.reconnectConn()
}

func (ac *amqpConnection) Channel() (AMQPChannel, error) {
	ch, err := ac.Connection.Channel()
	if err != nil {
		return nil, err
	}
	wrapperChan := newAMQPChannel(ch)
	go ac.reconnectChannel(wrapperChan)
	return wrapperChan, nil
}

// Private method to handle channel reconnect on error / close / timeout
func (ac *amqpConnection) reconnectChannel(c *amqpChannel) {
	chanClose := c.Channel.NotifyClose(make(chan *amqp.Error, 1))
	closeReason, open := <-chanClose
	if !open || c.IsClosed() {
		debug("AMQP channel closed with reason: %v", closeReason)
		if err := c.Close(); err != nil {
			debug("%s -> %v", "AMQP close error", err)
		}
		return
	}
	cb := breaker.New(3, 1, 5*time.Second)
	var (
		newChannel     *amqp.Channel
		hasReconnected = false
	)
	for !hasReconnected {
		result := cb.Run(func() (err error) {
			newChannel, err = ac.Connection.Channel()
			if err != nil {
				return err
			}
			return nil
		})

		switch result {
		case nil:
			debug("AMQP channel reconnected")
			hasReconnected = true
		case breaker.ErrBreakerOpen:
		default:
		}
	}
	c.Channel = newChannel
	// New connections set, rerun async reconnect
	go ac.reconnectChannel(c)
}

// AMQPChannel is an interface for amqp channel wrapper
type AMQPChannel interface {
	Close() error
	IsClosed() bool
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
	NotifyFlow(c chan bool) chan bool
	NotifyReturn(c chan amqp.Return) chan amqp.Return
	NotifyCancel(c chan string) chan string
	NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64)
	NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation
	Qos(prefetchCount, prefetchSize int, global bool) error
	Cancel(consumer string, noWait bool) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDeclarePassive(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueInspect(name string) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	QueueUnbind(name, key, exchange string, args amqp.Table) error
	QueuePurge(name string, noWait bool) (int, error)
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error)
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeDeclarePassive(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeDelete(name string, ifUnused, noWait bool) error
	ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error
	ExchangeUnbind(destination, key, source string, noWait bool, args amqp.Table) error
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	PublishWithDeferredConfirm(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) (*amqp.DeferredConfirmation, error)
	Get(queue string, autoAck bool) (msg amqp.Delivery, ok bool, err error)
	Tx() error
	TxCommit() error
	TxRollback() error
	Flow(active bool) error
	Confirm(noWait bool) error
	Recover(requeue bool) error
	Ack(tag uint64, multiple bool) error
	Nack(tag uint64, multiple bool, requeue bool) error
	Reject(tag uint64, requeue bool) error
	GetNextPublishSeqNo() uint64
}

// amqpChannel is AMQP chan wrapper struct
type amqpChannel struct {
	sync.Mutex
	*amqp.Channel
	closed uint32
}

// newAMQPChannel is wrapper for AMPQ channels constructor
func newAMQPChannel(channel *amqp.Channel) *amqpChannel {
	return &amqpChannel{
		Channel: channel,
		closed:  0,
	}
}

// IsClosed change internal channel state
func (ac *amqpChannel) IsClosed() bool {
	return atomic.LoadUint32(&ac.closed) == 1
}

// Close closes running channel and change internal state
func (ac *amqpChannel) Close() error {
	if ac.IsClosed() {
		return amqp.ErrClosed
	}
	if err := ac.Channel.Close(); err != nil {
		return err
	}
	atomic.StoreUint32(&ac.closed, 1)
	return nil
}
