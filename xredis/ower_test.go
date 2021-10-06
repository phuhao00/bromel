package xredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/phuhao00/bromel"
	"testing"
	"time"
)

//docker run -itd --name redis-test -p 6379:6379 redis

type testOwner struct {
	c *Client
}

func (t testOwner) Launch() {
	t.c.Launch()
}

func (t testOwner) Stop() {
	t.c.Stop()
}

func TestOwner(t *testing.T) {
	to := &testOwner{c: &Client{
		BaseComponent: bromel.NewBaseComponent(NewClient(context.Background(), &Config{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		})),
	}}
	to.c.Launch()
	fn := func(client bromel.Component) {
		real := client.GetReal()
		rdb := real.(*redis.Client)
		result := rdb.Ping(context.Background())
		fmt.Println(result)
	}
	op := bromel.Operation{CB: fn, Ret: make(chan interface{})}
	to.c.Resolve(op)
	<-op.Ret
	time.Sleep(time.Second * 5)
	to.c.Stop()
}
