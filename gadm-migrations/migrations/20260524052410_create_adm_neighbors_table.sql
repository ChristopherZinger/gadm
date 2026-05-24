-- +goose Up
CREATE TABLE IF NOT EXISTS gadm.adm_neighbors (
    n1 UUID NOT NULL REFERENCES gadm.adm(id) ON DELETE CASCADE,
    n2 UUID NOT NULL REFERENCES gadm.adm(id) ON DELETE CASCADE,
    CONSTRAINT relation_no_self   CHECK (n1 <> n2),
    CONSTRAINT relation_canonical CHECK (n1 < n2)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_adm_neighbors_n1_n2
    ON gadm.adm_neighbors (n1, n2);

CREATE INDEX IF NOT EXISTS idx_adm_neighbors_n2
    ON gadm.adm_neighbors (n2);

-- +goose Down
DROP TABLE IF EXISTS gadm.adm_neighbors;