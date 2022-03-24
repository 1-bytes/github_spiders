package user

import (
	"github.com/gocolly/colly/v2"
	configs "github_spiders/pkg/config"
	"github_spiders/spiders/types"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type auth struct {
	Index             uint32
	users             types.GitHubUser
	lock              sync.Mutex
	LastSwitchingTime int64
}

var (
	AuthInstance *auth
	once         sync.Once
)

// NewAuth 创建一个新的 Auth 对象.
func NewAuth() *auth {
	once.Do(func() {
		AuthInstance = &auth{}
		AuthInstance.initAuth()
	})
	return AuthInstance
}

// Init 结构体初始化.
func (a *auth) initAuth() {
	a.Index = 0
	a.users = a.loadGitHubUsers()
}

// getGitHubUsers 获取 GitHub 账户信息.
func (a *auth) loadGitHubUsers() types.GitHubUser {
	// token 可以去 https://github.com/settings/tokens 进行创建
	tokenStr := configs.GetString("spiders.github.users")
	userSplit := strings.Split(tokenStr, ",")
	var users []types.User
	for _, v := range userSplit {
		users = append(users, types.User{Token: v})
	}
	return types.GitHubUser{Users: users}
}

// AddToken 在 header 里面增加授权 token.
func (a *auth) AddToken(header *http.Header) *http.Header {
	a.lock.Lock()
	defer a.lock.Unlock()
	maxLen := a.GetTokenCount()
	if maxLen == 0 || int(a.Index) > maxLen {
		return header
	}
	if a.Index < 0 {
		rand.Seed(time.Now().UnixNano())
		a.Index = uint32(rand.Intn(maxLen)) - 1
	}
	header.Add("Authorization", "token "+a.users.Users[a.Index].Token)
	return header
}

// DelToken 删除 header 里面的授权 token.
func (a *auth) DelToken(header *http.Header) *http.Header {
	header.Del("Authorization")
	return header
}

// GetTokenCount 获取账号的数量.
func (a *auth) GetTokenCount() int {
	return len(a.users.Users)
}

// NextToken 切换下一条 token.
func (a *auth) NextToken() {
	a.lock.Lock()
	defer a.lock.Unlock()
	currTime := time.Now().Unix()
	if 10 > currTime-a.LastSwitchingTime {
		return
	}
	a.LastSwitchingTime = currTime
	maxCount := a.GetTokenCount()
	atomic.AddUint32(&a.Index, 1) // 注意这里协程安全，不能直接 a.Index++
	if int(a.Index) >= maxCount { // 检查是否到最后一条 Token 了
		a.Index = 0
	}
	log.Printf("\n\n-------------------------------------------\n"+
		"[TOKEN] change token, new index: %d"+
		"\n-------------------------------------------\n\n", a.Index)
}

// CheckTokenValidity 检查 token 是否被封禁
// 返回的第一个参数如果是 true 则代表 token 有效，反之无效.
// 返回的第二个参数是一个字符型格式化的时间，会显示 token 的解封日期，当 token 没有被封禁时，该值为空
func (a *auth) CheckTokenValidity(resp *colly.Response) (bool, string) {
	if resp.Headers == nil {
		return true, ""
	}

	rateLimit := resp.Headers.Get("X-RateLimit-Remaining")
	rateLimitRetimeStr := resp.Headers.Get("X-RateLimit-Reset")
	if resp.StatusCode != http.StatusOK && rateLimit == "0" {
		rateLimitRetime, _ := strconv.ParseInt(rateLimitRetimeStr, 10, 64)
		tm := time.Unix(rateLimitRetime, 0)
		recoveryTime := tm.Format("2006-01-02 15:04:05")
		return false, recoveryTime
	}
	return true, ""
}
