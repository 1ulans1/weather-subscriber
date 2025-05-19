package messaging

import (
	"context"
	"email-service/internal/core/domain"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	ch *amqp.Channel
}

func NewRabbitMQConsumer(url string) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	ch.QueueDeclare("confirm_queue", true, false, false, false, nil)
	ch.QueueDeclare("notify_queue", true, false, false, false, nil)
	return &RabbitMQConsumer{ch: ch}, nil
}

func (c *RabbitMQConsumer) ConsumeConfirm(ctx context.Context, handler func(msg *domain.ConfirmMessage) error) error {
	msgs, err := c.ch.Consume("confirm_queue", "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			var m domain.ConfirmMessage
			if err := json.Unmarshal(d.Body, &m); err != nil {
				continue
			}
			handler(&m)
		}
	}()
	return nil
}

func (c *RabbitMQConsumer) ConsumeNotification(ctx context.Context, handler func(msg *domain.NotifyMessage) error) error {
	msgs, err := c.ch.Consume("notify_queue", "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			var m domain.NotifyMessage
			if err := json.Unmarshal(d.Body, &m); err != nil {
				continue
			}
			handler(&m)
		}
	}()
	return nil
}
