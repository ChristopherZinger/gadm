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

const ADM_1_TABLE = "adm_1"

type Adm1_Table struct {
	FID      string
	GID0     string
	Country  string
	GID1     string
	Name1    string
	Varname1 string
	NlName1  string
	Type1    string
	Engtype1 string
	Cc1      string
	Hasc1    string
	Iso1     string
	Geometry string
}

var Adm1 = Adm1_Table{
	FID:      "fid",
	GID0:     "gid_0",
	Country:  "country",
	GID1:     "gid_1",
	Name1:    "name_1",
	Varname1: "varname_1",
	NlName1:  "nl_name_1",
	Type1:    "type_1",
	Engtype1: "engtype_1",
	Cc1:      "cc_1",
	Hasc1:    "hasc_1",
	Iso1:     "iso_1",
	Geometry: "geom",
}
