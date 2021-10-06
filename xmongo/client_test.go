package xmongo

import (
	"context"
	"fmt"
	"github.com/phuhao00/bromel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

type testOwner struct {
	c *Client
}

func (t testOwner) Launch() {
	t.c.Launch()
}

func (t testOwner) Stop() {
	t.c.Stop()
}

func TestClient(t *testing.T) {
	ctx := context.Background()
	to := &testOwner{}
	tc := &Client{
		BaseComponent: bromel.NewBaseComponent(NewClient(ctx, &Config{
			URI:         "mongodb://localhost:27017",
			MinPoolSize: 3,
			MaxPoolSize: 3000,
		})),
	}
	to.c = tc
	to.Launch()
	fn := func(client bromel.Component) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		real := client.GetReal()
		mongoCli := real.(*mongo.Client)
		res, err := InsertOne(ctx, GetCOll(mongoCli, "retu_test", "retu-test_collection"),
			bson.D{{"name", "pi"}, {"value", 3.14159}})
		if err != nil {
			fmt.Println(err)
		}
		id := res.InsertedID
		fmt.Println(id)
	}
	op := bromel.Operation{
		CB:  fn,
		Ret: make(chan interface{}),
	}
	to.c.Resolve(op)
	<-op.Ret
	fmt.Println("op success")
	time.Sleep(time.Second * 5)
	to.Stop()
}
