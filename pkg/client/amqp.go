package client

import (
	"fmt"
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
	config *AMQPManagerConfig
	connection *amqp.Connection
	closedConn chan *amqp.Error
	pool       *Pool
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
		pool:   NewPool(maxCap, nil, logger),
		closedConn:   make(chan *amqp.Error, 1),
		logger: logger,
	}
	err := manager.newAMQPConnection()
	if err != nil {
		return nil, err
	}
	manager.connection.NotifyClose(manager.closedConn)
	go manager.reconnectConn()
	return manager, nil
}

// SetPoolFactory adds func factory to pool
func (am *AMQPManager) SetPoolFactory(factory Factory) {
	am.pool.SetFactory(factory)
}

// AddClient adds a new Channel in pool
func (am *AMQPManager) AddChannel() error {
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
func (am *AMQPManager) GetChannel() (*AMQPWrapper, error) {
	sftpClient := am.pool.Get()
	if sftpClient == nil {
		return nil, helper.ErrorPrint("No SFTP client available")
	}
	return sftpClient.(*AMQPWrapper), nil
}

// PutClient add an existing SFTP client in pool
func (am *AMQPManager) PutClient(client *AMQPWrapper) {
	am.pool.Put(client)
}

// Close closes all SFTP connections
func (am *AMQPManager) Close() error {
	for {
		conn := am.pool.Get()
		if conn == nil {
			return nil
		}
		if err := conn.(*AMQPWrapper).Close(); err != nil {
			return err
		}
	}
}

// Private method to handle channel reconnect on error / close / timeout
func (am *AMQPManager) reconnectChannel(c *AMQPWrapper) {
	c.Channel.NotifyClose(c.closedChan)
	select {
	case resChan := <-c.closedChan:
		am.logger.Debugf("Channel closed, reconnecting: %s", resChan)
		cb := breaker.New(3, 1, 5*time.Second)
		var (
			newChannel *amqp.Channel
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
			case breaker.ErrBreakerOpen:
			default:
				am.logger.Errorf("Failed to reconnect: %v", result)
			}
		}

		atomic.AddUint64(&c.Reconnects, 1)
		c.Lock()
		c.Channel = newChannel
		c.closed = false
		c.Unlock()
		// New connections set, rerun async reconnect
		am.reconnectChannel(c)
	}
}

// Private method to handle connection reconnect on error / close / timeout
func (am *AMQPManager) reconnectConn() {
	select {
	case resConn := <-am.closedConn:
		am.logger.Debugf("Connection closed, reconnecting: %s", resConn)
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
			case breaker.ErrBreakerOpen:
			default:
				am.logger.Errorf("Failed to reconnect: %v", result)
			}
		}
	}
}

//------------------------------------------------------------------------------

// AMQPWrapper is SFTP client wrapper struct
type AMQPWrapper struct {
	sync.Mutex
	Channel     *amqp.Channel
	closedChan chan *amqp.Error
	shutdown   chan bool
	closed     bool
	Reconnects uint64
	logger     *logrus.Logger
}

// NewAMQPWrapper is AMQPWrapper constructor
func NewAMQPWrapper(channel *amqp.Channel) *AMQPWrapper {
	return &AMQPWrapper{
		Channel:    channel,
		shutdown:   make(chan bool, 1),
		closedChan:   make(chan *amqp.Error, 1),
		closed:     false,
		Reconnects: 0,
	}
}

// Close closes channel from AMQP => chan notify ssh connection to close
func (s *AMQPWrapper) Close() error {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return helper.ErrorPrint("Connection was already closed")
	}
	var err = s.Channel.Close()
	if err != nil {
		return helper.ErrorPrintf("unable to close ftp connection: %v", err)
	}
	s.shutdown <- true
	s.closed = true
	return nil
}
