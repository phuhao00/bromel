package producer

import (
	"flag"
	"log"
	"testing"
)

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchangeName = flag.String("exchange", "test-exchange", "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	routingKey   = flag.String("key", "test-key", "AMQP routing key")
	body         = flag.String("body", "foobar", "Body of message")
	reliable     = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")
)

func TestProduct(t *testing.T) {
	cfg := &Config{
		Uri:          "amqp://guest:guest@localhost:5672/",
		ExchangeName: "test-exchange",
		ExchangeType: "direct",
		RoutingKey:   "test-key",
		Body:         "foobar",
		Reliable:     true,
	}
	if err := publish(cfg); err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("published %dB OK", len(cfg.Body))
}
