package adm

import (
	"encoding/json"
	"fmt"
	"gadm-api/logger"

	geojson "github.com/paulmach/go.geojson"
)

func convertAdmsToFeatureCollection(adms []Adm) (*geojson.FeatureCollection, error) {
	fc := geojson.NewFeatureCollection()
	for _, adm := range adms {

		feature, err := convertAdmsToGeojson(adm)
		if err != nil {
			logger.Error("failed_to_convert_adm_to_geojson: adm_id=%s: %w", adm.ID, err)
			continue
		}
		fc.AddFeature(feature)
	}
	return fc, nil
}

func convertAdmsToGeojson(adm Adm) (*geojson.Feature, error) {
	if len(adm.Geom) == 0 {
		logger.Warning("missing_geometry_for_adm_fc_conversion: adm_id=%s", adm.ID)
		return nil, fmt.Errorf("missing_geometry_for_adm_fc_conversion: adm_id=%s", adm.ID)
	}

	geom, err := geojson.UnmarshalGeometry(adm.Geom)
	if err != nil {
		return nil, fmt.Errorf("failed_to_unmarshal_geometry: adm_id=%s: %w", adm.ID, err)
	}

	feature := geojson.NewFeature(geom)
	feature.ID = adm.ID
	feature.BoundingBox = adm.Bbox

	if len(adm.Metadata) > 0 {
		if err := json.Unmarshal(adm.Metadata, &feature.Properties); err != nil {
			return nil, fmt.Errorf("failed_to_unmarshal_metadata: adm_id=%s: %w", adm.ID, err)
		}
	}
	feature.SetProperty("lv", adm.Level)
	feature.SetProperty("geom_hash", adm.GeomHash)
	return feature, nil
}
