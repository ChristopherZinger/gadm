-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';


SELECT 'inserting adm_0';
INSERT INTO gadm.adm (lv, geom_hash, metadata)
SELECT 
0 AS lv, 
t.md5_geom_binary_hash AS geom_hash,
(to_jsonb(t) - 'geom') AS metadata
FROM gadm.adm_0 AS t;

SELECT 'inserting adm_1';
INSERT INTO gadm.adm (lv, geom_hash, metadata)
SELECT 
1 AS lv, 
t.md5_geom_binary_hash AS geom_hash, 
(to_jsonb(t) - 'geom') AS metadata
FROM gadm.adm_1 AS t;

SELECT 'inserting adm_2';
INSERT INTO gadm.adm (lv, geom_hash, metadata)
SELECT
2 AS lv,
t.md5_geom_binary_hash AS geom_hash,
(to_jsonb(t) - 'geom') AS metadata
FROM gadm.adm_2 AS t;

SELECT 'inserting adm_3';
INSERT INTO gadm.adm (lv, geom_hash, metadata)
SELECT
3 AS lv,
t.md5_geom_binary_hash AS geom_hash,
(to_jsonb(t) - 'geom') AS metadata
FROM gadm.adm_3 AS t;

SELECT 'inserting adm_4';
INSERT INTO gadm.adm (lv, geom_hash, metadata)
SELECT
4 AS lv,
t.md5_geom_binary_hash AS geom_hash,
(to_jsonb(t) - 'geom') AS metadata
FROM gadm.adm_4 AS t;

SELECT 'inserting adm_5';
INSERT INTO gadm.adm (lv, geom_hash, metadata)
SELECT
5 AS lv,
t.md5_geom_binary_hash AS geom_hash,
(to_jsonb(t) - 'geom') AS metadata
FROM gadm.adm_5 AS t;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
TRUNCATE TABLE gadm.adm;
-- +goose StatementEnd
