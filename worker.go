package asynctask

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"time"
)

var handlers = make(map[string]interface{})

func RegisteHandler(h map[string]interface{}) {
	for k, v := range h {
		if _, ok := handlers[k]; ok {
			panic(fmt.Sprintf("Duplicated handler registered: %s", k))
		} else {
			handlers[k] = v
		}
	}
}

type worker struct {
	broker  Broker
	backend Backend
	queue   string
}

func NewWorker(queue string, broker Broker, backend Backend) (*worker, error) {
	return &worker{queue: queue, broker: broker, backend: backend}, nil
}

func (t *worker) GetProducer() *Producer {
	return &Producer{broker: t.broker, queue: t.queue}
}

func (t *worker) consume(message []byte, errorMsg chan interface{}) {
	if len(message) == 0 {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			errorMsg <- err
		}
	}()
	msg := Message{Args: []TypeValue{}}
	panicIf(decode(message, &msg))
	t.broker.AckMessage(t.queue, msg.TaskId)
	if handler, ok := handlers[msg.Name]; ok {
		handlerFunc := reflect.ValueOf(handler)
		if handlerFunc.Kind() != reflect.Func {
			panic(errors.New("Invalid handler type"))
		}
		params := TypeValue2ReflectValue(msg.Args)
		result := handlerFunc.Call(params)
		if len(result) == 0 {
			panic(errors.New("The first result of your handler must be error"))
		}
		panicIf(t.backend.SetResult(msg.TaskId, result))
	} else {
		// TODO: If this message is not recognize, we will throw back to queue
		panicIf(t.GetProducer().Send(&msg))
	}
}

func (t *worker) CheckHealth() bool {
	return t.broker.CheckHealth() && t.backend.CheckHealth()
}

func (t *worker) Serve(concurrency int) {
	quit := make(chan os.Signal, 1) // 什么时候使用缓冲区，什么时候不使用
	signal.Notify(quit, os.Interrupt, os.Kill)

	var wg sync.WaitGroup

	healthChecker := time.NewTicker(15 * time.Second)
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
				logger.Error("%s", err)
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
				msg := t.broker.PopMessage(t.queue)
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
			logger.Info("worker Check Health Result: %t", t.CheckHealth())
		case msg, ok := <-msgPool:
			if !ok {
				break LOOP
			}
			pool <- struct{}{}
			poolWg.Add(1)
			go func() {
				defer poolWg.Done()
				t.consume(msg, errorMsg)
				<-pool
			}()
		}
	}
	poolWg.Wait()
	close(errorMsg)
	wg.Wait()
}
