package utils

import (
	//"gopkg.in/redis.v4"
	"time"
	"github.com/patrickmn/go-cache"
	"errors"
)

type LocalCache struct {
	Client *cache.Cache
}

//func GetCacheClient() *redis.Client {
//	return cache.New(5*time.Minute, 10*time.Minute)
//	client := redis.NewClient(&redis.Options{
//		Addr:     "127.0.0.1:6379",
//		Password: "", // no password set
//		DB:       0, // use default DB
//	})
//	return client
//}

func GetCacheClient() *cache.Cache {
	return cache.New(5*time.Minute, 10*time.Minute)
}


func (c *LocalCache) SetIntValueInCache(key string, value int64, time time.Duration)  {
	c.Client.Set(key, value, time)
}

func (c *LocalCache) GetIntValueFromCache(key string) (int64, error) {
	val, exist := c.Client.Get(key)
	if (exist) {
		return val.(int64), nil
	}
	return int64(0), errors.New("key not found")
}

func (c *LocalCache) Increment(key string, value int64) {
	c.Client.Increment(key, value)
}