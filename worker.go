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
	*AsyncBase
}

func (t *worker) consume(message *BrokerMessage, errorMsg chan interface{}) {
	if len(message.Body) == 0 {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			errorMsg <- err
		}
	}()
	msg := Message{Args: []TypeValue{}}
	panicIf(decode(message.Body, &msg))
	if handler, ok := handlers[msg.Name]; ok {
		t.broker.AckMessage(t.queue, message.Id)
		handlerFunc := reflect.ValueOf(handler)
		params := TypeValue2ReflectValue(msg.Args)
		funcResult := handlerFunc.Call(params)
		if len(funcResult) == 0 {
			panicIf(errors.New("The last result of your handler must be error"))
		}
		result := &Result{}
		// The last result must be error type
		if !funcResult[len(funcResult)-1].IsNil() {
			result.HasError = true
			err := funcResult[len(funcResult)-1].Interface().(error)
			result.Error = err.Error()
		}
		result.ReturnValues = make([]*ResultValue, len(funcResult)-1)
		for i := 0; i < len(result.ReturnValues); i++ {
			val := funcResult[i].Interface()
			result.ReturnValues[i] = &ResultValue{
				Type:  reflect.TypeOf(val).String(),
				Value: val,
			}
		}
		if resultBytes, err := encode(result); err != nil {
			panicIf(err)
		} else {
			panicIf(t.backend.SetResult(msg.TaskId, resultBytes))
		}
	}
}

func (t *worker) CheckHealth() bool {
	return t.broker.CheckHealth() && t.backend.CheckHealth()
}

func (t *worker) Serve(concurrency int) {
	quit := make(chan os.Signal, 1) // really need chan here?
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
	msgPool := make(chan *BrokerMessage, concurrency)
	go func() {
		for {
			select {
			case <-quit:
				close(msgPool)
				return
			default:
				msg, err := t.broker.PopMessage(t.queue)
				if err != nil {
					errorMsg <- err
				} else if msg != nil {
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
