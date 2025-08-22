-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE message_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    message_id UUID NOT NULL,
    file_url TEXT NOT NULL,
    file_type VARCHAR(50), -- 'image', 'video', 'pdf', etc.
    file_name VARCHAR(255),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message_attachments;
-- +goose StatementEnd
