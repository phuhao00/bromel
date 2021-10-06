package main

import (
	"fmt"
	"github.com/phuhao00/bromel/meshwork"
	"github.com/phuhao00/bromel/meshwork/examples/simple/common"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	// go printGoroutineNum()

	s := meshwork.NewServer(&meshwork.ServerOption{
		SocketReadBufferSize:  1024 * 1024,
		SocketWriteBufferSize: 1024 * 1024,
		ReadTimeout:           time.Second * 3,
		WriteTimeout:          time.Second * 3,
		ReqQueueSize:          -1,
		RespQueueSize:         -1,
		Packer:                meshwork.NewDefaultPacker(),
		Codec:                 nil,
	})
	s.OnSessionCreate = func(sess *meshwork.Session) {
		fmt.Printf("session created: %s", sess.ID())
	}
	s.OnSessionClose = func(sess *meshwork.Session) {
		fmt.Printf("session closed: %s", sess.ID())
	}

	// register global middlewares
	//s.Use(fixture.RecoverMiddleware(log), logMiddleware)

	// register a route
	s.AddRoute(common.MsgIdPingReq, func(c *meshwork.Context) (*meshwork.Entry, error) {
		return c.Response(common.MsgIdPingAck, "pong, pong, pong")
	})

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
	time.Sleep(time.Second)
}

func logMiddleware(next meshwork.HandlerFunc) meshwork.HandlerFunc {
	return func(c *meshwork.Context) (resp *meshwork.Entry, err error) {
		fmt.Printf("rec <<< | id:(%d) size:(%d) data: %s", c.Message().ID, len(c.Message().Data), c.Message().Data)
		defer func() {
			if err != nil || resp == nil {
				return
			}
			fmt.Printf("snd >>> | id:(%d) size:(%d) data: %s", resp.ID, len(resp.Data), resp.Data)
		}()
		return next(c)
	}
}

// nolint: deadcode, unused
func printGoroutineNum() {
	for {
		fmt.Println("goroutine num: ", runtime.NumGoroutine())
		time.Sleep(time.Second)
	}
}
