package model

import (
	"time"

	"github.com/google/uuid"
)

// Bảng: user_base
type UserBase struct {
	UserBaseID         uuid.UUID  `json:"user_base_id"`          // PK
	Passwords          string     `json:"passwords"`             // hashed password
	IsTwoFactorEnabled int16      `json:"is_two_factor_enabled"` // 0 or 1
	LoginTime          *time.Time `json:"login_time,omitempty"`  // nullable
	LogoutTime         *time.Time `json:"logout_time,omitempty"` // nullable
	LoginIP            string     `json:"login_ip"`              // default ''
	States             int16      `json:"states"`                // CHECK >= 0
	CreatedAt          time.Time  `json:"created_at"`            // default now()
	UpdatedAt          time.Time  `json:"updated_at"`            // default now()
}

// Bảng: users
type User struct {
	ID         uuid.UUID `json:"id"`                    // PK
	UserBaseID uuid.UUID `json:"user_base_id"`          // FK
	Names      string    `json:"user_names"`            // UNIQUE NOT NULL
	UserType   string    `json:"user_type"`             // CHECK in ('buyer','seller','both')
	ProfilePic *string   `json:"profile_pic,omitempty"` // nullable TEXT
	Bio        *string   `json:"bio,omitempty"`         // nullable TEXT
	CreatedAt  time.Time `json:"created_at"`            // default now()
	UpdatedAt  time.Time `json:"updated_at"`            // default now()
}

// Đóng gói tham số truyền vào transaction
type CreateUserTxParams struct {
	Email      string  `json:"email" binding:"required"`     // UNIQUE NOT NULL
	Password   string  `json:"password" binding:"required"`  // NOT NULL
	Names      string  `json:"names" binding:"required"`     // NOT NULL
	UserType   string  `json:"user_type" binding:"required"` // CHECK in ('buyer','seller','both')
	ProfilePic *string `json:"profile_pic,omitempty"`        // nullable
	Bio        *string `json:"bio,omitempty"`                // nullable
}

type UserToken struct {
	ID       uuid.UUID `json:"id"`
	UserType string    `json:"user_type"`
}

// Kết quả trả ra
type CreateUserTxResult struct {
	UserBaseID uuid.UUID
}
type LoginInPut struct {
	Email    string `json:"email" binding:"required,email"` // UNIQUE NOT NULL
	Password string `json:"password" binding:"required"`    // NOT NULL
	LoginIP  string `json:"login_ip,omitempty"`             // NOT NULL
}
type LoginOutput struct {
	Token        string `json:"token"`         // JWT token
	RefreshToken string `json:"refresh_token"` // Refresh token
}
type UpdateLogoutInfoParams struct {
	UserBaseID uuid.UUID `json:"user_base_id"`
}

type DeleteUserBaseParams struct {
	UserBaseID uuid.UUID `json:"user_base_id"`
}
type ListUsersByTypeParams struct {
	UserType string `json:"user_type"`
}
type UpdatePasswordRegisterParams struct {
	VerifyKey    string  `json:"verify_key"`
	UserPassword string  `json:"user_password"`         // FK
	Names        string  `json:"user_names"`            // UNIQUE NOT NULL
	UserType     string  `json:"user_type"`             // CHECK in ('buyer','seller','both')
	ProfilePic   *string `json:"profile_pic,omitempty"` // nullable TEXT
	Bio          *string `json:"bio,omitempty"`
}
