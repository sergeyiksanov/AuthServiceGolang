-- +goose Up
-- +goose StatementBegin
DELETE FROM credentials;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
