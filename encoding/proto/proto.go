// Package proto defines the protobuf codec. Importing this package will
// register the codec.
package proto

import (
	encoding2 "github.com/phuhao00/bromel/encoding"

	"google.golang.org/protobuf/proto"
)

// Name is the name registered for the proto compressor.
const Name = "proto"

func init() {
	encoding2.RegisterCodec(codec{})
}

// codec is a Codec implementation with protobuf. It is the default codec for Transport.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (codec) Name() string {
	return Name
}
