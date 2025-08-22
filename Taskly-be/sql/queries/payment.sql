-- name: CreateWallet :one
INSERT INTO wallet (user_id, balance_bigint)
VALUES ($1, $2)
RETURNING id, user_id, balance_bigint, updated_at, created_at;

-- name: GetWalletByUserID :one
SELECT id, user_id, balance_bigint, updated_at, created_at
FROM wallet
WHERE user_id = $1
LIMIT 1;

-- name: GetWalletForUpdateByUser :one
SELECT id, user_id, balance_bigint, updated_at, created_at
FROM wallet
WHERE user_id = $1
FOR UPDATE;

-- name: UpdateWalletBalance :one
UPDATE wallet
SET balance_bigint = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, balance_bigint, updated_at, created_at;


-- =========================
-- Topup / Payment intent
-- =========================

-- name: CreateTopupOrder :one
INSERT INTO topup_order (
  user_id,
  order_id,
  reference_code,
  idempotency_key,
  amount_bigint,
  currency,
  provider,
  provider_payment_url,
  status,
  expires_at
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
ON CONFLICT (idempotency_key) DO UPDATE
  SET updated_at = NOW()
RETURNING id, user_id, order_id, reference_code, idempotency_key, amount_bigint, currency, provider, provider_payment_url, status, created_at, updated_at, expires_at;



-- name: GetTopupByIdempotencyKeyAndOrder :one
SELECT *
FROM topup_order
WHERE idempotency_key = $1
  AND order_id = $2
  AND status = 'PENDING'
LIMIT 1;

-- name: GetLastPendingTopupByOrder :one
SELECT id, user_id, order_id, reference_code, idempotency_key, amount_bigint, currency, provider, provider_payment_url, status, created_at, updated_at, expires_at
FROM topup_order
WHERE order_id = $1 AND status = 'PENDING'
ORDER BY created_at DESC
LIMIT 1;

-- name: GetTopupOrderByReference :one
SELECT id, user_id, order_id, reference_code, idempotency_key, amount_bigint, currency, provider, provider_payment_url, status, created_at, updated_at, expires_at
FROM topup_order
WHERE reference_code = $1
LIMIT 1;


-- name: UpdateTopupPaymentInfo :one
UPDATE topup_order
SET provider_payment_url = $2, expires_at = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, provider_payment_url, expires_at, updated_at;

-- name: UpdateTopupOrderStatus :one
UPDATE topup_order
SET status = $2,  updated_at = NOW()
WHERE id = $1
RETURNING id, status, updated_at,order_id;


-- =========================
-- Payment Transaction (provider callbacks)
-- =========================

-- name: CreatePaymentTransaction :one
INSERT INTO payment_transaction (
  topup_order_id,
  provider,
  provider_tx_id,
  amount_bigint,
  status,
  remote_payload,
  signature,
  verified_at
)
VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7, $8)
RETURNING id, topup_order_id, provider, provider_tx_id, amount_bigint, status, remote_payload, signature, verified_at, created_at;

-- name: GetPaymentTransactionByProviderTx :one
SELECT id, topup_order_id, provider, provider_tx_id, amount_bigint, status, remote_payload, signature, verified_at, created_at
FROM payment_transaction
WHERE provider = $1 AND provider_tx_id = $2
LIMIT 1;


-- =========================
-- Ledger entries
-- =========================

-- name: CreateLedgerEntry :one
INSERT INTO ledger_entry (wallet_id, transaction_id, amount_bigint, balance_after, entry_type)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, wallet_id, transaction_id, amount_bigint, balance_after, entry_type, created_at;

-- name: GetLedgerEntriesByWallet :many
SELECT id, wallet_id, transaction_id, amount_bigint, balance_after, entry_type, created_at
FROM ledger_entry
WHERE wallet_id = $1
ORDER BY created_at DESC;
