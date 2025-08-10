package rabbitmq

import (
	"encoding/json"
	"fmt"
	"nusantara_service/configs"

	"github.com/streadway/amqp"
)

func PublishToQueue(exchangeName string, queueName string, payload interface{}) error {
	ch := configs.RabbitChannel

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName,
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.Publish(
		exchangeName, queueName, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
