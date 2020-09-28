package config

import (
	"context"
	"github.com/common-go/health"
	"github.com/common-go/kafka"
	"github.com/common-go/mongo"
	"github.com/common-go/mq"
	v "github.com/common-go/validator"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

type ApplicationContext struct {
	Consumer         mq.Consumer
	ConsumerCaller   mq.ConsumerCaller
	BatchWorker      mq.BatchWorker
	HealthController *health.HealthController
}

func NewApplicationContext(ctx context.Context, root Root) (*ApplicationContext, error) {
	mongoDb, err := mongo.SetupMongo(ctx, root.Mongo)
	if err != nil {
		logrus.Errorf("Can't connect mongoDB: Error: %s", err.Error())
		return nil, err
	}

	//_, err2 := kafka.NewConnect(root.KafkaConsumer, root.KafkaConsumer.Brokers[0])
	//if err2 != nil {
	//	logrus.Errorf("Can't connect Kafka: Error: %s", err.Error())
	//	return nil, err2
	//}

	consumer, err := kafka.NewConsumerByConfig(root.KafkaConsumer, true)
	if err != nil {
		logrus.Errorf("Can't new consumer: Error: %s", err.Error())
		return nil, err
	}
	producer, err := kafka.NewProducerByConfig(root.KafkaProducer, true)
	if err != nil {
		logrus.Errorf("Can't new producer: Error: %s", err.Error())
		return nil, err
	}

	userTypeOf := reflect.TypeOf(User{})
	bulkWriter := mongo.NewMongoBatchInsert(mongoDb, "users")
	//bulkWriter := mongo.NewMongoBatchUpdate(mongoDb, "users", userTypeOf)
	//bulkWriter := writer.NewMongoBatchWriter(mongoDb, "users", userTypeOf)
	batchHandler := mq.NewBatchHandler(userTypeOf, bulkWriter)

	retryService := mq.NewMqRetryService(producer)
	batchWorker := mq.NewDefaultBatchWorker(root.BatchWorkerConfig, batchHandler, retryService)

	v := NewUserValidator()
	validator := mq.NewValidator(userTypeOf, v)
	consumerCaller := mq.NewBatchConsumerCaller(batchWorker, validator)
	// consumerCaller := mq.NewBatchConsumerCaller(batchWorker, nil)

	mongoHealthService := mongo.NewDefaultMongoHealthService(mongoDb)
	subHealthService := kafka.NewKafkaHealthService(
		root.KafkaConsumer.Brokers,
		"kafka",
	)
	healthServices := []health.HealthService{mongoHealthService, subHealthService}
	healthController := health.NewHealthController(healthServices)
	return &ApplicationContext{
		Consumer:         consumer,
		ConsumerCaller:   consumerCaller,
		BatchWorker:      batchWorker,
		HealthController: healthController,
	}, nil
}

func NewUserValidator() v.Validator {
	validator := v.NewDefaultValidator()
	validator.CustomValidateList = append(validator.CustomValidateList, v.CustomValidate{Fn: CheckActive, Tag: "active"})
	return validator
}

func CheckActive(fl validator.FieldLevel) bool {
	return fl.Field().Bool()
}
