package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/olivere/elastic/v7"
	"github_spiders/pkg/collectors"
	"github_spiders/pkg/encoding/base64"
	"github_spiders/pkg/queued"
	"github_spiders/pkg/utils"
	"github_spiders/spiders/github_com"
	"github_spiders/spiders/github_com/common"
	"github_spiders/spiders/github_com/user"
	"github_spiders/spiders/types"
	"log"
	"net/http"
)

const TagUser = "user"

// ReposByUser 列出用户已加星标的存储库.
// GitHub API docs url:
// https://docs.github.com/cn/rest/reference/activity#list-repositories-starred-by-a-user
type ReposByUser struct {
}

// Callbacks 爬虫回调函数.
func (ru *ReposByUser) Callbacks() {
	// 初始化变量
	auth := user.NewAuth()
	collector := collectors.GetInstance(TagRepo)
	// elasticClient, err := elastic.NewClient(config.ElasticOptions...)
	// if err != nil {
	// 	panic(err)
	// }

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

		var (
			// repoID        string
			repoName string
			// repoURL       string
			// repoApiURL    string
			repoStarURL string
			// repoStarCount uint64
		)
		for _, repo := range repos {
			// 有关仓库的一些信息，想要什么值可以自己取
			// repoID = strconv.FormatUint(uint64(repo["id"].(float64)), 10)
			repoName = repo["full_name"].(string)
			// repoURL = repo["html_url"].(string)
			// repoApiURL = repo["url"].(string)
			repoStarURL = repo["stargazers_url"].(string)
			// repoStarCount = uint64(repo["stargazers_count"].(float64))
			log.Printf("【New Repo】 Name:%s, URL:%s", repoName, repoStarURL)

			// // 获取 Readme.md 的数据并存储
			// readme, url, _ := GetReadme(repoName)
			// err = SaveData(elasticClient, types.ElasticIndexConfig{
			// 	Index: "github_readme",
			// 	Item: types.Item{
			// 		RepoID:        repoID,
			// 		RepoName:      repoName,
			// 		RepoURL:       repoURL,
			// 		RepoApiURL:    repoApiURL,
			// 		RepoStarCount: repoStarCount,
			// 		ReadmeURL:     url,
			// 		Readme:        readme,
			// 	},
			// })
			// if err != nil {
			// 	log.Printf("Failed to store data in Elasticsearch, error: %s", err)
			// }

			repoStarURL = common.CheckUrl(repoStarURL)
			// _ = collectors.GetInstance(TagUser).Visit(repoStarURL)
			_ = queued.GetInstance(TagUser).AddURL(repoStarURL)
		}

		// 下一页
		url := common.GetNextPageUrl(resp.Request, repoLen)
		if url == "" {
			return
		}
		// _ = resp.Request.Visit(url)
		_ = queued.GetInstance(TagRepo).AddURL(url)
	})

	// 错误处理
	collector.OnError(func(resp *colly.Response, err error) {
		github_com.ErrLock.Lock()
		defer github_com.ErrLock.Unlock()
		validity, t := auth.CheckTokenValidity(resp)
		if !validity {
			log.Printf("Toktn or IP is temporarily blocked, "+
				"unblock time is: %s Trying to change Token", t)
			auth.NextToken()
			auth.DelToken(resp.Request.Headers)
			_ = resp.Request.Retry()
		}
	})
}

// SaveData 保存数据.
func SaveData(client *elastic.Client, cfg types.ElasticIndexConfig) error {
	return nil
	data, err := json.Marshal(cfg.Item)
	if err != nil {
		return err
	}
	_, err = client.Index().
		Index(cfg.Index).
		BodyJson(string(data)).
		Id(cfg.Item.RepoID).Do(context.Background())
	return err
}

// GetReadme 获取 Readme.md 文件的内容.
func GetReadme(fullName string) (string, string, error) {
	header := &http.Header{}
	auth := user.NewAuth()
	header.Add("Accept", "application/vnd.github.v3+json")
	header = auth.AddToken(header)

	url := fmt.Sprintf("https://api.github.com/repos/%s/readme", fullName)
	fetch, err := common.Fetcher(url, header)
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
