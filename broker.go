package asynctask

type BrokerMessage struct {
	Body []byte
	Id   uint64
}

// Broker used to manage the task queue
type Broker interface {
	CheckHealth() bool
	PopMessage(string) (*BrokerMessage, error)
	AckMessage(string, uint64) error
	PushMessage(string, []byte) error
	Close() error
}
