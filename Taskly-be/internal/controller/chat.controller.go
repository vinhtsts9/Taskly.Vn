package controller

import (
	"fmt"
	"net/http"
	"time"

	"Taskly.com/m/global"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/internal/service"
	"Taskly.com/m/package/utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatController struct {
	svc service.IChatService // Giả sử bạn đã có interface IChatService rồi
}

func NewChatController() *ChatController {
	return &ChatController{
		svc: service.GetChatService(), // Tạo từ singleton hoặc DI tùy bạn
	}
}

// 1. Tạo phòng chat
func (ctl *ChatController) CreateRoom(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
		User2ID string `json:"user2_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: user2_id is required"})
		return
	}

	userInfo := auth.GetUserFromContext(c)
	if userInfo.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user credentials"})
		return
	}
	content := req.Content
	user1ID := userInfo.ID
	user2ID, err := uuid.Parse(req.User2ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User2ID format"})
		return
	}

	// Ngăn người dùng tạo phòng với chính mình
	if user1ID == user2ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create a room with yourself"})
		return
	}

	room, err := ctl.svc.CreateRoom(c, user1ID, user2ID, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create room"})
		fmt.Println("error create room: ", err)
		return
	}

	c.JSON(http.StatusOK, room)
}

func (ctl *ChatController) GetRoomInfo(c *gin.Context) {
	roomIDStr := c.Param("room_id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	roomInfo, err := ctl.svc.GetRoomInfo(c, roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, roomInfo)
}

// 2. Gửi tin nhắn
func (ctl *ChatController) SetChatHistory(c *gin.Context) {
	var input model.SetChatInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	msg, err := ctl.svc.SetChatHistory(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, msg)
}

// 3. Lấy lịch sử chat theo room
func (ctl *ChatController) GetChatHistory(c *gin.Context) {
	roomIDStr := c.Param("room_id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Lấy cursor (sent_at của tin nhắn cuối cùng) từ query, nếu không có thì dùng thời gian hiện tại
	cursorStr := c.DefaultQuery("cursor", "")
	var cursor time.Time
	if cursorStr == "" {
		cursor = time.Now()
	} else {
		cursor, err = time.Parse(time.RFC3339, cursorStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor format"})
			return
		}
	}

	// Gọi service với các tham số mới
	messages, err := ctl.svc.GetChatHistory(c, roomID, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// 4. Lấy danh sách phòng theo user
func (ctl *ChatController) GetRoomChatByUserId(c *gin.Context) {
	userInfo := auth.GetUserFromContext(c)
	if userInfo.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user credentials"})
		return
	}

	rooms, err := ctl.svc.GetRoomChatByUserId(c, userInfo.ID)
	if err != nil {
		global.Logger.Sugar().Error("Failed to get rooms", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get rooms"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

// 5. Kiểm tra phòng đã hợp lệ chưa
func (ctl *ChatController) RoomExists(c *gin.Context) {
	var req struct {
		User2ID string `json:"user2_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: user2_id is required"})
		return
	}

	userInfo := auth.GetUserFromContext(c)
	if userInfo.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user credentials"})
		return
	}

	user1ID := userInfo.ID
	user2ID, err := uuid.Parse(req.User2ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User2ID format"})
		return
	}

	exists, err := ctl.svc.CheckRoomExist(c, user1ID, user2ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check room exists"})
		return
	}

	c.JSON(http.StatusOK, exists)
}
