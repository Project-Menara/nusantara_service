// package configs

// import (
// 	"log"
// 	"os"

// 	"github.com/streadway/amqp"
// )

// var RabbitConn *amqp.Connection
// var RabbitChannel *amqp.Channel

// func InitRabbitMQ() {
// 	amqpURL := os.Getenv("AMQP_URL")
// 	var err error

// 	RabbitConn, err = amqp.Dial(amqpURL)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
// 	}

// 	RabbitChannel, err = RabbitConn.Channel()
// 	if err != nil {
// 		log.Fatalf("Failed to open a channel: %v", err)
// 	}

// 	log.Println("RabbitMQ connected successfully")
// }

// func CloseConnections() {
// 	if RabbitChannel != nil {
// 		_ = RabbitChannel.Close()
// 	}
// 	if RabbitConn != nil {
// 		_ = RabbitConn.Close()
// 	}
// }

package configs

import (
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

var RabbitConn *amqp.Connection
var RabbitChannel *amqp.Channel

func InitRabbitMQ() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	var err error
	for i := 0; i < 5; i++ {
		RabbitConn, err = amqp.Dial(amqpURL)
		if err == nil {
			RabbitChannel, err = RabbitConn.Channel()
			if err == nil {
				log.Println("RabbitMQ connected successfully")
				return
			}
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in %d seconds: %v", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	log.Fatalf("Fatal: Could not connect to RabbitMQ after multiple retries")
}

func GetRabbitChannel() *amqp.Channel {
	// Periksa apakah koneksi RabbitMQ terputus. Jika ya, coba sambungkan kembali.
	if RabbitConn == nil || RabbitConn.IsClosed() {
		log.Println("RabbitMQ connection is closed. Reconnecting...")
		InitRabbitMQ()
	}

	// Periksa apakah channel RabbitMQ null. Jika ya, buat channel baru dari koneksi aktif.
	if RabbitChannel == nil {
		log.Println("RabbitMQ channel is not available. Opening a new channel...")
		var err error
		RabbitChannel, err = RabbitConn.Channel()
		if err != nil {
			log.Printf("Failed to open a new channel: %v", err)
			return nil
		}
	}
	return RabbitChannel
}

func CloseConnections() {
	if RabbitChannel != nil {
		_ = RabbitChannel.Close()
	}
	if RabbitConn != nil {
		_ = RabbitConn.Close()
	}
}
