package asynctask

import (
	"fmt"
	"sync"
)

var backends map[string]Backend

func RegisteBackend(key string, backend Backend) {
	backends[key] = backend
}

type Backend interface {
	Init() error
	CheckHealth() bool
	GetMessage() []string
}

func Serve(backendType string, concurrency int) {
	quit := make(chan struct{}, 1) // 什么时候使用缓冲区，什么时候不使用
	pool := make(chan struct{}, 10)
	errorMsg := make(chan string, 10)
	var wg sync.WaitGroup
	backend := backends[backendType]
	if err := backend.Init(); err != nil {
		panic(err)
	}
	for {
		select {
		case <-quit:
			goto STOP
		default:
			pool <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				workFunc(backend.GetMessage(), errorMsg)
				<-pool
			}()
		}
		fmt.Sprintf("%t is health\n", backend.CheckHealth())
		// TODO beat
	}
STOP:
	wg.Wait()
}

func workFunc(message []string, errorMsg chan string) {
	fmt.Println(message)
}
