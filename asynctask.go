package asynctask

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	backends = make(map[string]func() Backend)
	handlers = make(map[string]interface{})
)

func init() {

	log.SetPrefix("ERROR LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func RegisteBackend(key string, backendFunc func() Backend) {
	backends[key] = backendFunc
}

func RegisteHandler(_handlers map[string]interface{}) {
	for k, v := range _handlers {
		if _, ok := handlers[k]; ok {
			log.Panic("Duplicated handler registered", k)
		} else {
			handlers[k] = v
		}
	}
}

type Backend interface {
	Init() error
	CheckHealth() bool
	GetMessage(string) []byte
	CreateMessage(string, []byte) error
}

type AsyncTask struct {
	backend Backend
	queue   string
}

func (a *AsyncTask) Init(queue string, backendType string) error {
	a.queue = queue
	if backendFunc, ok := backends[backendType]; ok {
		a.backend = backendFunc()
	} else {
		return errors.New(fmt.Sprintf("backend %s not found", backendType))
	}
	return a.backend.Init()
}

func (a *AsyncTask) CheckHealth() bool {
	return a.backend.CheckHealth()
}
func (a *AsyncTask) Serve(concurrency int) {
	quit := make(chan os.Signal, 1) // 什么时候使用缓冲区，什么时候不使用
	signal.Notify(quit, os.Interrupt, os.Kill)

	var wg sync.WaitGroup

	healthChecker := time.NewTicker(5 * time.Second)
	pool := make(chan struct{}, concurrency)
	errorMsg := make(chan interface{}, concurrency)
	wg.Add(1)
	go func() {
		for {
			select {
			case err, ok := <-errorMsg:
				if err == nil && !ok {
					wg.Done()
					return
				}
				fmt.Println(err)
			}
		}
	}()
	msgPool := make(chan []byte, concurrency)
	go func() {
		for {
			select {
			case <-quit:
				close(msgPool)
				return
			default:
				msg := a.backend.GetMessage(a.queue)
				if msg != nil {
					msgPool <- msg
				}
			}
		}
	}()

	var poolWg sync.WaitGroup
LOOP:
	for {
		select {
		case <-healthChecker.C:
			// health
			fmt.Println(fmt.Sprintf("%t is health\n", a.backend.CheckHealth()))
		case msg, ok := <-msgPool:
			if msg == nil && !ok {
				break LOOP
			}
			pool <- struct{}{}
			poolWg.Add(1)
			go func() {
				defer poolWg.Done()
				a.consume(msg, errorMsg)
				<-pool
			}()
		}
	}
	poolWg.Wait()
	close(errorMsg)
	wg.Wait()
	fmt.Println("grace period timeout")
}

func (a *AsyncTask) consume(message []byte, errorMsg chan interface{}) {
	if len(message) == 0 {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			errorMsg <- err
		}
	}()
	msg := Message{Args: []Arg{}}
	if err := decode(message, &msg); err != nil {
		panic(err)
	}
	if _, ok := handlers[msg.Name]; ok {
		// 反射调用
		fmt.Println("Handle", msg)
	} else {
		err := a.Produce(&msg)
		if err != nil {
			panic(err)
		}
	}
}

func (a *AsyncTask) Produce(message *Message) error {
	if msg, err := encode(message); err != nil {
		return err
	} else {
		return a.backend.CreateMessage(a.queue, msg)
	}
}
