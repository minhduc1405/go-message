package internal

import (
	"github.com/common-go/log"
	"github.com/common-go/mongo"
	"github.com/common-go/mq"
	"github.com/common-go/pubsub"
)

type Root struct {
	Server            ServerConfig          `mapstructure:"server"`
	Log               log.Config            `mapstructure:"log"`
	Mongo             mongo.MongoConfig     `mapstructure:"mongo"`
	PubSubConsumer    pubsub.ConsumerConfig `mapstructure:"pubsub_consumer"`
	PubSubProducer    pubsub.ProducerConfig `mapstructure:"pubsub_producer"`
	BatchWorkerConfig mq.BatchWorkerConfig  `mapstructure:"batch_worker"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}
