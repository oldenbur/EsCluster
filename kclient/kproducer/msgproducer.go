package kproducer

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
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
	return []byte(fmt.Sprintf(`{"ts": "%s", "_id": "%s", "count": %d}`,
		time.Now().Format(tsFormat), uuid.NewV4(), j.counter))
}
