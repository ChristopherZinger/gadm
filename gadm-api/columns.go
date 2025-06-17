package main

const ADM_0_TABLE = "adm_0"

type Adm0_Table struct {
	FID      string
	GID0     string
	Country  string
	Geometry string
}

var Adm0 = Adm0_Table{
	FID:      "fid",
	GID0:     "gid_0",
	Country:  "country",
	Geometry: "geom",
}
