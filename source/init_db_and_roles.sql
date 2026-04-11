BEGIN ISOLATION LEVEL SERIALIZABLE;
-- database: gadm-db
-- schemas: common, gadm, postgis

-- create schemas as admin
CREATE SCHEMA gadm;
CREATE SCHEMA postgis;
CREATE SCHEMA common;

ALTER DATABASE "gadm-db" SET search_path TO common, gadm, postgis, public;


CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA postgis;

COMMIT;
