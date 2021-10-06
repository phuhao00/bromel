package main

import (
	"fmt"
	"github.com/phuhao00/bromel/encoding"
	"github.com/phuhao00/bromel/meshwork"
	"github.com/phuhao00/bromel/meshwork/examples/custom_packet/common"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	s := meshwork.NewServer(&meshwork.ServerOption{
		// specify codec and packer
		Codec:  encoding.GetCodec("json"),
		Packer: &meshwork.CustomPacker{},
	})

	s.AddRoute("json01-req", handler, meshwork.RecoverMiddleware(nil), logMiddleware)

	go func() {
		if err := s.Serve("0.0.0.0:8088"); err != nil {
			fmt.Printf("serve err: %s", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-sigCh
	if err := s.Stop(); err != nil {
		fmt.Printf("server stopped err: %s", err)
	}
}

func handler(ctx *meshwork.Context) (*meshwork.Entry, error) {
	var data common.Json01Req
	ctx.MustBind(&data)

	return ctx.Response("json01-resp", &common.Json01Resp{
		Success: true,
		Data:    fmt.Sprintf("%s:%d:%t", data.Key1, data.Key2, data.Key3),
	})
}

func logMiddleware(next meshwork.HandlerFunc) meshwork.HandlerFunc {
	return func(ctx *meshwork.Context) (resp *meshwork.Entry, err error) {
		fullSize := ctx.Message().MustGet("fullSize")
		fmt.Printf("recv request  | fullSize:(%d) id:(%v) dataSize(%d) data: %s", fullSize, ctx.Message().ID, len(ctx.Message().Data), ctx.Message().Data)

		defer func() {
			if err != nil {
				return
			}
			if resp != nil {
				fmt.Printf("send response | dataSize:(%d) id:(%v) data: %s", len(resp.Data), resp.ID, resp.Data)
			} else {
				fmt.Printf("don't send response since nil")
			}
		}()
		return next(ctx)
	}
}
