package service

import (
	"context"

	model "Taskly.com/m/internal/models"
	"github.com/google/uuid"
)

type IUserService interface {
	Register(ctx context.Context, verifyKey string, verifyType string) error

	VerifyOTP(ctx context.Context, VerifyKey string, OTP string) (err error)

	UpdatePasswordRegister(ctx context.Context, in model.UpdatePasswordRegisterParams) (err error)

	Login(ctx context.Context, in model.LoginInPut) (out *model.LoginOutput, err error)
	RefreshToken(ctx context.Context, refreshToken string) (out *model.LoginOutput, err error)
	// 4. Cập nhật logout
	UpdateLogoutInfo(ctx context.Context, userBaseID uuid.UUID) error

	// 5. Lấy thông tin người dùng đầy đủ (user + user_base)
	GetUserByID(ctx context.Context, userID uuid.UUID) (model.User, error)

	// 6. Xóa user base (xóa tài khoản)
	DeleteUserBase(ctx context.Context, userBaseID uuid.UUID) error

	// 7. Lấy danh sách user theo loại
	ListUsersByType(ctx context.Context, userType string) ([]model.User, error)
}

var (
	localUserService IUserService
)

func GetUserService() IUserService {

	if localUserService == nil {
		panic("implement localUserService notfound")
	}

	return localUserService
}

func InitUserService(i IUserService) {
	localUserService = i
}
