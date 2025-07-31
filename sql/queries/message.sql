-- name: CreateRoom :one
INSERT INTO rooms (user1_id, user2_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetRoomInfo :one
SELECT 
    r.id AS room_id,
    u1.id AS user1_id,
    u1.names AS user1_name,
    u1.profile_pic AS user1_profile_pic,
    u2.id AS user2_id,
    u2.names AS user2_name,
    u2.profile_pic AS user2_profile_pic
FROM rooms r
JOIN users u1 ON r.user1_id = u1.id
JOIN users u2 ON r.user2_id = u2.id
WHERE r.id = $1;

-- name: GetChatHistory :many
SELECT * FROM messages
WHERE room_id = $1 AND sent_at < $2
ORDER BY sent_at DESC
LIMIT 5;


-- name: SetChatHistory :one
INSERT INTO messages (room_id, sender_id, receiver_id, content)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetRoomChatByUserId :many
SELECT
    r.id,
    r.user1_id,
    r.user2_id,
    r.created_at,
    (SELECT content FROM messages WHERE room_id = r.id ORDER BY sent_at DESC LIMIT 1) AS last_message,
    (SELECT sent_at FROM messages WHERE room_id = r.id ORDER BY sent_at DESC LIMIT 1) AS last_message_time
FROM
    rooms r
WHERE
    r.user1_id = $1 OR r.user2_id = $1
ORDER BY
    last_message_time DESC; -- SẮP XẾP THEO THỜI GIAN CỦA TIN NHẮN CUỐI CÙNG
