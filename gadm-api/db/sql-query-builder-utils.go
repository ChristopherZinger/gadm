package db

import (
	"fmt"
	"gadm-api/logger"

	"github.com/Masterminds/squirrel"

	utils "gadm-api/utils"
)

var GADM_QUERY_ORDER_COLUMN_NAME = "fid"

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type SqlFilterParams struct {
	FilterColName string
	FilterVal     string
}

type SqlQueryParams struct {
	StartAtValue int
	LimitValue   int
	SqlFilterParams
}

func GetGadmBaseSelectBuilder(
	lv utils.GadmLevel,
	gidFilterValue string,
	filterColName string,
	startAtFid int,
	limit int,
) squirrel.SelectBuilder {
	baseQuery := psql.Select(
		"adm.metadata::jsonb - 'md5_geom_binary_hash' as properties",
		"g.bbox", "g.geom", "adm.geom_hash").
		From("adm").
		InnerJoin("adm_geometries g ON adm.geom_hash = g.geom_hash").
		Where(squirrel.Eq{"adm.lv": lv})

	if gidFilterValue != "" {
		baseQuery = baseQuery.Where(squirrel.Eq{fmt.Sprintf("adm.metadata::jsonb ->> '%s'", filterColName): gidFilterValue})
	}

	baseQuery = baseQuery.
		Where(squirrel.Gt{fmt.Sprintf("adm.metadata ->> '%s'", GADM_QUERY_ORDER_COLUMN_NAME): fmt.Sprintf("%d", startAtFid)}).
		OrderBy(fmt.Sprintf("adm.metadata ->> '%s'", GADM_QUERY_ORDER_COLUMN_NAME)).
		Limit(uint64(limit))

	return baseQuery
}

func GetGadmFeatureSelectBuilder(
	lv utils.GadmLevel,
	gidFilterValue string,
	filterColName string,
	startAtFid int,
	limit int,
) squirrel.SelectBuilder {
	baseGadmQuery := GetGadmBaseSelectBuilder(lv, gidFilterValue, filterColName, startAtFid, limit)

	featureQuery := psql.Select(`json_build_object(
		'type', 'Feature',
		'id', r.geom_hash,
		'geometry', ST_AsGeoJSON(r.geom)::json,
		'properties', r.properties,
		'bbox', ARRAY[
			ST_XMin(r.bbox),
			ST_YMin(r.bbox),
			ST_XMax(r.bbox),
			ST_YMax(r.bbox)
		]) as geojson`).
		FromSelect(baseGadmQuery, "r")

	return featureQuery
}

type GadmFeatureCollectionSelectBuilderParams struct {
	GadmLevel     utils.GadmLevel
	FilterValue   string
	FilterColName string
	StartAtFid    int
	PageSize      int
}

func BuildGadmFeatureCollectionSelectBuilder(
	params GadmFeatureCollectionSelectBuilderParams,
) squirrel.SelectBuilder {
	gadmFeatureQuery := GetGadmFeatureSelectBuilder(
		params.GadmLevel,
		params.FilterValue,
		params.FilterColName,
		params.StartAtFid,
		params.PageSize)

	result := psql.Select(
		`json_build_object(
			'type', 'FeatureCollection',
			'features', json_agg(geojson))`,
	).FromSelect(gadmFeatureQuery, "features")

	return result
}

func BuildGeojsonFeatureSqlQuery(
	lv utils.GadmLevel,
	gidFilterValue string,
	filterColName string,
	startAtFid int,
	limit int,
) (string, []interface{}, error) {
	gadmFeatureSelectBuilder := GetGadmFeatureSelectBuilder(

		lv,
		gidFilterValue,
		filterColName,
		startAtFid,
		limit,
	)

	sql, args, err := gadmFeatureSelectBuilder.ToSql()
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return "", nil, fmt.Errorf("failed to build sql query: %w", err)
	}
	return sql, args, nil
}

func GetNextFidSqlQuery(startAtFid int, pageSize int, filterParams SqlFilterParams) (string, []interface{}, error) {
	query := psql.Select(fmt.Sprintf("adm.metadata -> '%s'", GADM_QUERY_ORDER_COLUMN_NAME)).
		From(ADM_TABLE).
		Where(squirrel.GtOrEq{fmt.Sprintf("adm.metadata -> '%s'", GADM_QUERY_ORDER_COLUMN_NAME): startAtFid})

	if filterParams.FilterColName != "" {
		query = query.Where(squirrel.Eq{
			fmt.Sprintf("adm.metadata ->> '%s'", filterParams.FilterColName): filterParams.FilterVal})
	}

	query = query.OrderBy(fmt.Sprintf("adm.metadata -> '%s'", GADM_QUERY_ORDER_COLUMN_NAME)).
		Limit(uint64(pageSize)).
		Offset(uint64(pageSize))

	sql, args, err := query.ToSql()
	if err != nil {
		logger.Error("failed_to_get_next_fid_param", err)
		return "", nil, fmt.Errorf("failed to build next page check sql query: %w", err)
	}
	return sql, args, nil
}

func GetAccessTokenCreatedAtSqlQuery(token string) (string, []interface{}, error) {
	sql, args, err := psql.
		Select(AccessTokensTable.CreatedAt).
		From(ACCESS_TOKEN_TABLE).
		Where(squirrel.Eq{AccessTokensTable.Token: token}).
		ToSql()

	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func GetAccessTokenSqlQuery(token string) (string, []interface{}, error) {
	sql, args, err := psql.
		Select(AccessTokensTable.CreatedAt, AccessTokensTable.CanGenerateAccessTokens).
		From(ACCESS_TOKEN_TABLE).
		Where(squirrel.Eq{AccessTokensTable.Token: token}).
		ToSql()

	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func GetInsertAccessTokenWithReturningSqlQuery(email string) (string, []interface{}, error) {
	sql, args, err := psql.
		Insert(ACCESS_TOKEN_TABLE).
		Columns(AccessTokensTable.Email).
		Values(email).
		Suffix(fmt.Sprintf("RETURNING %s, %s", AccessTokensTable.Token, AccessTokensTable.CreatedAt)).
		ToSql()

	if err != nil {
		return "", nil, err
	}
	return sql, args, nil

}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func GetReverseGeocodeSqlQuery(point Point) (string, []interface{}, error) {
	withClause := ` 
		WITH input_point AS (
			SELECT ST_SetSRID(ST_MakePoint(?, ?), 4326)::geometry(Point,4326) AS pt
		),
		candidates as (
			SELECT g.geom_hash, g.geom, g.area_sq_m, ip.pt
			FROM gadm.adm_geometries AS g 
			INNER JOIN input_point ip 
			ON ST_Contains(g.bbox, ip.pt) 
			ORDER BY area_sq_m ASC
		),
		result_geometry as (
			SELECT c.geom_hash
			FROM candidates AS c 
			INNER JOIN input_point ip 
			ON ST_Contains(c.geom, ip.pt) 
			ORDER BY c.area_sq_m ASC 
			LIMIT 1
		)`

	query := psql.
		Select("(to_jsonb(adm.*) - 'id') as result").
		Prefix(withClause, point.Lng, point.Lat).
		From("adm").
		InnerJoin("result_geometry ON adm.geom_hash = result_geometry.geom_hash").
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}
