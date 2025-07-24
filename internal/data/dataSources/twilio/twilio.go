package twilio

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func SendWhatsAppOTP(to string, otp string) error {
	from := "+14155238886"
	sid := "AC75ae06a57ea07621e56d02c2f074e3b9"
	token := "792d46515d5f6ad2906ed51d932f3db3"

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
