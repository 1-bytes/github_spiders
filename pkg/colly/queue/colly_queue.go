package queue

import (
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	"sync"
)

var (
	collyQueueInstance map[int]*collyQueued
	lock               sync.Mutex
)

type collyQueued struct {
	Queue *queue.Queue
}

// init 初始化.
func init() {
	// 初始化单例的 Slice
	collyQueueInstance = make(map[int]*collyQueued)
}

// GetInstance 根据 ID 作为标识获取一个单例.
func GetInstance(id int) *collyQueued {
	if collyQueueInstance[id] != nil {
		return collyQueueInstance[id]
	}
	return func() *collyQueued {
		lock.Lock()
		defer lock.Unlock()
		if collyQueueInstance[id] == nil {
			collyQueueInstance[id] = &collyQueued{}
		}
		return collyQueueInstance[id]
	}()
}

// NewQueue 获取一个新的队列.
func (q *collyQueued) NewQueue(threads int, storage *redisstorage.Storage) {
	q.Queue, _ = queue.New(threads, storage)
}

// GetQueue 获取指定单例的队列.
func (q *collyQueued) GetQueue() *queue.Queue {
	return q.Queue
}
