package types

import "github.com/gocolly/colly/v2"

const DefaultPerPage = 30 // 默认每页数据条数
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
