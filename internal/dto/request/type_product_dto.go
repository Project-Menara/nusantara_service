package dto

type CreateTypeProductRequest struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type UpdateTypeProductRequest struct {
	Name string `json:"name"`
}

type UpdateStatusTypeProductRequest struct {
	Status int `json:"status"`
}
