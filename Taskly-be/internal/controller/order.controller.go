package controller

import (
	"net/http"

	"Taskly.com/m/global"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
	"Taskly.com/m/package/utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrderController struct {
	svc          service.IOrderService
	vnpayService service.IVNPayService
}

func NewOrderController() *OrderController {
	return &OrderController{
		svc:          service.GetOrderService(),
		vnpayService: service.GetVNPayService(),
	}
}

// 1. Tạo đơn hàng mới
func (ctl *OrderController) CreateOrder(c *gin.Context) {
	var input model.CreateOrderParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		global.Logger.Error("Invalid input for creating order", zap.Error(err))
		return
	}

	order,err := ctl.svc.CreateOrder(c, input, auth.GetUserFromContext(c).ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		global.Logger.Error("Failed to create order", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order": order})
}

// 1b. Tạo đơn hàng và lấy link thanh toán VNPay
// func (ctl *OrderController) CreateOrderAndGetVNPayURL(c *gin.Context) {
// 	var input model.CreateOrderParams
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	url, err := ctl.svc.CreateOrderAndGenerateVNPayURL(c, input, ctl.vnpayService)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order or generate VNPay URL"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"payment_url": url})
// }

// // 1c. Xử lý callback từ VNPay
// func (ctl *OrderController) HandleVNPayCallback(c *gin.Context) {
// 	if err := ctl.svc.HandleVNPayCallback(c, c.Request.URL.Query(), ctl.vnpayService); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Payment verified and order updated"})
// }

// 2. Lấy danh sách đơn theo user
func (ctl *OrderController) ListOrdersByUser(c *gin.Context) {
	// Lấy userID từ context của người dùng đã xác thực, thay vì từ URL param
	userInfo := auth.GetUserFromContext(c)
	if userInfo.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user credentials"})
		return
	}

	orders, err := ctl.svc.ListOrdersByUser(c, userInfo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// 3. Lấy chi tiết đơn
func (ctl *OrderController) GetOrderByID(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := ctl.svc.GetOrderByID(c, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// 4. Cập nhật trạng thái đơn (bất kỳ)
func (ctl *OrderController) UpdateOrderStatus(c *gin.Context) {
	var input model.UpdateOrderStatusParams
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := ctl.svc.UpdateOrderStatus(c, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

// 5. Freelancer giao hàng
func (ctl *OrderController) SubmitOrderDelivery(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := ctl.svc.SubmitOrderDelivery(c, orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit delivery"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order delivery submitted"})
}

// 6. Buyer xác nhận hoàn thành
func (ctl *OrderController) AcceptOrderCompletion(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := ctl.svc.AcceptOrderCompletion(c, orderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order marked as completed"})
}
