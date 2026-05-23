package adm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Masterminds/squirrel"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Adm struct {
	Metadata json.RawMessage `db:"metadata" json:"metadata"`
	ID       string          `db:"id" json:"id"`
	Level    int             `db:"lv" json:"lv"`
	GeomHash string          `db:"geom_hash" json:"geom_hash"`
}

type Repo struct {
	pgConn *pgxpool.Pool
}

func NewAdmRepo(pg *pgxpool.Pool) *Repo {
	return &Repo{pgConn: pg}
}

func (repo *Repo) GetAdmNeighbors(ctx context.Context, admId string) ([]Adm, error) {
	sql, args, err := getAdmNeighborsSqlQuery(admId)
	if err != nil {
		return nil, fmt.Errorf("failed_to_build_query: %w", err)
	}

	rows, err := repo.pgConn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"failed_to_query_database_for_adm_neighbors: sql_query: %s: %w",
			sql, err)
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[Adm])
	if err != nil {
		return nil, fmt.Errorf("failed_to_collect_rows: %w", err)
	}

	return result, nil
}

func getAdmNeighborsSqlQuery(admId string) (string, []interface{}, error) {
	withClause := `
		with ids as (SELECT DISTINCT id FROM adm_neighbors 
			JOIN ADM ON adm.id=n1 OR adm.id=n2
			WHERE n1=$1
			OR n2=$1)`

	query := psql.
		Select("adm.metadata", "adm.id", "adm.lv", "adm.geom_hash").
		Prefix(withClause, admId).
		From("ids").
		LeftJoin("adm ON ids.id=adm.id")

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}
