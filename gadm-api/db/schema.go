package db

type GidName string

const (
	GidName0 GidName = "gid_0"
	GidName1 GidName = "gid_1"
	GidName2 GidName = "gid_2"
	GidName3 GidName = "gid_3"
	GidName4 GidName = "gid_4"
	GidName5 GidName = "gid_5"
)

const ADM_TABLE = "adm"

type Adm_Table struct {
	Id       string
	Lv       string
	GeomHash string
	Metadata string
}

var Adm = Adm_Table{
	Id:       "id",
	Lv:       "lv",
	GeomHash: "geom_hash",
	Metadata: "Metadata",
}

const ACCESS_TOKEN_TABLE = "access_tokens"

type AccessTokens_Table struct {
	Id                      string
	Token                   string
	Email                   string
	CreatedAt               string
	UpdatedAt               string
	CanGenerateAccessTokens string
}

var AccessTokensTable = AccessTokens_Table{
	Id:                      "id",
	Token:                   "token",
	Email:                   "email",
	CreatedAt:               "created_at",
	UpdatedAt:               "updated_at",
	CanGenerateAccessTokens: "can_generate_access_tokens",
}
