package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/gukz/asynctask"
)

var _ asynctask.Backend = (*redisBackend)(nil)

type redisBackend struct {
	redisclient *redis.Client
	keyPrefix   string
	backendTTL  time.Duration
}

func NewBackend(host string, password string, dbNum int, keyPrefix string, ttl time.Duration) (asynctask.Backend, error) {
	t := &redisBackend{keyPrefix: keyPrefix, backendTTL: ttl}
	t.redisclient = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       dbNum,
	})
	_, err := t.redisclient.Ping().Result()
	return t, err
}

func (t *redisBackend) CheckHealth() bool {
	_, err := t.redisclient.Ping().Result()
	return err == nil
}
func (t *redisBackend) SetResult(taskId string, result []byte) error {
	res := t.redisclient.Set(t.keyPrefix+taskId, result, t.backendTTL)
	return res.Err()
}

func (t *redisBackend) GetResult(taskId string) ([]byte, error) {
	res := t.redisclient.Get(t.keyPrefix + taskId)
	return res.Bytes()
}
