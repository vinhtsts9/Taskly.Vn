package user

import (
	chatbot_controller "Taskly.com/m/internal/controller"
	middleware "Taskly.com/m/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type ChatBotRouter struct{}

func (r *ChatBotRouter) InitChatBotRouter(Router *gin.RouterGroup) {
	chatController := chatbot_controller.NewChatBotController()

	private := Router.Group("/ai")
	private.Use(middleware.AuthenMiddleware())
	{
		private.POST("/chat", chatController.ChatBotN8N)
	}
}
