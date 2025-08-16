package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID           uuid.UUID  `json:"id"`
	GigID        uuid.UUID  `json:"gig_id"`
	BuyerID      uuid.UUID  `json:"buyer_id"`
	SellerID     uuid.UUID  `json:"seller_id"`
	Status       string     `json:"status"`        // 'pending', 'active', 'delivered', 'completed', 'cancelled'
	TotalPrice   float64    `json:"total_price"`   // float64 vì FLOAT8
	OrderDate    time.Time  `json:"order_date"`    // mặc định là CURRENT_TIMESTAMP
	DeliveryDate *time.Time `json:"delivery_date"` // nullable trong DB
}

type CreateOrderParams struct {
	GigID        uuid.UUID     `json:"gig_id" binding:"required"`
	SellerID     uuid.UUID     `json:"seller_id" binding:"required"`
	DeliveryDate *time.Time    `json:"delivery_date,omitempty"`                                    // optional
	PackageTier  string        `json:"package_tier" binding:"required,oneof=basic standard premium"` // Tên gói dịch vụ
	Answers      []AnswerParam `json:"answers,omitempty"`                                            // New field for answers
}
type AnswerParam struct {
	QuestionID uuid.UUID `json:"question_id" binding:"required"`
	Answer     string    `json:"answer" binding:"required"`
}


type UpdateOrderStatusParams struct {
	ID     uuid.UUID `json:"id" binding:"required"`
	Status string    `json:"status" binding:"required,oneof=pending active delivered completed cancelled"`
}

type OrderResult struct {
	ID           uuid.UUID  `json:"id"`
	GigID        uuid.UUID  `json:"gig_id"`
	BuyerID      uuid.UUID  `json:"buyer_id"`
	SellerID     uuid.UUID  `json:"seller_id"`
	Status       string     `json:"status"`
	TotalPrice   float64    `json:"total_price"`
	OrderDate    time.Time  `json:"order_date"`
	DeliveryDate *time.Time `json:"delivery_date,omitempty"`
}

