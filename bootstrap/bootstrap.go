package bootstrap

import configs "github_spiders/config"

// Setup 初始化指定的服务，例如：Redis MySQL Logger 等模块.
func Setup() {
	autoLoader(
		configs.Initialize, // 配置文件
		SetupCollector,     // 初始化 Colly
		// SetupCollyRedis,    // 初始化 Colly.Redis(依赖于 collector)
		// SetupCollyQueue,    // 初始化 CollyQueue(依赖于 RedisStorage)
	)
}

// autoLoader 自动加载初始化.
func autoLoader(funcName ...func()) {
	for _, v := range funcName {
		v()
	}
}
