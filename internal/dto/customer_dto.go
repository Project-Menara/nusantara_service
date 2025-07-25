package dto

type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type NewPinRequest struct {
	Phone string `json:"phone"`
	PIN   string `json:"pin"`
}

type ConfirmPinRequest struct {
	Phone      string `json:"phone"`
	ConfirmPIN string `json:"confirm_pin"`
}
