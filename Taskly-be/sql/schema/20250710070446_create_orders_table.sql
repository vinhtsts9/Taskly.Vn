-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    gig_id UUID NOT NULL,
    buyer_id UUID NOT NULL,
    seller_id UUID NOT NULL,
    status VARCHAR(15) not null CHECK (status IN ('pending','paid', 'active', 'delivered', 'completed', 'cancelled')) DEFAULT 'pending',
    total_price FLOAT8 NOT NULL CHECK (total_price >= 0),
    order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    delivery_date TIMESTAMP,
    idempotency_key VARCHAR(150) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
