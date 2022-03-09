package callbacks

import (
	"github.com/gocolly/colly/v2"
	"github_spiders/pkg/collectors"
	"github_spiders/pkg/utils"
	"github_spiders/spiders/github_com/common"
	"github_spiders/spiders/github_com/user"
	"github_spiders/spiders/types"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// ReposByUser 列出用户已加星标的存储库.
// GitHub API docs url:
// https://docs.github.com/cn/rest/reference/activity#list-repositories-starred-by-a-user
type ReposByUser struct {
	tokenIndex    uint32
	tokenMaxCount int
	lock          sync.Mutex
}

// Callbacks 爬虫回调函数.
func (ru *ReposByUser) Callbacks() {
	// 初始化变量
	ru.tokenIndex = 0
	auth := user.NewAuth()
	ru.tokenMaxCount = auth.GetTokenCount()
	collector := collectors.GetInstance(types.TagsRepo)

	// 对要提交的请求进行处理
	collector.OnRequest(func(r *colly.Request) {
		// GitHub's docs:
		// By default, all requests to https://api.github.com receive the v3 version of the REST API.
		// We encourage you to explicitly request this version via the Accept header.
		r.Headers.Add("Accept", "application/vnd.github.v3+json")
		r.Headers = auth.AddToken(r.Headers, int(ru.tokenIndex))
	})

	//  返回数据处理
	collector.OnResponse(func(resp *colly.Response) {
		repos, err := utils.JsonUnmarshalBody(resp.Body)
		if err != nil {
			log.Printf("err: Failed to unmarshal the json: %s", err)
			return
		}
		// 检查该用户有没有 star 的存储库，或者是 star 的存储库列表到了最终页
		repoLen := len(repos)
		if repoLen == 0 {
			return
		}

		for _, repo := range repos {
			repoName := repo["full_name"]
			starViewUrl := repo["stargazers_url"]
			starCount := repo["stargazers_count"]
			log.Printf("【New Repo】 Name:%s, StarCount:%v, URL:%s",
				repoName, starCount, starViewUrl)
			_ = collectors.GetInstance(types.TagsUser).Visit(starViewUrl.(string))
		}

		// 下一页
		url := common.GetNextPageUrl(resp.Request, repoLen)
		if url == "" {
			return
		}
		_ = resp.Request.Visit(url)
	})

	//  错误处理
	collector.OnError(func(resp *colly.Response, err error) {
		rateLimit := resp.Headers.Get("X-RateLimit-Remaining")
		rateLimitRetimeStr := resp.Headers.Get("X-RateLimit-Reset")
		rateLimitRetime, _ := strconv.ParseInt(rateLimitRetimeStr, 10, 64)
		tm := time.Unix(rateLimitRetime, 0)
		recoveryTime := tm.Format("2006-01-02 15:04:05")
		ru.lock.Lock()
		if resp.StatusCode != http.StatusOK && rateLimit == "0" {
			log.Printf(
				"Toktn or IP is temporarily blocked, unblock time is: %s Trying to change Token",
				recoveryTime,
			)
			atomic.AddUint32(&ru.tokenIndex, 1) // 注意这里协程安全，不能直接 ru.tokenIndex++
			if int(ru.tokenIndex) == ru.tokenMaxCount {
				ru.tokenIndex = 0
			}
			_ = resp.Request.Retry()
		}
		ru.lock.Unlock()
	})
}
