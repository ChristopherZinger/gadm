package main

import (
	"testing"
)

func TestBuildJsonObjectExpression(t *testing.T) {
	columns := []string{"fid", "gid0", "country"}
	result := buildGeojsonFeaturePropertiesSqlExpression(columns...)

	expected := "json_build_object('fid', fid, 'gid0', gid0, 'country', country)"

	if result != expected {
		t.Errorf("buildJsonObjectExpression() = %v, want %v", result, expected)
	}
}
