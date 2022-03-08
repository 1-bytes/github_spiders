package types

import "github.com/gocolly/colly"

const DefaultPerPage = 30 // 默认每页数据条数

type User struct {
	Token string
}

type GitHubUser struct {
	Users []User
}

type GitHubCollector struct {
	ReposByUserC *colly.Collector
	UsersByRepoC *colly.Collector
}

// CollyRedisConfig colly redis 配置信息.
type CollyRedisConfig struct {
	Address  string
	Password string
	DB       int
	Prefix   string
}
