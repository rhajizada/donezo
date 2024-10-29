-- +goose Up
-- +goose StatementBegin
INSERT INTO boards (name)
VALUES (strftime('%m-%d-%Y', 'now'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM boards
WHERE name = strftime('%m-%d-%Y', 'now');
-- +goose StatementEnd

