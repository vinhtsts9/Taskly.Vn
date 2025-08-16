package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"Taskly.com/m/internal/database"
	"Taskly.com/m/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	svc service.IAdminUserService
}

func NewAdminUserController() *AdminUserController {
	return &AdminUserController{svc: service.GetAdminUserService()}
}

func (ctl *AdminUserController) ListUsers(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	// fields := c.DefaultQuery("fields", "basic") // hiện chưa dùng, để mở rộng

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	limit := int32(size)
	offset := int32((page - 1) * size)

	total, err := ctl.svc.CountUsers(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi đếm người dùng"})
		return
	}
	rows, err := ctl.svc.ListUsers(c, database.AdminListUsersParams{Column1: query, Limit: limit, Offset: offset})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "lỗi lấy danh sách người dùng"})
		return
	}
	fmt.Println("rows", rows)
	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page":  page,
		"size":  size,
		"data":  rows,
	})
}
