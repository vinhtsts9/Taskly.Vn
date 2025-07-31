package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"Taskly.com/m/global"
	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
	utils "Taskly.com/m/package/utils"
	"Taskly.com/m/package/utils/auth"
	"Taskly.com/m/package/utils/crypto"
	"Taskly.com/m/package/utils/random"
	"Taskly.com/m/package/utils/sendto"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const VerifyTypeEmail = "email"

type sUserService struct {
	store database.Store // Store là interface, có thể là SQLStore
}

func NewUserService(store database.Store) *sUserService {
	return &sUserService{
		store: store,
	}
}
func (s *sUserService) Register(ctx context.Context, verifyKey string, verifyType string) error {
	// 1. Băm key
	hashKey := crypto.GetHash(strings.ToLower(verifyKey))

	// 2. Kiểm tra người dùng đã tồn tại
	exists, err := s.store.CheckUserExist(ctx, verifyKey)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("user with %s already exists", verifyKey)
	}
	userKey := utils.GetUserKey(hashKey)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()

	switch {
	case err == redis.Nil:
		// Chưa có OTP nào → tiếp tục
		fmt.Println("Key not found in Redis, continue to create OTP")
	case err != nil:
		// Redis bị lỗi thật sự
		fmt.Println("Redis GET failed:", err)
		return err
	case otpFound != "":
		// OTP vẫn còn hạn sử dụng
		return fmt.Errorf("OTP has already been sent and is still valid")
	}

	// 4. Tạo OTP
	otp := random.GenerateSixDigitOtp()
	err = global.Rdb.SetEx(ctx, userKey, strconv.Itoa(otp), time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set OTP in redis: %w", err)
	}
	fmt.Println("OTP:", otp)

	// 5. Gửi OTP
	switch verifyType {
	case VerifyTypeEmail:
		err := sendto.SendTextEmail([]string{verifyKey}, "Vinhtiensinh17@gmail.com", strconv.Itoa(otp))
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	default:
		return fmt.Errorf("unsupported verify type: %s", verifyType)
	}

	// 6. Ghi vào bảng user_verify
	existsOTPVerify, err := s.store.CheckOTPVerifyExist(ctx, hashKey)
	if err != nil {
		return fmt.Errorf("check OTP verify failed: %w", err)
	}
	if existsOTPVerify {
		return fmt.Errorf("OTP verify entry already exists")
	}

	_, err = s.store.InsertOTPVerify(ctx, database.InsertOTPVerifyParams{
		VerifyOtp:     strconv.Itoa(otp),
		VerifyType:    database.VerifyTypeEnum(verifyType),
		VerifyKey:     verifyKey,
		VerifyHashKey: hashKey,
	})
	if err != nil {
		return fmt.Errorf("failed to insert OTP record: %w", err)
	}

	return nil
}

func (s *sUserService) VerifyOTP(ctx context.Context, VerifyKey string, OTP string) (err error) {
	hashKey := crypto.GetHash(strings.ToLower(VerifyKey))
	userKey := utils.GetUserKey(hashKey)

	otpFound, err := global.Rdb.Get(ctx, userKey).Result()
	if err != nil {
		return err
	}
	fmt.Println("otpFound:", otpFound)
	fmt.Println("OTP:", OTP)
	if OTP != otpFound {
		return fmt.Errorf("OTP not match")
	}

	// Xóa OTP khỏi Redis để tránh bị dùng lại
	if delErr := global.Rdb.Del(ctx, userKey).Err(); delErr != nil {
		// Không return luôn, chỉ log lỗi
		fmt.Printf("warning: failed to delete OTP from Redis: %v\n", delErr)
	}

	err = s.store.UpdateUserVerificationStatus(ctx, hashKey)
	if err != nil {
		return err
	}
	return nil
}

func (s *sUserService) UpdatePasswordRegister(ctx context.Context, in model.UpdatePasswordRegisterParams) (err error) {
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))
	infoOTP, err := s.store.GetInfoOTP(ctx, hashKey)

	if err != nil {
		return err
	}
	// check isVerified ok
	if infoOTP.IsVerified == false {
		return fmt.Errorf("user OTP not verified")
	}
	// update user_base table
	userBase := database.CreateUserBaseParams{}
	userBase.Email = infoOTP.VerifyKey
	if err != nil {
		return err
	}

	userSalt, err := crypto.GenerateSalt(16)
	if err != nil {
		return err
	}
	userBase.Passwords = crypto.HassPassword(in.UserPassword, userSalt)
	userBase.Salt = userSalt
	// add userBase to userBase table

	newUserBase, err := s.store.CreateUserBase(ctx, userBase)
	if err != nil {
		return err
	}
	err = s.store.CreateUserProfile(ctx, database.CreateUserProfileParams{
		UserBaseID: newUserBase,
		Names:      in.Names,
		UserType:   in.UserType,
		ProfilePic: utils.ToNullString(in.ProfilePic),
		Bio:        utils.ToNullString(in.Bio),
	})
	if err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}
	return nil
}

// Đăng nhập: cập nhật login info
func (s *sUserService) Login(ctx context.Context, req model.LoginInPut) (*model.LoginOutput, error) {
	userBase, err := s.store.GetUserBaseToCheckLogin(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user base: %w", err)
	}

	if !crypto.MatchingPassword(userBase.Passwords, req.Password, userBase.Salt) {
		return nil, errors.New("invalid email or password")
	}

	_ = s.store.UpdateLoginInfo(ctx, database.UpdateLoginInfoParams{
		Email:   req.Email,
		LoginIp: req.LoginIP,
	})

	userInfo, err := s.store.GetUserInfoToSetToken(ctx, userBase.UserBaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	userToken := model.UserToken{
		ID:       userInfo.ID,
		UserType: userInfo.UserType,
	}
	// Tạo token
	accessToken, refreshToken, err := auth.GenerateTokens(userToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}
	//Lưu access_token và refresh_token vào user_base table
	err = s.store.UpdateUserBaseToken(ctx, database.UpdateUserBaseTokenParams{
		UserBaseID:   userBase.UserBaseID,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user base: %w", err)
	}
	// Lưu thông tin người dùng vào Redis
	userInfoJson, err := json.Marshal(userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user info for Redis: %w", err)
	}

	// Lưu vào Redis với thời gian sống
	err = global.Rdb.Set(ctx, utils.GetUserKey(userInfo.ID.String()), userInfoJson, time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store user info in Redis: %w", err)
	}
	return &model.LoginOutput{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *sUserService) RefreshToken(ctx context.Context, refreshToken string) (*model.LoginOutput, error) {
	// 1. Xác thực refresh token và lấy user ID
	user_base_id, err := s.store.CheckRefreshToken(ctx, refreshToken)
	if user_base_id == uuid.Nil || err != nil {
		return nil, fmt.Errorf("failed to check refresh token")
	}

	// 2. Lấy thông tin user từ DB
	userInfo, err := s.store.GetUserByID(ctx, user_base_id)
	if err != nil {
		return nil, fmt.Errorf("user not found for refresh token: %w", err)
	}

	userToken := model.UserToken{
		ID:       userInfo.ID,
		UserType: userInfo.UserType,
	}

	// 3. Xoay vòng token: Tạo cặp token mới và thu hồi token cũ
	newAccessToken, newRefreshToken, err := auth.GenerateTokens(userToken)
	if err != nil {
		return nil, fmt.Errorf("failed to rotate refresh token: %w", err)
	}
	err = s.store.UpdateUserBaseToken(ctx, database.UpdateUserBaseTokenParams{
		UserBaseID:   user_base_id,
		RefreshToken: newRefreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user base: %w", err)
	}

	return &model.LoginOutput{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Đăng xuất
func (s *sUserService) UpdateLogoutInfo(ctx context.Context, userBaseID uuid.UUID) error {
	return s.store.UpdateLogoutInfo(ctx, userBaseID)
}

// Lấy user theo ID
func (s *sUserService) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	dbUser, err := s.store.GetUserByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:         dbUser.ID,
		UserBaseID: dbUser.UserBaseID,
		Names:      dbUser.Names,
		UserType:   dbUser.UserType,
		ProfilePic: utils.PtrIfValid(dbUser.ProfilePic),
		Bio:        utils.PtrIfValid(dbUser.Bio),
		CreatedAt:  dbUser.CreatedAt,
		UpdatedAt:  dbUser.UpdatedAt,
	}, nil
}

// Xóa user base (ví dụ rollback nếu step 2 fail)
func (s *sUserService) DeleteUserBase(ctx context.Context, userBaseID uuid.UUID) error {
	return s.store.DeleteUserBase(ctx, userBaseID)
}

// Lấy danh sách user theo type
func (s *sUserService) ListUsersByType(ctx context.Context, userType string) ([]model.User, error) {
	dbUsers, err := s.store.ListUsersByType(ctx, userType)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, len(dbUsers))
	for i, u := range dbUsers {
		users[i] = model.User{
			ID:         u.ID,
			UserBaseID: u.UserBaseID,
			Names:      u.Names,
			UserType:   u.UserType,
			ProfilePic: utils.PtrIfValid(u.ProfilePic),
			Bio:        utils.PtrIfValid(u.Bio),
			CreatedAt:  u.CreatedAt,
			UpdatedAt:  u.UpdatedAt,
		}
	}
	return users, nil
}
