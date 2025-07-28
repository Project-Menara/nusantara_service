package dto

type CreateBannerRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type UpdateBannerRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
