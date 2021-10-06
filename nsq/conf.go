package nsq

type ConsumerConfig struct {
	Topic   string `json:"topic"`
	Channel string `json:"channel"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

func defaultConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		Topic:   "test",
		Channel: "channelTest",
		Host:    "localhost",
		Port:    "4161",
	}
}

func NewConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{}
}

type ProducerConfig struct {
}

func NewProducerConfig() *ProducerConfig {
	return &ProducerConfig{}
}
