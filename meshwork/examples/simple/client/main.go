package main

import (
	"fmt"
	"github.com/phuhao00/bromel/meshwork"
	"github.com/phuhao00/bromel/meshwork/examples/simple/common"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:8088")
	if err != nil {
		panic(err)
	}
	packer := meshwork.NewDefaultPacker()
	go func() {
		// write loop
		for {
			time.Sleep(time.Second)
			rawData := []byte("ping, ping, ping")
			msg := &meshwork.Entry{
				ID:   common.MsgIdPingReq,
				Data: rawData,
			}
			packedMsg, err := packer.Pack(msg)
			if err != nil {
				panic(err)
			}
			if _, err := conn.Write(packedMsg); err != nil {
				panic(err)
			}
			fmt.Printf("snd >>> | id:(%d) size:(%d) data: %s", msg.ID, len(rawData), rawData)
		}
	}()
	go func() {
		// read loop
		for {
			msg, err := packer.Unpack(conn)
			if err != nil {
				panic(err)
			}
			fmt.Printf("rec <<< | id:(%d) size:(%d) data: %s", msg.ID, len(msg.Data), msg.Data)
		}
	}()
	select {}
}
