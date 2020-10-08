## async task client in go
This project is a package used to support async task client of golang, support multi backend and broker, include `Redis`, `AMQP` e.g.

## quick start
Here is a sample code to help you quick start, the api is easily to use, enjoy it!

## sample code
```go
package test

import (
	"fmt"
	"github.com/go-redis/redis"
	rBackend "github.com/gukz/asynctask/backend/redis"
	rBroker "github.com/gukz/asynctask/broker/redis"
)
func init() {
	asynctask.RegisteHandler(map[string]interface{}{
		"async_func": func(name string, a int64) (string, int64, error) {
			fmt.Println("function is running", name, a)
			return name + "is finished", a * a, nil
		},
	})
}

func main(){
	// defin a redis client
	client := redis.NewClient(&redis.Options{
		Addr:	  "127.0.0.1:6379",
		Password: "",
		DB:	      0,
	})

	backend, _ := rBackend.NewBackend(client, "backend_", 10*time.Minute)
	broker, _ := rBroker.NewBroker(client, "broker_")

	// create a async base object
	asyncBase := asynctask.NewAsyncTask("test_async_queue", broker, backend)

	go func() {
		var a = 1
		for {
			// build the async call
			msg := asynctask.NewMessage(
				"test_func", []asynctask.TypeValue{
					{Name: "name1", Type: "string", Value: "hello"},
					{Name: "age", Type: "int64", Value: a++},
				})
			// send this async call
			_ = asyncBase.GetProducer().Send(msg)
			// wait it to be process
			time.Sleep(2 * time.Second)
			// check the result of the async call
			res, _ := asyncBase.GetProducer().GetResult(msg.TaskId)
		}
	}()

	// start handle async call with concurrent of 10.
	asyncBase.GetWorker().Serve(10)
}
```

## more feature comes later...
- support middleware(used to monitor task failure).
- support task with context.
- support more `backend` and `broker`.
