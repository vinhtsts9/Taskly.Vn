-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS user_verify (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    verify_key varchar(255) NOT NULL UNIQUE,
    verify_hash_key varchar(255) NOT NULL UNIQUE,
    verify_otp varchar(255) NOT NULL,
    verify_type varchar(50)  NOT NULL,
    is_deleted boolean NOT NULL DEFAULT false,
    is_verified boolean NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_verify;
DROP TYPE IF EXISTS verify_type_enum;
-- +goose StatementEnd
