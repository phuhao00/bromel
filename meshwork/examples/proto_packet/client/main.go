package main

import (
	"fmt"
	"github.com/phuhao00/bromel/encoding"
	_ "github.com/phuhao00/bromel/encoding/proto"
	"github.com/phuhao00/bromel/meshwork"
	"github.com/phuhao00/bromel/meshwork/examples/proto_packet/common"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:8088")
	if err != nil {
		panic(err)
	}

	packer := &meshwork.CustomPacker{}
	codec := encoding.GetCodec("proto")

	go func() {
		for {
			var id = common.ID_FooReqID
			req := &common.FooReq{
				Bar: "bar",
				Buz: 22,
			}
			data, err := codec.Marshal(req)
			if err != nil {
				panic(err)
			}
			msg := &meshwork.Entry{ID: id, Data: data}
			packedMsg, err := packer.Pack(msg)
			if err != nil {
				panic(err)
			}
			if _, err := conn.Write(packedMsg); err != nil {
				panic(err)
			}
			fmt.Printf("send | id: %d; size: %d; data: %s", id, len(data), req.String())
			time.Sleep(time.Second)
		}
	}()

	for {
		msg, err := packer.Unpack(conn)
		if err != nil {
			panic(err)
		}
		var respData common.FooResp
		if err := codec.Unmarshal(msg.Data, &respData); err != nil {
			panic(err)
		}
		fmt.Printf("recv | id: %d; size: %d; data: %s", msg.ID, len(msg.Data), respData.String())
	}
}
