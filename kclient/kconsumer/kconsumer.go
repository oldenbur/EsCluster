package kconsumer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"io"
	"os"
	"os/signal"
	"syscall"
)

const (
	consumerGroup = "group1"
)

type Consumer interface {
	Consume(message *kafka.Message) error
}

type KConsumer struct {
	broker string
	topic  string
}

func NewKConsumer(broker, topic string) *KConsumer {
	return &KConsumer{}
}

func (c *KConsumer) Start() error {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    c.broker,
		"group.id":             consumerGroup,
		"session.timeout.ms":   6000,
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
	})
	if err != nil {
		fmt.Errorf("NewConsumer error: %s\n", err)
	}

	err = consumer.SubscribeTopics([]string{c.topic}, nil)
	if err != nil {
		fmt.Errorf("SubscribeTopics error: %s\n", err)
	}

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\nKey: %s  Value: %s\n", e.TopicPartition, string(e.Key), string(e.Value))
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	consumer.Close()

	return nil
}

func (c *KConsumer) Close() error {
	return nil
}
