package main

import (
	"fmt"
	"gadm-api/logger"

	"github.com/Masterminds/squirrel"
)

type GeoJsonFeatureSqlQueryParams struct {
	TableName              string
	FeaturePropertiesNames []string
	GeometryColumnName     string
	OrderByColumnName      string // This has to be a integer column!
	StartAtValue           int
	LimitValue             int
}

type FeatureCollectionQueryParams struct {
	TableName              string
	FeaturePropertiesNames []string
	GeometryColumnName     string
	OrderByColumnName      string // This has to be a integer column!
	StartAtValue           int
	LimitValue             int
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

func buildGeojsonFeatureSqlExpression(params struct {
	FeaturePropertiesNames []string
	GeometryColumnName     string
}) string {
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

func buildGeojsonFeatureCollectionSqlExpression(params struct {
	FeaturePropertiesNames []string
	GeometryColumnName     string
}) string {
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
	params GeoJsonFeatureSqlQueryParams,
) (string, []interface{}, error) {
	geojsonFeatureExpression := buildGeojsonFeatureSqlExpression(
		struct {
			FeaturePropertiesNames []string
			GeometryColumnName     string
		}{params.FeaturePropertiesNames, params.GeometryColumnName},
	)

	query := squirrel.Select(geojsonFeatureExpression).
		From(params.TableName).
		Where(squirrel.Expr(fmt.Sprintf("%s >= $1", params.OrderByColumnName), params.StartAtValue)).
		OrderBy(fmt.Sprintf("%s ASC", params.OrderByColumnName)).
		Limit(uint64(params.LimitValue))

	sql, args, err := query.ToSql()
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return "", nil, fmt.Errorf("failed to build sql query: %w", err)
	}
	return sql, args, nil
}

func buildFeatureCollectionSqlQuery(params FeatureCollectionQueryParams) (string, []interface{}, error) {
	featureCollectionSqlExpression := buildGeojsonFeatureCollectionSqlExpression(
		struct {
			FeaturePropertiesNames []string
			GeometryColumnName     string
		}{params.FeaturePropertiesNames, Adm0.Geometry},
	)

	columnNames := append(params.FeaturePropertiesNames, params.GeometryColumnName)
	subQuery := squirrel.Select(columnNames...).
		From(params.TableName).
		Where(squirrel.Expr(fmt.Sprintf("%s >= $1", params.OrderByColumnName), max(params.StartAtValue, MIN_FID))).
		OrderBy(fmt.Sprintf("%s ASC", params.OrderByColumnName)).
		Limit(uint64(params.LimitValue))

	mainQuery := squirrel.Select(featureCollectionSqlExpression).FromSelect(subQuery, "sub")

	sql, args, err := mainQuery.ToSql()
	if err != nil {
		logger.Error("failed_to_build_sql_query %v", err)
		return "", nil, fmt.Errorf("failed to build sql query: %w", err)
	}
	return sql, args, nil
}

func getNextFidSqlQuery(tableName string, orderByColumnName string, startAt int, pageSize int) (string, []interface{}, error) {
	query := squirrel.Select(orderByColumnName).
		From(tableName).
		Where(squirrel.Expr(fmt.Sprintf("%s >= $1", orderByColumnName), startAt)).
		OrderBy(fmt.Sprintf("%s ASC", orderByColumnName)).
		Limit(uint64(pageSize)).
		Offset(uint64(pageSize))

	sql, args, err := query.ToSql()
	if err != nil {
		logger.Error("failed_to_get_next_fid_param", err)
		return "", nil, fmt.Errorf("failed to build next page check sql query: %w", err)
	}
	return sql, args, nil
}
