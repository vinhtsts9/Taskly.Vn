-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE wallet (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    user_id UUID NOT NULL REFERENCES users(id),
    balance_bigint BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE topup_order (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    user_id UUID NOT NULL REFERENCES users(id),
    reference_code VARCHAR(100) UNIQUE NOT NULL,
    amount_bigint BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'VND',
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, VERIFIED, FAILED
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE payment_transaction (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    topup_order_id UUID REFERENCES topup_order(id),
    provider_tx_id VARCHAR(100) NOT NULL,
    amount_bigint BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL, -- PENDING, SUCCESS, FAILED
    remote_payload JSONB,
    signature TEXT,
    verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE ledger_entry (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    wallet_id UUID REFERENCES wallet(id),
    transaction_id UUID REFERENCES payment_transaction(id),
    amount_bigint BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    entry_type VARCHAR(50) NOT NULL, -- CREDIT, DEBIT
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ledger_entry;
DROP TABLE IF EXISTS payment_transaction;
DROP TABLE IF EXISTS topup_order;
DROP TABLE IF EXISTS wallet;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +goose StatementEnd
