package netpro

import (
	"reflect"
	"runtime"
)

// GetFunctionName return the name of the function
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
