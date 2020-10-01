package asynctask

type Producer struct {
	broker Broker
	queue  string
}

func NewProducer(queue string, broker Broker) *Producer {
	return &Producer{queue: queue, broker: broker}
}

func (t *Producer) CheckHealth() bool {
	if t.broker == nil {
		logger.Warn("The Producer's broker is nil")
		return false
	}
	return t.broker.CheckHealth()
}

func (t *Producer) Send(message *Message) error {
	if msg, err := encode(message); err != nil {
		return err
	} else {
		return t.broker.PushMessage(t.queue, msg)
	}
}
