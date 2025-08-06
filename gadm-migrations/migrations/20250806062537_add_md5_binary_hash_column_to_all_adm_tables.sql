-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

SELECT 'Adding md5_geom_binary_hash column to adm_0';
ALTER TABLE adm_0 ADD COLUMN md5_geom_binary_hash TEXT
    GENERATED ALWAYS AS (md5(ST_AsBinary(geom))) STORED;

SELECT 'Adding md5_geom_binary_hash column to adm_1';
ALTER TABLE adm_1 ADD COLUMN md5_geom_binary_hash TEXT
    GENERATED ALWAYS AS (md5(ST_AsBinary(geom))) STORED;

SELECT 'Adding md5_geom_binary_hash column to adm_2';
ALTER TABLE adm_2 ADD COLUMN md5_geom_binary_hash TEXT
    GENERATED ALWAYS AS (md5(ST_AsBinary(geom))) STORED;

SELECT 'Adding md5_geom_binary_hash column to adm_3';
ALTER TABLE adm_3 ADD COLUMN md5_geom_binary_hash TEXT
    GENERATED ALWAYS AS (md5(ST_AsBinary(geom))) STORED;

SELECT 'Adding md5_geom_binary_hash column to adm_4';
ALTER TABLE adm_4 ADD COLUMN md5_geom_binary_hash TEXT
    GENERATED ALWAYS AS (md5(ST_AsBinary(geom))) STORED;

SELECT 'Adding md5_geom_binary_hash column to adm_5';
ALTER TABLE adm_5 ADD COLUMN md5_geom_binary_hash TEXT
    GENERATED ALWAYS AS (md5(ST_AsBinary(geom))) STORED;

SELECT 'Creating unique indexes on md5_geom_binary_hash columns';
CREATE UNIQUE INDEX idx_adm_0_md5_geom_binary_hash ON adm_0(md5_geom_binary_hash);
CREATE UNIQUE INDEX idx_adm_1_md5_geom_binary_hash ON adm_1(md5_geom_binary_hash);
CREATE UNIQUE INDEX idx_adm_2_md5_geom_binary_hash ON adm_2(md5_geom_binary_hash);
CREATE UNIQUE INDEX idx_adm_3_md5_geom_binary_hash ON adm_3(md5_geom_binary_hash);
CREATE UNIQUE INDEX idx_adm_4_md5_geom_binary_hash ON adm_4(md5_geom_binary_hash);
CREATE UNIQUE INDEX idx_adm_5_md5_geom_binary_hash ON adm_5(md5_geom_binary_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

SELECT 'Dropping unique indexes on md5_geom_binary_hash columns';
DROP INDEX IF EXISTS idx_adm_0_md5_geom_binary_hash;
DROP INDEX IF EXISTS idx_adm_1_md5_geom_binary_hash;
DROP INDEX IF EXISTS idx_adm_2_md5_geom_binary_hash;
DROP INDEX IF EXISTS idx_adm_3_md5_geom_binary_hash;
DROP INDEX IF EXISTS idx_adm_4_md5_geom_binary_hash;
DROP INDEX IF EXISTS idx_adm_5_md5_geom_binary_hash;

SELECT 'Dropping md5_geom_binary_hash columns';
ALTER TABLE adm_0 DROP COLUMN md5_geom_binary_hash;
ALTER TABLE adm_1 DROP COLUMN md5_geom_binary_hash;
ALTER TABLE adm_2 DROP COLUMN md5_geom_binary_hash;
ALTER TABLE adm_3 DROP COLUMN md5_geom_binary_hash;
ALTER TABLE adm_4 DROP COLUMN md5_geom_binary_hash;
ALTER TABLE adm_5 DROP COLUMN md5_geom_binary_hash;
-- +goose StatementEnd
