package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/data/repositories"
	"nusantara_service/internal/workers/payload"

	"github.com/streadway/amqp"
)

func SendImageQueue() {
	err := rabbitmq.ConsumeQueueManual(rabbitmq.SendImageQueueName, func(msg amqp.Delivery) {
		log.Printf("Received a message from %s: %s", rabbitmq.SendImageQueueName, string(msg.Body))
		var payload payload.ImageSendTask
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			log.Printf("Unmarshal error for %s: %v", rabbitmq.SendImageQueueName, err)
			msg.Nack(false, false)
			return
		}

		cashierRepo := repositories.NewCashierRepositoryImpl(configs.DB)

		cld, err := cloudinary.NewCloudinaryService()
		if err != nil {
			log.Println("cloudinary init:", err)
			msg.Nack(false, true)
			return
		}

		ctx := context.Background()
		log.Printf("image upload for %s/%s", payload.Folder, payload.Filename)

		uploadResult, err := cld.UploadImageBytes(ctx, bytes.NewReader(payload.FileBytes), payload.Folder, payload.Filename)
		if err != nil {
			log.Printf("Cloudinary upload failed for %s: %v", payload.Filename, err)
			msg.Nack(false, true)
			return
		}

		cashier, err := cashierRepo.FindById(ctx, payload.UserID)
		if err != nil {
			log.Printf("Failed to find cashier with ID %s: %v", payload.UserID, err)
			msg.Nack(false, true)
			return
		}

		cashier.Photo = &uploadResult.URL
		err = cashierRepo.Update(ctx, payload.UserID, cashier)
		if err != nil {
			log.Printf("Failed to update cashier photo for ID %s: %v", payload.UserID, err)
			msg.Nack(false, true)
			return
		}

		log.Printf("Successfully processed image task for user %s. Photo URL: %s", payload.UserID, uploadResult.URL)

		msg.Ack(false)
	})
	if err != nil {
		log.Fatalf("Failed to consume queue %s: %v", rabbitmq.SendImageQueueName, err)
	}

	log.Printf("Consumer for %s started. Waiting for messages...", rabbitmq.SendImageQueueName)
	forever := make(chan bool)
	<-forever
}
