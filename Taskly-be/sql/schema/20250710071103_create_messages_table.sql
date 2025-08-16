-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Tạo bảng rooms

CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    user1_id UUID NOT NULL,
    user2_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
--CREATE UNIQUE INDEX ux_rooms_user_pair
--ON rooms (LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id));

-- Thêm index cho truy vấn theo user
CREATE INDEX idx_rooms_user1_id ON rooms(user1_id);
CREATE INDEX idx_rooms_user2_id ON rooms(user2_id);

-- 2. Tạo bảng messages, gắn với room
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    room_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    receiver_id UUID NOT NULL,
    content TEXT NOT NULL,
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Thêm index cho truy vấn tin nhắn
CREATE INDEX idx_messages_room_id ON messages(room_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_room_id_sent_at ON messages(room_id, sent_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Xóa index trước khi xóa bảng
DROP INDEX IF EXISTS idx_messages_room_id_sent_at;
DROP INDEX IF EXISTS idx_messages_sender_id;
DROP INDEX IF EXISTS idx_messages_room_id;
DROP INDEX IF EXISTS idx_rooms_user2_id;
DROP INDEX IF EXISTS idx_rooms_user1_id;

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS rooms;
-- +goose StatementEnd
