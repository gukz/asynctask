package main

import (
	"fmt"
	"github.com/gukz/asynctask"
	rBackend "github.com/gukz/asynctask/backend/redis"
	rBroker "github.com/gukz/asynctask/broker/redis"
	"time"
)

func init() {
	asynctask.RegisteHandler(map[string]interface{}{
		"test_func": func(name string, a int64) error { fmt.Println("handling message", name, a, "\n"); return nil },
	})
}

func main() {
	redisHost := "127.0.0.1:6379"
	redisPassword := ""
	redisDbNum := 0
	queue := "Main_async_queue"
	backend, _ := rBackend.NewBackend(redisHost, redisPassword, redisDbNum, "", 10*time.Minute)
	broker, _ := rBroker.NewBroker(redisHost, redisPassword, redisDbNum)
	worker, _ := asynctask.NewWorker(queue, broker, backend)
	var a = 1
	go func() {
		for {
			msg := asynctask.NewMessage(
				"test_func", []asynctask.TypeValue{
					{Name: "name1", Type: "string", Value: "hello"},
					{Name: "age", Type: "int64", Value: a},
				})
			a += 1
			if err := worker.GetProducer().Send(msg); err != nil {
				fmt.Println(err, msg)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	worker.Serve(10)
}
