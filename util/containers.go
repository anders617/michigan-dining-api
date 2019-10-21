package util

import (
	"reflect"

	"github.com/golang/glog"
)

// Values - Returns the value from the given map
func Values(input interface{}) *[]interface{} {
	v := reflect.ValueOf(input)
	values := make([]interface{}, 0, v.Len())
	if v.Kind() == reflect.Map {
		iter := v.MapRange()
		for iter.Next() {
			values = append(values, iter.Value().Interface())
		}
	} else {
		glog.Fatalf("Non Map Type Passed to Values %v", input)
	}
	return &values
}

func AsSliceType(s interface{}, outSlice interface{}) interface{} {
	// Some reflection magic that takes input array and converts to a different type
	inV := reflect.Indirect(reflect.ValueOf(s))

	t := reflect.TypeOf(outSlice)
	elemT := t.Elem()
	out := reflect.MakeSlice(t, 0, inV.Len())

	for idx := 0; idx < inV.Len(); idx++ {
		out = reflect.Append(out, inV.Index(idx).Elem().Convert(elemT))
	}
	return out.Interface()
}
