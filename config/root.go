package config

import (
	"github.com/common-go/kafka"
	"github.com/common-go/log"
	"github.com/common-go/mongo"
	"github.com/common-go/mq"
)

type Root struct {
	Server            ServerConfig          `mapstructure:"server"`
	Log               log.Config            `mapstructure:"log"`
	Mongo             mongo.MongoConfig     `mapstructure:"mongo"`
	KafkaConsumer     kafka.ConsumerConfig  `mapstructure:"kafka_consumer"`
	KafkaProducer     kafka.ProducerConfig  `mapstructure:"kafka_producer"`
	BatchWorkerConfig mq.BatchWorkerConfig  `mapstructure:"batch_worker"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}
