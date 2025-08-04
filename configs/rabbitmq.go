package configs

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func InitRabbitMQ() {
	amqpURL := os.Getenv("AMQP_URL")
	var err error

	RabbitConn, err = amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	RabbitChannel, err = RabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	log.Println("RabbitMQ connected successfully")
}

func CloseConnections() {
	if RabbitChannel != nil {
		_ = RabbitChannel.Close()
	}
	if RabbitConn != nil {
		_ = RabbitConn.Close()
	}
}
