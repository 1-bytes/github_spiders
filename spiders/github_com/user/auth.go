package user

import (
	configs "github_spiders/pkg/config"
	"github_spiders/spiders/types"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type auth struct {
	Index int
	users types.GitHubUser
}

var AuthInstance *auth

// NewAuth 创建一个新的 Auth 对象.
func NewAuth() *auth {
	if AuthInstance == nil { // 单例模式
		AuthInstance = &auth{}
		AuthInstance.initAuth()
	}
	return AuthInstance
}

// Init 结构体初始化.
func (a *auth) initAuth() {
	a.users = a.loadGitHubUsers()
}

// getGitHubUsers 获取 GitHub 账户信息.
func (a *auth) loadGitHubUsers() types.GitHubUser {
	// token 可以去 https://github.com/settings/tokens 进行创建
	tokenStr := configs.GetString("spiders.github.users")
	userSplit := strings.Split(tokenStr, ",")
	users := []types.User{}
	for _, v := range userSplit {
		users = append(users, types.User{Token: v})
	}
	return types.GitHubUser{Users: users}
}

// AddToken 在 header 里面增加授权 token.
func (a *auth) AddToken(header *http.Header, userIndex int) *http.Header {
	maxLen := len(a.users.Users)
	if maxLen == 0 || userIndex > maxLen {
		return header
	}
	if userIndex < 0 {
		rand.Seed(time.Now().UnixNano())
		userIndex = rand.Intn(maxLen) - 1
	}
	header.Add("Authorization", "token "+a.users.Users[userIndex].Token)
	return header
}
