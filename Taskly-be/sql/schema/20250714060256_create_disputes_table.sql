-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'dispute_status') THEN
        CREATE TYPE dispute_status AS ENUM ('pending', 'under_review', 'resolved', 'refunded', 'rejected');
    END IF;
END
$$;
CREATE TABLE disputes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    order_id UUID NOT NULL ,
    user_id UUID NOT NULL ,
    reason TEXT NOT NULL,
    status dispute_status not null DEFAULT 'pending', -- 'pending', 'resolved', 'rejected' , 'under _review', 'refunded'
    created_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists disputes;
DROP TYPE IF EXISTS dispute_status;
-- +goose StatementEnd
