package asynctask

type Producer struct {
	broker Broker
	queue  string
}

func NewProducer(brokerType string, queue string) (*Producer, error) {
	p := &Producer{queue: queue}
	if brokerFunc, ok := brokers[brokerType]; ok {
		p.broker = brokerFunc()
	} else {
		return nil, logger.Error("broker %s is not found", brokerType)
	}
	if err := p.broker.Init(); err != nil {
		return nil, err
	}
	return p, nil
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
		return t.broker.CreateMessage(t.queue, msg)
	}
}
