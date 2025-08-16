package service

import (
	"context"
	"time"

	model "Taskly.com/m/internal/models"
	"github.com/google/uuid"
)

type IChatService interface {

	// 1. Tạo phòng chat giữa 2 user
	CreateRoom(ctx context.Context, user1ID, user2ID uuid.UUID, content string) (model.Room, error)
	GetRoomInfo(ctx context.Context, roomID uuid.UUID) (model.RoomInfo, error)

	// 2. Gửi tin nhắn vào phòng chat
	SetChatHistory(ctx context.Context, input model.SetChatInput) (model.Message, error)

	// 3. Lấy lịch sử chat trong phòng (phân trang)
	GetChatHistory(ctx context.Context, roomID uuid.UUID, sentAt time.Time) ([]model.Message, error)

	// 4. Lấy tất cả phòng chat của user
	GetRoomChatByUserId(ctx context.Context, userID uuid.UUID) ([]model.RoomWithLastMessage, error)

	// 5.Kiểm tra phòng tồn tại
	CheckRoomExist(ctx context.Context, user1ID, user2ID uuid.UUID) (model.Room, error)
}

var (
	localChatService IChatService
)

func GetChatService() IChatService {
	if localChatService == nil {
		panic("implement localChatService not found")
	}
	return localChatService
}

func InitChatService(i IChatService) {
	localChatService = i
}
