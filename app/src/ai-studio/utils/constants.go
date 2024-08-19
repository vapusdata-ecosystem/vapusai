package utils

import "time"

const (
	//default port for datamesh server
	DEFAULTPORT          = 8900
	DEFAULTSERVERKEYTLS  = "x509/server_key.pem"
	DEFAULTSERVERCERTTLS = "x509/server_cert.pem"
	DEFAULT_CONFIG_TYPE  = "toml"
	DOT                  = "."
	EMPTYSTR             = ""
	HASHICORPVAULT       = "hashicorpVault"
	AWS_SECRETMANAGER    = "awsSecretManager"
)

// context Contexts

type ctxKeys string

const (
	// ContextKey for the context
	PARSED_CLAIMS ctxKeys = "parsedClaims"
	ACCESS_TOKEN  ctxKeys = "accessToken"
	USERID_CTX    ctxKeys = "userId"
)

// REDIS action and keys that are constant
const (
	LIST              = "list"
	ADD               = "Add"
	EXISTS            = "exists"
	COUNT             = "count"
	DEL               = "del"
	MADD              = "add-mulitple"
	ACCOUNT_KEY       = "account"
	ACCOUNT_DM_CF_KEY = "account:datamesh"
	DM_IDENTIFIER     = "datamesh"
	MESHID            = "meshId"
	DMNID             = "dmnId"
	DNID              = "dnId"
)

const (
	// Vault Secret engines for different resources
	DATAMESH_SE_VAULT      = "vapus-datamesh"
	DataMeshNODES_SE_VAULT = "vapus-DataMeshnodes"
	DATANODES_SE_VAULT     = "vapus-datanodes"
	VPR_SE_VAULT           = "vapus-vpr"
	VDR_SE_VAULT           = "vapus-vdr"
)

// User constants

const (
	DEFAULT_USER_INVITE_EXPIRE_TIME = time.Hour * 24 * 30
	DEFAULT_PLATFORM_AT_VALIDITY    = time.Hour * 24
)

// REDIS action and keys that are constant
const (
	DATAMESH_NODE_MAP        = "datamesh_node_map"
	DOMAINNODES_TABLE        = "domainnodes-table"
	DATANODES_TABLE          = "datanodes-table"
	DATAMESH_TABLE           = "datamesh-table"
	DATANODES_METADATA_TABLE = "datamesh-metadata"
)

// ES Indexes
const (
	DATANODE_METADATA_ES_INDEX = "vapusdata-datasource-metadata"
	DOMAINNODES_INDEX          = "vapusdata-domain"
	DATANODES_INDEX            = "vapusdata-datasources"
	DATAMESH_INDEX             = "vapusdata-datamesh"
	ACCOUNT_INDEX              = "vapusdata-accounts"
)

// DB action and keys that are configurable
var (
	DataSourcePackages   = "schema::datasource::packages::%v"
	DataSourceDatabase   = "schema::datasource::database::%v"
	DATANODE_MEADATA_KEY = "datanode-scheme-%v"
)
