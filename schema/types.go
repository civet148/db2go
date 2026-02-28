package schema

import (
	"encoding/json"
	"strings"

	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
)

const (
	TAG_NAME_SQLCA      = "sqlca"
	TAG_VALUE_IS_NULL   = "isnull"
	TAG_VALUE_READ_ONLY = "readonly"
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
	DB_COLUMN_TYPE_BIT:       "uint8",
	DB_COLUMN_TYPE_DECIMAL:   "float64",
	DB_COLUMN_TYPE_REAL:      "float32",
	DB_COLUMN_TYPE_DOUBLE:    "float64",
	DB_COLUMN_TYPE_FLOAT:     "float32",
	DB_COLUMN_TYPE_NUMERIC:   "float64",
}

// 数据库字段类型与go语言类型对照表
var db2goTypes = map[string]string{

	DB_COLUMN_TYPE_BIGINT:     "int64",
	DB_COLUMN_TYPE_INT:        "int32",
	DB_COLUMN_TYPE_INTEGER:    "int32",
	DB_COLUMN_TYPE_MEDIUMINT:  "int32",
	DB_COLUMN_TYPE_SMALLINT:   "int16",
	DB_COLUMN_TYPE_TINYINT:    "int8",
	DB_COLUMN_TYPE_BIT:        "int8",
	DB_COLUMN_TYPE_DECIMAL:    "float64",
	DB_COLUMN_TYPE_REAL:       "float64",
	DB_COLUMN_TYPE_DOUBLE:     "float64",
	DB_COLUMN_TYPE_FLOAT:      "float64",
	DB_COLUMN_TYPE_NUMERIC:    "float64",
	DB_COLUMN_TYPE_BOOL:       "bool",
	DB_COLUMN_TYPE_BOOLEAN:    "bool",
	DB_COLUMN_TYPE_DATETIME:   "time.Time",
	DB_COLUMN_TYPE_YEAR:       "string",
	DB_COLUMN_TYPE_DATE:       "string",
	DB_COLUMN_TYPE_TIME:       "string",
	DB_COLUMN_TYPE_TIMESTAMP:  "time.Time",
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
	DB_COLUMN_TYPE_JSON:       "struct{}",
	DB_COLUMN_TYPE_JSONB:      "string",
	DB_COLUMN_TYPE_POINT:      "sqlca.Point", //暂定
	DB_COLUMN_TYPE_POLYGON:    "string",      //暂定
}

// 数据库字段类型与protobuf类型对照表
var db2protoTypes = map[string]string{

	DB_COLUMN_TYPE_BIGINT:     "sint64",
	DB_COLUMN_TYPE_INT:        "sint32",
	DB_COLUMN_TYPE_INTEGER:    "sint32",
	DB_COLUMN_TYPE_MEDIUMINT:  "sint32",
	DB_COLUMN_TYPE_SMALLINT:   "sint32",
	DB_COLUMN_TYPE_TINYINT:    "sint32",
	DB_COLUMN_TYPE_BIT:        "sint32",
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

// 数据库字段类型与protobuf类型对照表(无符号)
var db2protoTypesUnsigned = map[string]string{
	DB_COLUMN_TYPE_BIGINT:    "uint64",
	DB_COLUMN_TYPE_INT:       "uint32",
	DB_COLUMN_TYPE_INTEGER:   "uint32",
	DB_COLUMN_TYPE_MEDIUMINT: "uint32",
	DB_COLUMN_TYPE_SMALLINT:  "uint32",
	DB_COLUMN_TYPE_TINYINT:   "uint32",
	DB_COLUMN_TYPE_DECIMAL:   "double",
	DB_COLUMN_TYPE_REAL:      "float",
	DB_COLUMN_TYPE_DOUBLE:    "double",
	DB_COLUMN_TYPE_FLOAT:     "float",
	DB_COLUMN_TYPE_NUMERIC:   "double",
}

const (
	SCHEME_MYSQL             = "mysql"
	SCHEME_POSTGRES          = "postgres"
	SCHEME_MSSQL             = "mssql"
	SCHEME_OPEN_GAUSS        = "opengauss"
	JSON_PROPERTY_OMIT_EMTPY = "omitempty"
	DAO_SUFFIX               = "dao"
	TABLE_ALL                = "__all_tables__"
)

const (
	IMPORT_GOGO_PROTO = `import "github.com/gogo/protobuf/gogoproto/gogo.proto";`
	SQLCA_V2_PKG      = `github.com/civet148/sqlca/v2`
	SQLCA_V3_PKG      = `github.com/civet148/sqlca/v3`
	IMPORT_SQLCA_V3   = `import "github.com/civet148/sqlca/v3"`
	IMPORT_SQLCA_V2   = `import "github.com/civet148/sqlca/v2"`
)

const (
	JSON_STYLE_DEFAULT         = "default"
	JSON_STYLE_SMALL_CAMELCASE = "smallcamel"
	JSON_STYLE_BIG_CAMELCASE   = "bigcamel"
)

const (
	HookType_Sqlca = "sqlca"
	HookType_Gorm  = "gorm"
)

type SpecType struct {
	Table   string            `json:"table"`
	Column  string            `json:"column"`
	Type    string            `json:"type"`
	IsPtr   bool              `json:"is_ptr"`
	Package map[string]string `json:"package"`
}

type CommTagType struct {
	Table    string `json:"table"`
	Column   string `json:"column"`
	TagName  string `json:"tag_name"`
	TagValue string `json:"tag_value"`
}

type BaseModel struct {
	Columns []string          `json:"columns"`
	Type    string            `json:"type"`
	Package map[string]string `json:"package"`
}

type CmdFlags struct {
	ConnUrl        string
	Database       string
	Tables         []string
	Without        []string
	ReadOnly       []string
	ExtraTags      []string
	ExcludeTables  []string
	Scheme         string
	Host           string
	User           string
	Password       string
	Charset        string
	OutDir         string
	Prefix         string
	Suffix         string
	PackageName    string
	Protobuf       bool
	EnableDecimal  bool
	OneFile        bool
	GogoOptions    []string
	DAO            string
	ImportModels   string
	OmitEmpty      bool
	TinyintAsBool  []string
	Engine         *sqlca.Engine
	JsonProperties string
	JsonStyle      string
	SSH            string
	SpecTypes      []*SpecType
	ImportVer      string
	SqlcaPkg       string
	Debug          bool
	ExportDDL      string
	TagTypes       []*CommTagType
	ProtoOptions   map[string]string
	FieldStyle     FieldStyle
	BaseModel      *BaseModel
}

func NewCmdFlags() *CmdFlags {
	return &CmdFlags{
		ProtoOptions: make(map[string]string),
	}
}

func (c *CmdFlags) String() string {
	data, _ := json.Marshal(c)
	return string(data)
}

func (c *CmdFlags) GoString() string {
	return c.String()
}

func (c *CmdFlags) IsBaseColumn(strColumnName string) bool {
	if c.BaseModel == nil || len(c.BaseModel.Columns) == 0 {
		return false
	}
	for _, v := range c.BaseModel.Columns {
		if v == strColumnName {
			return true
		}
	}
	return false
}

func (c *CmdFlags) ParseSpecTypes(strSpecType string) {
	var sts []*SpecType
	if strSpecType == "" {
		return
	}
	ss := strings.Split(strSpecType, ",")
	for _, v := range ss {
		v = strings.TrimSpace(v)
		tt := strings.Split(v, "=")
		if len(tt) != 2 {
			log.Errorf("spec type [%s] format illegal", v)
			continue
		}
		var strTableName, strColumnName string
		tcs := strings.Split(tt[0], ".")
		if len(tcs) == 0 {
			continue
		}
		if len(tcs) == 1 {
			strTableName = TABLE_ALL
			strColumnName = tcs[0]
		} else {
			strTableName = tcs[0]
			strColumnName = tcs[1]
		}
		strSpecType = tt[1]
		idx := strings.LastIndex(strSpecType, ".")

		var isPtr bool
		var pack = make(map[string]string)
		var strPackage, strAlias string
		if idx > 0 {
			strPackage = strSpecType[:idx]
			strSpecType = strSpecType[idx+1:]
			if strings.HasPrefix(strPackage, "*") {
				strPackage = strPackage[1:]
				isPtr = true
			}
			strAlias = strings.ReplaceAll(strPackage, "/", "_")
			strAlias = strings.ReplaceAll(strAlias, "-", "_")
			strAlias = strings.ReplaceAll(strAlias, ".", "_")
			pack[strAlias] = strPackage
		}

		sts = append(sts, &SpecType{
			Table:   strTableName,
			Column:  strColumnName,
			Type:    strSpecType,
			IsPtr:   isPtr,
			Package: pack,
		})
	}
	c.SpecTypes = sts
	return
}

func (c *CmdFlags) ParseBaseModel(strBaseModel string) {
	strBaseModel = strings.TrimSpace(strBaseModel)
	if strBaseModel == "" {
		return
	}
	tt := strings.Split(strBaseModel, "=")
	if len(tt) != 2 {
		log.Warnf("base model [%s] format illegal", strBaseModel)
		return
	}

	c.BaseModel = &BaseModel{}

	var strSpecType string
	var strColumns string
	var strPackage, strAlias string
	strSpecType = tt[0]
	strColumns = tt[1]

	var columns = strings.Split(strColumns, ",")
	for i, col := range columns {
		columns[i] = strings.TrimSpace(col)
	}

	idx := strings.LastIndex(strSpecType, ".")
	if idx > 0 {
		strPackage = strSpecType[:idx]
		strSpecType = strSpecType[idx+1:]
		strAlias = strings.ReplaceAll(strPackage, "/", "_")
		strAlias = strings.ReplaceAll(strAlias, "-", "_")
		strAlias = strings.ReplaceAll(strAlias, ".", "_")
		c.BaseModel.Package = map[string]string{
			strAlias: strPackage,
		}
	}
	c.BaseModel.Columns = columns
	c.BaseModel.Type = strSpecType
	return
}

func (c *CmdFlags) GetJsonPropertiesSlice() (jsonProps []string) {
	var dup bool

	if c.JsonProperties != "" {
		jsonProps = strings.Split(c.JsonProperties, ",")
	}

	if c.OmitEmpty {
		for _, v := range jsonProps {
			if v == JSON_PROPERTY_OMIT_EMTPY {
				dup = true
				break
			}
		}
		if !dup {
			jsonProps = append(jsonProps, JSON_PROPERTY_OMIT_EMTPY)
		}
	}
	//log.Debugf("jsonProps [%+v]", jsonProps)
	return
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

type CreateDatabaseDDL struct {
	Database  string `db:"Database"`
	CreateSQL string `db:"Create Database"`
}
