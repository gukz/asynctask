package redis

import (
	"github.com/gukz/asynctask"
)

func init() {
	asynctask.RegisteBackend("redis", func() asynctask.Backend { return &RedisBackend{} })
}

type RedisBackend struct {
}

func (t *RedisBackend) Init() error {
	return nil
}
func (t *RedisBackend) CheckHealth() bool {
	return true
}
func (t *RedisBackend) SetResult(result interface{}) error {
	return nil
}
