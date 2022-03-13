package types

import "github.com/gocolly/colly/v2"

const DefaultPage = 30 // 默认每页数据条数
const MaxPerPage = 100 // 每页数据最大条数
const TagsRepo = "repo"
const TagsUser = "user"

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

type ElasticIndexConfig struct {
	Index string
	Item  Item
}

type Item struct {
	RepoID        string
	RepoName      string
	RepoURL       string
	RepoApiURL    string
	ReadmeURL     string
	RepoStarCount uint64
	Readme        string
}
