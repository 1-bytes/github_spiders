package callbacks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github_spiders/pkg/elastic"
	"github_spiders/spiders/types"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
)

var Tags = map[string]string{
	TagRepo: "repo",
	TagUser: "user",
}

type BasicCallback struct {
}

func (*BasicCallback) CheckUrl(u string) string {
	parse, _ := url.Parse(u)
	params, _ := url.ParseQuery(parse.RawQuery)
	if params.Has("per_page") {
		params.Del("per_page")
	}
	params.Add("per_page", strconv.Itoa(types.MaxPerPage))
	if !params.Has("page") {
		params.Add("page", "1")
	}
	parse.RawQuery = params.Encode()
	return parse.String()
}

// GetNextPageUrl 获取下一页要请求的链接.
func (*BasicCallback) GetNextPageUrl(r *colly.Request, dataLen int) string {
	params := r.URL.Query()
	perPage, err := strconv.Atoi(params.Get("per_page"))
	if err != nil {
		perPage = types.DefaultPage
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

// Fetcher 获取指定网页内容.
func (*BasicCallback) Fetcher(url string, headers *http.Header) ([]byte, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = *headers
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"the returned status code is not the expected status: %d",
			resp.StatusCode,
		)
	}
	return ioutil.ReadAll(resp.Body)
}

// SaveData 保存数据.
func (*BasicCallback) SaveData(index string, id string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = elastic.GetInstance().Index().
		Index(index).
		BodyJson(string(j)).
		Id(id).Do(context.Background())
	return err
}
