-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE gigs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    user_id UUID NOT NULL,
    title VARCHAR(150) NOT NULL,
    category_id INT[] NOT NULL,
    image_url TEXT[] not null,
    description text not null,
    pricing_mode TEXT NOT NULL DEFAULT 'single'
        CHECK (pricing_mode IN ('single', 'triple')),
    status VARCHAR(10) not null CHECK (status IN ('active', 'paused', 'draft')) DEFAULT 'draft',
    created_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS gigs;
-- +goose StatementEnd
