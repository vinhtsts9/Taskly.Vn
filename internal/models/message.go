package model

import (
	"time"

	"github.com/google/uuid"
)

// ==== Room ====

type Room struct {
	ID        uuid.UUID `json:"id"`
	User1ID   uuid.UUID `json:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}

// RoomWithLastMessage bao gồm thông tin room và tin nhắn cuối cùng
type RoomWithLastMessage struct {
	ID              uuid.UUID `json:"id"`
	User1ID         uuid.UUID `json:"user1_id"`
	User2ID         uuid.UUID `json:"user2_id"`
	CreatedAt       time.Time `json:"created_at"`
	LastMessage     string    `json:"last_message"`
	LastMessageTime time.Time `json:"last_message_time"`
}

type CreateRoomParams struct {
	User1ID uuid.UUID `json:"user1_id" binding:"required"`
	User2ID uuid.UUID `json:"user2_id" binding:"required"`
}

// UserInfo định nghĩa thông tin cơ bản của người dùng trong phòng chat
type UserInfo struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	ProfilePic *string   `json:"profile_pic,omitempty"`
}

// RoomInfo định nghĩa thông tin chi tiết của phòng chat, bao gồm cả hai người dùng
type RoomInfo struct {
	RoomID uuid.UUID `json:"room_id"`
	User1  UserInfo  `json:"user1"`
	User2  UserInfo  `json:"user2"`
}

type RoomResult struct {
	ID        uuid.UUID `json:"id"`
	User1ID   uuid.UUID `json:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ==== Message ====

type Message struct {
	ID         uuid.UUID `json:"id"`
	RoomID     uuid.UUID `json:"room_id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
}

type SetChatInput struct {
	RoomID     uuid.UUID `json:"room_id" binding:"required"`
	SenderID   uuid.UUID `json:"sender_id" binding:"required"`
	ReceiverID uuid.UUID `json:"receiver_id" binding:"required"`
	Content    string    `json:"content" binding:"required"`
}
type MessageResult struct {
	ID         uuid.UUID `json:"id"`
	RoomID     uuid.UUID `json:"room_id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
}
