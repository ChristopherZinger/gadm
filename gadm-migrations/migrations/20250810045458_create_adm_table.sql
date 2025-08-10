-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE gadm.adm (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
lv INTEGER,
geom_hash TEXT,
metadata JSONB
);
CREATE INDEX idx_adm_geom_hash ON gadm.adm (geom_hash);
CREATE INDEX idx_adm_lv ON gadm.adm (lv);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX IF EXISTS gadm.idx_adm_geom_hash;
DROP INDEX IF EXISTS gadm.idx_adm_lv;
DROP TABLE IF EXISTS gadm.adm;
-- +goose StatementEnd
