package asynctask

// Backend used to manage task's status and result
type Backend interface {
	CheckHealth() bool
	SetResult(taskId string, result []byte) error
	GetResult(taskId string) ([]byte, error)
	Close() error
}
