-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access_tokens (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL DEFAULT gen_random_uuid (),
    email TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_access_tokens_token ON access_tokens (token) INCLUDE (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_access_tokens_token;
DROP TABLE IF EXISTS access_tokens;
-- +goose StatementEnd