# Ingest data from GADM to postgres

## Setup postgres

User `docker-compose.yml` file with `postgis/postgis:16-3.5` image.

This image comes with preinstalled postgis extension. Unfortunately it is
installed on public schema.

## Ingest GADM database

First download geopackage GADM database - the one that includes all layers.

Second - instal ogr2ogr tool. On mac it is a part of gamd

```
    $ brew install gdal
    $ export PATH=$PATH:/opt/homebrew/opt/gdal/bin
```

Next use `ingest-gadm-pkg-to-pg.sh` script to move data to postgres instance.
Make sure you have `psql` installed and that it matches version of postgres
image.

## rename tables ?

## add indexes ?
