package service

import (
	"context"

	"Taskly.com/m/internal/database"
	"github.com/google/uuid"
)

// RBACService định nghĩa các phương thức cho tầng service.
type IRBACService interface {
	// Quản lý permissions
	CreatePermission(ctx context.Context, arg database.CreatePermissionParams) (database.Permission, error)
	GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]database.Permission, error)
	GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]database.Permission, error)

	// Quản lý roles
	CreateRole(ctx context.Context, name string) (database.Role, error)
	GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]database.Role, error)
	AddPermissionToRole(ctx context.Context, arg database.AddPermissionToRoleParams) error
	RemovePermissionFromRole(ctx context.Context, arg database.RemovePermissionFromRoleParams) error

	// Quản lý user-roles
	AddRoleToUser(ctx context.Context, arg database.AddRoleToUserParams) error
	RemoveRoleFromUser(ctx context.Context, arg database.RemoveRoleFromUserParams) error

	// Kiểm tra quyền
	CheckPermission(ctx context.Context, userID uuid.UUID, resource string, action string) (bool, error)
}

var localRBACService IRBACService

func InitRBACService(i IRBACService) {
	localRBACService = i
}

func GetRbacService() IRBACService {
	if localRBACService == nil {
		panic("implement rbacService failed")
	}
	return localRBACService
}
