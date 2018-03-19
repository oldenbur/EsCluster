package kproducer

import (
	"fmt"
	"time"
)

const (
	tsFormat = "2006-01-02T03:04:05.0000000"
)

type MessageProducer interface {
	Produce() []byte
}

type jsonProducer struct {
	counter int
}

func (j *jsonProducer) Produce() []byte {
	j.counter += 1
	return []byte(fmt.Sprintf(`{"ts": "%s", "id1": %d}`, time.Now().Format(tsFormat), j.counter))
}
