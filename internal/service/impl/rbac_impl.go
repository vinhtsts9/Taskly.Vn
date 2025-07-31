package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Taskly.com/m/internal/database"
	"Taskly.com/m/internal/service"
	"github.com/google/uuid"
)

type rbacService struct {
	store database.Store
}

func NewRBACService(store database.Store) service.IRBACService {
	return &rbacService{
		store: store,
	}
}

// CreatePermission implements service.RBACService.
func (s *rbacService) CreatePermission(ctx context.Context, arg database.CreatePermissionParams) (database.Permission, error) {
	return s.store.CreatePermission(ctx, arg)
}

// GetPermissionsByRoleID implements service.RBACService.
func (s *rbacService) GetPermissionsByRoleID(ctx context.Context, roleID int32) ([]database.Permission, error) {
	return s.store.GetPermissionsByRoleID(ctx, roleID)
}

// GetPermissionsByUserID implements service.RBACService.
func (s *rbacService) GetPermissionsByUserID(ctx context.Context, userID uuid.UUID) ([]database.Permission, error) {
	return s.store.GetPermissionsByUserID(ctx, userID)
}

// CreateRole implements service.RBACService.
func (s *rbacService) CreateRole(ctx context.Context, name string) (database.Role, error) {
	return s.store.CreateRole(ctx, name)
}

// GetRolesByUserID implements service.RBACService.
func (s *rbacService) GetRolesByUserID(ctx context.Context, userID uuid.UUID) ([]database.Role, error) {
	return s.store.GetRolesByUserID(ctx, userID)
}

// AddPermissionToRole implements service.RBACService.
func (s *rbacService) AddPermissionToRole(ctx context.Context, arg database.AddPermissionToRoleParams) error {
	return s.store.AddPermissionToRole(ctx, arg)
}

// RemovePermissionFromRole implements service.RBACService.
func (s *rbacService) RemovePermissionFromRole(ctx context.Context, arg database.RemovePermissionFromRoleParams) error {
	return s.store.RemovePermissionFromRole(ctx, arg)
}

// AddRoleToUser implements service.RBACService.
func (s *rbacService) AddRoleToUser(ctx context.Context, arg database.AddRoleToUserParams) error {
	return s.store.AddRoleToUser(ctx, arg)
}

// RemoveRoleFromUser implements service.RBACService.
func (s *rbacService) RemoveRoleFromUser(ctx context.Context, arg database.RemoveRoleFromUserParams) error {
	return s.store.RemoveRoleFromUser(ctx, arg)
}

// CheckPermission implements service.RBACService.
func (s *rbacService) CheckPermission(ctx context.Context, userID uuid.UUID, resource string, action string) (bool, error) {
	permissions, err := s.store.GetPermissionsByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // Không có quyền nào, coi như không có quyền truy cập
		}
		return false, fmt.Errorf("lỗi khi lấy quyền của người dùng: %w", err)
	}

	for _, p := range permissions {
		if p.Resource == resource && p.Action == action {
			return true, nil // Tìm thấy quyền hợp lệ
		}
	}

	return false, nil // Không tìm thấy quyền phù hợp
}
