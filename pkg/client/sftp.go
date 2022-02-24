package client

import (
	"fmt"
	syncPkg "github.com/buyco/funicular/pkg/sync"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/xerrors"
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
		return nil, xerrors.Errorf("unable to connect to [%s]: %v", addr, err)
	}
	return conn, nil
}

// NewSFTPClient is SFTP client constructor
func NewSFTPClient(sshClient *ssh.Client) (*sftp.Client, error) {
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, xerrors.Errorf("unable to start sftp subsytem: %v", err)
	}
	return client, nil
}

//------------------------------------------------------------------------------

// SFTPManager is a struct to manage SFTP connections
type SFTPManager struct {
	host      string
	port      uint32
	pool      *syncPkg.Pool
	sshConfig *ssh.ClientConfig
}

// NewSFTPManager is SFTPManager constructor
func NewSFTPManager(host string, port uint32, sshConfig *ssh.ClientConfig, maxCap uint) *SFTPManager {
	return &SFTPManager{
		host:      host,
		port:      port,
		pool:      syncPkg.NewPool(maxCap, nil),
		sshConfig: sshConfig,
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
		debug("Error: %v", err)
	}
	return nil
}

// GetClient get a new SFTP client in pool
func (sm *SFTPManager) GetClient() (*SFTPWrapper, error) {
	sftpClient := sm.pool.Get()
	if sftpClient == nil {
		return nil, xerrors.New("no SFTP client available")
	}
	return sftpClient.(*SFTPWrapper), nil
}

// PutClient add an existing SFTP client in pool
func (sm *SFTPManager) PutClient(client *SFTPWrapper) {
	err := sm.pool.Put(client)
	if err != nil {
		debug("Error: %v", err)
		err = client.Close()
		if err != nil {
			debug("%s -> %v", "An error occurred while closing connection", err)
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

	return sshConn, sftpConn, nil
}

// Private method to handle reconnect on error / close / timeout
func (sm *SFTPManager) reconnect(c *SFTPWrapper) {
	closed := make(chan error, 1)
	go func() {
		closed <- c.client.Wait()
	}()

	select {
	case <-c.shutdown:
		err := c.client.Close()
		if err != nil {
			debug("SFTP client cannot be closed: %v", err)
		}
		err = c.connection.Close()
		if err != nil {
			debug("SSH session cannot be closed: %v", err)
		}
		break
	case res := <-closed:
		c.mutex.Lock()
		debug("SFTP connection closed, reconnecting: %s", res)
		// Force SSH close & wait connection even if already closed
		c.connection.Close()
		c.connection.Wait()
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
				debug("SFTP failed to reconnect: %v", result)
			}
		}

		c.connection = sshConn
		c.client = sftpConn
		c.closed = false
		c.mutex.Unlock()
		atomic.AddUint64(&c.Reconnects, 1)
		debug("SFTP connection recovered")

		// New connections set, rerun async reconnect
		go sm.reconnect(c)
	}
}

//------------------------------------------------------------------------------

// SFTPWrapper is SFTP client wrapper struct
type SFTPWrapper struct {
	mutex      sync.Mutex
	connection *ssh.Client
	client     *sftp.Client
	shutdown   chan bool
	closed     bool
	Reconnects uint64
}

// NewSFTPWrapper is SFTPWrapper constructor
func NewSFTPWrapper(sshClient *ssh.Client, sftpClient *sftp.Client) *SFTPWrapper {
	return &SFTPWrapper{
		connection: sshClient,
		client:     sftpClient,
		shutdown:   make(chan bool, 1),
		closed:     false,
		Reconnects: 0,
	}
}

// Close closes connection from SFTP => chan notify ssh connection to close
func (s *SFTPWrapper) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.closed {
		return xerrors.New("SFTP connection was already closed")
	}
	s.shutdown <- true
	s.client.Wait()
	s.connection.Wait()
	s.closed = true

	return nil
}

// Client ensures that client can be fetched and is not reconnecting
func (s *SFTPWrapper) Client() *sftp.Client {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.client
}
