package queued

import (
	"github.com/gocolly/colly/v2/queue"
	"github_spiders/pkg/storage"
	"sync"
)

var (
	instance map[string]*queue.Queue
	lock     sync.Mutex
)

func init() {
	instance = make(map[string]*queue.Queue)
}

// GetInstance 获取指定的队列..
func GetInstance(tag string) *queue.Queue {
	if instance[tag] != nil {
		return instance[tag]
	}
	return func(tag string) *queue.Queue {
		var err error
		lock.Lock()
		defer lock.Unlock()
		if instance[tag] == nil {
			redisStorage := storage.GetRedisStorage(tag)
			instance[tag], err = queue.New(1, &redisStorage)
			if err != nil {
				panic(err)
			}
		}
		return instance[tag]
	}(tag)
}
