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

type AMQPClient interface {
	LocalAddr() net.Addr
	ConnectionState() tls.ConnectionState
	NotifyClose(receiver chan *amqp.Error) chan *amqp.Error
	NotifyBlocked(receiver chan amqp.Blocking) chan amqp.Blocking
	Close() error
	IsClosed() bool
	Channel() (*amqp.Channel, error)
}

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

// AMQPConnection is a struct to handle AMQP connection
type AMQPConnection struct {
	connection AMQPClient
	config     *AMQPConnectionConfig
}

type AMQPClientCreateFunc func() (AMQPClient, error)

// NewAMQPConnection is AMQPConnection constructor without auto reconnect
func NewAMQPConnection(client AMQPClient) (*AMQPConnection, error) {
	conn := &AMQPConnection{
		connection: client,
	}

	return conn, nil
}

// NewAMQPConnectionWithReconnect is AMQPConnection constructor with auto reconnect
func NewAMQPConnectionWithReconnect(clientFunc AMQPClientCreateFunc) (*AMQPConnection, error) {
	amqpClient, err := clientFunc()
	if err != nil {
		return nil, err
	}

	conn := &AMQPConnection{
		connection: amqpClient,
	}

	go conn.reconnectConn(clientFunc)
	return conn, nil
}

// Dial is AMQP client connection
func Dial(config *AMQPConnectionConfig) AMQPClientCreateFunc {
	return func() (AMQPClient, error) {
		var connection *amqp.Connection
		var err error

		url := fmt.Sprintf(
			"%s://%s:%s@%s:%d/",
			"amqp",
			config.user,
			config.password,
			config.host,
			config.port,
		)
		if config.config != nil {
			connection, err = amqp.DialConfig(
				url,
				*config.config,
			)
		} else {
			connection, err = amqp.Dial(url)
		}
		if err != nil {
			return nil, xerrors.Errorf("unable to start AMQP subsystem: %v", err)
		}
		return connection, err
	}
}

func (ac *AMQPConnection) LocalAddr() net.Addr {
	return ac.connection.LocalAddr()
}

func (ac *AMQPConnection) ConnectionState() tls.ConnectionState {
	return ac.connection.ConnectionState()
}

func (ac *AMQPConnection) NotifyClose(receiver chan *amqp.Error) chan *amqp.Error {
	return ac.connection.NotifyClose(receiver)
}

func (ac *AMQPConnection) NotifyBlocked(receiver chan amqp.Blocking) chan amqp.Blocking {
	return ac.connection.NotifyBlocked(receiver)
}

func (ac *AMQPConnection) Close() error {
	return ac.connection.Close()
}

func (ac *AMQPConnection) IsClosed() bool {
	return ac.connection.IsClosed()
}

// Private method to handle connection reconnect on error / close / timeout
func (ac *AMQPConnection) reconnectConn(reconnect AMQPClientCreateFunc) {
	connClose := ac.connection.NotifyClose(make(chan *amqp.Error, 1))

	closeReason, open := <-connClose
	if !open {
		debug("AMQP connection closed with reason: %v", closeReason)
		return
	}
	cb := breaker.New(3, 1, 5*time.Second)
	var hasReconnected = false
	for !hasReconnected {
		result := cb.Run(func() (err error) {
			connection, err := reconnect()
			if err != nil {
				return err
			}
			ac.connection = connection
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
	go ac.reconnectConn(reconnect)
}

// ------------------------------------------------------------------------------

// AMQPChannel is AMQP chan wrapper struct
type AMQPChannel struct {
	sync.Mutex
	*amqp.Channel
	closed uint32
}

// NewAMQPChannel is wrapper for AMPQ channels constructor
func NewAMQPChannel(channel *amqp.Channel) *AMQPChannel {
	return &AMQPChannel{
		Channel: channel,
		closed:  0,
	}
}

func (ac *AMQPConnection) Channel() (*AMQPChannel, error) {
	ch, err := ac.connection.Channel()
	if err != nil {
		return nil, err
	}
	wrapperChan := NewAMQPChannel(ch)
	go ac.reconnectChannel(wrapperChan)
	return wrapperChan, nil
}

// Private method to handle channel reconnect on error / close / timeout
func (ac *AMQPConnection) reconnectChannel(c *AMQPChannel) {
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
			newChannel, err = ac.connection.Channel()
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

// IsClosed change internal channel state
func (ac *AMQPChannel) IsClosed() bool {
	return atomic.LoadUint32(&ac.closed) == 1
}

// Close closes running channel and change internal state
func (ac *AMQPChannel) Close() error {
	if ac.IsClosed() {
		return amqp.ErrClosed
	}
	if err := ac.Channel.Close(); err != nil {
		return err
	}
	atomic.StoreUint32(&ac.closed, 1)
	return nil
}
