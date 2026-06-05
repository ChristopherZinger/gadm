package adm

import (
	"fmt"
	"gadm-api/logger"
	"gadm-api/utils"

	"github.com/Masterminds/squirrel"
)

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

func getAdmForPointSqlQuery(point utils.Point) (string, []interface{}, error) {
	withClause := `
		WITH input_point AS (
			SELECT ST_SetSRID(ST_MakePoint(?, ?), 4326)::geometry(Point,4326) AS pt
		),
		candidates AS (
			SELECT g.geom_hash, g.geom, g.area_sq_m, ip.pt
			FROM gadm.adm_geometries AS g
			INNER JOIN input_point ip
			ON ST_Contains(g.bbox, ip.pt)
			ORDER BY area_sq_m ASC
		),
		result_geometry AS (
			SELECT c.geom_hash
			FROM candidates AS c
			INNER JOIN input_point ip
			ON ST_Contains(c.geom, ip.pt)
			ORDER BY c.area_sq_m ASC
			LIMIT 1
		)`

	query := psql.
		Select("adm.metadata", "adm.id", "adm.lv", "adm.geom_hash").
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

func getSelectOneAdmByIdSqlQuery(admId string) (string, []interface{}, error) {
	query := psql.
		Select("adm.metadata", "adm.id", "adm.lv", "adm.geom_hash").
		From("adm").
		Where("adm.id = $1", admId).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func getSelectAdmsSqlQuery(options admQueryOpts) (string, []interface{}, error) {
	fields := []string{"adm.metadata", "adm.id", "adm.lv", "adm.geom_hash"}
	if options.includeGeometry {
		fields = append(
			fields,
			"ST_AsGeoJSON(g.geom, 6) as geom",
			"ST_AsGeoJSON(g.bbox, 6) as bbox",
		)
	}
	query := psql.
		Select(fields...).
		From("adm")

	if options.startAfterId != nil {
		query = query.Where("adm.id > $1", *options.startAfterId)
		query = query.OrderBy("adm.id")
	}

	if options.lv != nil {
		if *options.lv < 0 || *options.lv > 5 {
			logger.Error("invalid_lv_when_building_adm_query: %d", *options.lv)
		} else {
			query = query.Where("adm.lv = $1", *options.lv)
		}
	}

	if options.includeGeometry {
		query = query.Join("adm.adm_geometries g on adm.geom_hash = g.geom_hash")
	}

	if options.startAfterFid != nil {
		query = query.Where(squirrel.Gt{"adm.metadata ->> 'fid'": *options.startAfterFid})
		query = query.OrderBy("adm.metadata ->> 'fid'")
	}

	query = query.Limit(uint64(options.batchSize))

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func getUpsertAdmTreeSqlQuery(parentId string, childIds []string) (string, []interface{}, error) {
	if len(childIds) == 0 {
		return "", nil, nil
	}

	query := psql.
		Insert("gadm.adm_tree").
		Columns("parent", "child").
		Suffix("ON CONFLICT (parent, child) DO NOTHING")

	for _, childId := range childIds {
		query = query.Values(parentId, childId)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func getSelectAdmDirectChildrenForIdSqlQuery(admId string, lv int) (string, []interface{}, error) {
	query := psql.
		Select("adm.metadata", "adm.id", "adm.lv", "adm.geom_hash").
		From("adm").
		Where(fmt.Sprintf("adm.metadata->>'gid_%d' = $1", lv), admId).
		Where("adm.lv = $2", lv+1)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func getNeighborsSqlQuery(admId string) (string, []interface{}, error) {
	withClause := `
		WITH seed AS (
			SELECT g.geom, g.bbox
			FROM gadm.adm a
			JOIN gadm.adm_geometries g ON a.geom_hash = g.geom_hash
			WHERE a.id = ?
		)`

	query := psql.
		Select("adm.metadata", "adm.id", "adm.lv", "adm.geom_hash").
		Prefix(withClause, admId).
		From("seed").
		InnerJoin("gadm.adm_geometries cg ON cg.bbox && seed.bbox AND ST_Touches(cg.geom, seed.geom)").
		InnerJoin("gadm.adm ON adm.geom_hash = cg.geom_hash").
		Where("adm.id > ?", admId).
		Where("NOT EXISTS (SELECT 1 FROM gadm.adm_tree t WHERE t.parent = adm.id)")

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func getSelectLeafAdmsSqlQuery(startAfterId string, batchSize int) (string, []interface{}, error) {
	query := psql.
		Select("adm.metadata", "adm.id", "adm.lv", "adm.geom_hash").
		From("adm").
		LeftJoin("adm_tree atree ON adm.id = atree.parent").
		Where("atree.parent IS NULL")

	if startAfterId != "" {
		query = query.Where("adm.id > ?", startAfterId)
	}

	query = query.OrderBy("adm.id ASC").Limit(uint64(batchSize))

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func getUpsertAdmNeighborsSqlQuery(n1Id, n2Id string) (string, []interface{}, error) {
	query := psql.
		Insert("gadm.adm_neighbors").
		Columns("n1", "n2").
		Values(
			squirrel.Expr("LEAST(?::uuid, ?::uuid)", n1Id, n2Id),
			squirrel.Expr("GREATEST(?::uuid, ?::uuid)", n1Id, n2Id),
		).
		Suffix("ON CONFLICT (n1, n2) DO NOTHING")

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}

func getUpsertAdmNeighborsBatchSqlQuery(admId string, neighborIds []string) (string, []interface{}, error) {
	if len(neighborIds) == 0 {
		return "", nil, nil
	}

	query := psql.
		Insert("gadm.adm_neighbors").
		Columns("n1", "n2").
		Suffix("ON CONFLICT (n1, n2) DO NOTHING")

	for _, neighborId := range neighborIds {
		query = query.Values(
			squirrel.Expr("LEAST(?::uuid, ?::uuid)", admId, neighborId),
			squirrel.Expr("GREATEST(?::uuid, ?::uuid)", admId, neighborId),
		)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, nil
}
