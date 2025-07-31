package model

import (
	"github.com/google/uuid"
)

type Gig struct {
	ID           uuid.UUID   `json:"id"`
	UserID       uuid.UUID   `json:"user_id"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	CategoryID   int32       `json:"category_id"`
	Price        float64     `json:"price"`
	DeliveryTime int32       `json:"delivery_time"`
	ImageURL     *string     `json:"image_url"`
	Status       string      `json:"status"`
	CreatedAt    interface{} `json:"created_at"`
}

type SellerInfo struct {
	UserID         uuid.UUID `json:"user_id"`
	UserName       string    `json:"user_name"`
	UserProfilePic *string   `json:"user_profile_pic"`
}

type CategoryInfo struct {
	CategoryName string `json:"category_name"`
}

type GigDetailDTO struct {
	Gig
	SellerInfo
	CategoryInfo
}

type ListServicesParams struct {
	Search     *string `json:"search"`
	CategoryID *int32  `json:"category_id"`
	Status     string  `json:"status"`
	Limit      int32   `json:"limit"`
	Offset     int32   `json:"offset"`
}

type CreateServiceParams struct {
	UserID       uuid.UUID `json:"user_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CategoryID   int32     `json:"category_id"`
	Price        float64   `json:"price"`
	DeliveryTime int32     `json:"delivery_time"`
	ImageURL     *string   `json:"image_url"` // viết hoa cho đồng bộ với struct Gig
	Status       string    `json:"status"`
}

type UpdateServiceParams struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CategoryID   int32     `json:"category_id"`
	Price        float64   `json:"price"`
	DeliveryTime int32     `json:"delivery_time"`
	ImageURL     *string   `json:"image_url"`
	Status       string    `json:"status"`
}
