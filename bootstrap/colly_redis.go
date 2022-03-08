package bootstrap

import (
	"github_spiders/pkg/config"
	"github_spiders/spiders/types"
)

var redisStorageCfg types.CollyRedisConfig

// SetupCollyRedis 读取 redis 配置.
func SetupCollyRedis() {
	redisStorageCfg = types.CollyRedisConfig{
		Address:  config.GetString("redis.github.address"),
		Password: config.GetString("redis.github.password"),
		DB:       config.GetInt("redis.github.db"),
		Prefix:   config.GetString("redis.github.prefix"),
	}
}

// GetCollyRedisConfig 返回 colly redis 配置.
func GetCollyRedisConfig() types.CollyRedisConfig {
	return redisStorageCfg
}
