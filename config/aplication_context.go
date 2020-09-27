package config

import (
	"context"
	"reflect"
	"time"

	"github.com/common-go/health"
	"github.com/common-go/mongo"
	"github.com/common-go/mq"
	"github.com/common-go/pubsub"
	v "github.com/common-go/validator"
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

	consumer, err := pubsub.NewConsumerByConfig(ctx, root.PubSubConsumer, true)
	if err != nil {
		logrus.Errorf("Can't new consumer: Error: %s", err.Error())
		return nil, err
	}
	producer, err := pubsub.NewProducerByConfig(ctx, root.PubSubProducer)
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
	subHealthService := pubsub.NewPubSubHealthService(
		"pubsubSubscriber",
		consumer.Client,
		10 * time.Second,
		pubsub.PermissionSubscribe,
		root.PubSubConsumer.SubscriptionId,
	)
	pubHealthService := pubsub.NewPubSubHealthService(
		"pubsubPublisher",
		producer.Client,
		10 * time.Second,
		pubsub.PermissionPublish,
		root.PubSubConsumer.SubscriptionId,
	)
	healthServices := []health.HealthService{mongoHealthService, subHealthService, pubHealthService}
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
