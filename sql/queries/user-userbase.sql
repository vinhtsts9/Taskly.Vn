-- name: CreateUserBase :one
INSERT INTO user_base (passwords,email,salt)
VALUES ($1,$2,$3)
RETURNING user_base_id;

-- name: GetUserBaseToCheckLogin :one
SELECT user_base_id, salt, passwords
FROM user_base
WHERE email = $1;

-- name: CheckUserExist :one
SELECT EXISTS(SELECT 1 FROM user_base WHERE email = $1);

-- name: CreateUserProfile :exec
INSERT INTO users (user_base_id, names, user_type, profile_pic, bio)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateLoginInfo :exec
UPDATE user_base
SET login_time = CURRENT_TIMESTAMP,
    login_ip = $2,
    states = 1
WHERE email = $1;

-- name: UpdateLogoutInfo :exec
UPDATE user_base
SET logout_time = CURRENT_TIMESTAMP,
    states = 3
WHERE user_base_id = $1;

-- name: GetUserInfoToSetToken :one
SELECT 
    u.id, u.names, u.user_type, u.profile_pic, u.bio,
    u.created_at, u.updated_at
FROM users u
JOIN user_base b ON u.user_base_id = b.user_base_id
WHERE b.user_base_id = $1;

-- name: DeleteUserBase :exec
DELETE FROM user_base WHERE user_base_id = $1;

-- name: ListUsersByType :many
SELECT * FROM users WHERE user_type = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserBaseToken :exec
UPDATE user_base
SET refresh_token = $2
WHERE user_base_id = $1;


-- name: CheckRefreshToken :one
SELECT user_base_id FROM user_base WHERE refresh_token = $1;