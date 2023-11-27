package main

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
)

type Config struct {
	Url string
}

type Producer struct {
	conn   *nats.Conn
	stream jetstream.JetStream
}

func NewProducer(cfg Config) (*Producer, error) {
	conn, err := nats.Connect(cfg.Url)
	if err != nil {
		return &Producer{}, errors.Wrap(err, "fail to init nats connection")
	}

	stream, err := jetstream.New(conn)
	if err != nil {
		return &Producer{}, errors.Wrap(err, "fail to init nats-jetstream")
	}

	return &Producer{
		conn:   conn,
		stream: stream,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, subject string, msg []byte) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		if _, err := p.stream.Publish(ctx, subject, msg); err != nil {
			return errors.Wrap(err, "fail to publish message")
		}
		return nil
	}
}
