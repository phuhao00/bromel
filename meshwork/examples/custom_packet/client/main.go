package main

import (
	"fmt"
	"github.com/phuhao00/bromel/encoding"
	_ "github.com/phuhao00/bromel/encoding/json"
	"github.com/phuhao00/bromel/meshwork"
	"github.com/phuhao00/bromel/meshwork/examples/custom_packet/common"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:8088")
	if err != nil {
		panic(err)
	}
	codec := encoding.GetCodec("json")
	packer := &meshwork.CustomPacker{}
	go func() {
		// write loop
		for {
			time.Sleep(time.Second)
			req := &common.Json01Req{
				Key1: "hello",
				Key2: 10,
				Key3: true,
			}
			data, err := codec.Marshal(req)
			if err != nil {
				panic(err)
			}
			msg := &meshwork.Entry{
				ID:   "json01-req",
				Data: data,
			}
			packedMsg, err := packer.Pack(msg)
			if err != nil {
				panic(err)
			}
			if _, err := conn.Write(packedMsg); err != nil {
				panic(err)
			}
		}
	}()
	go func() {
		// read loop
		for {
			msg, err := packer.Unpack(conn)
			if err != nil {
				panic(err)
			}
			fullSize := msg.MustGet("fullSize")
			fmt.Printf("ack received | fullSize:(%d) id:(%v) dataSize:(%d) data: %s", fullSize, msg.ID, len(msg.Data), msg.Data)
		}
	}()
	select {}
}
