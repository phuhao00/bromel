package pubsub

import (
	"context"
	"flag"
	"os"
	"testing"
)

var url = flag.String("url", "amqp:///", "AMQP url for both the publisher and subscriber")

// exchange binds the publishers to the subscribers

func TestPubSub(t *testing.T) {

	ctx, done := context.WithCancel(context.Background())
	go func() {
		publish(redial(ctx, "amqp://guest:guest@localhost:5672/"), read(os.Stdin))
		done()
	}()

	go func() {
		subscribe(redial(ctx, "amqp://guest:guest@localhost:5672/"), write(os.Stdout))
		done()
	}()

	<-ctx.Done()
}
