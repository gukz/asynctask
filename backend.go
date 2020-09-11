package asynctask

var backends = make(map[string]func() Backend)

func RegisteBackend(key string, backendFunc func() Backend) {
	backends[key] = backendFunc
}

// Backend used to manage task's status and result
type Backend interface {
	Init() error
	CheckHealth() bool
	SetResult(result interface{}) error
}
