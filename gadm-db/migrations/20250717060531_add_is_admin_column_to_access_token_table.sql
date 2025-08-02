-- +goose Up
-- +goose StatementBegin
ALTER TABLE access_tokens ADD COLUMN can_generate_access_tokens BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE access_tokens DROP COLUMN can_generate_access_tokens;
-- +goose StatementEnd
