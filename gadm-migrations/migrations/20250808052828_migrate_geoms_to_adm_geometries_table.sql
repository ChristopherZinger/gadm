-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

INSERT INTO gadm.adm_geometries (geom_hash, geom)
    SELECT md5_geom_binary_hash, geom
    FROM adm_0
    ON CONFLICT (geom_hash) DO NOTHING;

INSERT INTO gadm.adm_geometries (geom_hash, geom)
    SELECT md5_geom_binary_hash, geom
    FROM adm_1
    ON CONFLICT (geom_hash) DO NOTHING;

INSERT INTO gadm.adm_geometries (geom_hash, geom)
    SELECT md5_geom_binary_hash, geom
    FROM adm_2
    ON CONFLICT (geom_hash) DO NOTHING;

INSERT INTO gadm.adm_geometries (geom_hash, geom)
    SELECT md5_geom_binary_hash, geom
    FROM adm_3
    ON CONFLICT (geom_hash) DO NOTHING;

INSERT INTO gadm.adm_geometries (geom_hash, geom)
    SELECT md5_geom_binary_hash, geom
    FROM adm_4
    ON CONFLICT (geom_hash) DO NOTHING;

INSERT INTO gadm.adm_geometries (geom_hash, geom)
    SELECT md5_geom_binary_hash, geom
    FROM adm_5
    ON CONFLICT (geom_hash) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
TRUNCATE TABLE gadm.adm_geometries;
-- +goose StatementEnd
