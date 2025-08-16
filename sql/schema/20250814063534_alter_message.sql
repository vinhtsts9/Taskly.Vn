-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX ux_rooms_user_pair
ON rooms (LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id));
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ux_rooms_user_pair;
-- +goose StatementEnd
