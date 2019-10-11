package util

import (
	"io/ioutil"

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
