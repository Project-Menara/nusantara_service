package eventconsumer

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/cloudinary"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/data/repositories"
	"nusantara_service/internal/domain/entities"
	"nusantara_service/internal/workers/payload"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func StartImageEventUploadConsumer() {
	qn := rabbitmq.SendImageEventQueueName

	err := rabbitmq.ConsumeQueueAuto(qn, func(msg amqp.Delivery) {
		var payload payload.ImageEventUploadPayload
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			msg.Nack(false, false)
			return
		}

		eventRepo := repositories.NewEventRepositoryImpl(configs.DB)

		cld, err := cloudinary.NewCloudinaryService()
		if err != nil {
			msg.Nack(false, true)
			return
		}

		upload, err := cld.UploadImageBytes(context.Background(), bytes.NewReader(payload.FileBytes), payload.Folder, payload.Filename)
		if err != nil {
			msg.Nack(false, true)
			return
		}

		img := &entities.ImageEntity{
			ID:        uuid.New(),
			ImagePath: upload.URL,
		}

		if err := eventRepo.CreateImage(context.Background(), img); err != nil {
			msg.Nack(false, true)
			return
		}

		switch payload.Type {
		case "cover":
			if err := eventRepo.UpdateEventCover(context.Background(), payload.EventID, upload.URL); err != nil {
				log.Printf("[event-image] update cover error: %s", payload.Type)
			}
		default:
			log.Printf("[event-image] unknown type: %s", payload.Type)
		}
	})
	if err != nil {
		log.Fatalf("Failed to consume queue %s: %v", rabbitmq.SendImageShopQueueName, err)
	}

	log.Printf("Consumer for %s started. Waiting for messages...", rabbitmq.SendImageShopQueueName)
	forever := make(chan bool)
	<-forever
}
