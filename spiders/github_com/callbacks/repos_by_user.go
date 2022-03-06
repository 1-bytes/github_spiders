package callbacks

import (
	"github.com/gocolly/colly"
	"github_spiders/pkg/utils"
	"github_spiders/spiders/github_com/user"
	"github_spiders/spiders/types"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

// ReposByUser 列出用户已加星标的存储库.
// GitHub API docs url:
// https://docs.github.com/cn/rest/reference/activity#list-repositories-starred-by-a-user
type ReposByUser struct {
	Colly types.GitHubCollector
	Index int
}

// Callbacks 爬虫回调函数.
func (ru *ReposByUser) Callbacks() {
	ru.Index = 0
	auth := user.NewAuth()
	collector := ru.Colly
	collector.ReposByUserC.OnRequest(func(r *colly.Request) {
		// GitHub's docs:
		// By default, all requests to https://api.github.com receive the v3 version of the REST API.
		// We encourage you to explicitly request this version via the Accept header.
		r.Headers.Add("Accept", "application/vnd.github.v3+json")
		r.Headers = auth.AddToken(r.Headers, ru.Index)
	})

	collector.ReposByUserC.OnResponse(func(resp *colly.Response) {
		if resp.StatusCode != http.StatusOK {
			// TODO:// 回头再处理这个问题，状态码为非200状态时，有可能是需要更换帐号的 token 了 ..
		}
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
			log.Printf("【New Repositorie】 Name:%s, URL:%s",
				repo["full_name"], repo["html_url"])
		}

		// 下一页
		url := ru.getNextPageUrl(resp.Request, repoLen)
		if url == "" {
			return
		}
		_ = collector.ReposByUserC.Visit(url)
	})
}

// getNextPageUrl 获取下一页要请求的链接.
func (ru *ReposByUser) getNextPageUrl(r *colly.Request, dataLen int) string {
	params := r.URL.Query()
	perPage, err := strconv.Atoi(params.Get("per_page"))
	if err != nil {
		perPage = types.DefaultPerPage
	}
	if dataLen < perPage {
		return ""
	}
	pageStr := params.Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, _ := strconv.ParseInt(pageStr, 10, 64)
	atomic.AddInt64(&page, 1)
	params.Set("page", strconv.FormatInt(page, 10))
	r.URL.RawQuery = params.Encode()
	return r.URL.String()
}
