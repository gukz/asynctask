package redis

import (
	"github.com/go-redis/redis"
	"github.com/gukz/asynctask"
	"time"
)

func init() {
	asynctask.RegisteBackend("redis", &RedisBackend{})
}

type RedisBackend struct {
	redisclient *redis.Client
}

func (b *RedisBackend) Init() error {
	host := "127.0.0.1"
	password := ""
	dbNum := 0
	b.redisclient = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       dbNum,
	})
	_, err := b.redisclient.Ping().Result()
	return err
}
func (b *RedisBackend) CheckHealth() bool {
	_, err := b.redisclient.Ping().Result()
	return err != nil
}
func (b *RedisBackend) GetMessage() []string {
	queue := "gukz_asynctask"
	if msg, err := b.redisclient.BLPop(10*time.Second, queue).Result(); err != nil {
		panic(err)
	} else {
		return msg
	}
}
