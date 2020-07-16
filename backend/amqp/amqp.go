package amqp

import (
	"github.com/gukz/asynctask"
)

func init() {
	asynctask.RegisteBackend("amqp", func() asynctask.Backend { return &RedisBackend{} })
}

type AmqpBackend struct {
}

func (a *AmqpBackend) Init() error {
	return nil
}
func (a *AmqpBackend) CheckHealth() bool {
	return false
}
func (a *AmqpBackend) GetMessage(queue string) []byte {
	return nil
}
func (a *AmqpBackend) CreateMessage(queue string, data []byte) error {
	return nil
}
