package router

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
)

type RouterGroup struct {
	//middleware
	preHandler     HandlerFunc
	postHandler    HandlerFunc
	handlers       map[uint16]HandlerFunc
	littleEndian   bool
	defaultHandler HandlerFunc
}

func NewRouterGroup() *RouterGroup {
	return &RouterGroup{handlers: make(map[uint16]HandlerFunc)}
}

type HandlerFunc func(msgID uint16, data []byte)

func (rg *RouterGroup) Register(msgID uint16, handlerFunc HandlerFunc) {
	rg.handlers[msgID] = handlerFunc
}

func (rg *RouterGroup) Unmarshal(data []byte, msg interface{}) error {
	if len(data) < 2 {
		return errors.New("protobuf data too short")
	}
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("msg is not protobuf message")
	}
	return proto.UnmarshalMerge(data[2:], pbMsg)
}

func (rg *RouterGroup) Marshal(msgID uint16, msg interface{}) ([]byte, error) {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return []byte{}, fmt.Errorf("msg is not protobuf message")
	}
	// data
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		return data, err
	}
	// 4byte = len(flag)[2byte] + len(msgID)[2byte]
	buf := make([]byte, 4+len(data))
	if rg.littleEndian {
		binary.LittleEndian.PutUint16(buf[0:2], 0)
		binary.LittleEndian.PutUint16(buf[2:], msgID)
	} else {
		binary.BigEndian.PutUint16(buf[0:2], 0)
		binary.BigEndian.PutUint16(buf[2:], msgID)
	}
	copy(buf[4:], data)
	return buf, err
}

func (rg *RouterGroup) Route(data []byte) (uint16, error) {

	if len(data) < 2 {
		return 0, errors.New("protobuf data too short")
	}

	var msgID uint16
	if rg.littleEndian {
		msgID = binary.LittleEndian.Uint16(data)
	} else {
		msgID = binary.BigEndian.Uint16(data)
	}

	handler, ok := rg.handlers[msgID]

	if ok && handler != nil {
		if rg.preHandler != nil {
			rg.preHandler(msgID, data)
		}
		handler(msgID, data)
		if rg.postHandler != nil {
			rg.postHandler(msgID, data)
		}
	} else {
		if rg.defaultHandler != nil {
			rg.defaultHandler(msgID, data)
			return msgID, nil
		}
		return msgID, errors.New("unknow msg msgID:" + strconv.Itoa(int(msgID)))
	}

	return msgID, nil
}
