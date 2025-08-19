package db

import "github.com/jackc/pgx/v5/pgxpool"

type PgConn struct {
	Db *pgxpool.Pool
}

func CreatePgConnector(db *pgxpool.Pool) *PgConn {
	return &PgConn{Db: db}
}
