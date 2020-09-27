package config

import (
	"context"
	"reflect"

	"github.com/common-go/amq"
	"github.com/common-go/health"
	"github.com/common-go/mongo"
	"github.com/common-go/mq"
	v "github.com/common-go/validator"
	"github.com/go-stomp/stomp"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
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

	consumer, err := amq.NewConsumerByConfig(root.Amq, stomp.AckAuto, true)
	//consumer, err := pubsub.NewConsumerByConfig(ctx, root.PubSubConsumer, true)
	if err != nil {
		logrus.Errorf("Can't new consumer: Error: %s", err.Error())
		return nil, err
	}
	producer, err := amq.NewProducerByConfig(root.Amq, "")
	// producer, err := pubsub.NewProducerByConfig(ctx, root.PubSubProducer)
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
	consumerHealthService := amq.NewAMQHealthService(consumer.Conn, "amqConsumer", "")
	producerHealthService := amq.NewAMQHealthService(consumer.Conn, "amqProducer", "")
	healthServices := []health.HealthService{mongoHealthService, consumerHealthService, producerHealthService}
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
	validator.CustomValidateList = append(validator.CustomValidateList, v.CustomValidate{Fn: CheckFlag1, Tag: "active"})
	return validator
}

func CheckFlag1(fl validator.FieldLevel) bool {
	return fl.Field().Bool()
}
