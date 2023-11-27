package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dany-ykl/logger"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"wb_test_task/consumer/internal/common"
	"wb_test_task/consumer/internal/config"
	"wb_test_task/consumer/internal/domain"
	"wb_test_task/libs/model"
)

type Msg struct {
	Subject string
	Data    []byte
}

type Consumer struct {
	conn           *nats.Conn
	js             jetstream.JetStream
	stream         jetstream.Stream
	consumer       jetstream.Consumer
	onMessage      func(ctx context.Context, msg *Msg)
	orderService   orderService
	countConsumers int
}

func New(cfg config.NatsConsumer, service orderService) (*Consumer, error) {
	nc, err := nats.Connect(
		cfg.Url,
		nats.RetryOnFailedConnect(cfg.RetryOfFailedConnect),
	)
	if err != nil {
		return &Consumer{}, errors.Wrap(err, "fail to connect to nats")
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return &Consumer{}, errors.Wrap(err, "fail to create jetstream")
	}

	stream, err := js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name:     cfg.StreamName,
		Subjects: cfg.Subjects,
	})
	if err != nil {
		return &Consumer{}, errors.Wrap(err, "fail to create stream")
	}

	consumer, err := stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Durable:   fmt.Sprintf("%sconsumer", cfg.StreamName),
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return &Consumer{}, errors.Wrap(err, "fail to create consumer")
	}

	return &Consumer{
		conn:           nc,
		js:             js,
		stream:         stream,
		consumer:       consumer,
		orderService:   service,
		countConsumers: cfg.CountConsumers,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	g, ctxG := errgroup.WithContext(ctx)

	for i := 0; i < c.countConsumers; i++ {
		id := i
		g.Go(func() error {
			logger.Info("nats consumer is starting", zap.Int("id", id))
			cc, err := c.consumer.Consume(func(msg jetstream.Msg) {
				if err := c.OnMessage(ctxG, &Msg{
					Subject: msg.Subject(),
					Data:    msg.Data(),
				}); err != nil {
					var wrapErr common.WrapError
					if errors.As(err, &wrapErr) {
						logger.Warn("fail to handle message", zap.String("errorType", wrapErr.Error()), zap.String("error", wrapErr.Message()))
					} else {
						logger.Warn("fail to handle message", zap.Error(err))
					}
				}

				if err := msg.Ack(); err != nil {
					logger.Warn("fail to ack message", zap.Error(err))
				}
			})
			if err != nil {
				return errors.Wrap(err, "fail to start consume")
			}
			defer cc.Stop()

			select {
			case <-ctxG.Done():
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		return errors.Wrap(err, "fail to wait")
	}

	return nil
}

func (c *Consumer) OnMessage(ctx context.Context, msg *Msg) error {
	switch msg.Subject {
	case "order.create":
		return c.orderCreateHandler(ctx, msg)
	default:
		return errors.Wrap(errors.New("fail to onMessage msg"), msg.Subject)
	}
}

//go:generate mockgen -source=consumer.go -destination=mocks/mock.go
type orderService interface {
	Create(ctx context.Context, request *domain.OrderCreateRequest) (*model.Order, error)
}

func (c *Consumer) orderCreateHandler(ctx context.Context, msg *Msg) error {
	var orderCrateRequest domain.OrderCreateRequest
	if err := json.Unmarshal(msg.Data, &orderCrateRequest); err != nil {
		return errors.Wrap(err, "fail to unmarshal msg")
	}

	if _, err := c.orderService.Create(ctx, &orderCrateRequest); err != nil {
		return errors.Wrap(err, "fail to create order")
	}

	return nil
}

func (c *Consumer) Shutdown() error {
	c.conn.Close()
	return nil
}
