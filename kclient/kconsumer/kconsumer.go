package kconsumer

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sync"
)

const pollTimeout = 100

type MessageConsumer interface {
	Consume(message *kafka.Message) error
}

type KConsumer struct {
	broker string
	topic  string
	group  string

	closeChan chan bool
	wg        *sync.WaitGroup

	messageConsumer MessageConsumer
	kafkaConsumer   *kafka.Consumer
}

func NewKConsumer(broker, topic, group string, messageConsumer MessageConsumer) *KConsumer {
	return &KConsumer{
		broker:          broker,
		topic:           topic,
		group:           group,
		closeChan:       make(chan bool),
		wg:              &sync.WaitGroup{},
		messageConsumer: messageConsumer,
	}
}

func (c *KConsumer) Start() error {

	seelog.Infof("KConsumer starting - broker: %s  topic: %s  group: %s", c.broker, c.topic, c.group)

	var err error
	c.kafkaConsumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    c.broker,
		"group.id":             c.group,
		"session.timeout.ms":   6000,
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"},
	})
	if err != nil {
		return fmt.Errorf("NewConsumer error: %s\n", err)
	}

	err = c.kafkaConsumer.SubscribeTopics([]string{c.topic}, nil)
	if err != nil {
		return fmt.Errorf("SubscribeTopics error: %s\n", err)
	}

	go func() {
		c.wg.Add(1)
		defer c.wg.Done()

		c.pollCycle()
	}()

	seelog.Info("KConsumer started")

	return nil
}

func (c *KConsumer) pollCycle() {

	for {
		select {
		case <-c.closeChan:
			seelog.Debugf("KConsumer close signal")
			return

		default:
			ev := c.kafkaConsumer.Poll(pollTimeout)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				seelog.Infof("KConsumer message: %s  key: %s  value: %s",
					e.TopicPartition, string(e.Key), string(e.Value))

				if err := c.messageConsumer.Consume(e); err != nil {
					fmt.Errorf("messageConsumer error: %v", err)
				}

			case kafka.PartitionEOF:
				seelog.Infof("KConsumer PartitionEOF: %v", e)

			case kafka.Error:
				seelog.Warnf("KConsumer error: %v", e)

			default:
				seelog.Warnf("KConsumer unknown Poll result: %v", e)
			}
		}
	}
}

func (c *KConsumer) Close() error {

	seelog.Infof("KConsumer closing\n")
	defer seelog.Infof("KConsumer closed\n")

	c.closeChan <- true
	c.wg.Wait()

	if c.kafkaConsumer != nil {
		if err := c.kafkaConsumer.Close(); err != nil {
			return err
		}
	}

	return nil
}
