package config

import (
	"github.com/common-go/kafka"
	"github.com/common-go/log"
	"github.com/common-go/mq"
)

type Root struct {
	Server            ServerConfig         `mapstructure:"server"`
	Log               log.Config           `mapstructure:"log"`
	Sql               SqlConfig            `mapstructure:"sql"`
	KafkaConsumer     kafka.ConsumerConfig `mapstructure:"kafka_consumer"`
	KafkaProducer     kafka.ProducerConfig `mapstructure:"kafka_producer"`
	BatchWorkerConfig mq.BatchWorkerConfig `mapstructure:"batch_worker"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type SqlConfig struct {
	Dialect string `mapstructure:"dialect"`
	Uri     string `mapstructure:"uri"`
}
