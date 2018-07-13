package cache

import "github.com/go-redis/redis"

var redisAddr string
var redisPass string

func InitRedis(addr string, pass string) {
	redisAddr = addr
	redisPass = pass
}

func NewRedisClient(db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       db,
	})

	return client
}
