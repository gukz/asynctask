package redis

import (
	"github.com/go-redis/redis"
	"github.com/gukz/asynctask"
	"time"
)

func init() {
	asynctask.RegisteBackend("redis", func() asynctask.Backend { return &RedisBackend{} })
}

type RedisBackend struct {
	redisclient *redis.Client
}

func (b *RedisBackend) Init() error {
	host := "127.0.0.1:6379"
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
	return err == nil
}
func (b *RedisBackend) GetMessage(queue string) []byte {
	if msg, err := b.redisclient.BLPop(2*time.Second, queue).Result(); err != nil {
		if err != redis.Nil {
			panic(err)
		}
	} else if len(msg) > 1 {
		return ([]byte)(msg[1])
	}
	return nil
}
func (b *RedisBackend) CreateMessage(queue string, data []byte) error {
	res := b.redisclient.RPush(queue, string(data))
	return res.Err()
}
