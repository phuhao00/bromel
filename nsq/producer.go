package nsq

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

type XProducer struct {
	*nsq.Producer
	*ProducerConfig
}

func NewProducer() *XProducer {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}
	return &XProducer{producer, NewProducerConfig()}
}

//Pub ...
func (p XProducer) Pub(topicName string, messageBody string) (err error) {
	if messageBody == "" {
		return errors.New("message is empty")
	}
	if err = p.Publish(topicName, []byte(messageBody)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// 延迟消息
func (p *XProducer) deferredPublish(topic string, delay time.Duration, message string) (err error) {
	if message == "" {
		return errors.New("message is empty")
	}
	if err = p.DeferredPublish(topic, delay, []byte(message)); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//Stop ...
func (p *XProducer) Stop() {
	p.Producer.Stop()

}
