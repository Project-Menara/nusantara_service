package producer

import "nusantara_service/internal/workers/consumer"

func StartWorkers() {
	go consumer.ConsumeOTPQueue()
	go consumer.ConsumeVerifiedQueue()
	go consumer.ConsumeLinkForgotPINQueue()

	go consumer.ConsumeImageDeleteQueue()
	go consumer.ConsumeCacheInvalidateQueue()
	go consumer.SendImageQueue()
}
