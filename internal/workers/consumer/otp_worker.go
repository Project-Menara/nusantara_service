package consumer

import (
	"encoding/json" // Import errors package
	"log"
	"nusantara_service/configs"
	"nusantara_service/internal/data/dataSources/rabbitmq"
	"nusantara_service/internal/data/dataSources/twilio" // Pastikan jalur ini benar
	"strings"                                            // Import strings package
)

// Asumsi: Jika package twilio Anda tidak mengembalikan error yang terstruktur,
// kita akan mencoba mencari substring error code di pesan error.
// Cara yang lebih robust adalah dengan memodifikasi package twilio
// untuk mengembalikan struct error kustom (contoh: twilio.Error)
// yang berisi code dan message, sehingga bisa dicek langsung.
const twilioRateLimitErrorCode = "63038"

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
	log.Println("[*] Waiting for messages in otp_queue")
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
			msg.Nack(false, false) // Nack pesan tanpa requeue jika gagal unmarshal (pesan rusak)
			continue
		}

		log.Printf("Received OTP message for phone: %s, code: %s", payload.Phone, payload.Code)
		if err := twilio.SendWhatsAppOTP(payload.Phone, payload.Code); err != nil {
			log.Printf("Failed to send OTP via WhatsApp for %s: %v", payload.Phone, err)

			// Cek apakah error adalah karena rate limit Twilio
			// Ini adalah asumsi, jika package twilio Anda mengembalikan struct error kustom,
			// gunakan cara yang lebih spesifik untuk mengeceknya.
			if strings.Contains(err.Error(), twilioRateLimitErrorCode) {
				log.Printf("Twilio rate limit exceeded for %s. Nacking message without requeue.", payload.Phone)
				msg.Nack(false, false) // Nack tanpa requeue jika Twilio limit
			} else {
				// Untuk error lain (misal, masalah koneksi sementara), Nack dengan requeue
				log.Printf("Other error occurred for %s. Nacking message with requeue.", payload.Phone)
				msg.Nack(false, true) // Nack dan requeue jika gagal kirim OTP (error transient)
			}
		} else {
			log.Printf("OTP sent successfully to %s", payload.Phone)
			msg.Ack(false) // Ack pesan setelah berhasil
		}
	}
}

func ConsumeVerifiedQueue() {
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
			msg.Nack(false, false) // Nack pesan tanpa requeue jika gagal unmarshal
			continue
		}

		log.Printf("Received Verified message for phone: %s", payload.Phone)
		message := "Nomor kamu telah berhasil diverifikasi."
		if err := twilio.SendWhatsAppMessage(payload.Phone, message); err != nil {
			log.Printf("Failed to send verification message via WhatsApp for %s: %v", payload.Phone, err)

			// Cek apakah error adalah karena rate limit Twilio
			if strings.Contains(err.Error(), twilioRateLimitErrorCode) {
				log.Printf("Twilio rate limit exceeded for %s. Nacking message without requeue.", payload.Phone)
				msg.Nack(false, false) // Nack tanpa requeue jika Twilio limit
			} else {
				log.Printf("Other error occurred for %s. Nacking message with requeue.", payload.Phone)
				msg.Nack(false, true) // Nack dan requeue jika gagal kirim pesan verifikasi (error transient)
			}
		} else {
			log.Printf("Verification message sent successfully to %s", payload.Phone)
			msg.Ack(false) // Ack pesan setelah berhasil
		}
	}
}

func ConsumeLinkForgotPINQueue() {
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
			msg.Nack(false, false) // Nack pesan tanpa requeue jika gagal unmarshal
			continue
		}

		log.Printf("Received link message for phone: %s, code: %s", payload.Phone, payload.Link)
		if err := twilio.SendWhatsAppMessage(payload.Phone, payload.Link); err != nil {
			log.Printf("Failed to send link via WhatsApp for %s: %v", payload.Phone, err)

			// Cek apakah error adalah karena rate limit Twilio
			if strings.Contains(err.Error(), twilioRateLimitErrorCode) {
				log.Printf("Twilio rate limit exceeded for %s. Nacking message without requeue.", payload.Phone)
				msg.Nack(false, false) // Nack tanpa requeue jika Twilio limit
			} else {
				log.Printf("Other error occurred for %s. Nacking message with requeue.", payload.Phone)
				msg.Nack(false, true) // Nack dan requeue jika gagal kirim OTP (error transient)
			}
		} else {
			log.Printf("Link sent successfully to %s", payload.Phone)
			msg.Ack(false) // Ack pesan setelah berhasil
		}
	}
}
