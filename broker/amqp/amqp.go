package amqp

import (
	"github.com/gukz/asynctask"
)

func init() {
	asynctask.RegisteBroker("amqp", func() asynctask.Broker { return &AmqpBroker{} })
}

type AmqpBroker struct {
}

func (a *AmqpBroker) Init() error {
	return nil
}
func (a *AmqpBroker) CheckHealth() bool {
	return false
}
func (a *AmqpBroker) GetMessage(queue string) []byte {
	return nil
}
func (a *AmqpBroker) CreateMessage(queue string, data []byte) error {
	return nil
}
