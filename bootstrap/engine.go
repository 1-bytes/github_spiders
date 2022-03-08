package bootstrap

import (
	"github_spiders/pkg/colly/queue"
	"github_spiders/pkg/colly/redis"
	"github_spiders/pkg/config"
	"github_spiders/spiders/types"
)

func SetupEngine(c types.GitHubCollector) {
	cfg := types.CollyRedisConfig{
		Address:  config.GetString("redis.github.address"),
		Password: config.GetString("redis.github.password"),
		DB:       config.GetInt("redis.github.db"),
		Prefix:   config.GetString("redis.github.prefix"),
	}
	redisStorage := redis.NewRedisStorage(c.UsersByRepoC, cfg)
	queue.GetInstance(int(c.UsersByRepoC.ID)).NewQueue(3, redisStorage)
	redisStorage = redis.NewRedisStorage(c.ReposByUserC, cfg)
	queue.GetInstance(int(c.ReposByUserC.ID)).NewQueue(3, redisStorage)
	SetupCallback(c)
}
