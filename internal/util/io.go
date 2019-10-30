package util

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

//
// ReadProtoFromFile - Reads a proto from the given file path
//
func ReadProtoFromFile(path string, p proto.Message) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Error("Failed to read in text proto, %v", err)
		return err
	}
	proto.UnmarshalText(string(data), p)
	return nil
}

func MapToString(m interface{}) string {
	v := reflect.ValueOf(m)
	b := strings.Builder{}
	mapToString(v, 0, &b)
	return b.String()
}

func mapToString(v reflect.Value, indent int, builder *strings.Builder) {
	if v.Kind() == reflect.Map {
		iter := v.MapRange()
		for iter.Next() {
			for i := 0; i < indent*4; i++ {
				builder.WriteString(" ")
			}
			builder.WriteString(fmt.Sprint(iter.Key()))
			builder.WriteString(": ")
			if iter.Value().Kind() == reflect.Map {
				builder.WriteString("\n")
				mapToString(iter.Value(), indent+1, builder)
			} else {
				builder.WriteString(fmt.Sprint(iter.Value()))
			}
			builder.WriteString("\n")
		}
	} else {
		glog.Fatalf("Non Map Type Passed to Values %v", v)
	}
}
