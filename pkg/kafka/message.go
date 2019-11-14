package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Topic string

	Partition int
	Offset    int64
	Key       []byte
	Value     []byte
	Headers   []Header

	Time time.Time
}

type Header struct {
	Key   string
	Value []byte
}

func messageFromKafka(msg kafka.Message) Message {
	h := make([]Header, len(msg.Headers))

	for hk, hv := range msg.Headers {
		h[hk] = Header{
			Key:   hv.Key,
			Value: hv.Value,
		}
	}

	return Message{
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Key:       msg.Key,
		Value:     msg.Value,
		Headers:   h,
		Time:      msg.Time,
	}
}

func messageToKafka(msg Message) kafka.Message {
	h := make([]kafka.Header, len(msg.Headers))
	for hk, hv := range msg.Headers {
		h[hk] = kafka.Header{
			Key:   hv.Key,
			Value: hv.Value,
		}
	}

	return kafka.Message{
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
		Key:       msg.Key,
		Value:     msg.Value,
		Headers:   h,
		Time:      msg.Time,
	}
}
