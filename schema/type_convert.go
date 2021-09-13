package schema

import (
	"github.com/civet148/log"
	"strings"
)

const (
	POSTGRES_COLUMN_INTEGER           = "integer"
	POSTGRES_COLUMN_BIT               = "bit"
	POSTGRES_COLUMN_BOOLEAN           = "boolean"
	POSTGRES_COLUMN_BOX               = "box"
	POSTGRES_COLUMN_BYTEA             = "bytea"
	POSTGRES_COLUMN_CHARACTER         = "character"
	POSTGRES_COLUMN_CIDR              = "cidr"
	POSTGRES_COLUMN_CIRCLE            = "circle"
	POSTGRES_COLUMN_DATE              = "date"
	POSTGRES_COLUMN_NUMERIC           = "numeric"
	POSTGRES_COLUMN_REAL              = "real"
	POSTGRES_COLUMN_DOUBLE            = "double"
	POSTGRES_COLUMN_INET              = "inet"
	POSTGRES_COLUMN_SMALLINT          = "smallint"
	POSTGRES_COLUMN_BIGINT            = "bigint"
	POSTGRES_COLUMN_INTERVAL          = "interval"
	POSTGRES_COLUMN_JSON              = "json"
	POSTGRES_COLUMN_JSONB             = "jsonb"
	POSTGRES_COLUMN_LINE              = "line"
	POSTGRES_COLUMN_LSEG              = "lseg"
	POSTGRES_COLUMN_MACADDR           = "macaddr"
	POSTGRES_COLUMN_MONEY             = "money"
	POSTGRES_COLUMN_PATH              = "path"
	POSTGRES_COLUMN_POINT             = "point"
	POSTGRES_COLUMN_POLYGON           = "polygon"
	POSTGRES_COLUMN_TEXT              = "text"
	POSTGRES_COLUMN_TIME              = "time"
	POSTGRES_COLUMN_TIMESTAMP         = "timestamp"
	POSTGRES_COLUMN_TSQUERY           = "tsquery"
	POSTGRES_COLUMN_TSVECTOR          = "tsvector"
	POSTGRES_COLUMN_TXID_SNAPSHOT     = "txid_snapshot"
	POSTGRES_COLUMN_UUID              = "uuid"
	POSTGRES_COLUMN_BIT_VARYING       = "bit varying"
	POSTGRES_COLUMN_CHARACTER_VARYING = "character varying"
	POSTGRES_COLUMN_XML               = "xml"
)

const (
	DB_COLUMN_TYPE_BIGINT     = "bigint"
	DB_COLUMN_TYPE_INT        = "int"
	DB_COLUMN_TYPE_INTEGER    = "integer"
	DB_COLUMN_TYPE_MEDIUMINT  = "mediumint"
	DB_COLUMN_TYPE_SMALLINT   = "smallint"
	DB_COLUMN_TYPE_TINYINT    = "tinyint"
	DB_COLUMN_TYPE_BIT        = "bit"
	DB_COLUMN_TYPE_BOOL       = "bool"
	DB_COLUMN_TYPE_BOOLEAN    = "boolean"
	DB_COLUMN_TYPE_DECIMAL    = "decimal"
	DB_COLUMN_TYPE_REAL       = "real"
	DB_COLUMN_TYPE_DOUBLE     = "double"
	DB_COLUMN_TYPE_FLOAT      = "float"
	DB_COLUMN_TYPE_NUMERIC    = "numeric"
	DB_COLUMN_TYPE_DATETIME   = "datetime"
	DB_COLUMN_TYPE_YEAR       = "year"
	DB_COLUMN_TYPE_DATE       = "date"
	DB_COLUMN_TYPE_TIME       = "time"
	DB_COLUMN_TYPE_TIMESTAMP  = "timestamp"
	DB_COLUMN_TYPE_ENUM       = "enum"
	DB_COLUMN_TYPE_SET        = "set"
	DB_COLUMN_TYPE_VARCHAR    = "varchar"
	DB_COLUMN_TYPE_NVARCHAR   = "nvarchar"
	DB_COLUMN_TYPE_CHAR       = "char"
	DB_COLUMN_TYPE_TEXT       = "text"
	DB_COLUMN_TYPE_TINYTEXT   = "tinytext"
	DB_COLUMN_TYPE_MEDIUMTEXT = "mediumtext"
	DB_COLUMN_TYPE_LONGTEXT   = "longtext"
	DB_COLUMN_TYPE_BLOB       = "blob"
	DB_COLUMN_TYPE_TINYBLOB   = "tinyblob"
	DB_COLUMN_TYPE_MEDIUMBLOB = "mediumblob"
	DB_COLUMN_TYPE_LONGBLOB   = "longblob"
	DB_COLUMN_TYPE_BINARY     = "binary"
	DB_COLUMN_TYPE_VARBINARY  = "varbinary"
	DB_COLUMN_TYPE_JSON       = "json"
	DB_COLUMN_TYPE_JSONB      = "jsonb"
	DB_COLUMN_TYPE_POINT      = "point"
	DB_COLUMN_TYPE_POLYGON    = "polygon"
)

var db2goTypesUnsigned = map[string]string{
	DB_COLUMN_TYPE_BIGINT:    "uint64",
	DB_COLUMN_TYPE_INT:       "uint32",
	DB_COLUMN_TYPE_INTEGER:   "uint32",
	DB_COLUMN_TYPE_MEDIUMINT: "uint32",
	DB_COLUMN_TYPE_SMALLINT:  "uint16",
	DB_COLUMN_TYPE_TINYINT:   "uint8",
}

//数据库字段类型与go语言类型对照表
var db2goTypes = map[string]string{

	DB_COLUMN_TYPE_BIGINT:     "int64",
	DB_COLUMN_TYPE_INT:        "int32",
	DB_COLUMN_TYPE_INTEGER:    "int32",
	DB_COLUMN_TYPE_MEDIUMINT:  "int32",
	DB_COLUMN_TYPE_SMALLINT:   "int16",
	DB_COLUMN_TYPE_TINYINT:    "int8",
	DB_COLUMN_TYPE_BIT:        "int8",
	DB_COLUMN_TYPE_BOOL:       "bool",
	DB_COLUMN_TYPE_BOOLEAN:    "bool",
	DB_COLUMN_TYPE_DECIMAL:    "float64",
	DB_COLUMN_TYPE_REAL:       "float64",
	DB_COLUMN_TYPE_DOUBLE:     "float64",
	DB_COLUMN_TYPE_FLOAT:      "float64",
	DB_COLUMN_TYPE_NUMERIC:    "float64",
	DB_COLUMN_TYPE_DATETIME:   "string",
	DB_COLUMN_TYPE_YEAR:       "string",
	DB_COLUMN_TYPE_DATE:       "string",
	DB_COLUMN_TYPE_TIME:       "string",
	DB_COLUMN_TYPE_TIMESTAMP:  "string",
	DB_COLUMN_TYPE_ENUM:       "string",
	DB_COLUMN_TYPE_SET:        "string",
	DB_COLUMN_TYPE_VARCHAR:    "string",
	DB_COLUMN_TYPE_NVARCHAR:   "string",
	DB_COLUMN_TYPE_CHAR:       "string",
	DB_COLUMN_TYPE_TEXT:       "string",
	DB_COLUMN_TYPE_TINYTEXT:   "string",
	DB_COLUMN_TYPE_MEDIUMTEXT: "string",
	DB_COLUMN_TYPE_LONGTEXT:   "string",
	DB_COLUMN_TYPE_BLOB:       "string",
	DB_COLUMN_TYPE_TINYBLOB:   "string",
	DB_COLUMN_TYPE_MEDIUMBLOB: "string",
	DB_COLUMN_TYPE_LONGBLOB:   "string",
	DB_COLUMN_TYPE_BINARY:     "string",
	DB_COLUMN_TYPE_VARBINARY:  "string",
	DB_COLUMN_TYPE_JSON:       "string",
	DB_COLUMN_TYPE_JSONB:      "string",
	DB_COLUMN_TYPE_POINT:      "string", //暂定
	DB_COLUMN_TYPE_POLYGON:    "string", //暂定
}

//数据库字段类型与protobuf类型对照表
var db2protoTypes = map[string]string{

	DB_COLUMN_TYPE_BIGINT:     "int64",
	DB_COLUMN_TYPE_INT:        "int32",
	DB_COLUMN_TYPE_INTEGER:    "int32",
	DB_COLUMN_TYPE_MEDIUMINT:  "int32",
	DB_COLUMN_TYPE_SMALLINT:   "int32",
	DB_COLUMN_TYPE_TINYINT:    "int32",
	DB_COLUMN_TYPE_BIT:        "int32",
	DB_COLUMN_TYPE_BOOL:       "bool",
	DB_COLUMN_TYPE_BOOLEAN:    "bool",
	DB_COLUMN_TYPE_DECIMAL:    "double",
	DB_COLUMN_TYPE_REAL:       "float",
	DB_COLUMN_TYPE_DOUBLE:     "double",
	DB_COLUMN_TYPE_FLOAT:      "float",
	DB_COLUMN_TYPE_NUMERIC:    "double",
	DB_COLUMN_TYPE_DATETIME:   "string",
	DB_COLUMN_TYPE_YEAR:       "string",
	DB_COLUMN_TYPE_DATE:       "string",
	DB_COLUMN_TYPE_TIME:       "string",
	DB_COLUMN_TYPE_TIMESTAMP:  "string",
	DB_COLUMN_TYPE_ENUM:       "string",
	DB_COLUMN_TYPE_SET:        "string",
	DB_COLUMN_TYPE_VARCHAR:    "string",
	DB_COLUMN_TYPE_CHAR:       "string",
	DB_COLUMN_TYPE_TEXT:       "string",
	DB_COLUMN_TYPE_TINYTEXT:   "string",
	DB_COLUMN_TYPE_MEDIUMTEXT: "string",
	DB_COLUMN_TYPE_LONGTEXT:   "string",
	DB_COLUMN_TYPE_BLOB:       "string",
	DB_COLUMN_TYPE_TINYBLOB:   "string",
	DB_COLUMN_TYPE_MEDIUMBLOB: "string",
	DB_COLUMN_TYPE_LONGBLOB:   "string",
	DB_COLUMN_TYPE_BINARY:     "string",
	DB_COLUMN_TYPE_VARBINARY:  "string",
	DB_COLUMN_TYPE_JSON:       "string",
	DB_COLUMN_TYPE_JSONB:      "string",
	DB_COLUMN_TYPE_POINT:      "string", //暂定
	DB_COLUMN_TYPE_POLYGON:    "string", //暂定
}

func ConvertPostgresColumnType(table *TableSchema) (err error) {

	for i, v := range table.Columns {
		if _, ok := db2goTypes[v.DataType]; !ok {
			convertPostgresType(&table.Columns[i])
			log.Infof("postgres column [%v] data type [%v] converted to [%v]", v.Name, v.DataType, table.Columns[i].DataType)
		}
	}

	return
}

func ConvertMssqlColumnType(table *TableSchema) (err error) {
	return
}

func convertPostgresType(column *TableColumn) {
	column.DataType = getFamiliarType(column.DataType)
}

func getFamiliarType(strDataType string) (strType string) {

	if strings.Contains(strDataType, POSTGRES_COLUMN_BIT_VARYING) || strings.Contains(strDataType, POSTGRES_COLUMN_BIT) {
		return DB_COLUMN_TYPE_BIT
	} else if strings.Contains(strDataType, POSTGRES_COLUMN_BOX) || strings.Contains(strDataType, POSTGRES_COLUMN_CIRCLE) {
		return DB_COLUMN_TYPE_POLYGON
	} else if strings.Contains(strDataType, POSTGRES_COLUMN_MONEY) || strings.Contains(strDataType, POSTGRES_COLUMN_NUMERIC) {
		return DB_COLUMN_TYPE_DECIMAL
	}
	return DB_COLUMN_TYPE_TEXT
}
