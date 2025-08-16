package impl

import (
	"context"
	"database/sql"
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
func (s *sChatService) CreateRoom(ctx context.Context, user1ID, user2ID uuid.UUID, content string) (model.Room, error) {
	var room model.Room
	err := s.store.ExecTx(ctx, func(q *database.Queries) error {
		var err error
		roomRs, err := s.store.CreateRoom(ctx, database.CreateRoomParams{
			Column1: user1ID,
			Column2: user2ID,
		})
		if err != nil {
			return err
		}
		room = model.Room{
			ID:        roomRs.ID,
			User1ID:   roomRs.User1ID,
			User2ID:   roomRs.User2ID,
			CreatedAt: roomRs.CreatedAt,
		}

		_, err = s.store.SetChatHistory(ctx, database.SetChatHistoryParams{
			RoomID:     room.ID,
			SenderID:   user1ID,
			ReceiverID: user2ID,
			Content:    content,
		})
		if err != nil {
			return err
		}

		return nil
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
			ProfilePic: utils.PtrStringIfValid(roomInfo.User1ProfilePic),
		},
		User2: model.UserInfo{
			ID:         roomInfo.User2ID,
			Name:       roomInfo.User2Name,
			ProfilePic: utils.PtrStringIfValid(roomInfo.User2ProfilePic),
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
			User1Name:       r.User1Name,
			User1ProfilePic: utils.PtrStringIfValid(r.User1ProfilePic),
			User2ID:         r.User2ID,
			User2Name:       r.User2Name,
			User2ProfilePic: utils.PtrStringIfValid(r.User2ProfilePic),
			CreatedAt:       r.CreatedAt,
			LastMessage:     r.LastMessage,
			LastMessageTime: r.LastMessageTime,
		})
	}
	return result, nil
}

// 5. Kiểm tra phòng tồn tại
func (s *sChatService) CheckRoomExist(ctx context.Context, user1ID, user2ID uuid.UUID) (model.Room, error) {
	exists, err := s.store.CheckRoomExist(ctx, database.CheckRoomExistParams{
		Column1: user1ID,
		Column2: user2ID,
	})
	if err == sql.ErrNoRows {
		return model.Room{}, nil
	} else if err != nil {
		return model.Room{}, err
	} else {
		return model.Room{
			ID:        exists.ID,
			User1ID:   exists.User1ID,
			User2ID:   exists.User2ID,
			CreatedAt: exists.CreatedAt,
		}, nil
	}
}
