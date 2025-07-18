# GADM api

This project exposes GADM open-source database through REST api endpoints.

## Running Migrations

Migrations should be run by one user who is the owner of all tables.
Scripts `run-migration.sh` and `rollback-migrations.sh` will run or rollback
a single migration at the time.

## Teardown and setup the database

The `init_database.sh` script will reset old database and create new one.
It will setup appropriate roles and ingest GADM geopackage. It requires
GDAL and ogr2ogr to be per installed. The script has to be run from /gadm-app
directory.
