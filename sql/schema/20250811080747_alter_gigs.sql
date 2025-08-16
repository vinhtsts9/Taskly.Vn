-- +goose Up
-- +goose StatementBegin

ALTER TABLE gigs
    add category_id INT[] NOT NULL,
    add image_url TEXT[] not null,
    add updated_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
    add description text not null,
    ADD COLUMN pricing_mode TEXT NOT NULL DEFAULT 'single'
        CHECK (pricing_mode IN ('single', 'triple'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER table gigs
DROP COLUMN category_id,
DROP COLUMN image_url,
DROP COLUMN updated_at,
drop COLUMN description,
drop COLUMN pricing_mode;
-- +goose StatementEnd
