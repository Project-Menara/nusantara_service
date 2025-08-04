package producer

import "nusantara_service/internal/workers/consumer"

func StartWorkers() {
	go consumer.ConsumeOTPQueue()
	go consumer.ConsumeVerifiedQueue()
	go consumer.ConsumeLinkForgotPINQueue()
}
