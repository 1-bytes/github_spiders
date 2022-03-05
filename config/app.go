package config

import "github_spiders/pkg/config"

func init() {
	config.Add("app", config.StrMap{
		// 应用名称
		"name": config.Env("APP_NAME", "spiders"),
		// 当前环境，用以区分多环境
		"env": config.Env("APP_ENV", "develop"),
		// APP 安全密钥，务必去创建一个自己的 GUID 作为密钥：https://www.guidgen.com
		"key": config.Env("APP_KEY", "b2581f25-99a2-4dd2-826d-753a6702903e"),
		// 是否开启调试模式
		"debug": config.Env("APP_DEBUG", false),
		// 日志存放目录
		"logger": config.Env("APP_LOGGER_PATH", "./runtime/logs"),
	})
}
