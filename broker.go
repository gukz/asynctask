package asynctask

// Broker used to manage the task queue
type Broker interface {
	CheckHealth() bool
	PopMessage(string) []byte
	AckMessage(string, string) error
	PushMessage(string, []byte) error
}
