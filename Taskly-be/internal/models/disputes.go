package model

import (
	"time"

	"github.com/google/uuid"
)

type Dispute struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	UserID    uuid.UUID `json:"user_id"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateDisputeParams struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
	UserID  uuid.UUID `json:"user_id" binding:"required"`
	Reason  string    `json:"reason" binding:"required"`
}

type UpdateDisputeStatusParams struct {
	ID     uuid.UUID `json:"id" binding:"required"`
	Status string    `json:"status" binding:"required,oneof=refunded rejected resolved under_review"`
}
