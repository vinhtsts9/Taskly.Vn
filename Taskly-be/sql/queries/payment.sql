-- name:CreatePayment
INSERT INTO topup_order (user_id, reference_code, amount_bigint, status)
VALUES (:user_id, :reference_code, :amount, 'PENDING')
RETURNING id;
