package producer

import (
	"nusantara_service/internal/workers/consumer"
	shopconsumer "nusantara_service/internal/workers/consumer/Shop_Consumer"
)

func StartWorkers() {
	go consumer.ConsumeOTPQueue()
	go consumer.ConsumeVerifiedQueue()
	go consumer.ConsumeLinkForgotPINQueue()

	go consumer.ConsumeImageDeleteQueue()
	go consumer.ConsumeCacheInvalidateQueue()
	go consumer.SendImageQueue()
	go shopconsumer.StartImageUploadConsumer()
}
