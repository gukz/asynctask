package asynctask

var brokers = make(map[string]func() Broker)

func RegisteBroker(key string, brokerFunc func() Broker) {
	brokers[key] = brokerFunc
}

// Broker used to manage the task queue
type Broker interface {
	Init() error
	CheckHealth() bool
	GetMessage(string) []byte
	CreateMessage(string, []byte) error
}
