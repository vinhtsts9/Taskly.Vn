-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS user_base (
    user_base_id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    email VARCHAR(255) UNIQUE NOT NULL,
    salt VARCHAR(255) NOT NULL,
    passwords VARCHAR(255) NOT NULL,
    is_two_factor_enabled SMALLINT NOT NULL DEFAULT 0,
    login_time TIMESTAMP DEFAULT NULL,
    logout_time TIMESTAMP DEFAULT NULL,
    login_ip VARCHAR(255) NOT NULL DEFAULT '',
    states SMALLINT NOT NULL DEFAULT 3 CHECK (states >= 0),
    refresh_token VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_base;
-- +goose StatementEnd
