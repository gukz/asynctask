package redis

import (
	"github.com/go-redis/redis"
	"github.com/gukz/asynctask"
	"time"
)

func init() {
	asynctask.RegisteBroker("redis", func() asynctask.Broker { return &RedisBroker{} })
}

type RedisBroker struct {
	redisclient *redis.Client
}

func (b *RedisBroker) Init() error {
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
func (b *RedisBroker) CheckHealth() bool {
	_, err := b.redisclient.Ping().Result()
	return err == nil
}
func (b *RedisBroker) GetMessage(queue string) []byte {
	if msg, err := b.redisclient.BLPop(2*time.Second, queue).Result(); err != nil {
		if err != redis.Nil {
			panic(err)
		}
	} else if len(msg) > 1 {
		return ([]byte)(msg[1])
	}
	return nil
}
func (b *RedisBroker) CreateMessage(queue string, data []byte) error {
	res := b.redisclient.RPush(queue, string(data))
	return res.Err()
}
