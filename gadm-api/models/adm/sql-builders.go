package adm

import (
	"fmt"
	"gadm-api/utils"
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

func getSelectAdmsSqlQuery(startAfterId string, batchSize int) (string, []interface{}, error) {
	query := psql.
		Select("adm.metadata", "adm.id", "adm.lv", "adm.geom_hash").
		From("adm")

	if startAfterId != "" {
		query = query.Where("adm.id > $1", startAfterId)
	}

	query = query.OrderBy("adm.id").OrderBy("adm.id").Limit(uint64(batchSize))

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
