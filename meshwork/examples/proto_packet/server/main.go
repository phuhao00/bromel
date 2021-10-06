package main

import (
	"fmt"
	"github.com/phuhao00/bromel/encoding"
	"github.com/phuhao00/bromel/meshwork"
	"github.com/phuhao00/bromel/meshwork/examples/proto_packet/common"
	"google.golang.org/protobuf/proto"
)

func main() {
	srv := meshwork.NewServer(&meshwork.ServerOption{
		Packer: &common.CustomPacker{},
		Codec:  encoding.GetCodec("proto"),
	})

	srv.AddRoute(common.ID_FooReqID, handle, logTransmission(&common.FooReq{}, &common.FooResp{}))

	if err := srv.Serve("0.0.0.0:8088"); err != nil {
		fmt.Printf("serve err: %s", err)
	}
}

func handle(c *meshwork.Context) (*meshwork.Entry, error) {
	var reqData common.FooReq
	c.MustBind(&reqData)
	return c.Response(common.ID_FooRespID, &common.FooResp{
		Code:    2,
		Message: "success",
	})
}

func logTransmission(req, resp proto.Message) meshwork.MiddlewareFunc {
	return func(next meshwork.HandlerFunc) meshwork.HandlerFunc {
		return func(c *meshwork.Context) (*meshwork.Entry, error) {
			if err := c.Bind(req); err == nil {
				fmt.Printf("recv | id: %d; size: %d; data: %s", c.Message().ID, len(c.Message().Data), req)
			}

			respEntry, err := next(c)

			if err == nil && respEntry != nil {
				c.MustDecodeTo(respEntry.Data, resp)
				fmt.Printf("send | id: %d; size: %d; data: %s", respEntry.ID, len(respEntry.Data), resp)
			}
			return respEntry, err
		}
	}
}
