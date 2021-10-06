package main

import (
	"fmt"
	"github.com/phuhao00/bromel/meshwork"
	"github.com/phuhao00/bromel/meshwork/examples/broadcast/common"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s := meshwork.NewServer(&meshwork.ServerOption{
		Packer: meshwork.NewDefaultPacker(),
	})

	s.Use(meshwork.RecoverMiddleware(nil), logMiddleware)

	s.AddRoute(common.MsgIdBroadCastReq, func(ctx *meshwork.Context) (*meshwork.Entry, error) {
		reqData := ctx.Message().Data

		// broadcasting
		go meshwork.Sessions().Range(func(id string, sess *meshwork.Session) (next bool) {
			if ctx.Session().ID() == id {
				return true // next iteration
			}
			msg, err := ctx.Response(common.MsgIdBroadCastAck, fmt.Sprintf("%s (broadcast from %s)", reqData, ctx.Session().ID()))
			if err != nil {
				fmt.Printf("create response err: %s", err)
				return true
			}
			if err := sess.SendResp(msg); err != nil {
				fmt.Printf("broadcast err: %s", err)
			}
			return true
		})

		return ctx.Response(common.MsgIdBroadCastAck, "broadcast done")
	})

	go func() {
		if err := s.Serve("0.0.0.0:8088"); err != nil {
			fmt.Println(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-sigCh
	if err := s.Stop(); err != nil {
		fmt.Printf("server stopped err: %s", err)
	}
	time.Sleep(time.Second)
}

func logMiddleware(next meshwork.HandlerFunc) meshwork.HandlerFunc {
	return func(ctx *meshwork.Context) (resp *meshwork.Entry, err error) {
		fmt.Printf("recv request | %s", ctx.Message().Data)
		// defer func() {
		// 	if err != nil || resp == nil {
		// 		return
		// 	}
		// 	log.Infof("send response | id: %d; size: %d; data: %s", resp.ID, len(resp.Data), resp.Data)
		// }()
		return next(ctx)
	}
}
