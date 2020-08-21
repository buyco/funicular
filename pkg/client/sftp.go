package client

import (
	"fmt"
	syncPkg "github.com/buyco/funicular/pkg/sync"
	"github.com/buyco/keel/pkg/helper"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"gopkg.in/eapache/go-resiliency.v1/breaker"
	"sync"
	"sync/atomic"
	"time"
)

// NewSSHConfig is SSH ClientConfig constructor
func NewSSHConfig(user string, password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         2 * time.Second,
	}
}

// NewSSHClient is SSH client constructor
func NewSSHClient(host string, port uint32, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)

	if err != nil {
		err = helper.ErrorPrintf("unable to connect to [%s]: %v", addr, err)
		return nil, err
	}
	return conn, err
}

// NewSFTPClient is SFTP client constructor
func NewSFTPClient(sshClient *ssh.Client) (*sftp.Client, error) {
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		err = helper.ErrorPrintf("unable to start sftp subsytem: %v", err)
		return nil, err
	}
	return client, err
}

//------------------------------------------------------------------------------

// SFTPManager is a struct to manage SFTP connections
type SFTPManager struct {
	host      string
	port      uint32
	user      string
	password  string
	pool      *syncPkg.Pool
	sshConfig *ssh.ClientConfig
	logger    *logrus.Logger
	sync.Mutex
}

// NewSFTPManager is SFTPManager constructor
func NewSFTPManager(host string, port uint32, sshConfig *ssh.ClientConfig, maxCap uint, logger *logrus.Logger) *SFTPManager {
	return &SFTPManager{
		host:      host,
		port:      port,
		pool:      syncPkg.NewPool(maxCap, nil, logger),
		sshConfig: sshConfig,
		logger:    logger,
	}
}

// SetPoolFactory adds func factory to pool
func (sm *SFTPManager) SetPoolFactory(factory syncPkg.Factory) {
	sm.pool.SetFactory(factory)
}

// AddClient adds a new SFTP client in pool
func (sm *SFTPManager) AddClient() error {
	sshConn, sftpConn, err := sm.newConnections()
	if err != nil {
		return err
	}
	sftpStruct := NewSFTPWrapper(sshConn, sftpConn)
	go sm.reconnect(sftpStruct)
	err = sm.pool.Put(sftpStruct)
	if err != nil {
		sm.logger.Warn(err)
	}
	return nil
}

// GetClient get a new SFTP client in pool
func (sm *SFTPManager) GetClient() (*SFTPWrapper, error) {
	sftpClient := sm.pool.Get()
	if sftpClient == nil {
		return nil, helper.ErrorPrint("no SFTP client available")
	}
	return sftpClient.(*SFTPWrapper), nil
}

// PutClient add an existing SFTP client in pool
func (sm *SFTPManager) PutClient(client *SFTPWrapper) {
	err := sm.pool.Put(client)
	if err != nil {
		sm.logger.Warn(err)
		err := client.Close()
		if err != nil {
			sm.logger.WithError(err).Warn("An error occurred while closing connection")
		}
	}
}

// Close closes all SFTP connections
func (sm *SFTPManager) Close() error {
	for {
		conn := sm.pool.Get()
		if conn == nil {
			return nil
		}
		if err := conn.(*SFTPWrapper).Close(); err != nil {
			return err
		}
	}
}

// Private method to create ssh and sftp clients
func (sm *SFTPManager) newConnections() (*ssh.Client, *sftp.Client, error) {
	sshConn, err := NewSSHClient(sm.host, sm.port, sm.sshConfig)
	if err != nil {
		return sshConn, nil, err
	}
	sftpConn, err := NewSFTPClient(sshConn)
	if err != nil {
		return nil, sftpConn, err
	}

	return sshConn, sftpConn, err
}

// Private method to handle reconnect on error / close / timeout
func (sm *SFTPManager) reconnect(c *SFTPWrapper) {
	closed := make(chan error, 1)
	go func() {
		closed <- c.connection.Wait()
	}()

	select {
	case <-c.shutdown:
		_ = c.connection.Close()
		break
	case res := <-closed:
		sm.logger.Debugf("SFTP connection closed, reconnecting: %s", res)
		cb := breaker.New(3, 1, 5*time.Second)
		var (
			sshConn        *ssh.Client
			sftpConn       *sftp.Client
			hasReconnected = false
		)
		for !hasReconnected {
			result := cb.Run(func() error {
				var err error
				sshConn, sftpConn, err = sm.newConnections()
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
				sm.logger.Errorf("SFTP failed to reconnect: %v", result)
			}
		}

		atomic.AddUint64(&c.Reconnects, 1)
		c.Lock()
		c.connection = sshConn
		c.Client = sftpConn
		c.closed = false
		c.Unlock()
		// New connections set, rerun async reconnect
		sm.reconnect(c)
	}
}

//------------------------------------------------------------------------------

// SFTPWrapper is SFTP client wrapper struct
type SFTPWrapper struct {
	sync.Mutex
	connection *ssh.Client
	Client     *sftp.Client
	shutdown   chan bool
	closed     bool
	Reconnects uint64
}

// NewSFTPWrapper is SFTPWrapper constructor
func NewSFTPWrapper(sshClient *ssh.Client, sftpClient *sftp.Client) *SFTPWrapper {
	return &SFTPWrapper{
		connection: sshClient,
		Client:     sftpClient,
		shutdown:   make(chan bool, 1),
		closed:     false,
		Reconnects: 0,
	}
}

// Close closes connection from SFTP => chan notify ssh connection to close
func (s *SFTPWrapper) Close() error {
	s.Lock()
	defer s.Unlock()
	if s.closed {
		return helper.ErrorPrint("SFTP connection was already closed")
	}
	var err = s.Client.Close()
	if err != nil {
		return helper.ErrorPrintf("unable to close sftp connection: %v", err)
	}
	s.shutdown <- true
	s.closed = true
	return s.Client.Wait()
}
