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

func NewBackend(client *redis.Client, keyPrefix string, ttl time.Duration) (asynctask.Backend, error) {
    t := &redisBackend{redisclient: client, keyPrefix: keyPrefix, backendTTL: ttl}
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

func (t *redisBackend) Close() error {
	return t.redisclient.Close()
}
