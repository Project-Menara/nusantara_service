package twilio

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func SendWhatsAppOTP(to string, otp string) error {
	from := os.Getenv("TWILIO_PHONE_NUMBER")
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	token := os.Getenv("TWILIO_AUTH_TOKEN")

	// Prepare the message
	msg := fmt.Sprintf("Kode OTP kamu adalah: %s\nHanya Berlaku Selama 1 menit", otp)
	msgData := url.Values{}
	msgData.Set("To", "whatsapp:"+to)
	msgData.Set("From", "whatsapp:"+from)
	msgData.Set("Body", msg)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Request to Twilio
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", sid)
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(sid, token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Optional: check response
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("failed to send WhatsApp message: %s", string(body))
}
