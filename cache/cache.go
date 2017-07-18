package cache

import (
	"log"
	"strconv"

	"github.com/go-redis/redis"
)

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

func (c Cache) GetUserId(session string) (int, error) {
	val, err := c.client.Get(session).Result()

	if err != nil {
		return -1, err
	}

	return strconv.Atoi(val)
}
