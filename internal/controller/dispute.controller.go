package controller

import (
	"net/http"

	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
	"github.com/gin-gonic/gin"
)

type DisputeController struct {
	svc service.IDisputeService // Interface, cần được define
}

func NewDisputeController() *DisputeController {
	return &DisputeController{
		svc: service.GetDisputeService(), // singleton hoặc inject
	}
}

// 1. Tạo tranh chấp
func (ctl *DisputeController) CreateDispute(c *gin.Context) {
	var req model.CreateDisputeParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	dispute, err := ctl.svc.CreateDispute(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dispute)
}

// 2. Lấy danh sách tranh chấp
func (ctl *DisputeController) ListDisputes(c *gin.Context) {
	disputes, err := ctl.svc.ListDisputes(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch disputes"})
		return
	}

	c.JSON(http.StatusOK, disputes)
}

// 3. Cập nhật trạng thái tranh chấp
func (ctl *DisputeController) UpdateDisputeStatus(c *gin.Context) {
	var req model.UpdateDisputeStatusParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ctl.svc.UpdateDisputeStatus(c, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dispute status updated successfully"})
}
