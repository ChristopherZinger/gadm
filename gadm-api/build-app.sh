go build -tags='no_clickhouse  \
    no_libsql \
    no_mssql \
    no_mysql \
    no_postgres \
    no_sqlite3 \
    no_vertica \
    no_ydb' -o gadm-api ./cmd/gadm-api

