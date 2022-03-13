package bootstrap

import (
	"fmt"
	"github.com/gocolly/redisstorage"
	"github_spiders/pkg/collectors"
	"github_spiders/pkg/config"
	"github_spiders/spiders/types"
)

// SetupCollyRedis 初始化 Colly Redis.
func SetupCollyRedis() {
	storage := &redisstorage.Storage{
		Address: fmt.Sprintf("%s:%s",
			config.GetString("redis.github.host"),
			config.GetString("redis.github.port"),
		),
		Password: config.GetString("redis.github.password"),
		DB:       config.GetInt("redis.github.db"),
		Prefix:   config.GetString("redis.github.prefix"),
	}

	err := collectors.GetInstance(types.TagsRepo).SetStorage(storage)
	if err != nil {
		panic(err)
	}
}
