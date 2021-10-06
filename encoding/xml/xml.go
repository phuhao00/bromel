package xml

import (
	"encoding/xml"
	encoding2 "github.com/phuhao00/bromel/encoding"
)

// Name is the name registered for the xml codec.
const Name = "xml"

func init() {
	encoding2.RegisterCodec(codec{})
}

// codec is a Codec implementation with xml.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
