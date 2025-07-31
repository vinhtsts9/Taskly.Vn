package controller

import (
	"net/http"

	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
	"github.com/gin-gonic/gin"
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
	var req model.CreateRoleRequest
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
	var req model.CreatePermissionRequest
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
	var req model.AddRoleToUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dữ liệu không hợp lệ"})
		return
	}

	arg := database.AddRoleToUserParams{
		UserID: req.UserID,
		RoleID: req.RoleID,
	}

	err := ctl.svc.AddRoleToUser(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi gán vai trò cho người dùng"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "gán vai trò cho người dùng thành công"})
}

// AddPermissionToRole gán một quyền cho vai trò
func (ctl *RBACController) AddPermissionToRole(c *gin.Context) {
	var req model.AddPermissionToRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dữ liệu không hợp lệ"})
		return
	}

	arg := database.AddPermissionToRoleParams{
		RoleID:       req.RoleID,
		PermissionID: req.PermissionID,
	}

	err := ctl.svc.AddPermissionToRole(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi khi gán quyền cho vai trò"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "gán quyền cho vai trò thành công"})
}
