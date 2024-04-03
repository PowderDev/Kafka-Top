package main

import (
	"PowderDev/kafka-top/consumer"
	"PowderDev/kafka-top/websockets"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/valyala/fasthttp"
)

var wsHub = websockets.NewHub()

func main() {
	var config, envErr = getConfig()
	if envErr != nil {
		log.Fatalf("Error loading environment variables: %v\n", envErr)
	}

	var cons, consErr = kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers": config.KafkaHost,
			"group.id":          "top_events_consumer",
			"auto.offset.reset": "latest",
		},
	)
	if consErr != nil {
		log.Fatalf("Error creating Kafka consumer: %v\n", consErr)
	}

	r := getRouter()

	s := &fasthttp.Server{
		Handler: r.Handler,
	}

	go s.ListenAndServe("0.0.0.0:4000")
	go wsHub.Run()
	go consumer.Consume(cons, []string{"top-events"}, wsHub.HandleKafkaMessage)

	waitForSignals(s, wsHub, cons)
}

func waitForSignals(server *fasthttp.Server, wsHub *websockets.Hub, cons *kafka.Consumer) {
	signalCh := make(chan os.Signal, 1024)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGQUIT)

	for range signalCh {
		s := <-signalCh
		log.Printf("%v signal received.\n", s)

		if s == syscall.SIGQUIT || s == syscall.SIGINT {
			wsHub.Stop()
			cons.Close()
			err := server.Shutdown()

			if err != nil {
				log.Printf("Shutdown error: %v\n", err)
			}
		}
	}
}
