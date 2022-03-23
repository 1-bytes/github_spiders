package callbacks

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github_spiders/pkg/collectors"
	"github_spiders/pkg/encoding/base64"
	"github_spiders/pkg/queued"
	"github_spiders/pkg/utils"
	"github_spiders/spiders/github_com"
	"github_spiders/spiders/github_com/user"
	"github_spiders/spiders/types"
	"log"
	"net/http"
	"strconv"
)

const TagUser = "user"

// ReposByUser 列出用户已加星标的存储库.
// GitHub API docs url:
// https://docs.github.com/cn/rest/reference/activity#list-repositories-starred-by-a-user
type ReposByUser struct {
	BasicCallback
}

// Callbacks 爬虫回调函数.
func (r *ReposByUser) Callbacks() {
	// 初始化变量
	auth := user.NewAuth()
	collector := collectors.GetInstance(TagUser)

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
		repos := types.JsonRepos{}
		err := json.Unmarshal(resp.Body, &repos)
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
			// 有关仓库的一些信息
			readme, readmeAPI, err := r.GetReadme(repo.FullName)
			if err != nil {
				continue
			}
			repo.OtherData = struct {
				ReadmeAPI string `json:"readme_api"`
				Readme    string `json:"readme"`
			}{
				ReadmeAPI: readmeAPI,
				Readme:    readme,
			}

			id := strconv.FormatInt(repo.ID, 10)
			if err = r.SaveData("github_repos", id, repo); err != nil {
				log.Printf("Failed to store data in Elasticsearch, error: %s", err)
			} else {
				log.Printf("【New Repo】 Name:%s, URL:%s", repo.FullName, repo.StargazersURL)
			}
			_ = queued.GetInstance(TagRepo).AddURL(r.CheckUrl(repo.StargazersURL))
		}

		// 下一页
		if url := r.GetNextPageUrl(resp.Request, repoLen); url != "" {
			_ = queued.GetInstance(TagUser).AddURL(url)
		}
	})

	// 错误处理
	collector.OnError(func(resp *colly.Response, err error) {
		github_com.ErrLock.Lock()
		defer github_com.ErrLock.Unlock()
		instance := queued.GetInstance(TagUser)
		validity, t := auth.CheckTokenValidity(resp)
		if !validity {
			log.Printf("Toktn or IP is temporarily blocked, "+
				"unblock time is: %s Trying to change Token", t)
			auth.NextToken()
			auth.DelToken(resp.Request.Headers)
			_ = instance.AddRequest(resp.Request)
		}
		log.Println("Failed to initiate request, " +
			"task has been sent back to queue, waiting for retry.")
		_ = instance.AddRequest(resp.Request)
	})
}

// GetReadme 获取 Readme.md 文件的内容.
func (r *ReposByUser) GetReadme(fullName string) (string, string, error) {
	header := &http.Header{}
	auth := user.NewAuth()
	header.Add("Accept", "application/vnd.github.v3+json")
	header = auth.AddToken(header)

	url := fmt.Sprintf("https://api.github.com/repos/%s/readme", fullName)
	fetch, err := r.Fetcher(url, header)
	if err != nil {
		return "", "", err
	}
	body, err := utils.JsonUnmarshalBody(fetch)
	if err != nil {
		return "", "", err
	}

	var readme string
	if len(body) != 0 {
		var ok bool
		if _, ok = body[0]["content"].(string); ok {
			readme = body[0]["content"].(string)
			url = body[0]["url"].(string)
		}
	}
	return string(base64.Decoding(readme)), url, nil
}
