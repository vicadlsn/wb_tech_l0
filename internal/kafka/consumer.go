package kafka

import (
	"context"
	"log/slog"

	"webtechl0/internal/config"

	"github.com/segmentio/kafka-go"
)

type Handler interface {
	HandleMessage(ctx context.Context, data []byte) error
}

type Consumer struct {
	reader  *kafka.Reader
	handler Handler
	lg      *slog.Logger
}

func NewConsumer(cfg config.Kafka, handler Handler, lg *slog.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})
	return &Consumer{reader: reader, handler: handler, lg: lg}
}

func (c *Consumer) Start(ctx context.Context) error {
	lg := c.lg.With(slog.String("op", "Consumer.Start"))

	for {
		msg, err := c.reader.FetchMessage(ctx)

		if err != nil {
			if ctx.Err() != nil {
				lg.Info("Context cancelled, stopping consumer")
				return ctx.Err()
			}

			lg.Error("Failed to fetch message", slog.Any("error", err))
			continue
		}

		lg.Debug("Fetched message", slog.String("topic", msg.Topic), slog.Int64("offset", msg.Offset), slog.Int("partition", msg.Partition), slog.String("key", string(msg.Key)))

		if err := c.handler.HandleMessage(ctx, msg.Value); err != nil {
			lg.Error("Failed to handle message", slog.Any("error", err))
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			lg.Error("Failed to commit message", slog.Any("error", err))
		}
	}
}

func (c *Consumer) Stop() error {
	c.lg.Info("Closing kafka reader")
	return c.reader.Close()
}
