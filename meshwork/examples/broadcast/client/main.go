package main

import (
	"fmt"
	"github.com/phuhao00/bromel/meshwork"
	"github.com/phuhao00/bromel/meshwork/examples/broadcast/common"
	"net"
	"time"
)

var packer meshwork.Packer

func init() {

	packer = meshwork.NewDefaultPacker()
}

func main() {
	senderClient()
	for i := 0; i < 10; i++ {
		readerClient(i)
	}

	select {}
}

func establish() (net.Conn, error) {
	return net.Dial("tcp", "0.0.0.0:8088")
}

func senderClient() {
	conn, err := establish()
	if err != nil {
		fmt.Println(err)
		return
	}
	// send
	go func() {
		for {
			time.Sleep(time.Second)
			data := []byte(fmt.Sprintf("hello everyone @1111111111111111111111111111111111111%d", time.Now().Unix()))
			msg := &meshwork.Entry{
				ID:   common.MsgIdBroadCastReq,
				Data: data,
			}

			packedMsg, _ := packer.Pack(msg)
			if _, err := conn.Write(packedMsg); err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	// read
	go func() {
		for {
			msg, err := packer.Unpack(conn)
			if err != nil {
				fmt.Println(err)

				return
			}
			fmt.Printf("sender | recv ack | %s", msg.Data)
		}
	}()
}

func readerClient(id int) {
	conn, err := establish()
	if err != nil {
		fmt.Println(err)

		return
	}

	go func() {
		for {
			msg, err := packer.Unpack(conn)
			if err != nil {
				fmt.Println(err)

				return
			}
			fmt.Printf("reader %03d | recv broadcast | %s", id, msg.Data)
		}
	}()
}
