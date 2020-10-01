package asynctask

import (
	"github.com/google/uuid"
	"reflect"
)

type TypeValue struct {
	Name  string      `bson:"name"`
	Type  string      `bson:"type"`
	Value interface{} `bson:"value"`
}
type Message struct {
	TaskId string
	Name   string
	Args   []TypeValue
}

func NewMessage(name string, args []TypeValue) *Message {
	return &Message{
		TaskId: uuid.New().String(),
		Name:   name,
		Args:   args,
	}
}

func TypeValue2ReflectValue(data []TypeValue) []reflect.Value {
	res := make([]reflect.Value, len(data))
	for i, arg := range data {
		val, err := ReflectValue(arg.Type, arg.Value)
		panicIf(err)
		res[i] = val
	}
	return res
}

func ReflectValue2TypeValue(reflectData []reflect.Value, data []TypeValue) {

}
