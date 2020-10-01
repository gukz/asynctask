package amqp

import (
	"github.com/gukz/asynctask"
)

var _ asynctask.Broker = (*amqpBroker)(nil)

type amqpBroker struct {
}

func NewBroker() (asynctask.Broker, error) {
	a := &amqpBroker{}
	return a, nil
}

func (a *amqpBroker) CheckHealth() bool {
	return false
}

func (a *amqpBroker) AckMessage(queue string, taskId string) error {
	return nil
}

func (a *amqpBroker) PopMessage(queue string) []byte {
	return nil
}

func (a *amqpBroker) PushMessage(queue string, data []byte) error {
	return nil
}
