package impl

import (
	"context"
	"time"

	"Taskly.com/m/internal/database"
	model "Taskly.com/m/internal/models"
	"Taskly.com/m/package/utils"

	"github.com/google/uuid"
)

type sChatService struct {
	store database.Store
}

func NewChatService(store database.Store) *sChatService {
	return &sChatService{store: store}
}

// 1. Tạo phòng chat (CreateRoom)
func (s *sChatService) CreateRoom(ctx context.Context, user1ID, user2ID uuid.UUID) (model.Room, error) {
	room, err := s.store.CreateRoom(ctx, database.CreateRoomParams{
		User1ID: user1ID,
		User2ID: user2ID,
	})
	if err != nil {
		return model.Room{}, err
	}
	return model.Room{
		ID:        room.ID,
		User1ID:   room.User1ID,
		User2ID:   room.User2ID,
		CreatedAt: room.CreatedAt,
	}, nil
}

func (s *sChatService) GetRoomInfo(ctx context.Context, roomID uuid.UUID) (model.RoomInfo, error) {
	roomInfo, err := s.store.GetRoomInfo(ctx, roomID)
	if err != nil {
		return model.RoomInfo{}, err
	}

	return model.RoomInfo{
		RoomID: roomInfo.RoomID,
		User1: model.UserInfo{
			ID:         roomInfo.User1ID,
			Name:       roomInfo.User1Name,
			ProfilePic: utils.PtrIfValid(roomInfo.User1ProfilePic),
		},
		User2: model.UserInfo{
			ID:         roomInfo.User2ID,
			Name:       roomInfo.User2Name,
			ProfilePic: utils.PtrIfValid(roomInfo.User2ProfilePic),
		},
	}, nil
}

// 2. Gửi tin nhắn (SetChatHistory)
func (s *sChatService) SetChatHistory(ctx context.Context, input model.SetChatInput) (model.Message, error) {
	msg, err := s.store.SetChatHistory(ctx, database.SetChatHistoryParams{
		RoomID:     input.RoomID,
		SenderID:   input.SenderID,
		ReceiverID: input.ReceiverID,
		Content:    input.Content,
	})
	if err != nil {
		return model.Message{}, err
	}
	return model.Message{
		ID:         msg.ID,
		RoomID:     msg.RoomID,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		SentAt:     msg.SentAt,
	}, nil
}

// 3. Lấy lịch sử chat theo room (GetChatHistory)
func (s *sChatService) GetChatHistory(ctx context.Context, roomID uuid.UUID, sentAt time.Time) ([]model.Message, error) {
	msgs, err := s.store.GetChatHistory(ctx, database.GetChatHistoryParams{
		RoomID: roomID,
		SentAt: sentAt,
	})
	if err != nil {
		return nil, err
	}
	var result []model.Message
	for _, m := range msgs {
		result = append(result, model.Message{
			ID:         m.ID,
			RoomID:     m.RoomID,
			SenderID:   m.SenderID,
			ReceiverID: m.ReceiverID,
			Content:    m.Content,
			SentAt:     m.SentAt,
		})
	}
	return result, nil
}

// 4. Lấy danh sách phòng theo user (GetRoomChatByUserId)
func (s *sChatService) GetRoomChatByUserId(ctx context.Context, userID uuid.UUID) ([]model.RoomWithLastMessage, error) {
	rooms, err := s.store.GetRoomChatByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []model.RoomWithLastMessage
	for _, r := range rooms {
		result = append(result, model.RoomWithLastMessage{
			ID:              r.ID,
			User1ID:         r.User1ID,
			User2ID:         r.User2ID,
			CreatedAt:       r.CreatedAt,
			LastMessage:     r.LastMessage,
			LastMessageTime: r.LastMessageTime,
		})
	}
	return result, nil
}
