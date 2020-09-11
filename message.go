package asynctask

import (
	"github.com/google/uuid"
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
