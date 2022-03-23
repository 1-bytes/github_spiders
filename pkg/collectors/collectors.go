package collectors

import (
	"github.com/gocolly/colly/v2"
	"github_spiders/pkg/storage"
	"github_spiders/spiders/github_com"
	"sync"
)

var (
	collectorInstance map[string]*colly.Collector
	lock              sync.Mutex
	cfg               *github_com.Spider
)

// SetConfig 设置 collector 配置.
func SetConfig(config *github_com.Spider) {
	cfg = config
}

func CloneToTag(tag string, toTag string) {
	lock.Lock()
	defer lock.Unlock()
	collectorInstance[toTag] = collectorInstance[tag].Clone()
}

func init() {
	collectorInstance = make(map[string]*colly.Collector)
}

// GetInstance 获取指定的 collector (伪单例模式)..
func GetInstance(tag string) *colly.Collector {
	if collectorInstance[tag] != nil {
		return collectorInstance[tag]
	}
	return func(tag string) *colly.Collector {
		var err error
		lock.Lock()
		defer lock.Unlock()
		if collectorInstance[tag] == nil {
			collectorInstance[tag] = cfg.Create()
			redisStorage := storage.GetRedisStorage(tag)
			err = collectorInstance[tag].SetStorage(&redisStorage)
			if err != nil {
				panic(err)
			}
		}
		return collectorInstance[tag]
	}(tag)
}
