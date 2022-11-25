package util

import (
	"fmt"
	"reflect"
)

func CallbackDataSerialize[T any](t T, data any) string {
	typeName := reflect.TypeOf(t).Name()
	res := fmt.Sprintf("%s||%v", typeName, data)
	return res
}
