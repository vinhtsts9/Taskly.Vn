


-- name: CreateService :one
INSERT INTO gigs (
  user_id, title, description, category_id, price, delivery_time, image_url, status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;


-- name: UpdateService :one
UPDATE gigs
SET
  title = $2,
  description = $3,
  category_id = $4,
  price = $5,
  delivery_time = $6,
  image_url = $7,
  status = $8
WHERE id = $1
RETURNING *;


-- name: DeleteService :exec
DELETE FROM gigs WHERE id = $1;

-- name: GetService :one
SELECT
    g.id,
    g.user_id,
    g.title,
    g.description,
    g.price,
    g.delivery_time,
    g.image_url,
    g.status,
    g.created_at,
    c.name AS category_name,
    u.names AS user_name,
    u.profile_pic AS user_profile_pic
FROM
    gigs g
JOIN
    users u ON g.user_id = u.id
JOIN
    categories c ON g.category_id = c.id
WHERE
    g.id = $1;
