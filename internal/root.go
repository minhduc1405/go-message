package internal

import (
	"github.com/common-go/amq"
	"github.com/common-go/log"
	"github.com/common-go/mongo"
	"github.com/common-go/mq"
)

type Root struct {
	Server            ServerConfig          `mapstructure:"server"`
	Log               log.Config            `mapstructure:"log"`
	Mongo             mongo.MongoConfig     `mapstructure:"mongo"`
	Amq               amq.Config            `mapstructure:"amq"`
	BatchWorkerConfig mq.BatchWorkerConfig  `mapstructure:"batch_worker"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}
