package meshwork

import (
	"encoding/binary"
	"fmt"
	"github.com/phuhao00/pdproto/go/messageId"
	"io"
	"net"
	"time"
)

type CustomPacker struct{}

func (p *CustomPacker) Pack(entry *Entry) ([]byte, error) {
	buffer := make([]byte, 4+4+len(entry.Data))
	p.byteOrder().PutUint32(buffer[0:4], uint32(len(buffer)))                    // write totalSize
	p.byteOrder().PutUint32(buffer[4:8], uint32(entry.ID.(messageId.MessageID))) // write id
	copy(buffer[8:], entry.Data)                                                 // write data
	return buffer, nil
}

func (p *CustomPacker) Unpack(reader io.Reader) (*Entry, error) {
	reader.(*net.TCPConn).SetReadDeadline(time.Now().Add(time.Second * 10))
	headerBuffer := make([]byte, 4+4)
	if _, err := io.ReadFull(reader, headerBuffer); err != nil {
		//fmt.Println(reflect.TypeOf(p).Name())
		return nil, err
	}
	totalSize := p.byteOrder().Uint32(headerBuffer[:4])               // read totalSize
	id := messageId.MessageID(p.byteOrder().Uint32(headerBuffer[4:])) // read id

	// read data
	dataSize := totalSize - 4 - 4
	data := make([]byte, dataSize)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, fmt.Errorf("read data from reader err: %s", err)
	}
	entry := &Entry{
		ID:   id,
		Data: data,
	}
	return entry, nil
}

func (*CustomPacker) byteOrder() binary.ByteOrder {
	return binary.BigEndian
}
