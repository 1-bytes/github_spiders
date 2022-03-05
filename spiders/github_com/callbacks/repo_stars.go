package callbacks

import (
	"fmt"
	"github.com/gocolly/colly"
	"github_spiders/pkg/utils"
	"github_spiders/spiders/github_com/user"
	"log"
)

// RepoStars 项目主页.
type RepoStars struct {
	Colly *colly.Collector
	Index int
}

// Callbacks 爬取回调.
func (rs *RepoStars) Callbacks() {
	rs.Index = 0
	auth := user.NewAuth()
	rs.Colly.OnRequest(func(r *colly.Request) {
		// GitHub's docs:
		// By default, all requests to https://api.github.com receive the v3 version of the REST API.
		// We encourage you to explicitly request this version via the Accept header.
		r.Headers.Add("Accept", "application/vnd.github.v3+json")
		r.Headers = auth.AddToken(r.Headers, rs.Index)
	})

	rs.Colly.OnResponse(func(resp *colly.Response) {
		log.Printf("%d, %s", resp.StatusCode, resp.Body)
		// body := make(map[string]interface{})
		body, err := utils.JsonUnmarshalBody(resp.Body)
		// err := json.Unmarshal(resp.Body, &body)
		if err != nil {
			log.Printf("err: Failed to unmarshal the json. %s", err)
			return
		}

		for k, v := range body {
			fmt.Println(k, v)
		}
	})
}
