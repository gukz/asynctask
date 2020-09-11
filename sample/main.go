package main

import (
	"fmt"
	"github.com/gukz/asynctask"
	_ "github.com/gukz/asynctask/backend/redis"
	_ "github.com/gukz/asynctask/broker/redis"
	"time"
)

func init() {
	asynctask.RegisteHandler(map[string]interface{}{
		"test_func": func(name string, a int64) error { fmt.Println("handling message", name, a, "\n"); return nil },
	})
}

func main() {
	queue := "async_queue"
	worker, _ := asynctask.NewWorker("redis", "redis", queue)
	var a = 1
	go func() {
		for {
			msg := asynctask.NewMessage(
				"test_func", []asynctask.Arg{
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
