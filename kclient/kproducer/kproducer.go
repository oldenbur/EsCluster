package kproducer

import (
	"context"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/time/rate"
	"os"
	"sync"
	"time"
)

const (
	burstSize = 1
)

type KProducer struct {
	doneChan chan bool
	wg       *sync.WaitGroup

	mps    uint
	broker string
	topic  string
	group  string
}

func NewKProducer(mps int, broker, topic, group string) *KProducer {
	return &KProducer{
		doneChan: make(chan bool),
		wg:       &sync.WaitGroup{},
		mps:      uint(mps),
		broker:   broker,
		topic:    topic,
		group:    group,
	}
}

func (p *KProducer) Start() error {

	seelog.Info("KProducer starting")
	prod, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": p.broker})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	go func() {

		p.wg.Add(1)
		defer p.wg.Done()

		last := time.Now()
		period := rate.Every(time.Second / time.Duration(p.mps))
		seelog.Debugf("converted mps: %d  to period: %v", p.mps, period)
		lim := rate.NewLimiter(period, burstSize)

		for {
			select {
			case <-p.doneChan:
				seelog.Info("KProducer complete")
				return

			case ev := <-prod.Events():
				p.handleKafkaEvent(ev)

			default:
				seelog.Debug("Waiting")
				ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
				if err := lim.Wait(ctx); err != nil {
					seelog.Warnf("KProducer Wait returned error: %v", err)
				}
				seelog.Infof("KProducer Wait returned - delay: %v", time.Since(last))
				last = time.Now()

				prod.ProduceChannel() <- &kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
					Value:          []byte(`howdy`),
				}
			}
		}
	}()

	return nil
}

func (p *KProducer) handleKafkaEvent(evInt kafka.Event) {

	var ev *kafka.Message
	var ok bool
	if ev, ok = evInt.(*kafka.Message); !ok {
		seelog.Warnf(`unexpected producer event type '%T' for event: %s`, evInt, evInt)
		return
	}

	if ev.TopicPartition.Error != nil {
		fmt.Printf("kafka delivery failed: %v\n", ev.TopicPartition.Error)
	} else {
		fmt.Printf("kafka delivered message to topic %s [%d] at offset %v\n",
			*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
	}
}

func (p *KProducer) Close() error {
	seelog.Info("KProducer closing")
	p.doneChan <- true
	p.wg.Wait()
	seelog.Info("KProducer closed")

	return nil
}
