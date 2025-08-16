package model

import (
	"time"

	"github.com/google/uuid"
)

// Role model
type Role struct {
	ID        uuid.UUID `json:"id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Permission model
type Permission struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Resource  string    `json:"resource"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRole (bảng trung gian user_roles)
type UserRole struct {
	UserID    uuid.UUID `json:"user_id"`
	RoleID    uuid.UUID `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// RolePermission (bảng trung gian role_permissions)
type RolePermission struct {
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
}

// Structs cho Roles
type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

// Structs cho Permissions
type CreatePermissionRequest struct {
	Name     string `json:"name" binding:"required"`
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required"`
}

// Structs cho việc gán quyền
type AddRoleToUserRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

type AddPermissionToRoleRequest struct {
	RoleID       uuid.UUID `json:"role_id" binding:"required"`
	PermissionID uuid.UUID `json:"permission_id" binding:"required"`
}
