package redis

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
	"github_spiders/spiders/types"
	"strconv"
	"sync"
)

var (
	collyRedisInstance map[int]*collyRedis
	lock               sync.Mutex
)

type collyRedis struct {
	redisStorage *redisstorage.Storage
}

// init 初始化.
func init() {
	collyRedisInstance = make(map[int]*collyRedis)
}

// GetInstance 根据 ID 作为标识获取一个单例.
func GetInstance(id int) *collyRedis {
	if collyRedisInstance[id] != nil {
		return collyRedisInstance[id]
	}
	return func() *collyRedis {
		lock.Lock()
		defer lock.Unlock()
		if collyRedisInstance[id] == nil {
			collyRedisInstance[id] = &collyRedis{}
		}
		return collyRedisInstance[id]
	}()
}

// NewRedisStorage 创建一个新的 redis storage.
func (r *collyRedis) NewRedisStorage(c *colly.Collector, cfg types.CollyRedisConfig) *redisstorage.Storage {
	r.redisStorage = &redisstorage.Storage{
		Address:  cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
		Prefix:   cfg.Prefix + strconv.Itoa(int(c.ID)),
	}
	_ = c.SetStorage(r.redisStorage)
	return r.redisStorage
}

// GetRedisStorage 获取指定单例的 redis storage.
func (r *collyRedis) GetRedisStorage() *redisstorage.Storage {
	return r.redisStorage
}
