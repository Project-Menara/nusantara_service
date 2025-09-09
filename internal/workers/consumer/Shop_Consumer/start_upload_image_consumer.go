package shopconsumer

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

func StartImageUploadConsumer() {
	qn := rabbitmq.SendImageShopQueueName

	err := rabbitmq.ConsumeQueueAuto(qn, func(msg amqp.Delivery) {
		var payload payload.ImageUploadPayload
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			log.Printf("[Shop image] unmarshall error: %v", err)
			msg.Nack(false, false)
			return
		}

		shopRepo := repositories.NewShopRepositoryImpl(configs.DB)

		cld, err := cloudinary.NewCloudinaryService()
		if err != nil {
			log.Println("cloudinary init:", err)
			msg.Nack(false, true)
			return
		}

		upload, err := cld.UploadImageBytes(context.Background(), bytes.NewReader(payload.FileBytes), payload.Folder, payload.Filename)
		if err != nil {
			log.Printf("[shop-image] upload error: %v", err)
			msg.Nack(false, true)
			return
		}

		img := &entities.ImageEntity{
			ID:        uuid.New(),
			ImagePath: upload.URL,
		}

		if err := shopRepo.CreateImage(context.Background(), img); err != nil {
			log.Printf("[shop-image] create image error: %v", err)
			msg.Nack(false, true)
			return
		}

		switch payload.Type {
		case "cover":
			if err := shopRepo.UpdateShopCover(context.Background(), payload.ShopID, upload.URL); err != nil {
				log.Printf("[shop-image] update cover error: %v", err)
			}
		case "gallery":
			if err := shopRepo.CreateGallery(context.Background(), payload.ShopID, img.ID, payload.Filename); err != nil {
				log.Printf("[shop-image] create gallery error: %v", err)
			}
		default:
			log.Printf("[shop-image] unknown type: %s", payload.Type)
		}
	})
	if err != nil {
		log.Fatalf("Failed to consume queue %s: %v", rabbitmq.SendImageShopQueueName, err)
	}

	log.Printf("Consumer for %s started. Waiting for messages...", rabbitmq.SendImageShopQueueName)
	forever := make(chan bool)
	<-forever
}
