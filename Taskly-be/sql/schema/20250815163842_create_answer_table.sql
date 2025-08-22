-- +goose Up
-- +goose StatementBegin
CREATE TABLE answers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
  order_id UUID NOT NULL,
  user_id UUID NOT NULL,
  question_id UUID NOT NULL,
  answer TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if EXISTS answers;
-- +goose StatementEnd
