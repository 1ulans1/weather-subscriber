package messaging

import (
	"encoding/json"
	"subscription-service/internal/core/ports"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	ch *amqp.Channel
}

func NewRabbitMQPublisher(url string) (ports.EmailPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare("confirm_queue", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare("notify_queue", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &RabbitMQPublisher{ch: ch}, nil
}

func (p *RabbitMQPublisher) PublishConfirm(email, token string) error {
	msg := ConfirmMessage{
		Email: email,
		Token: token,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.ch.Publish("", "confirm_queue", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

type ConfirmMessage struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (p *RabbitMQPublisher) PublishNotification(email, city, temperature, condition, unsubToken string) error {
	msg := NotifyMessage{
		Email: email,
		City:  city,
		Weather: struct {
			Temperature string `json:"temperature"`
			Condition   string `json:"condition"`
		}{
			Temperature: temperature,
			Condition:   condition,
		},
		UnsubToken: unsubToken,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.ch.Publish("", "notify_queue", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

type NotifyMessage struct {
	Email   string `json:"email"`
	City    string `json:"city"`
	Weather struct {
		Temperature string `json:"temperature"`
		Condition   string `json:"condition"`
	} `json:"weather"`
	UnsubToken string `json:"unsub_token"`
}
