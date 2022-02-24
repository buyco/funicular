package client

import (
	"fmt"
	"github.com/streadway/amqp"
	"golang.org/x/xerrors"
	"gopkg.in/eapache/go-resiliency.v1/breaker"
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

//------------------------------------------------------------------------------

// AMQPConnection is a struct to manage AMQP connections
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
	*amqp.Connection
	config *AMQPConnectionConfig
}

// NewAMQPConnection is AMQPConnection constructor
func NewAMQPConnection(config *AMQPConnectionConfig) (*AMQPConnection, error) {
	conn := &AMQPConnection{
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
func (ac *AMQPConnection) createAMQPConnection() (err error) {
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
func (ac *AMQPConnection) reconnectConn() {
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

//------------------------------------------------------------------------------

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
	ch, err := ac.Connection.Channel()
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
