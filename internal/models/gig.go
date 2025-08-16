package model

import (
	"github.com/google/uuid"
)

type Gig struct {
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"user_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CategoryID  []int32     `json:"category_id"`
	ImageURL    *[]string   `json:"image_url"`
	PricingMode string      `json:"pricing_mode" validate:"oneof=single triple"`
	Status      string      `json:"status"`
	CreatedAt   interface{} `json:"created_at"`
	UpdatedAt   interface{} `json:"updated_at"`
}

type SellerInfo struct {
	UserName       string  `json:"user_name"`
	UserProfilePic *string `json:"user_profile_pic"`
}

type CategoryInfo struct {
	CategoryName []string `json:"category_name"`
}

type GigDetailDTO struct {
	Gig
	SellerInfo
	CategoryInfo
	GigPackage []GigPackage
	Question   []Question
}

type ListServicesParams struct {
	Search     *string  `json:"search"`
	CategoryID *[]int32 `json:"category_id"`
	Status     string   `json:"status"`
	Limit      int32    `json:"limit"`
	Offset     int32    `json:"offset"`
}

type UpdateServiceParams struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CategoryID  []int32   `json:"category_id"`
	ImageURL    []string  `json:"image_url"`
	Status      string    `json:"status"`
}
type GigPackageOptions struct {
	Revisions int32 `json:"revisions"`
	Files     int32 `json:"files"`
}

type GigPackage struct {
	Tier         string            `json:"tier" binding:"required"`
	Price        float64           `json:"price" binding:"required,min=0"`
	DeliveryDays int32             `json:"delivery_days" binding:"required,min=1"`
	Options      GigPackageOptions `json:"options"`
}

type Question struct {
	ID       uuid.UUID `json:"id"`
	Question string `json:"question" binding:"required"`
	Required bool   `json:"required"`
}
type CreateServiceParams struct {
	UserID       uuid.UUID    `json:"user_id"`
	Title        string       `json:"title" binding:"required,min=10"`
	Description  string       `json:"description" binding:"required,min=30"`
	CategoryID   []int32      `json:"category_id" binding:"required"`
	ImageURL     []string     `json:"image_url"`
	PricingMode  string       `json:"pricing_mode" binding:"required,oneof=single triple"`
	Packages     []GigPackage `json:"packages" binding:"required,min=1"`
	Requirements struct {
		Questions []Question `json:"questions"`
	} `json:"requirements"`
	Status string `json:"status" binding:"required,oneof=draft active"`
}
type GetCategoriesRs struct {
	ParentID     int32   `json:"parent_id"`
	ParentName   string  `json:"parent_name"`
	ChildrenID   *int32  `json:"children_id"`
	ChildrenName *string `json:"children_name"`
}

type SearchGigParams struct {
	SearchTerm  string    `json:"search_term" form:"search_term"`
	MinPrice    *float64  `json:"min_price" form:"min_price"`
	MaxPrice    *float64  `json:"max_price" form:"max_price"`
	CategoryIDs []int32   `json:"category_ids" form:"category_ids"`
	LastGigID   uuid.UUID `json:"last_gig_id" form:"last_gig_id"`
}

type SearchGigDTO struct {
	ID          uuid.UUID   `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ImageURL    *[]string   `json:"image_url"`
	PricingMode string      `json:"pricing_mode"`
	CreatedAt   interface{} `json:"created_at"`
	UpdatedAt   interface{} `json:"updated_at"`
	BasicPrice  float64     `json:"basic_price"`
}
