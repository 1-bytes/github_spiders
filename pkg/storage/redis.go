package storage

import (
	"fmt"
	"github.com/gocolly/redisstorage"
	"github_spiders/pkg/config"
)

// GetRedisStorage 获取 redis storage 实例.
func GetRedisStorage(tag string) redisstorage.Storage {
	return redisstorage.Storage{
		Address: fmt.Sprintf("%s:%s",
			config.GetString("redis.github.host"),
			config.GetString("redis.github.port"),
		),
		Password: config.GetString("redis.github.password"),
		DB:       config.GetInt("redis.github.db"),
		Prefix: fmt.Sprintf("%s_%s",
			config.GetString("redis.github.prefix"),
			tag,
		),
	}
}
