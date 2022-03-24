package callbacks

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"github_spiders/pkg/collectors"
	"github_spiders/pkg/queued"
	"github_spiders/spiders/github_com"
	"github_spiders/spiders/github_com/user"
	"github_spiders/spiders/types"
	"log"
	"strconv"
	"time"
)

const TagRepo = "repo"

// UsersByRepo 列出已为存储库加注星标的人员.
// GitHub API docs url:
// https://docs.github.com/cn/rest/reference/activity#list-stargazers
type UsersByRepo struct {
	BasicCallback
}

// Callbacks 爬虫回调函数.
func (u *UsersByRepo) Callbacks() {
	// 初始化变量
	auth := user.NewAuth()
	collector := collectors.GetInstance(TagRepo)

	// 对要提交的请求进行处理
	collector.OnRequest(func(r *colly.Request) {
		// GitHub's docs:
		// By default, all requests to https://api.github.com receive the v3 version of the REST API.
		// We encourage you to explicitly request this version via the Accept header.
		r.Headers.Add("Accept", "application/vnd.github.v3+json")
		r.Headers = auth.AddToken(r.Headers)
	})

	//  返回数据处理
	collector.OnResponse(func(resp *colly.Response) {
		users := types.JsonUsers{}
		err := json.Unmarshal(resp.Body, &users)
		if err != nil {
			log.Printf("err: Failed to unmarshal the json: %s", err)
			return
		}

		// 检查该用户有没有 star 的存储库，或者是 star 的存储库列表到了最终页
		userLen := len(users)
		if userLen == 0 {
			return
		}

		for _, user := range users {
			// 有关用户的一些信息，想要什么值可以自己取
			id := strconv.FormatInt(user.ID, 10)
			if err = u.SaveData("github_users", id, user); err != nil {
				log.Printf("Failed to store data in Elasticsearch, error: %s", err)
			} else {
				time.Sleep(time.Second)
				log.Printf("【New User】 Name:%s, URL:%s", user.Login, user.URL+"/starred")
			}
			_ = queued.GetInstance(TagUser).AddURL(u.CheckUrl(user.URL + "/starred"))
		}

		// 下一页
		if url := u.GetNextPageUrl(resp.Request, userLen); url != "" {
			_ = queued.GetInstance(TagRepo).AddURL(url)
		}
	})

	// 错误处理
	collector.OnError(func(resp *colly.Response, err error) {
		github_com.ErrLock.Lock()
		defer github_com.ErrLock.Unlock()
		instance := queued.GetInstance(TagRepo)
		validity, t := auth.CheckTokenValidity(resp)
		if !validity {
			log.Printf("Token[%d] or IP is temporarily blocked, "+
				"unblock time is: %s Trying to change Token", auth.Index, t)
			auth.NextToken()
			auth.DelToken(resp.Request.Headers)
			_ = instance.AddRequest(resp.Request)
			return
		}
		log.Println("Failed to initiate request, " +
			"task has been sent back to queue, waiting for retry.")
		_ = instance.AddRequest(resp.Request)
	})
}
