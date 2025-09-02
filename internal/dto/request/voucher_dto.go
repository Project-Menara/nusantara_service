package dto

import "time"

type CreateVoucherRequest struct {
	Code            string    `json:"code"`
	DiscountAmount  int       `json:"discount_amount"`
	DiscountPercent int       `json:"discount_percent"`
	MinimumSpend    int       `json:"minimum_spend"`
	PointCost       int       `json:"point_cost"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Quota           int       `json:"quota"`
	Description     string    `json:"description"`
	DiscountType    string    `json:"discount_type"`
	Status          int       `json:"status"`
}

type UpdateVoucherRequest struct {
	Code            string    `json:"code"`
	DiscountAmount  int       `json:"discount_amount"`
	DiscountPercent int       `json:"discount_percent"`
	MinimumSpend    int       `json:"minimum_spend"`
	PointCost       int       `json:"point_cost"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Quota           int       `json:"quota"`
	Description     string    `json:"description"`
	DiscountType    string    `json:"discount_type"`
}

type UpdateStatusVoucherRequest struct {
	Status int `json:"status"`
}
