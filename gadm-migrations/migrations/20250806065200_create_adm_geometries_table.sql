-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE gadm.adm_geometries (
    geom_hash TEXT PRIMARY KEY,
    geom GEOMETRY(MultiPolygon, 4326) NOT NULL,
    bbox GEOMETRY(Polygon, 4326) 
        GENERATED ALWAYS AS (ST_Envelope(geom)) STORED,
    area_sq_m double precision 
        GENERATED ALWAYS AS (ST_Area(geom::geography)) STORED
);
CREATE INDEX idx_adm_geometries_geom ON gadm.adm_geometries USING GIST (geom);
CREATE INDEX idx_adm_geometries_bbox ON gadm.adm_geometries USING GIST (bbox);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX IF EXISTS gadm.idx_adm_geometries_geom;
DROP INDEX IF EXISTS gadm.idx_adm_geometries_bbox;
DROP TABLE gadm.adm_geometries;
-- +goose StatementEnd
