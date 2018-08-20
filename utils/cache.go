package utils

import (
	"gopkg.in/redis.v4"
	"time"
	"github.com/patrickmn/go-cache"
)

type Cache struct {
	Client *redis.Client
}

func GetCacheClient() *redis.Client {
	cache.New(5*time.Minute, 10*time.Minute)
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0, // use default DB
	})
	return client
}

func (c *Cache) SetIntValueInCache(key string, value int64, time time.Duration) error {
	return c.Client.Set(key, value, time).Err()
}

func (c *Cache) GetIntValueFromCache(key string) (int64, error) {
	return c.Client.Get(key).Int64()
}

func (c *Cache) Increment(key string, value int64) {
	c.Client.IncrBy(key, value)
}