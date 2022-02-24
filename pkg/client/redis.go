package client

import (
	"github.com/go-redis/redis/v7"
	"golang.org/x/xerrors"
	"net"
	"strconv"
	"sync"
)

// RedisConfig is a struct definition for Redis Client
type RedisConfig struct {
	Host string
	Port uint16
	DB   uint8
}

// ToOption is used by RedisConfig to generate Option struct from client
func (rc *RedisConfig) ToOption() *redis.Options {
	return &redis.Options{
		Addr: net.JoinHostPort(rc.Host, strconv.Itoa(int(rc.Port))),
		DB:   int(rc.DB),
	}
}

//------------------------------------------------------------------------------

// RedisManager is a struct to manage Redis clients
type RedisManager struct {
	config  RedisConfig
	Clients map[string]*redis.Client
	sync.RWMutex
}

// NewRedisManager is RedisManager struct constructor
func NewRedisManager(config RedisConfig) *RedisManager {
	return &RedisManager{
		config:  config,
		Clients: make(map[string]*redis.Client),
	}
}

// AddClient pushes a new client in manager
func (rw *RedisManager) AddClient(category string) (*redis.Client, error) {
	rw.Lock()
	defer rw.Unlock()
	if category == "" {
		return nil, xerrors.New("category must be filled")
	}
	client := redis.NewClient(rw.config.ToOption())
	return rw.set(client, category), nil
}

// GetCategories fetches categories of client available
func (rw *RedisManager) GetCategories() (clientsCat []string) {
	rw.RLock()
	defer rw.RUnlock()
	for key := range rw.Clients {
		clientsCat = append(clientsCat, key)
	}
	return
}

// Close closes all connections
func (rw *RedisManager) Close() error {
	var err error
	if len(rw.Clients) > 0 {
		for _, redisClient := range rw.Clients {
			err = redisClient.Close()
			if err != nil {
				return xerrors.Errorf("an error occurred while closing client connection pool: %v", err)
			}
		}
	} else {
		err = xerrors.New("manager have no clients to close")
	}
	return err
}

func (rw *RedisManager) set(client *redis.Client, category string) *redis.Client {
	content := rw.Clients[category]
	if content == nil {
		rw.Clients[category] = client
	} else {
		debug("Redis client already set for category [%s]", category)
	}
	return rw.Clients[category]
}
