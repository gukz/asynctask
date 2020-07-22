package main

import (
	"fmt"
	"github.com/gukz/asynctask"
	_ "github.com/gukz/asynctask/backend/redis"
	"time"
)

func init() {
	asynctask.RegisteHandler(map[string]interface{}{
		"handler1": func(name string, a int64) error { fmt.Println(name, a); return nil },
	})
}

func main() {
	queue := "async_queue"
	async := &asynctask.AsyncTask{}
	async.Init(queue, "redis")
	go func() {
		for i := 0; i < 5; i++ {
			msg := asynctask.NewMessage("handler1", []asynctask.Arg{{Name: "name1", Type: "string", Value: "hello"}, {Name: "age", Type: "int64", Value: i}})
			if err := async.Produce(msg); err != nil {
				fmt.Println(err, msg)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	async.Serve(10)
}
