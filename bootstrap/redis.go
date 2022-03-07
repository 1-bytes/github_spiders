package bootstrap

import (
	"github.com/gocolly/redisstorage"
	"github_spiders/pkg/config"
)

// SetupCollyRedis 初始化 Colly Redis.
func SetupCollyRedis() {
	storage := &redisstorage.Storage{
		Address:  config.GetString("redis.github.address"),
		Password: config.GetString("redis.github.password"),
		DB:       config.GetInt("redis.github.db"),
		Prefix:   config.GetString("redis.github.prefix"),
	}
	err := GetCollector().SetStorage(storage)
	if err != nil {
		panic(err)
	}
}
