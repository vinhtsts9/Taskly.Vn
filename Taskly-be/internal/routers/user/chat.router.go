package user

import (
	chat_controller "Taskly.com/m/internal/controller"
	middleware "Taskly.com/m/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type ChatRouter struct{}

func (r *ChatRouter) InitChatRouter(Router *gin.RouterGroup) {
	chatController := chat_controller.NewChatController()

	private := Router.Group("/chat")
	private.Use(middleware.AuthenMiddleware())
	{
		private.GET("/rooms", chatController.GetRoomChatByUserId)
		private.GET("/rooms/:room_id/info", chatController.GetRoomInfo) // Endpoint mới
		private.POST("/create-room", chatController.CreateRoom)
		private.POST("/send", chatController.SetChatHistory)
		private.GET("/history/:room_id", chatController.GetChatHistory)
		private.POST("/room-exists", chatController.RoomExists)
	}
}
