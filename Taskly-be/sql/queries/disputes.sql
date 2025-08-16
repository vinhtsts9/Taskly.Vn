-- name: CreateDispute :one
INSERT INTO disputes (order_id, user_id, reason)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetDisputeByOrderID :one
SELECT * FROM disputes WHERE order_id = $1;

-- name: ListDisputes :many
SELECT * FROM disputes ORDER BY created_at DESC;

-- name: UpdateDisputeStatus :exec
UPDATE disputes SET status = $2 WHERE id = $1;

-- name: GetDisputeByID :one
SELECT * FROM disputes
WHERE id = $1;
