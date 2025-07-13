package main

import (
	"fmt"
	"gadm-api/logger"

	"github.com/Masterminds/squirrel"
)

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

func buildGeojsonFeaturePropertiesSqlExpression(columns ...string) string {
	v := "json_build_object("
	len := len(columns)
	for i, column := range columns {
		v += fmt.Sprintf("'%s', %s", column, column)
		if i < len-1 {
			v += ", "
		}
	}
	v += ")"
	return v
}

func buildGeojsonFeatureSqlExpression(params GeoJSONFeatureConfig) string {
	featurePropertiesSqlExpr := buildGeojsonFeaturePropertiesSqlExpression(
		params.FeaturePropertiesNames...,
	)
	jsonBuildObjectFeature := fmt.Sprintf(
		`json_build_object(
			'type', 'Feature',
			'geometry', ST_AsGeoJSON(%[1]s)::json,
			'properties', %[2]s
		)`, params.GeometryColumnName, featurePropertiesSqlExpr,
	)
	return jsonBuildObjectFeature
}

func buildGeojsonFeatureCollectionSqlExpression(params GeoJSONFeatureConfig) string {
	geojsonFeatureExpression := buildGeojsonFeatureSqlExpression(params)
	jsonBuildObject := fmt.Sprintf(
		`json_build_object(
				'type', 'FeatureCollection',
				'features', json_agg(%s)
			)`,
		geojsonFeatureExpression,
	)
	return jsonBuildObject
}

func buildGeojsonFeatureSqlQuery(
	gadmLevel GadmLevel,
	queryParams SqlQueryParams,
) (string, []interface{}, error) {
	handlerConfig := geojsonHandlerQueryConfig[gadmLevel]

	geojsonFeatureExpression := buildGeojsonFeatureSqlExpression(
		GeoJSONFeatureConfig{
			FeaturePropertiesNames: handlerConfig.FeaturePropertiesNames,
			GeometryColumnName:     handlerConfig.GeometryColumnName,
		},
	)

	query := psql.Select(geojsonFeatureExpression).
		From(handlerConfig.TableName).
		Where(squirrel.GtOrEq{handlerConfig.OrderByColumnName: max(queryParams.StartAtValue, MIN_FID)})

	if queryParams.FilterColName != "" {
		query = query.Where(squirrel.Eq{queryParams.FilterColName: queryParams.FilterVal})
	}

	query = query.
		OrderBy(fmt.Sprintf("%s ASC", handlerConfig.OrderByColumnName)).
		Limit(uint64(queryParams.LimitValue))

	sql, args, err := query.ToSql()
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return "", nil, fmt.Errorf("failed to build sql query: %w", err)
	}
	return sql, args, nil
}

func buildFeatureCollectionSqlQuery(
	gadmLevel GadmLevel,
	queryParams SqlQueryParams,
) (string, []interface{}, error) {
	queryConfig := featureCollectionHandlerQueryConfig[gadmLevel]

	featureCollectionSqlExpression := buildGeojsonFeatureCollectionSqlExpression(
		GeoJSONFeatureConfig{queryConfig.FeaturePropertiesNames, Adm0.Geometry},
	)

	columnNames := append(queryConfig.FeaturePropertiesNames, queryConfig.GeometryColumnName)
	subQuery := psql.Select(columnNames...).
		From(queryConfig.TableName).
		Where(squirrel.GtOrEq{queryConfig.OrderByColumnName: max(queryParams.StartAtValue, MIN_FID)})

	if queryParams.SqlFilterParams.FilterColName != "" {
		subQuery = subQuery.
			Where(
				squirrel.Eq{
					queryParams.SqlFilterParams.FilterColName: queryParams.SqlFilterParams.FilterVal,
				})
	}

	subQuery = subQuery.OrderBy(fmt.Sprintf("%s ASC", queryConfig.OrderByColumnName)).
		Limit(uint64(queryParams.LimitValue))

	mainQuery := psql.Select(featureCollectionSqlExpression).FromSelect(subQuery, "sub")

	sql, args, err := mainQuery.ToSql()
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return "", nil, fmt.Errorf("failed to build sql query: %w", err)
	}
	return sql, args, nil
}

func getNextFidSqlQuery(tableName string, orderByColumnName string, startAt int, pageSize int, filterParams SqlFilterParams) (string, []interface{}, error) {
	query := psql.Select(orderByColumnName).
		From(tableName).
		Where(squirrel.GtOrEq{orderByColumnName: startAt})

	if filterParams.FilterColName != "" {
		query = query.Where(squirrel.Eq{filterParams.FilterColName: filterParams.FilterVal})
	}

	query = query.OrderBy(fmt.Sprintf("%s ASC", orderByColumnName)).
		Limit(uint64(pageSize)).
		Offset(uint64(pageSize))

	sql, args, err := query.ToSql()
	if err != nil {
		logger.Error("failed_to_get_next_fid_param", err)
		return "", nil, fmt.Errorf("failed to build next page check sql query: %w", err)
	}
	return sql, args, nil
}

func getAccessTokenCreatedAtSqlQuery(token string) (string, []interface{}, error) {
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
