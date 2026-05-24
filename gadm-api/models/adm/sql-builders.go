package adm

import "gadm-api/utils"

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
