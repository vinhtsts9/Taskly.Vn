-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE gigs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    user_id UUID NOT NULL,
    title VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    category_id INT NOT NULL,
    price FLOAT8 NOT NULL CHECK (price >= 0),
    delivery_time INT NOT NULL CHECK (delivery_time > 0),
    image_url TEXT,
    status VARCHAR(10) not null CHECK (status IN ('active', 'paused', 'draft')) DEFAULT 'draft',
    created_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS gigs;
-- +goose StatementEnd
