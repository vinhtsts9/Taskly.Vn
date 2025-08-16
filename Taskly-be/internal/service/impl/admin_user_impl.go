package impl

import (
	"context"

	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
)

type adminUserService struct{ store database.Store }

func NewAdminUserService(store database.Store) service.IAdminUserService {
	return &adminUserService{store: store}
}

func (s *adminUserService) ListUsers(ctx context.Context, params database.AdminListUsersParams) ([]model.AdminListUsersResponse, error) {

	rows, err := s.store.AdminListUsers(ctx, params)
	if err != nil {
		return nil, err
	}
	rs := []model.AdminListUsersResponse{}
	for _, row := range rows {
		rs = append(rs, model.AdminListUsersResponse{
			ID:         row.ID,
			Names:      row.Names,
			ProfilePic: row.ProfilePic,
			CreatedAt:  row.CreatedAt,
			Email:      row.Email,
			States:     row.States,
			RoleName:   row.RoleName,
		})
	}
	return rs, nil
}

func (s *adminUserService) CountUsers(ctx context.Context, query string) (int64, error) {
	return s.store.AdminCountUsers(ctx, query)
}
