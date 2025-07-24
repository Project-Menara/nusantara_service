package dto

type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
