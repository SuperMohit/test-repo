package main

import (
	"log"
	"os"
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/cloudevents/sdk-go"
)

type envConfig struct {
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}
	os.Exit(_main(os.Args[1:], env))
}

func _main(args []string, env envConfig) int {
	ctx := context.Background()

	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithPort(env.Port),
		cloudevents.WithPath(env.Path),
	)
	if err != nil {
		log.Printf("failed to create transport, %v", err)
		return 1
	}
	c, err := cloudevents.NewClient(t)
	if err != nil {
		log.Printf("failed to create client, %v", err)
		return 1
	}

	log.Printf("will listen on :%d%s\n", env.Port, env.Path)
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, Event))

	return 0
}