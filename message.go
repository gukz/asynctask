package asynctask

import (
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack"
)

type Arg struct {
	Name  string      `bson:"name"`
	Type  string      `bson:"type"`
	Value interface{} `bson:"value"`
}
type Message struct {
	TaskId string
	Name   string
	Args   []Arg
}

func NewMessage(name string, args []Arg) *Message {
	return &Message{
		TaskId: uuid.New().String(),
		Name:   name,
		Args:   args,
	}
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
