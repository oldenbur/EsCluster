package main

import (
	"github.com/cihub/seelog"
	"github.com/urfave/cli"
	"github.com/vrecan/death"
	"io"
	"os"
)

const (
	cmdFlagKafkaBroker = "broker"
	cmdEnvKafkaBroker  = "KCLIENT_BROKER"

	cmdFlagKafkaTopic = "topic"
	cmdEnvKafkaTopic  = "KCLIENT_TOPIC"
)

func main() {

	configureLogger()
	seelog.Info("kclient started")

	app := cli.NewApp()
	app.Name = "ConfigControl"
	app.Usage = "Tool to control merging config file changes during an upgrade"
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
			},
			Action: consume,
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

var injector KClientInjector

type KClientInjector interface {
	NewKConsumer(broker, topic string) StartCloser
}

type kclientInjector struct {
	newKConsumer func(broker, topic string) StartCloser
}

func (i *kclientInjector) NewKConsumer(broker, topic string) StartCloser {
	return i.newKConsumer(broker, topic)
}

func consume(args argsProducer) error {

	broker := args.String(cmdFlagKafkaBroker)
	if broker == "" {
		return seelog.Errorf("%s not defined", cmdFlagKafkaBroker)
	}

	topic := args.String(cmdFlagKafkaTopic)
	if topic == "" {
		return seelog.Errorf("%s not defined", cmdFlagKafkaTopic)
	}

	death := death.NewDeath(sys.SIGINT, sys.SIGTERM) //pass the signals you want to end your application
	objects := make([]io.Closer, 0)

	kconsumer := injector.NewKConsumer(broker, topic)
	err := kconsumer.Start()
	if err != nil {
		return seelog.Errorf("NewKConsumer error: %v", err)
	}

	objects = append(objects, kconsumer)

	death.WaitForDeath(objects...) // this will finish when a signal of your type is sent to your application
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
