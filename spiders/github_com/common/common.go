package common

import (
	"github.com/gocolly/colly/v2"
	"github_spiders/spiders/types"
	"net/url"
	"strconv"
	"sync/atomic"
)

func CheckUrl(u string) string {
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
func GetNextPageUrl(r *colly.Request, dataLen int) string {
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
