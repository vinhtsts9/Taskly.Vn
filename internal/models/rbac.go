package model

import "github.com/google/uuid"

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
	RoleID int32     `json:"role_id" binding:"required"`
}

type AddPermissionToRoleRequest struct {
	RoleID       int32 `json:"role_id" binding:"required"`
	PermissionID int32 `json:"permission_id" binding:"required"`
}
