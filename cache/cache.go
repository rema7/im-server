package cache

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

var (
	instance *Cache
	once     sync.Once
)

func GetCache() *Cache {
	once.Do(func() {
		redis := Cache{}
		redis.Init()
		instance = &redis
	})
	return instance
}

type Cache struct {
	client *redis.Client
}

func (c *Cache) Init() {
	c.client = redis.NewClient(&redis.Options{
		Addr:     "localhost:7379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := c.client.Ping().Result()
	if err != nil {
		log.Panic(err)
	}
}

func (c Cache) GetUserId(session string) (int64, error) {
	val, err := c.client.Get(session).Result()

	if err != nil {
		return -1, err
	}

	return strconv.ParseInt(val, 10, 64)
}

func (c Cache) GetChats(userId string) ([]int, error) {
	val, err := c.client.Get(userId).Result()

	if err != nil {
		return nil, err
	}

	var ints []int
	json.Unmarshal([]byte(val), &ints)
	return ints, nil
}
