package asynctask

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var typesMap = map[string]reflect.Type{
	// base types
	"bool":    reflect.TypeOf(true),
	"int":     reflect.TypeOf(int(1)),
	"int8":    reflect.TypeOf(int8(1)),
	"int16":   reflect.TypeOf(int16(1)),
	"int32":   reflect.TypeOf(int32(1)),
	"int64":   reflect.TypeOf(int64(1)),
	"uint":    reflect.TypeOf(uint(1)),
	"uint8":   reflect.TypeOf(uint8(1)),
	"uint16":  reflect.TypeOf(uint16(1)),
	"uint32":  reflect.TypeOf(uint32(1)),
	"uint64":  reflect.TypeOf(uint64(1)),
	"float32": reflect.TypeOf(float32(0.5)),
	"float64": reflect.TypeOf(float64(0.5)),
	"string":  reflect.TypeOf(string("")),
	// slices
	"[]bool":    reflect.TypeOf(make([]bool, 0)),
	"[]int":     reflect.TypeOf(make([]int, 0)),
	"[]int8":    reflect.TypeOf(make([]int8, 0)),
	"[]int16":   reflect.TypeOf(make([]int16, 0)),
	"[]int32":   reflect.TypeOf(make([]int32, 0)),
	"[]int64":   reflect.TypeOf(make([]int64, 0)),
	"[]uint":    reflect.TypeOf(make([]uint, 0)),
	"[]uint8":   reflect.TypeOf(make([]uint8, 0)),
	"[]uint16":  reflect.TypeOf(make([]uint16, 0)),
	"[]uint32":  reflect.TypeOf(make([]uint32, 0)),
	"[]uint64":  reflect.TypeOf(make([]uint64, 0)),
	"[]float32": reflect.TypeOf(make([]float32, 0)),
	"[]float64": reflect.TypeOf(make([]float64, 0)),
	"[]byte":    reflect.TypeOf(make([]byte, 0)),
	"[]string":  reflect.TypeOf([]string{""}),
}

func typeError(valueType string, value interface{}) error {
	return errors.New(fmt.Sprintf("Unsupported value type: %s, value: %s", typesMap[valueType].String(), value))
}

func ReflectValue(valueType string, value interface{}) (reflect.Value, error) {
	if strings.HasPrefix(valueType, "[]") {
		return reflectSlice(valueType, value)
	}
	return reflectValue(valueType, value)
}
func getBoolValue(valueType string, value interface{}) (bool, error) {
	if boolValue, ok := value.(bool); ok {
		return boolValue, nil
	} else {
		return boolValue, typeError(valueType, value)
	}
}
func getIntValue(valueType string, value interface{}) (int64, error) {
	if strings.HasPrefix(fmt.Sprintf("%T", value), "json.Number") {
		n, ok := value.(json.Number)
		if !ok {
			return 0, typeError(valueType, value)
		}
		return n.Int64()
	}
	if intVal, ok := value.(int64); ok {
		return intVal, nil
	}
	return 0, typeError(valueType, value)
}

func getUintValue(theType string, value interface{}) (uint64, error) {
	// We use https://golang.org/pkg/encoding/json/#Decoder.UseNumber when unmarshaling signatures.
	// This is because JSON only supports 64-bit floating point numbers and we could lose precision
	// when converting from float64 to unsigned integer
	if strings.HasPrefix(fmt.Sprintf("%T", value), "json.Number") {
		n, ok := value.(json.Number)
		if !ok {
			return 0, typeError(theType, value)
		}
		intVal, err := n.Int64()
		if err != nil {
			return 0, err
		}
		return uint64(intVal), nil
	}

	var n uint64
	switch value.(type) {
	case uint64:
		n = value.(uint64)
	case uint8:
		n = uint64(value.(uint8))
	default:
		return 0, typeError(theType, value)
	}
	return n, nil
}

func getFloatValue(theType string, value interface{}) (float64, error) {
	// We use https://golang.org/pkg/encoding/json/#Decoder.UseNumber when unmarshaling signatures.
	// This is because JSON only supports 64-bit floating point numbers and we could lose precision
	if strings.HasPrefix(fmt.Sprintf("%T", value), "json.Number") {
		n, ok := value.(json.Number)
		if !ok {
			return 0, typeError(theType, value)
		}

		return n.Float64()
	}

	f, ok := value.(float64)
	if !ok {
		return 0, typeError(theType, value)
	}

	return f, nil
}

func getStringValue(theType string, value interface{}) (string, error) {
	s, ok := value.(string)
	if !ok {
		return "", typeError(theType, value)
	}

	return s, nil
}

func reflectValue(valueType string, value interface{}) (reflect.Value, error) {
	curType, ok := typesMap[valueType]
	if !ok {
		return reflect.Value{}, typeError(valueType, value)
	}
	curValue := reflect.New(curType)
	if curType.String() == "bool" {
		b, err := getBoolValue(valueType, value)
		curValue.Elem().SetBool(b)
		return curValue.Elem(), err
	}
	if strings.HasPrefix(curType.String(), "int") {
		i, err := getIntValue(valueType, value)
		curValue.Elem().SetInt(i)
		return curValue.Elem(), err
	}
	if strings.HasPrefix(curType.String(), "uint") {
		ui, err := getUintValue(valueType, value)
		curValue.Elem().SetUint(ui)
		return curValue.Elem(), err
	}

	if strings.HasPrefix(curType.String(), "float") {
		f, err := getFloatValue(valueType, value)
		curValue.Elem().SetFloat(f)
		return curValue.Elem(), err
	}
	if curType.String() == "string" {
		s, err := getStringValue(valueType, value)
		curValue.Elem().SetString(s)
		return curValue.Elem(), err
	}
	return curValue, typeError(valueType, value)
}

func reflectSlice(valueType string, value interface{}) (reflect.Value, error) {
	curType, ok := typesMap[valueType]
	if !ok {
		return reflect.Value{}, errors.New(fmt.Sprintf("Unsupported value type %s, %s", valueType, value))
	}
	if value == nil {
		return reflect.MakeSlice(curType, 0, 0), nil
	}

	var theValue reflect.Value
	// Booleans
	if curType.String() == "[]bool" {
		bools := reflect.ValueOf(value)

		theValue = reflect.MakeSlice(curType, bools.Len(), bools.Len())
		for i := 0; i < bools.Len(); i++ {
			boolValue, err := getBoolValue(strings.Split(curType.String(), "[]")[1], bools.Index(i).Interface())
			if err != nil {
				return reflect.Value{}, err
			}

			theValue.Index(i).SetBool(boolValue)
		}

		return theValue, nil
	}

	// Integers
	if strings.HasPrefix(curType.String(), "[]int") {
		ints := reflect.ValueOf(value)

		theValue = reflect.MakeSlice(curType, ints.Len(), ints.Len())
		for i := 0; i < ints.Len(); i++ {
			intValue, err := getIntValue(strings.Split(curType.String(), "[]")[1], ints.Index(i).Interface())
			if err != nil {
				return reflect.Value{}, err
			}

			theValue.Index(i).SetInt(intValue)
		}

		return theValue, nil
	}

	// Unsigned integers
	if strings.HasPrefix(curType.String(), "[]uint") || curType.String() == "[]byte" {

		// Decode the base64 string if the value type is []uint8 or it's alias []byte
		// See: https://golang.org/pkg/encoding/json/#Marshal
		// > Array and slice values encode as JSON arrays, except that []byte encodes as a base64-encoded string
		if reflect.TypeOf(value).String() == "string" {
			output, err := base64.StdEncoding.DecodeString(value.(string))
			if err != nil {
				return reflect.Value{}, err
			}
			value = output
		}

		uints := reflect.ValueOf(value)

		theValue = reflect.MakeSlice(curType, uints.Len(), uints.Len())
		for i := 0; i < uints.Len(); i++ {
			uintValue, err := getUintValue(strings.Split(curType.String(), "[]")[1], uints.Index(i).Interface())
			if err != nil {
				return reflect.Value{}, err
			}

			theValue.Index(i).SetUint(uintValue)
		}

		return theValue, nil

	}
	return theValue, typeError(valueType, value)
}
