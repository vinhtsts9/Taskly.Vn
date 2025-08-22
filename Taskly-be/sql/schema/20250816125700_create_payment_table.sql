-- +goose Up
-- +goose StatementBegin

-- Enable uuid generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Wallet: store per-user balance
CREATE TABLE IF NOT EXISTS wallet (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    user_id UUID NOT NULL REFERENCES users(id),
    balance_bigint BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Topup order: merchant-side representation of a fund top-up request
CREATE TABLE IF NOT EXISTS topup_order (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    user_id UUID NOT NULL REFERENCES users(id),
    order_id UUID NULL, -- optional link to an order if this topup is for an order
    reference_code VARCHAR(150) UNIQUE NOT NULL,
    idempotency_key VARCHAR(150) UNIQUE NOT NULL, -- optional idempotency key for retry safety
    amount_bigint BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'VND',
    provider VARCHAR(50) NOT NULL DEFAULT 'vnpay',
    provider_payment_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING','COMPLETED','FAILED','EXPIRED','CANCELLED')),
    
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE
);

-- Payment transaction: records provider transaction callbacks/events
CREATE TABLE IF NOT EXISTS payment_transaction (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    topup_order_id UUID REFERENCES topup_order(id) ON DELETE SET NULL,
    provider VARCHAR(50) NOT NULL DEFAULT 'vnpay',
    provider_tx_id VARCHAR(150) NOT NULL,
    amount_bigint BIGINT NOT NULL,
    status          VARCHAR(20) NOT NULL
                     CHECK (status IN ('PENDING','COMPLETED','FAILED','REFUNDED')),
    remote_payload JSONB,
    signature TEXT,
    verified_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Ledger entries: record wallet movements
CREATE TABLE IF NOT EXISTS ledger_entry (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    wallet_id UUID REFERENCES wallet(id) ON DELETE CASCADE,
    transaction_id UUID REFERENCES payment_transaction(id) ON DELETE SET NULL,
    amount_bigint BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    entry_type      VARCHAR(10) NOT NULL
                     CHECK (entry_type IN ('CREDIT','DEBIT')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes and constraints
-- Unique idempotency key when provided
CREATE UNIQUE INDEX IF NOT EXISTS uq_topup_order_idempotency_key ON topup_order(idempotency_key)
    WHERE idempotency_key IS NOT NULL;

-- Indexes for lookup performance
CREATE INDEX IF NOT EXISTS idx_topup_order_user_status ON topup_order(user_id, status);
CREATE INDEX IF NOT EXISTS idx_topup_order_reference ON topup_order(reference_code);
CREATE INDEX IF NOT EXISTS idx_topup_order_order_id ON topup_order(order_id);

-- Ensure provider transaction uniqueness per provider
CREATE UNIQUE INDEX IF NOT EXISTS uq_payment_transaction_provider_tx
    ON payment_transaction(provider, provider_tx_id)
    WHERE provider_tx_id IS NOT NULL;

-- Trigger function to auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trg_topup_updated_at
BEFORE UPDATE ON topup_order
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE TRIGGER trg_wallet_updated_at
BEFORE UPDATE ON wallet
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Optional: helper view to see pending topups
CREATE OR REPLACE VIEW view_pending_topups AS
SELECT t.id, t.user_id, t.reference_code, t.amount_bigint, t.currency, t.provider, t.status, t.created_at, t.expires_at
FROM topup_order t
WHERE t.status = 'PENDING';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP VIEW IF EXISTS view_pending_topups;

DROP TRIGGER IF EXISTS trg_topup_updated_at ON topup_order;
DROP TRIGGER IF EXISTS trg_wallet_updated_at ON wallet;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS uq_payment_transaction_provider_tx;
DROP INDEX IF EXISTS idx_topup_order_order_id;
DROP INDEX IF EXISTS idx_topup_order_reference;
DROP INDEX IF EXISTS idx_topup_order_user_status;
DROP INDEX IF EXISTS uq_topup_order_idempotency_key;

DROP TABLE IF EXISTS ledger_entry;
DROP TABLE IF EXISTS payment_transaction;
DROP TABLE IF EXISTS topup_order;
DROP TABLE IF EXISTS wallet;

-- Note: do NOT drop users table here (if it exists in your app). If you created a users table in the UP block,
-- consider commenting out the DROP for it, or handle it explicitly.

-- +goose StatementEnd
