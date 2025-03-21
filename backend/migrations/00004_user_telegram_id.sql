-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN telegram_id BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN telegram_id;
-- +goose StatementEnd
