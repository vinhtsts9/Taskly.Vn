package controller

import (
	"net/http"

	"Taskly.com/m/internal/database"
	"Taskly.com/m/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RBACController struct {
	svc service.IRBACService
}

func NewRBACController() *RBACController {
	return &RBACController{
		svc: service.GetRbacService(),
	}
}

// CreateRole tạo một vai trò mới
func (ctl *RBACController) CreateRole(c *gin.Context) {
	type reqBody struct {
		Name string `json:"name" binding:"required"`
	}
	var req reqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dữ liệu không hợp lệ"})
		return
	}

	role, err := ctl.svc.CreateRole(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi tạo vai trò"})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// CreatePermission tạo một quyền mới
func (ctl *RBACController) CreatePermission(c *gin.Context) {
	type reqBody struct {
		Name     string `json:"name" binding:"required"`
		Resource string `json:"resource" binding:"required"`
		Action   string `json:"action" binding:"required"`
	}
	var req reqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dữ liệu không hợp lệ"})
		return
	}

	arg := database.CreatePermissionParams{
		Name:     req.Name,
		Resource: req.Resource,
		Action:   req.Action,
	}

	permission, err := ctl.svc.CreatePermission(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi tạo quyền"})
		return
	}

	c.JSON(http.StatusCreated, permission)
}

// AddRoleToUser gán một vai trò cho người dùng
func (ctl *RBACController) AddRoleToUser(c *gin.Context) {
	type reqBody struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
		RoleID uuid.UUID `json:"role_id" binding:"required"`
	}
	var req reqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dữ liệu không hợp lệ"})
		return
	}

	arg := database.AddRoleToUserParams{
		UserID: req.UserID,
		RoleID: req.RoleID,
	}

	if err := ctl.svc.AddRoleToUser(c, arg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi gán vai trò cho người dùng"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "gán vai trò cho người dùng thành công"})
}

// AddPermissionToRole gán một quyền cho vai trò
func (ctl *RBACController) AddPermissionToRole(c *gin.Context) {
	type reqBody struct {
		RoleID       uuid.UUID `json:"role_id" binding:"required"`
		PermissionID uuid.UUID `json:"permission_id" binding:"required"`
	}
	var req reqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dữ liệu không hợp lệ"})
		return
	}

	arg := database.AddPermissionToRoleParams{
		RoleID:       req.RoleID,
		PermissionID: req.PermissionID,
	}

	if err := ctl.svc.AddPermissionToRole(c, arg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi gán quyền cho vai trò"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "gán quyền cho vai trò thành công"})
}

// RemovePermissionFromRole hủy gán quyền khỏi vai trò
func (ctl *RBACController) RemovePermissionFromRole(c *gin.Context) {
	type reqBody struct {
		RoleID       uuid.UUID `json:"role_id" binding:"required"`
		PermissionID uuid.UUID `json:"permission_id" binding:"required"`
	}
	var req reqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dữ liệu không hợp lệ"})
		return
	}
	arg := database.RemovePermissionFromRoleParams{
		RoleID:       req.RoleID,
		PermissionID: req.PermissionID,
	}
	if err := ctl.svc.RemovePermissionFromRole(c, arg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi hủy gán quyền khỏi vai trò"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "đã hủy gán quyền khỏi vai trò"})
}

// RemoveRoleFromUser hủy gán vai trò khỏi người dùng
func (ctl *RBACController) RemoveRoleFromUser(c *gin.Context) {
	type reqBody struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
		RoleID uuid.UUID `json:"role_id" binding:"required"`
	}
	var req reqBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dữ liệu không hợp lệ"})
		return
	}
	arg := database.RemoveRoleFromUserParams{
		UserID: req.UserID,
		RoleID: req.RoleID,
	}
	if err := ctl.svc.RemoveRoleFromUser(c, arg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi hủy gán vai trò khỏi người dùng"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "đã hủy gán vai trò khỏi người dùng"})
}

// GetPermissionsByRoleID lấy quyền theo role ID
func (ctl *RBACController) GetPermissionsByRoleID(c *gin.Context) {
	roleIDParam := c.Param("role_id")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_id không hợp lệ"})
		return
	}
	perms, err := ctl.svc.GetPermissionsByRoleID(c, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi lấy quyền theo role"})
		return
	}
	c.JSON(http.StatusOK, perms)
}

// GetPermissionsByUserID lấy quyền theo user ID
func (ctl *RBACController) GetPermissionsByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id không hợp lệ"})
		return
	}
	perms, err := ctl.svc.GetPermissionsByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi lấy quyền theo user"})
		return
	}
	c.JSON(http.StatusOK, perms)
}

// GetRolesByUserID lấy role theo user ID
func (ctl *RBACController) GetRolesByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id không hợp lệ"})
		return
	}
	roles, err := ctl.svc.GetRolesByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi lấy vai trò của người dùng"})
		return
	}
	c.JSON(http.StatusOK, roles)
}
