package redis

import (
	"github.com/go-redis/redis"
	"github.com/gukz/asynctask"
	"time"
)

var _ asynctask.Broker = (*redisBroker)(nil)

type redisBroker struct {
	redisclient *redis.Client
}

func NewBroker(host string, password string, dbNum int) (asynctask.Broker, error) {
	b := &redisBroker{}
	b.redisclient = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       dbNum,
	})
	_, err := b.redisclient.Ping().Result()
	return b, err
}

func (b *redisBroker) CheckHealth() bool {
	_, err := b.redisclient.Ping().Result()
	return err == nil
}

func (b *redisBroker) AckMessage(queue string, taskId string) error {
	return nil
}

func (b *redisBroker) PopMessage(queue string) []byte {
	if msg, err := b.redisclient.BLPop(2*time.Second, queue).Result(); err != nil {
		if err != redis.Nil {
			panic(err)
		}
	} else if len(msg) > 1 {
		return ([]byte)(msg[1])
	}
	return nil
}

func (b *redisBroker) PushMessage(queue string, data []byte) error {
	res := b.redisclient.RPush(queue, string(data))
	return res.Err()
}
