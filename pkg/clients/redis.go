package clients

import (
	"crypto/sha256"
	"fmt"
	"github.com/buyco/keel/pkg/utils"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

type RedisConfig struct {
	Host string
	Port uint16
	DB   uint8
}

func (rc *RedisConfig) ToOption() *redis.Options {
	return &redis.Options{
		Addr: net.JoinHostPort(rc.Host, strconv.Itoa(int(rc.Port))),
		DB:   int(rc.DB),
	}
}

//------------------------------------------------------------------------------

type RedisManager struct {
	config  RedisConfig
	Clients map[string][]*RedisWrapper
	logger  *logrus.Logger
	sync.RWMutex
}

func NewRedisManager(config RedisConfig, logger *logrus.Logger) *RedisManager {
	return &RedisManager{
		config:  config,
		Clients: make(map[string][]*RedisWrapper),
		logger:  logger,
	}
}

func (rw *RedisManager) AddClient(category string, channel string, consumerName string) (*RedisWrapper, error) {
	rw.Lock()
	defer rw.Unlock()
	if category == "" {
		return nil, utils.ErrorPrint("category must be filled")
	}
	if channel == "" {
		channel = category
	}
	client, _ := NewRedisWrapper(rw.config, channel, consumerName)
	rw.add(client, category)
	return client, nil
}

func (rw *RedisManager) GetCategories() (clientsCat []string) {
	rw.RLock()
	defer rw.RUnlock()
	for key := range rw.Clients {
		clientsCat = append(clientsCat, key)
	}
	return
}

func (rw *RedisManager) Close() error {
	var err error
	if len(rw.Clients) > 0 {
		for _, clients := range rw.Clients {
			for _, client := range clients {
				if client.Closed {
					rw.logger.Debug("Ignore closing client. Already closed")
					continue
				}
				err = client.Close()
				if err != nil {
					return utils.ErrorPrintf("an error occurred while closing client connexion pool: %v", err)
				}
			}
		}
	} else {
		err = utils.ErrorPrint("manager have no clients to close")
	}
	return err
}

func (rw *RedisManager) add(redisWrapper *RedisWrapper, category string) {
	mm, ok := rw.Clients[category]
	if !ok {
		mm = make([]*RedisWrapper, 0)
	}
	mm = append(mm, redisWrapper)
	rw.Clients[category] = mm
}

//------------------------------------------------------------------------------

type RedisWrapper struct {
	client       *redis.Client
	config       *RedisConfig
	channel      string
	consumerName string
	Closed       bool
	sync.RWMutex
}

func NewRedisWrapper(config RedisConfig, channel string, consumerName string) (*RedisWrapper, error) {
	if channel == "" {
		return nil, utils.ErrorPrint("channel must be filled")
	}
	if consumerName == "" {
		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("%f", rand.Float64())))
		consumerName = fmt.Sprintf("%x", h.Sum(nil))
	}
	client := redis.NewClient(config.ToOption())
	return &RedisWrapper{
			client:       client,
			config:       &config,
			channel:      channel,
			consumerName: consumerName,
			Closed:       false,
		},
		nil
}

func (w *RedisWrapper) Reconnect() error {
	if !w.Closed {
		return utils.ErrorPrint("client is not closed")
	}

	w.Lock()
	w.client = redis.NewClient(w.config.ToOption())
	w.Closed = false
	w.Unlock()

	return nil
}

func (w *RedisWrapper) AddMessage(data map[string]interface{}) (string, error) {
	xAddArgs := &redis.XAddArgs{
		Stream: w.channel,
		Values: data,
	}
	result := w.client.XAdd(xAddArgs)
	return result.Result()
}

func (w *RedisWrapper) ReadMessage(lastId string, count int64, block time.Duration) ([]redis.XStream, error) {
	var channels = make([]string, 0)
	channels = append(channels, w.channel)
	// Not implemented explicitly in the lib but works the way the code is written
	channels = append(channels, lastId)
	xReadArgs := &redis.XReadArgs{
		Streams: channels,
		Count:   count,
		Block:   block,
	}
	result := w.client.XRead(xReadArgs)
	return result.Result()
}

func (w *RedisWrapper) ReadGroupMessage(group string, count int64, block time.Duration, extraIds ...string) ([]redis.XStream, error) {
	var channels = make([]string, 0)
	channels = append(channels, w.channel)
	// The special > ID, which means that the consumer want to receive only messages that were never delivered to any other consumer
	channels = append(channels, ">")
	// Any other ids if given
	if len(extraIds) > 0 {
		for _, id := range extraIds {
			channels = append(channels, id)
		}
	}
	xReadGroupArgs := &redis.XReadGroupArgs{
		Group:    group,
		Streams:  channels,
		Consumer: w.consumerName,
		Count:    count,
		Block:    block,
		NoAck:    false,
	}
	result := w.client.XReadGroup(xReadGroupArgs)
	return result.Result()
}

func (w *RedisWrapper) GetChannel() string {
	return w.channel
}

func (w *RedisWrapper) ReadRangeMessage(start string, stop string) ([]redis.XMessage, error) {
	result := w.client.XRange(w.channel, start, stop)
	return result.Result()
}

func (w *RedisWrapper) DeleteMessage(ids ...string) (int64, error) {
	result := w.client.XDel(w.channel, ids...)
	return result.Result()
}

func (w *RedisWrapper) CreateGroup(group string, start string) (string, error) {
	// MKSTREAM is not documented in Redis and allow to create stream if it is not created beforehand
	result := w.client.XGroupCreateMkStream(w.channel, group, start)
	return result.Result()
}

func (w *RedisWrapper) DeleteGroup(group string) (int64, error) {
	result := w.client.XGroupDestroy(w.channel, group)
	return result.Result()
}

func (w *RedisWrapper) PendingMessage(group string) (*redis.XPending, error) {
	result := w.client.XPending(w.channel, group)
	return result.Result()
}

func (w *RedisWrapper) AckMessage(group string, ids ...string) (int64, error) {
	result := w.client.XAck(w.channel, group, ids...)
	return result.Result()
}

func (w *RedisWrapper) DeleteGroupConsumer(group string) (int64, error) {
	result := w.client.XGroupDelConsumer(w.channel, group, w.consumerName)
	return result.Result()
}

func (w *RedisWrapper) Close() error {
	w.Lock()
	defer w.Unlock()
	w.Closed = true
	return w.client.Close()
}

func (w *RedisWrapper) FlushAll() (string, error) {
	result := w.client.FlushAll()
	return result.Result()
}

func (w *RedisWrapper) FlushAllAsync() (string, error) {
	result := w.client.FlushAllAsync()
	return result.Result()
}

func (w *RedisWrapper) FlushDB() (string, error) {
	result := w.client.FlushDB()
	return result.Result()
}

func (w *RedisWrapper) FlushDBAsync() (string, error) {
	result := w.client.FlushDBAsync()
	return result.Result()
}
