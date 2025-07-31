package controller

import (
	"net/http"
	"strconv"

	"Taskly.com/m/global"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type GigController struct {
	svc service.IGigService
}

func NewGigController() *GigController {
	return &GigController{
		svc: service.GetGigService(),
	}
}

// 2. Tạo dịch vụ mới
func (ctl *GigController) CreateService(c *gin.Context) {
	var input model.CreateServiceParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	gig, err := ctl.svc.CreateService(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusOK, gig)
}

// 3. Lấy chi tiết dịch vụ theo ID
func (ctl *GigController) GetServiceByID(c *gin.Context) {
	idStr := c.Param("id")
	global.Logger.Info("Attempting to get service with ID string", zap.String("id", idStr))

	id, err := uuid.Parse(idStr)
	if err != nil {
		global.Logger.Error("Failed to parse service ID", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}
	global.Logger.Info("Successfully parsed service ID", zap.String("parsed_uuid", id.String()))

	gig, err := ctl.svc.GetServiceByID(c, id)
	if err != nil {
		global.Logger.Error("Service layer returned an error for GetServiceByID", zap.String("id", id.String()), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	global.Logger.Info("Successfully retrieved service", zap.Any("gig", gig))
	c.JSON(http.StatusOK, gig)
}

// 4. Cập nhật dịch vụ
func (ctl *GigController) UpdateService(c *gin.Context) {
	var input model.UpdateServiceParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	gig, err := ctl.svc.UpdateService(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, gig)
}

// 5. Xoá dịch vụ
func (ctl *GigController) DeleteService(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID"})
		return
	}

	if err := ctl.svc.DeleteService(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

// Helper parseInt (an toàn)
func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
