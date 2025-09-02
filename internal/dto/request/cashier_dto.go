package dto

import "mime/multipart"

type CreateCashierRequest struct {
	Name     string                `form:"name" json:"name"`
	Username string                `form:"username" json:"username"`
	Email    string                `form:"email" json:"email"`
	Password string                `form:"password" json:"password"`
	Photo    *multipart.FileHeader `form:"photo" json:"photo"`
	Status   int                   `form:"status" json:"status"`
}

type UpdateCashierRequest struct {
	Name     *string               `form:"name" json:"name"`
	Username *string               `form:"username" json:"username"`
	Email    *string               `form:"email" json:"email"`
	Photo    *multipart.FileHeader `form:"photo" json:"photo"`
	Status   *int                  `form:"status" json:"status"`
}
