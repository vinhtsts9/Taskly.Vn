package controller

import (
	"fmt"
	"net/http"

	"Taskly.com/m/global" // Thêm import global
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
	"Taskly.com/m/package/utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	// Thêm import time
)

type UserController struct {
	svc service.IUserService
}

func NewUserController() *UserController {
	return &UserController{svc: service.GetUserService()}
}

// Đăng ký tài khoản (gửi OTP)
func (ctl *UserController) Register(c *gin.Context) {
	var req struct {
		VerifyKey  string `json:"verify_key" binding:"required"`
		VerifyType string `json:"verify_type" binding:"required"` // email / phone
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing verify_key or verify_type"})
		return
	}

	if err := ctl.svc.Register(c.Request.Context(), req.VerifyKey, req.VerifyType); err != nil {
		// Ghi log lỗi chi tiết ra terminal/file log
		global.Logger.Sugar().Errorf("Failed to register user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred while processing the registration."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// Xác minh OTP
func (ctl *UserController) VerifyOTP(c *gin.Context) {
	var req struct {
		VerifyKey string `json:"verify_key" binding:"required"`
		OTP       string `json:"otp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing verify_key or otp"})
		return
	}

	if err := ctl.svc.VerifyOTP(c.Request.Context(), req.VerifyKey, req.OTP); err != nil {
		global.Logger.Sugar().Errorf("Failed to verify OTP: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

// Tạo mật khẩu sau khi xác minh OTP
func (ctl *UserController) UpdatePasswordRegister(c *gin.Context) {
	var req model.UpdatePasswordRegisterParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ctl.svc.UpdatePasswordRegister(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete registration"})
		fmt.Println("Failed to complete registration", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Lấy thông tin của người dùng đang đăng nhập (dựa vào cookie)
func (ctl *UserController) Me(c *gin.Context) {
	userInfo := auth.GetUserFromContext(c)
	if userInfo.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or session"})
		return
	}

	// Trả về thông tin người dùng
	c.JSON(http.StatusOK, userInfo)
}

// Đăng nhập
func (ctl *UserController) Login(c *gin.Context) {
	var req model.LoginInPut
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login input"})
		return
	}

	out, err := ctl.svc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
		fmt.Println("Login failed", err.Error())
		return
	}

	// Thiết lập cookie cho Access Token (sống ngắn)
	accessMaxAge := 3600 // 1 giờ
	c.SetCookie("token", out.Token, accessMaxAge, "/", "", false, true)

	// Thiết lập cookie cho Refresh Token (sống dài)
	refreshMaxAge := 3600 * 24 * 7 // 7 ngày
	c.SetCookie("refresh_token", out.RefreshToken, refreshMaxAge, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Đăng xuất
func (ctl *UserController) Logout(c *gin.Context) {
	userInfo := auth.GetUserFromContext(c)
	if userInfo.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user credentials"})
		return
	}

	if err := ctl.svc.UpdateLogoutInfo(c.Request.Context(), userInfo.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	// Xóa cả hai cookie
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Lấy thông tin user theo ID
func (ctl *UserController) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := ctl.svc.GetUserByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Xóa tài khoản
func (ctl *UserController) DeleteUser(c *gin.Context) {
	userInfo := auth.GetUserFromContext(c)
	if userInfo.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user credentials"})
		return
	}

	if err := ctl.svc.DeleteUserBase(c.Request.Context(), userInfo.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (ctl *UserController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	newTokens, err := ctl.svc.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh token", "details": err.Error()})
		return
	}

	// Thiết lập lại cả 2 cookie với token mới
	accessMaxAge := 60 // 1 giờ
	c.SetCookie("token", newTokens.Token, accessMaxAge, "/", "", false, true)

	refreshMaxAge := 3600 * 24 * 7 // 7 ngày
	c.SetCookie("refresh_token", newTokens.RefreshToken, refreshMaxAge, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully"})
}
