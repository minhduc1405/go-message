package main

import (
	"context"
	"github.com/common-go/config"
	"github.com/common-go/health"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"

	c "go-message/internal"
)

func main() {

	parentPath := "go-message"
	resource := "configs"
	env := os.Getenv("ENV")

	var conf c.Root
	er1 := config.LoadConfig(parentPath, resource, env, &conf, "config")
	if er1 != nil {
		panic(er1)
	}
	logrus.SetLevel(logrus.DebugLevel)
	ctx := context.Background()

	app, er2 := c.NewApplicationContext(ctx, conf)
	if er2 != nil {
		panic(er2)
	}

	go serve(ctx, conf.Server, app.HealthController)

	app.BatchWorker.RunScheduler(ctx)
	app.Consumer.Consume(ctx, app.ConsumerCaller)
}

// Start a http server to serve HTTP requests
func serve(ctx context.Context, config c.ServerConfig, healthController *health.HealthController) {
	server := ""
	if config.Port > 0 {
		server = ":" + strconv.Itoa(config.Port)
	}
	http.HandleFunc("/health", healthController.Check)
	http.HandleFunc("/", healthController.Check)
	http.ListenAndServe(server, nil)
}
