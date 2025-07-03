package main

import (
	"fmt"
	"gadm-api/logger"

	"github.com/Masterminds/squirrel"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type SqlQueryParams struct {
	StartAtValue int
	LimitValue   int
	FilterValue  string
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

	query := squirrel.Select(geojsonFeatureExpression).
		From(handlerConfig.TableName).
		Where(squirrel.GtOrEq{handlerConfig.OrderByColumnName: max(queryParams.StartAtValue, MIN_FID)}).
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
		Where(squirrel.GtOrEq{queryConfig.OrderByColumnName: max(queryParams.StartAtValue, MIN_FID)}).
		OrderBy(fmt.Sprintf("%s ASC", queryConfig.OrderByColumnName)).
		Limit(uint64(queryParams.LimitValue))

	mainQuery := psql.Select(featureCollectionSqlExpression).FromSelect(subQuery, "sub")

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
		Where(squirrel.GtOrEq{orderByColumnName: startAt}).
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
