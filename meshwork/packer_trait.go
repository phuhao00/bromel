package meshwork

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Packer is a generic interface to pack and unpack message packet.
type Packer interface {
	// Pack packs Message into the packet to be written.
	// Pack(msg Message) ([]byte, error)
	Pack(entry *Entry) ([]byte, error)

	// Unpack unpacks the message packet from reader,
	// returns the Message interface, and error if error occurred.
	Unpack(reader io.Reader) (*Entry, error)
}

var _ Packer = &DefaultPacker{}

// NewDefaultPacker create a *DefaultPacker with initial field value.
func NewDefaultPacker() *DefaultPacker {
	return &DefaultPacker{MaxSize: 1024 * 1024}
}

// DefaultPacker is the default Packer used in session.
// DefaultPacker treats the packet with the format:
// 	(size)(id)(data):
// 		size: uint32 | took 4 bytes, only the size of `data`
// 		id:   uint32 | took 4 bytes
// 		data: []byte | took `size` bytes
type DefaultPacker struct {
	MaxSize int
}

func (d *DefaultPacker) bytesOrder() binary.ByteOrder {
	return binary.BigEndian
}

func (d *DefaultPacker) assertID(id interface{}) (uint32, bool) {
	switch v := id.(type) {
	case uint32:
		return v, true
	case *uint32:
		return *v, true
	default:
		return 0, false
	}
}

// Pack implements the Packer Pack method.
func (d *DefaultPacker) Pack(entry *Entry) ([]byte, error) {
	dataSize := len(entry.Data)
	buffer := make([]byte, 4+4+dataSize)
	d.bytesOrder().PutUint32(buffer[:4], uint32(dataSize)) // write dataSize
	id, err := ToUint32E(entry.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid type of entry.ID: %s", err)
	}
	d.bytesOrder().PutUint32(buffer[4:8], id) // write id
	copy(buffer[8:], entry.Data)              // write data
	return buffer, nil
}

// Unpack implements the Packer Unpack method.
func (d *DefaultPacker) Unpack(reader io.Reader) (*Entry, error) {
	headerBuffer := make([]byte, 4+4)
	if _, err := io.ReadFull(reader, headerBuffer); err != nil {
		return nil, fmt.Errorf("read size and id err: %s", err)
	}
	dataSize := d.bytesOrder().Uint32(headerBuffer[:4])
	if d.MaxSize > 0 && int(dataSize) > d.MaxSize {
		return nil, fmt.Errorf("the dataSize %d is beyond the max: %d", dataSize, d.MaxSize)
	}
	id := d.bytesOrder().Uint32(headerBuffer[4:8])
	data := make([]byte, dataSize)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, fmt.Errorf("read data err: %s", err)
	}
	entry := &Entry{
		ID:   int(id),
		Data: data,
	}
	return entry, nil
}
