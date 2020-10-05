package redis

import (
	"github.com/go-redis/redis"
	"github.com/gukz/asynctask"
	"time"
)

var _ asynctask.Broker = (*redisBroker)(nil)

type redisBroker struct {
	redisclient *redis.Client
	prefix      string
}

func NewBroker(client *redis.Client, prefix string) (asynctask.Broker, error) {
	b := &redisBroker{redisclient: client, prefix: prefix}
	_, err := b.redisclient.Ping().Result()
	return b, err
}

func (b *redisBroker) CheckHealth() bool {
	_, err := b.redisclient.Ping().Result()
	return err == nil
}

func (b *redisBroker) AckMessage(queue string, messageId uint64) error {
	return nil
}

func (b *redisBroker) PopMessage(queue string) (*asynctask.BrokerMessage, error) {
	if msg, err := b.redisclient.BLPop(2*time.Second, queue).Result(); err != nil {
		if err != redis.Nil {
			panic(err)
		}
	} else if len(msg) > 1 {
		return &asynctask.BrokerMessage{Body: ([]byte)(msg[1]), Id: 0}, nil
	}
	return nil, nil
}

func (b *redisBroker) PushMessage(queue string, data []byte) error {
	res := b.redisclient.RPush(queue, string(data))
	return res.Err()
}

func (b *redisBroker) Close() error {
	return b.redisclient.Close()
}
