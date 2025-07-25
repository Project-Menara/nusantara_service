package dto

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type UpdateRoleRequest struct {
	Name string `json:"name"`
}

type LoginCustomerRequest struct {
	Phone string `json:"phone"`
	Pin   string `json:"pin"`
}
