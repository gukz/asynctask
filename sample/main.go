package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gukz/asynctask"
	rBackend "github.com/gukz/asynctask/backend/redis"
	mqBroker "github.com/gukz/asynctask/broker/amqp"
	rBroker "github.com/gukz/asynctask/broker/redis"
	"time"
)

func init() {
	asynctask.RegisteHandler(map[string]interface{}{
		"test_func": func(name string, a int64) (string, int64, error) {
			return name + "is finished", a * a, nil
		},
	})
}

func main() {
	test_mq()
	test_redis()
}

func test_mq() {
	queue := "Main_async_queue"
	broker, err := mqBroker.NewBroker("", "", "", "")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	backend, _ := rBackend.NewBackend(client, "backend_", 10*time.Minute)
	asyncBase := asynctask.NewAsyncTask(queue, broker, backend)
	test_base(asyncBase)
}

func test_redis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	queue := "Main_async_queue"
	backend, _ := rBackend.NewBackend(client, "backend_", 10*time.Minute)
	broker, _ := rBroker.NewBroker(client, "broker_")
	asyncBase := asynctask.NewAsyncTask(queue, broker, backend)
	test_base(asyncBase)
}

func test_base(asyncBase *asynctask.AsyncBase) {
	worker := asyncBase.GetWorker()
	producer := asyncBase.GetProducer()
	var a = 1
	go func() {
		for {
			msg := asynctask.NewMessage(
				"test_func", []asynctask.TypeValue{
					{Name: "name1", Type: "string", Value: "hello"},
					{Name: "age", Type: "int64", Value: a},
				})
			a += 1
			if err := producer.Send(msg); err != nil {
				fmt.Println(err, msg)
			}
			time.Sleep(2 * time.Second)
			if res, err := producer.GetResult(msg.TaskId); err != nil {
				fmt.Println(err, *res)
			} else {
				fmt.Println("\nResult: ")
				for i := 0; i < len(res.ReturnValues); i++ {
					fmt.Print(res.ReturnValues[i].Type, ": ", res.ReturnValues[i].Value, " ")
				}
			}
		}
	}()

	worker.Serve(10)
}
