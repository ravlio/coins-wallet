package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// kafka reader interface for possibility of mocking
type Writer interface {
	WriteMessages(ctx context.Context, msgs ...Message) error
	Close() error
}

type KafkaWriter struct {
	kw *kafka.Writer
}

func NewWriter(cfg kafka.WriterConfig) Writer {
	kw := kafka.NewWriter(cfg)
	return &KafkaWriter{kw: kw}
}

func (kw *KafkaWriter) WriteMessages(ctx context.Context, msgs ...Message) error {
	m := make([]kafka.Message, len(msgs))

	for k, v := range msgs {
		m[k] = messageToKafka(v)
	}
	return kw.kw.WriteMessages(ctx, m...)
}

func (kw *KafkaWriter) Close() error {
	return kw.kw.Close()
}
