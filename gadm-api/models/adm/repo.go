package adm

import (
	"context"
	"encoding/json"
	"fmt"
	"gadm-api/utils"

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

func (repo *Repo) GetAdmForPoint(ctx context.Context, point utils.Point) (Adm, error) {
	sql, args, err := getAdmForPointSqlQuery(point)
	if err != nil {
		return Adm{}, fmt.Errorf("failed_to_build_query: %w", err)
	}

	var adm Adm
	err = repo.pgConn.QueryRow(ctx, sql, args...).
		Scan(&adm.Metadata, &adm.ID, &adm.Level, &adm.GeomHash)
	if err != nil {
		return Adm{}, fmt.Errorf(
			"failed_to_query_database_for_adm_for_lat_lng: sql_query: %s: %w",
			sql, err)
	}
	return adm, nil
}

func (repo *Repo) GetAdmById(ctx context.Context, admId string) (Adm, error) {
	sql, args, err := getSelectOneAdmByIdSqlQuery(admId)
	if err != nil {
		return Adm{}, fmt.Errorf("failed_to_build_query: %w", err)
	}

	var adm Adm
	err = repo.pgConn.QueryRow(ctx, sql, args...).
		Scan(&adm.Metadata, &adm.ID, &adm.Level, &adm.GeomHash)
	if err != nil {
		return Adm{}, fmt.Errorf("failed_to_query_database_for_adm_by_id: sql_query: %s: %w", sql, err)
	}
	return adm, nil
}

func (repo *Repo) GetAdmsDirectChildrenForId(ctx context.Context, admId string, lv int) ([]Adm, error) {
	sql, args, err := getSelectAdmDirectChildrenForIdSqlQuery(admId, lv)
	if err != nil {
		return nil, fmt.Errorf("failed_to_build_query: %w", err)
	}

	rows, err := repo.pgConn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed_to_query_database_for_adm_direct_children_for_id: sql_query: %s: %w", sql, err)
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[Adm])
	if err != nil {
		return nil, fmt.Errorf("failed_to_collect_rows: %w", err)
	}

	return result, nil
}

func (repo *Repo) UpsertAdmTreeRelationships(ctx context.Context, parentId string, childIds []string) error {
	sql, args, err := getUpsertAdmTreeSqlQuery(parentId, childIds)
	if err != nil {
		return fmt.Errorf("failed_to_build_query: %w", err)
	}

	_, err = repo.pgConn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed_to_upsert_adm_tree_relationships: sql_query: %s: %w", sql, err)
	}
	return nil
}

func (repo *Repo) GetAdms(ctx context.Context, startAfterId string, batchSize int) ([]Adm, error) {
	sql, args, err := getSelectAdmsSqlQuery(startAfterId, batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed_to_build_query: %w", err)
	}

	rows, err := repo.pgConn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed_to_query_database_for_adm_tree_relationships: sql_query: %s: %w", sql, err)
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[Adm])
	if err != nil {
		return nil, fmt.Errorf("failed_to_collect_rows: %w", err)
	}

	return result, nil
}

func (repo *Repo) GetNeighbors(ctx context.Context, admId string) ([]Adm, error) {
	sql, args, err := getNeighborsSqlQuery(admId)
	if err != nil {
		return nil, fmt.Errorf("failed_to_build_query: %w", err)
	}

	rows, err := repo.pgConn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf(
			"failed_to_query_database_for_neighbors: sql_query: %s: %w",
			sql, err)
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[Adm])
	if err != nil {
		return nil, fmt.Errorf("failed_to_collect_rows: %w", err)
	}

	return result, nil
}

func (repo *Repo) UpsertAdmNeighbors(ctx context.Context, n1Id, n2Id string) error {
	sql, args, err := getUpsertAdmNeighborsSqlQuery(n1Id, n2Id)
	if err != nil {
		return fmt.Errorf("failed_to_build_query: %w", err)
	}

	_, err = repo.pgConn.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed_to_upsert_adm_neighbors: sql_query: %s: %w", sql, err)
	}
	return nil
}

func (repo *Repo) GetLeafAdms(ctx context.Context, startAfterId string, batchSize int) ([]Adm, error) {
	sql, args, err := getSelectLeafAdmsSqlQuery(startAfterId, batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed_to_build_query: %w", err)
	}

	rows, err := repo.pgConn.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed_to_query_database_for_leaf_adms: sql_query: %s: %w", sql, err)
	}
	defer rows.Close()

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[Adm])
	if err != nil {
		return nil, fmt.Errorf("failed_to_collect_rows: %w", err)
	}

	return result, nil
}

func (repo *Repo) IterateLeafAdms(
	ctx context.Context,
	batchSize int,
	fn func(ctx context.Context, batch []Adm) error,
) error {
	var startAfterId string
	for {
		batch, err := repo.GetLeafAdms(ctx, startAfterId, batchSize)
		if err != nil {
			return fmt.Errorf("failed_to_get_leaf_adms_batch: %w", err)
		}
		if len(batch) == 0 {
			return nil
		}

		if err := fn(ctx, batch); err != nil {
			return err
		}

		if len(batch) < batchSize {
			return nil
		}
		startAfterId = batch[len(batch)-1].ID
	}
}
