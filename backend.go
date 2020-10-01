package asynctask

/*
var backends = make(map[string]func() Backend)

func RegisteBackend(key string, backendFunc func() Backend) {
	backends[key] = backendFunc
}
*/

// Backend used to manage task's status and result
type Backend interface {
	CheckHealth() bool
	SetResult(taskId string, result interface{}) error
	GetResult(taskId string) (interface{}, error)
}
