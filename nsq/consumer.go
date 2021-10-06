package nsq

import (
	"log"

	"github.com/nsqio/go-nsq"
)

type XConsumer struct {
	*nsq.Consumer
	conf *ConsumerConfig
}

type OptionConsumer func(consumer *XConsumer)

func WithConfig(config *ConsumerConfig) OptionConsumer {
	return func(consumer *XConsumer) {
		consumer.conf = config
	}
}

//NewXConsumer ...
func NewXConsumer(opts ...OptionConsumer) XConsumer {
	c := &XConsumer{conf: defaultConsumerConfig()}
	for _, opt := range opts {
		opt(c)
	}
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(c.conf.Topic, c.conf.Channel, config)
	if err != nil {
		log.Fatal(err)
	}
	consumer.AddHandler(&NormalMsgHandle{})
	err = consumer.ConnectToNSQLookupd(c.conf.Host + ":" + c.conf.Port)
	if err != nil {
		log.Fatal(err)
	}
	return XConsumer{consumer, NewConsumerConfig()}
}

//Stop ...
func (c *XConsumer) Stop() {
	c.Consumer.Stop()
}
