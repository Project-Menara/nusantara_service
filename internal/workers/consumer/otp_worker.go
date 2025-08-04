package consumer

import (
	"encoding/json"
	"log"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/data/dataSources/twilio" // Pastikan jalur ini benar
	// Impor amqp untuk Ack/Nack
)

type OTPPayload struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type VerifiedPayload struct {
	Phone string `json:"phone"`
}

type LinkForgotPINPayload struct {
	Phone string `json:"phone"`
	Link  string `json:"link"`
}

func ConsumeOTPQueue() {
	// Pastikan antrean dideklarasikan di sisi konsumen juga, jika belum di deklarasi di produser
	_, err := configs.RabbitChannel.QueueDeclare(
		"otp_queue", // Nama antrean
		true,        // Durable
		false,       // Delete when unused
		false,       // Exclusive
		false,       // No-wait
		nil,         // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare otp_queue: %v", err)
	}

	msgs, err := configs.RabbitChannel.Consume(
		"otp_queue", // queue
		"",          // consumer
		false,       // auto-ack (ubah ke false)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer for otp_queue: %v", err)
	}

	for msg := range msgs {
		var payload OTPPayload
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			log.Printf("Error unmarshaling OTP message: %v. Message body: %s", err, msg.Body)
			msg.Nack(false, false) // Nack pesan jika gagal unmarshal
			continue
		}

		log.Printf("Received OTP message for phone: %s, code: %s", payload.Phone, payload.Code)
		if err := twilio.SendWhatsAppOTP(payload.Phone, payload.Code); err != nil {
			log.Printf("Failed to send OTP via WhatsApp for %s: %v", payload.Phone, err)
			msg.Nack(false, true) // Nack dan requeue jika gagal kirim OTP
		} else {
			log.Printf("OTP sent successfully to %s", payload.Phone)
			msg.Ack(false) // Ack pesan setelah berhasil
		}
	}
}

func ConsumeVerifiedQueue() {
	// Pastikan antrean dideklarasikan di sisi konsumen juga
	_, err := configs.RabbitChannel.QueueDeclare(
		"verified_queue", // Nama antrean yang BENAR
		true,             // Durable
		false,            // Delete when unused
		false,            // Exclusive
		false,            // No-wait
		nil,              // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare verified_queue: %v", err)
	}

	msgs, err := configs.RabbitChannel.Consume(
		"verified_queue", // queue yang BENAR
		"",               // consumer
		false,            // auto-ack (ubah ke false)
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer for verified_queue: %v", err)
	}

	for msg := range msgs {
		var payload VerifiedPayload
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			log.Printf("Error unmarshaling verified message: %v. Message body: %s", err, msg.Body)
			msg.Nack(false, false) // Nack pesan jika gagal unmarshal
			continue
		}

		log.Printf("Received Verified message for phone: %s", payload.Phone)
		message := "Nomor kamu telah berhasil diverifikasi."
		if err := twilio.SendWhatsAppMessage(payload.Phone, message); err != nil {
			log.Printf("Failed to send verification message via WhatsApp for %s: %v", payload.Phone, err)
			msg.Nack(false, true) // Nack dan requeue jika gagal kirim pesan verifikasi
		} else {
			log.Printf("Verification message sent successfully to %s", payload.Phone)
			msg.Ack(false) // Ack pesan setelah berhasil
		}
	}
}

func ConsumeLinkForgotPINQueue() {
	// Pastikan antrean dideklarasikan di sisi konsumen juga, jika belum di deklarasi di produser
	_, err := configs.RabbitChannel.QueueDeclare(
		rabbitmq.LinkForgotPINQueueName, // Nama antrean
		true,                            // Durable
		false,                           // Delete when unused
		false,                           // Exclusive
		false,                           // No-wait
		nil,                             // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare link forgot pin: %v", err)
	}

	msgs, err := configs.RabbitChannel.Consume(
		rabbitmq.LinkForgotPINQueueName, // queue
		"",                              // consumer
		false,                           // auto-ack (ubah ke false)
		false,                           // exclusive
		false,                           // no-local
		false,                           // no-wait
		nil,                             // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer for link forgot pin: %v", err)
	}

	for msg := range msgs {
		var payload LinkForgotPINPayload
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			log.Printf("Error unmarshaling link message: %v. Message body: %s", err, msg.Body)
			msg.Nack(false, false) // Nack pesan jika gagal unmarshal
			continue
		}

		log.Printf("Received link message for phone: %s, code: %s", payload.Phone, payload.Link)
		if err := twilio.SendWhatsAppMessage(payload.Phone, payload.Link); err != nil {
			log.Printf("Failed to send link via WhatsApp for %s: %v", payload.Phone, err)
			msg.Nack(false, true) // Nack dan requeue jika gagal kirim OTP
		} else {
			log.Printf("Link sent successfully to %s", payload.Phone)
			msg.Ack(false) // Ack pesan setelah berhasil
		}
	}
}
