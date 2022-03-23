package elastic

import (
	"github.com/olivere/elastic/v7"
	"sync"
)

var (
	Options       []elastic.ClientOptionFunc
	once          sync.Once
	elasticClient *elastic.Client
)

// GetInstance 获取一个 elastic 的客户端单例.
func GetInstance() *elastic.Client {
	once.Do(func() {
		var err error
		elasticClient, err = elastic.NewClient(Options...)
		if err != nil {
			panic(err)
		}
	})
	return elasticClient
}
