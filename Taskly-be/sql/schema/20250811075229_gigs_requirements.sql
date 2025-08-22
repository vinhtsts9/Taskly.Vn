-- +goose Up
-- +goose StatementBegin
CREATE TABLE gig_requirements (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
  gig_id UUID NOT NULL,
  question TEXT NOT NULL,
  required BOOLEAN not null DEFAULT true,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS gig_requirements;
-- +goose StatementEnd
