package service

import (
	"context"

	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
)

type IAdminUserService interface {
	ListUsers(ctx context.Context, params database.AdminListUsersParams) ([]model.AdminListUsersResponse, error)
	CountUsers(ctx context.Context, query string) (int64, error)
}

var localAdminUserService IAdminUserService

func InitAdminUserService(s IAdminUserService) { localAdminUserService = s }
func GetAdminUserService() IAdminUserService {
	if localAdminUserService == nil {
		panic("admin user service not initialized")
	}
	return localAdminUserService
}
