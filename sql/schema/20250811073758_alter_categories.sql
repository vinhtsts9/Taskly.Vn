-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories
ADD COLUMN parent_id INT REFERENCES categories(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories
DROP COLUMN parent_id;
-- +goose StatementEnd
