package asynctask

import (
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack"
	"log"
)

var logger = &Logger{}

func init() {
	log.SetPrefix("ERROR LOG: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

type Logger struct {
}

func (t *Logger) Info(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (t *Logger) Warn(format string, args ...interface{}) {
	t.Info(format, args...)
}

func (t *Logger) Error(format string, args ...interface{}) error {
	t.Info(format, args...)
	return errors.New(fmt.Sprintf(format, args...))
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

type AsyncBase struct {
	queue   string
	backend Backend
	broker  Broker
}

func NewAsyncTask(queue string, broker Broker, backend Backend) *AsyncBase {
	return &AsyncBase{queue: queue, broker: broker, backend: backend}
}

func (t *AsyncBase) GetWorker() *worker {
	return &worker{AsyncBase: t}
}

func (t *AsyncBase) GetProducer() *producer {
	return &producer{AsyncBase: t}
}

func encode(value interface{}) ([]byte, error) {
	jsonBytes, err := msgpack.Marshal(&value)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func decode(data []byte, out interface{}) error {
	return msgpack.Unmarshal(data, out)
}
