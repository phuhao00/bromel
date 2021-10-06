package consumer

import (
	"log"
	"testing"
	"time"
)

//var (
//	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
//	exchange     = flag.String("exchange", "test-exchange", "Durable, non-auto-deleted AMQP exchange name")
//	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
//	queue        = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
//	bindingKey   = flag.String("key", "test-key", "AMQP binding key")
//	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
//	lifetime     = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
//)
//
//func init() {
//	flag.Parse()
//}

func TestConsumer(t *testing.T) {
	cfg := &Config{
		Uri:          "amqp://guest:guest@localhost:5672/",
		Exchange:     "test-exchange",
		ExchangeType: "direct",
		Queue:        "test-queue",
		BindingKey:   "test-key",
		ConsumerTag:  "simple-consumer",
		LifeTime:     5 * time.Second,
	}
	c, err := NewConsumer(cfg)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if cfg.LifeTime > 0 {
		log.Printf("running for %s", cfg.LifeTime)
		time.Sleep(cfg.LifeTime)
	} else {
		log.Printf("running forever")
		select {}
	}

	log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}
}
