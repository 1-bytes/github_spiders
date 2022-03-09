package types

import "github.com/gocolly/colly/v2"

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
