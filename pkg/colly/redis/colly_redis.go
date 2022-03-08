package redis

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/redisstorage"
	"github_spiders/spiders/types"
	"strconv"
)

// NewRedisStorage 创建一个新的 Redis Storage.
func NewRedisStorage(c *colly.Collector, cfg types.CollyRedisConfig) *redisstorage.Storage {
	storage := &redisstorage.Storage{
		Address:  cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
		Prefix:   cfg.Prefix + strconv.Itoa(int(c.ID)),
	}
	_ = c.SetStorage(storage)
	return storage
}
