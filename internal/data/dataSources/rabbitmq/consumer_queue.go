package rabbitmq

import (
	"log"
	"nusantara_service/configs"

	"github.com/streadway/amqp"
)

func ConsumeQueueAuto(queueName string, handler func(amqp.Delivery)) error {
	msgs, err := configs.RabbitChannel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg)
		}
	}()

	log.Printf("[*] Waiting for messages in %s", queueName)

	return nil
}

func ConsumeQueueManual(queueName string, handler func(amqp.Delivery)) error {
	msgs, err := configs.RabbitChannel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg)
		}
	}()

	log.Printf("[*] Waiting for messages in %s", queueName)

	return nil
}
