package bootstrap

//
// var storage *redisstorage.Storage
//
// // SetupCollyRedis 初始化 Colly Redis.
// func SetupCollyRedis() {
// 	storage = &redisstorage.Storage{
// 		Address:  config.GetString("redis.github.address"),
// 		Password: config.GetString("redis.github.password"),
// 		DB:       config.GetInt("redis.github.db"),
// 		Prefix:   config.GetString("redis.github.prefix"),
// 	}
// 	defer func(c *redis.Client) {
// 		_ = c.Close()
// 	}(storage.Client)
//
// 	err := GetCollector().SetStorage(storage)
// 	if err != nil {
// 		panic(err)
// 	}
// }
//
// // GetCollyRedisStorage 获取 CollyRedis.
// func GetCollyRedisStorage() *redisstorage.Storage {
// 	return storage
// }
