package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// kafka reader interface for possibility of mocking
type Reader interface {
	ReadMessage(ctx context.Context) (Message, error)
	FetchMessage(ctx context.Context) (Message, error)
	CommitMessages(ctx context.Context, msgs ...Message) error
	Close() error
}

// concrete segmentio kafka reader
type KafkaReader struct {
	kr *kafka.Reader
}

func NewReader(cfg kafka.ReaderConfig) Reader {
	return &KafkaReader{kr: kafka.NewReader(cfg)}
}

// read message from kafka
func (kr *KafkaReader) ReadMessage(ctx context.Context) (Message, error) {
	// proxy call to segmentio method
	msg, err := kr.kr.ReadMessage(ctx)
	if err != nil {
		return Message{}, err
	}

	return messageFromKafka(msg), nil
}

// fetch messafe from kafka
func (kr *KafkaReader) FetchMessage(ctx context.Context) (Message, error) {
	msg, err := kr.kr.FetchMessage(ctx)
	if err != nil {
		return Message{}, err
	}

	return messageFromKafka(msg), nil
}

// commit messages into kafka
func (kr *KafkaReader) CommitMessages(ctx context.Context, msgs ...Message) error {
	m := make([]kafka.Message, len(msgs))

	for k, v := range msgs {
		m[k] = messageToKafka(v)
	}

	return kr.kr.CommitMessages(ctx, m...)
}

// close reader
func (kr *KafkaReader) Close() error {
	return kr.kr.Close()
}
