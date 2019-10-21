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
	inT := reflect.TypeOf(s).Elem()
	inElemT := inT.Elem()
	inV := reflect.Indirect(reflect.ValueOf(s))

	t := reflect.TypeOf(outSlice)
	outElemT := t.Elem()
	out := reflect.MakeSlice(t, 0, inV.Len())

	for idx := 0; idx < inV.Len(); idx++ {
		if outElemT.Kind() == reflect.Interface && inElemT.Kind() == reflect.Ptr {
			// If out elem type is Interface and in elem type is ptr
			// then no need to call .Elem() on inV since interfaces are ptrs to values
			out = reflect.Append(out, inV.Index(idx))
		} else if outElemT.Kind() == reflect.Interface && inElemT.Kind() == reflect.Struct {
			out = reflect.Append(out, inV.Index(idx))
		} else {
			out = reflect.Append(out, inV.Index(idx).Elem())
		}
	}
	return out.Interface()
}
