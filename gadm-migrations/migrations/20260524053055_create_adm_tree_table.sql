-- +goose Up
CREATE TABLE IF NOT EXISTS gadm.adm_tree (
    parent UUID NOT NULL REFERENCES gadm.adm(id) ON DELETE CASCADE,
    child  UUID          REFERENCES gadm.adm(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_adm_tree_parent_child
    ON gadm.adm_tree (parent, child) NULLS NOT DISTINCT;

CREATE INDEX IF NOT EXISTS idx_adm_tree_child
    ON gadm.adm_tree (child);

-- +goose Down
DROP TABLE IF EXISTS gadm.adm_tree;
