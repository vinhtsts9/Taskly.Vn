package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"Taskly.com/m/package/utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatBotController struct{}

func NewChatBotController() *ChatBotController {
	return &ChatBotController{}
}

func (ctl *ChatBotController) ChatBotN8N(c *gin.Context) {
	var req struct {
		Question string `json:"question" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: question is required"})
		return
	}

	userInfo := auth.GetUserFromContext(c)
	if userInfo.ID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user credentials"})
		return
	}

	payload := map[string]string{
		"sessionId": userInfo.ID.String(),
		"question":  req.Question,
	}
	jsonPayload, _ := json.Marshal(payload)

	// Call n8n webhook
	resp, err := http.Post(
		"https://n8n-js.onrender.com/webhook/test_voice_message_elevenlabs",
		"application/json",
		bytes.NewReader(jsonPayload), // có thể dùng NewReader thay Buffer
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to n8n"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var n8nResp map[string]interface{}
	_ = json.Unmarshal(body, &n8nResp)

	c.JSON(http.StatusOK, gin.H{"answer": n8nResp["output"]})
}
