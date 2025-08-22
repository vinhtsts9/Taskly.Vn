-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE disputes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    order_id UUID NOT NULL ,
    user_id UUID NOT NULL ,
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'UNDER_REVIEW', 'RESOLVED', 'REFUND', 'REJECT')),
    created_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists disputes;
-- +goose StatementEnd
