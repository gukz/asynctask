package asynctask

/*
var brokers = make(map[string]func() Broker)

func RegisteBroker(key string, brokerFunc func() Broker) {
	brokers[key] = brokerFunc
}
*/

// Broker used to manage the task queue
type Broker interface {
	CheckHealth() bool
	PopMessage(string) []byte
	// TODO Add ack to the message
	AckMessage(string, string) error
	PushMessage(string, []byte) error
}
