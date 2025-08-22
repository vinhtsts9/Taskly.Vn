-- +goose Up
-- +goose StatementBegin
CREATE TABLE gig_packages (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
  gig_id UUID NOT NULL,
  tier VARCHAR(50) NOT NULL,
  price FLOAT8 NOT NULL CHECK (price >= 0),
  delivery_time INT NOT NULL CHECK (delivery_time > 0),
  options JSONB DEFAULT '{}'::jsonb,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS gig_packages;
-- +goose StatementEnd
