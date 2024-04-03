package main

import (
	"errors"
	"os"
)

type Config struct {
	KafkaHost string
}

func getConfig() (*Config, error) {
	host := os.Getenv("KAFKA_HOST")

	if host == "" {
		return nil, errors.New("KAFKA_HOST is not set")
	}

	return &Config{
		KafkaHost: host,
	}, nil
}
