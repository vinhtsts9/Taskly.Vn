package model

import (
	"time"

	"github.com/google/uuid"
)

type AdminListUsersResponse struct {
	ID         uuid.UUID `json:"id"`
	Names      string    `json:"names"`
	ProfilePic string    `json:"profile_pic"`
	CreatedAt  time.Time `json:"created_at"`
	Email      string    `json:"email"`
	States     int16     `json:"states"`
	RoleName   string    `json:"role_name"`
}

type AdminListUsersRequest struct {
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Query string `json:"query"`
}
