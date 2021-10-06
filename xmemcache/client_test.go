package xmemcache

import (
	"fmt"
	"github.com/phuhao00/bromel"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

type TestOwner struct {
	c bromel.Component
}

func (t TestOwner) Launch() {
	t.c.Launch()
}

func (t TestOwner) Stop() {
	t.c.Stop()
}

func TestClient(t *testing.T) {
	c := &Client{
		BaseComponent: bromel.NewBaseComponent(NewClient(&Config{ServerAddress: "127.0.0.1:11211"})),
	}
	to := TestOwner{c: c}
	to.c.Launch()
	op := bromel.Operation{
		CB: func(component bromel.Component) {
			real := component.GetReal()
			mc := real.(*memcache.Client)
			err := mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})
			if err != nil {
				fmt.Println(err)
			}
			it, err := mc.Get("foo")
			fmt.Println(string(it.Value), err)
		},
		Ret: make(chan interface{}),
	}
	to.c.Resolve(op)
	<-op.Ret

}
