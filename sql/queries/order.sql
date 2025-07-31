-- name: CreateOrder :one
INSERT INTO orders (gig_id, buyer_id, seller_id, total_price, delivery_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListOrdersByUser :many
SELECT * FROM orders
WHERE buyer_id = $1 OR seller_id = $1
ORDER BY order_date DESC;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;

-- name: SubmitOrderDelivery :exec
UPDATE orders
SET status = 'delivered'
WHERE id = $1 AND status = 'active';

-- name: AcceptOrderCompletion :exec
UPDATE orders
SET status = 'completed'
WHERE id = $1 AND status = 'delivered';
