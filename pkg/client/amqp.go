package client

import (
	"fmt"
	syncPkg "github.com/buyco/funicular/pkg/sync"
	"github.com/buyco/keel/pkg/helper"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
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

// AMQPManager is a struct to manage AMQP connections
type AMQPManagerConfig struct {
	address  string
	user     string
	password string
	config   *amqp.Config
}

// NewAMQPClient is AMQP client constructor
func NewAMQPManagerConfig(address, user, password string, config *amqp.Config) *AMQPManagerConfig {
	return &AMQPManagerConfig{
		address:  address,
		user:     user,
		password: password,
		config:   config,
	}
}

// AMQPManager is a struct to manage AMQP connections
type AMQPManager struct {
	config     *AMQPManagerConfig
	connection *amqp.Connection
	pool       *syncPkg.Pool
	shutdown   bool
	Reconnects uint64
	logger     *logrus.Logger
	sync.Mutex
}

// NewAMQPClient is AMQP client constructor
func (am *AMQPManager) newAMQPConnection() (err error) {
	var connection *amqp.Connection
	url := fmt.Sprintf(
		"%s://%s:%s@%s/",
		"amqp",
		am.config.user,
		am.config.password,
		am.config.address,
	)
	if am.config.config != nil {
		connection, err = amqp.DialConfig(
			url,
			*am.config.config,
		)
	} else {
		connection, err = amqp.Dial(url)
	}
	if err != nil {
		return helper.ErrorPrintf("unable to start AMQP subsytem: %v", err)
	}
	am.Lock()
	am.connection = connection
	am.Unlock()
	return nil
}

// NewAMQPManager is AMQPManager constructor
func NewAMQPManager(config *AMQPManagerConfig, maxCap uint, logger *logrus.Logger) (*AMQPManager, error) {
	manager := &AMQPManager{
		config: config,
		pool:   syncPkg.NewPool(maxCap, nil, logger),
		logger: logger,
	}
	err := manager.newAMQPConnection()
	if err != nil {
		return nil, err
	}
	go manager.reconnectConn()
	return manager, nil
}

// SetPoolFactory adds func factory to pool
func (am *AMQPManager) SetPoolFactory(factory syncPkg.Factory) {
	am.pool.SetFactory(factory)
}

// AddClient adds a new Channel in pool
func (am *AMQPManager) AddClient() error {
	channel, err := am.connection.Channel()
	if err != nil {
		return err
	}
	wrapper := NewAMQPWrapper(channel)
	go am.reconnectChannel(wrapper)
	am.pool.Put(wrapper)
	return nil
}

// GetChannel get a new AMQP connection channel in pool
func (am *AMQPManager) GetClient() (*AMQPWrapper, error) {
	amqpWrapper := am.pool.Get()
	if amqpWrapper == nil {
		return nil, helper.ErrorPrint("no AMQP client available")
	}
	return amqpWrapper.(*AMQPWrapper), nil
}

// PutClient add an existing AMQP client in pool
func (am *AMQPManager) PutClient(client *AMQPWrapper) {
	err := am.pool.Put(client)
	if err != nil {
		am.logger.Warn(err)
		err = client.Channel.Close()
		if err != nil {
			am.logger.WithError(err).Warn("An error occurred while closing channel")
		}
	}
}

// Close closes AMQP connection and channels
func (am *AMQPManager) Close() error {
	am.shutdown = true
	err := am.connection.Close()
	if err != nil {
		am.shutdown = false
		return err
	}
	return nil
}

// Private method to handle channel reconnect on error / close / timeout
func (am *AMQPManager) reconnectChannel(c *AMQPWrapper) {
	chanClose := c.Channel.NotifyClose(make(chan *amqp.Error, 1))
	select {
	case resChan := <-chanClose:
		if am.shutdown && am.connection.IsClosed() {
			am.logger.Debugf("AMQP connection is closing, stopping channel auto-reconnect loop")
			break
		}
		am.logger.Debugf("Channel closed, reconnecting: %v", resChan)
		cb := breaker.New(3, 1, 5*time.Second)
		var (
			newChannel     *amqp.Channel
			hasReconnected = false
		)
		for !hasReconnected {
			result := cb.Run(func() (err error) {
				newChannel, err = am.connection.Channel()
				if err != nil {
					return err
				}
				return nil
			})

			switch result {
			case nil:
				hasReconnected = true
				am.logger.Info("Channel recreated")
			case breaker.ErrBreakerOpen:
			default:
				am.logger.Errorf("Failed to reconnect channel: %v", result)
			}
		}

		atomic.AddUint64(&c.Reconnects, 1)
		c.Lock()
		c.Channel = newChannel
		c.Unlock()
		// New connections set, rerun async reconnect
		am.reconnectChannel(c)
	}
}

// Private method to handle connection reconnect on error / close / timeout
func (am *AMQPManager) reconnectConn() {
	connClose := am.connection.NotifyClose(make(chan *amqp.Error, 1))
	select {
	case resConn := <-connClose:
		if am.shutdown && am.connection.IsClosed() {
			am.logger.Debugf("Stop AMQP auto-reconnect loop")
			break
		}
		am.logger.Debugf("AMQP connection closed, reconnecting: %s", resConn)
		cb := breaker.New(3, 1, 5*time.Second)
		var hasReconnected = false
		for !hasReconnected {
			result := cb.Run(func() (err error) {
				err = am.newAMQPConnection()
				if err != nil {
					return err
				}
				return nil
			})

			switch result {
			case nil:
				hasReconnected = true
				am.logger.Info("AMQP connection recreated")
			case breaker.ErrBreakerOpen:
			default:
				am.logger.Errorf("Failed to reconnect AMQP connection: %v", result)
			}
		}
		atomic.AddUint64(&am.Reconnects, 1)
		// New connections set, rerun async reconnect
		am.reconnectConn()
	}
}

//------------------------------------------------------------------------------

// AMQPWrapper is AMQP client wrapper struct
type AMQPWrapper struct {
	sync.Mutex
	Channel    *amqp.Channel
	Reconnects uint64
}

// NewAMQPWrapper is AMQPWrapper constructor
func NewAMQPWrapper(channel *amqp.Channel) *AMQPWrapper {
	return &AMQPWrapper{
		Channel:    channel,
		Reconnects: 0,
	}
}
