package main

import (
	"github.com/cihub/seelog"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/oldenbur/esCluster/kclient/kconsumer"
	"github.com/oldenbur/esCluster/kclient/kproducer"
	"github.com/urfave/cli"
	"github.com/vrecan/death"
	"io"
	"os"
	"syscall"
)

const (
	cmdFlagKafkaBroker = "broker"
	cmdEnvKafkaBroker  = "KCLIENT_BROKER"

	cmdFlagKafkaTopic = "topic"
	cmdEnvKafkaTopic  = "KCLIENT_TOPIC"

	cmdFlagKafkaGroup = "group"
	cmdEnvKafkaGroup  = "KCLIENT_GROUP"
)

func main() {

	configureLogger()
	seelog.Info("kclient started")

	injector = &kclientInjector{
		newKConsumer: kconsumer.NewKConsumer,
	}

	app := cli.NewApp()
	app.Name = "KClient"
	app.Usage = "Kafka client with configurable roles and behavior"
	app.Commands = []cli.Command{
		{
			Name:  "consume",
			Usage: "register as a kafka consumer",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   cmdFlagKafkaBroker,
					Usage:  "kafka broker endpoint",
					EnvVar: cmdEnvKafkaBroker,
				},
				cli.StringFlag{
					Name:   cmdFlagKafkaTopic,
					Usage:  "kafka topic",
					EnvVar: cmdEnvKafkaTopic,
				},
				cli.StringFlag{
					Name:   cmdFlagKafkaGroup,
					Usage:  "kafka group",
					EnvVar: cmdEnvKafkaGroup,
				},
			},
			Action: consumeAction,
		},
		{
			Name:   "produce",
			Usage:  "produce stuff",
			Action: produceAction,
		},
	}
	app.Run(os.Args)

	seelog.Info("kclient complete")
}

type argsProducer interface {
	String(string) string
}

type StartCloser interface {
	Start() error
	io.Closer
}

func consumeAction(context *cli.Context) error { return consume(context) }
func produceAction(context *cli.Context) error { kproducer.Produce(); return nil }

func consume(args argsProducer) error {

	broker := args.String(cmdFlagKafkaBroker)
	if broker == "" {
		return seelog.Errorf("%s not defined", cmdFlagKafkaBroker)
	}

	topic := args.String(cmdFlagKafkaTopic)
	if topic == "" {
		return seelog.Errorf("%s not defined", cmdFlagKafkaTopic)
	}

	group := args.String(cmdFlagKafkaGroup)
	if group == "" {
		return seelog.Errorf("%s not defined", cmdFlagKafkaGroup)
	}

	kconsumer := injector.NewKConsumer(broker, topic, group, &nilConsumer{})
	err := kconsumer.Start()
	if err != nil {
		return seelog.Errorf("NewKConsumer error: %v", err)
	}

	death := death.NewDeath(syscall.SIGINT, syscall.SIGTERM)
	death.WaitForDeath(kconsumer)

	return nil
}

func configureLogger() {

	testConfig := `
        <seelog type="sync" minlevel="debug">
            <outputs formatid="main"><console/></outputs>
            <formats><format id="main" format="%Date %Time [%LEVEL] %Msg%n"/></formats>
        </seelog>`

	logger, err := seelog.LoggerFromConfigAsBytes([]byte(testConfig))
	if err != nil {
		panic(err)
	}
	err = seelog.ReplaceLogger(logger)
	if err != nil {
		panic(err)
	}

}

var injector KClientInjector

type KClientInjector interface {
	NewKConsumer(broker, topic, group string, c kconsumer.MessageConsumer) StartCloser
}

type kclientInjector struct {
	newKConsumer func(broker, topic, group string, c kconsumer.MessageConsumer) *kconsumer.KConsumer
}

func (i *kclientInjector) NewKConsumer(broker, topic, group string, c kconsumer.MessageConsumer) StartCloser {
	return i.newKConsumer(broker, topic, group, c)
}

type nilConsumer struct{}

func (*nilConsumer) Consume(message *kafka.Message) error {
	return nil
}
