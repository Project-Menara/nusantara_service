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

type UpdateCustomerRequest struct {
	Name        string `json:"name,omitempty" form:"name"`
	Username    string `json:"username,omitempty" form:"username"`
	Email       string `json:"email,omitempty" form:"email"`
	Gender      string `json:"gender,omitempty" form:"gender"`
	DateOfBirth string `json:"date_of_birth,omitempty" form:"date_of_birth"`
}

type VerifyPINCustomerRequest struct {
	PIN string `json:"pin"`
}

type NewPINCustomer struct {
	NewPIN string `json:"new_pin"`
}

type ConfirmNewPINCustomer struct {
	ConfirmPIN string `json:"confirm_pin"`
}

type NewPhoneCustomerRequest struct {
	Phone string `json:"new_phone"`
}

type VerifyOTPCustomerUpdateRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}
