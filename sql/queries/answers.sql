-- name: CreateAnswer :one
INSERT INTO answers (
  gig_id,
  user_id,
  question_id,
  answer
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetAnswersByOrderID :many
SELECT 
  a.id, 
  a.gig_id, 
  a.user_id, 
  a.question_id, 
  a.answer, 
  a.created_at, 
  a.updated_at
FROM answers a
JOIN orders o ON a.gig_id = o.gig_id AND a.user_id = o.buyer_id
WHERE o.id = $1;

