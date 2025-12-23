package schema

import (
	"encoding/json"
	"fmt"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
	"os"
	"strings"
)

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

type SpecType struct {
	Table   string            `json:"table"`
	Column  string            `json:"column"`
	Type    string            `json:"type"`
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
	ExportTo       string
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
		var pack = make(map[string]string)
		var strPackage, strAliase string
		if idx > 0 {
			strPackage = strSpecType[:idx]
			strSpecType = strSpecType[idx+1:]
			strAliase = strings.ReplaceAll(strPackage, "/", "_")
			strAliase = strings.ReplaceAll(strAliase, "-", "_")
			strAliase = strings.ReplaceAll(strAliase, ".", "_")
			pack[strAliase] = strPackage
		}

		sts = append(sts, &SpecType{
			Table:   strTableName,
			Column:  strColumnName,
			Type:    strSpecType,
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
	log.Infof("base model [%+v]", c.BaseModel)
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

type TableSchema struct {
	SchemeName         string        `json:"table_schema" db:"table_schema"`   //database name
	TableName          string        `json:"table_name" db:"table_name"`       //table name
	TableEngine        string        `json:"engine" db:"engine"`               //database engine
	TableComment       string        `json:"table_comment" db:"table_comment"` //comment of table schema
	SchemeDir          string        `json:"schema_dir" db:"schema_dir"`       //output path
	PkName             string        `json:"pk_name" db:"pk_name"`             //primary key column name
	StructName         string        `json:"struct_name" db:"struct_name"`     //struct name
	StructDAO          string        `json:"struct_dao" db:"struct_dao"`       //struct DAO name
	OutDir             string        `json:"out_dir" db:"out_dir"`             //output directory
	OutFilePath        string        `json:"file_name" db:"file_name"`         //output file path
	Columns            []TableColumn `json:"table_columns" db:"table_columns"` //columns with database and golang
	TableNameCamelCase string        `json:"-"`                                //table name in camel case
	TableCreateSQL     string        `json:"-"`                                //table create SQL
}

type TableColumn struct {
	Name         string `json:"column_name" db:"column_name"`
	DataType     string `json:"data_type" db:"data_type"`
	ColumnType   string `json:"column_type" db:"column_type"`
	Key          string `json:"column_key" db:"column_key"`
	Extra        string `json:"extra" db:"extra"`
	Comment      string `json:"column_comment" db:"column_comment"`
	IsNullable   string `json:"is_nullable" db:"is_nullable"`
	IsPrimaryKey bool   // is primary key
	IsDecimal    bool   // is decimal type
	IsReadOnly   bool   // is read only
	GoName       string // column name in golang
	GoType       string // column type in golang
}

type Exporter interface {
	ExportGo() (err error)
	ExportProto() (err error)
}

type Instance func(cmd *CmdFlags, e *sqlca.Engine) Exporter

var instances = make(map[string]Instance, 1)

func Register(strScheme string, inst Instance) {
	instances[strScheme] = inst
}

func NewExporter(cmd *CmdFlags, e *sqlca.Engine) Exporter {
	var ok bool
	var inst Instance
	if inst, ok = instances[cmd.Scheme]; !ok {
		log.Errorf("scheme [%v] instance not registered", cmd.Scheme)
		return nil
	}
	return inst(cmd, e)
}

func IsInSlice(in string, s []string) bool {
	for _, v := range s {
		if v == in {
			return true
		}
	}
	return false
}

func MakeTags(cmd *CmdFlags, strColName, strColType, strTagValue, strComment string, strAppends string) string {
	strComment = ReplaceCRLF(strComment)
	var strJsonValue string
	var strJsonProperties = cmd.GetJsonPropertiesSlice()
	strJsonValue = strTagValue
	if cmd.JsonStyle == JSON_STYLE_SMALL_CAMELCASE {
		strJsonValue = SmallCamelCase(strTagValue)
	} else if cmd.JsonStyle == JSON_STYLE_BIG_CAMELCASE {
		strJsonValue = BigCamelCase(strTagValue)
	}
	if len(strJsonProperties) > 0 {
		strJsonValue += fmt.Sprintf(",%s", strings.Join(strJsonProperties, ","))
	}
	return fmt.Sprintf("	%v %v `json:\"%v\" db:\"%v\" %v` //%v \n",
		strColName, strColType, strJsonValue, strTagValue, strAppends, strComment)
}

func ReplaceColumnType(cmd *CmdFlags, strTableName, strColName, strColType string) string {

	for _, st := range cmd.SpecTypes {
		if (st.Table == strTableName || st.Table == TABLE_ALL) && strColName == st.Column {
			if len(st.Package) != 0 {
				for k, _ := range st.Package {
					strColType = fmt.Sprintf("%s.%s", k, st.Type)
				}
			} else {
				strColType = st.Type
			}
		}
	}
	return strColType
}

func MakeGetter(strStructName, strColName, strColType string) (strGetter string) {

	return fmt.Sprintf("func (do *%v) Get%v() %v { return do.%v } \n", strStructName, strColName, strColType, strColName)
}

func MakeSetter(strStructName, strColName, strColType string) (strSetter string) {

	return fmt.Sprintf("func (do *%v) Set%v(v %v) { do.%v = v } \n", strStructName, strColName, strColType, strColName)
}

func ReplaceCRLF(strIn string) (strOut string) {
	strOut = strings.ReplaceAll(strIn, "\r", "")
	strOut = strings.ReplaceAll(strOut, "\n", "")
	return
}

func CreateOutputFile(cmd *CmdFlags, table *TableSchema, strFileSuffix string, append bool) (file *os.File, err error) {

	var strOutDir = cmd.OutDir
	var strPackageName = cmd.PackageName
	var strNamePrefix = cmd.Prefix
	var strNameSuffix = cmd.Suffix

	_, errStat := os.Stat(strOutDir)
	if errStat != nil && os.IsNotExist(errStat) {

		log.Info("mkdir [%v]", strOutDir)
		if err = os.MkdirAll(strOutDir, os.ModePerm); err != nil {
			log.Error("mkdir [%v] error (%v)", strOutDir, err.Error())
			return
		}
	}

	table.OutDir = strOutDir

	if strPackageName == "" {
		//mkdir by output dir + scheme name
		strPackageName = table.SchemeName
		if strings.LastIndex(strOutDir, fmt.Sprintf("%v", os.PathSeparator)) == -1 {
			table.SchemeDir = fmt.Sprintf("%v/%v", strOutDir, strPackageName)
		} else {
			table.SchemeDir = fmt.Sprintf("%v%v", strOutDir, strPackageName)
		}
	} else {
		table.SchemeDir = fmt.Sprintf("%v/%v", strOutDir, strPackageName) //mkdir by package name
	}

	_, errStat = os.Stat(table.SchemeDir)

	if errStat != nil && os.IsNotExist(errStat) {

		log.Info("mkdir [%v]", table.SchemeDir)
		if err = os.MkdirAll(table.SchemeDir, os.ModePerm); err != nil {
			log.Errorf("mkdir path name [%v] error (%v)", table.SchemeDir, err.Error())
			return
		}
	}

	if strNamePrefix != "" {
		strNamePrefix = fmt.Sprintf("%v_", strNamePrefix)
	}
	if strNameSuffix != "" {
		strNameSuffix = fmt.Sprintf("_%v", strNameSuffix)
	}

	var flag = os.O_CREATE | os.O_RDWR | os.O_TRUNC

	if append {
		flag = os.O_CREATE | os.O_RDWR | os.O_APPEND
	}

	if cmd.OneFile { //数据库名称作为文件名
		table.OutFilePath = fmt.Sprintf("%v/%v%v%v.%v", table.SchemeDir, strNamePrefix, table.SchemeName, strNameSuffix, strFileSuffix)
	} else { //数据表名作为文件名
		table.OutFilePath = fmt.Sprintf("%v/%v%v%v.%v", table.SchemeDir, strNamePrefix, table.TableName, strNameSuffix, strFileSuffix)
	}

	file, err = os.OpenFile(table.OutFilePath, flag, os.ModePerm)
	if err != nil {
		log.Errorf("open file [%v] error (%v)", table.OutFilePath, err.Error())
		return
	}
	log.Infof("generate table [%s] protobuf schema to file [%v] successfully", table.TableName, table.OutFilePath)
	return
}

func BigCamelCase(strIn string) (strOut string) {
	var idxUnderLine = int(-1)
	for i, v := range strIn {
		strChr := string(v)
		if i == 0 {
			strOut += strings.ToUpper(strChr)
		} else {
			if v == '_' {
				idxUnderLine = i //ignore
			} else {
				if i == idxUnderLine+1 {
					strOut += strings.ToUpper(strChr)
				} else {
					strOut += strChr
				}
			}
		}
	}
	return strOut
}

func SmallCamelCase(strIn string) (strOut string) {
	var idxUnderLine = int(-1)
	for i, v := range strIn {
		strChr := string(v)
		if i == 0 {
			strOut += strings.ToLower(strChr)
		} else {
			if v == '_' {
				idxUnderLine = i //ignore
			} else {
				if i == idxUnderLine+1 {
					strOut += strings.ToUpper(strChr)
				} else {
					strOut += strChr
				}
			}
		}
	}
	return
}

func Split(s string) (ss []string) {
	if strings.Contains(s, ",") {
		ss = strings.Split(s, ",")
	} else {
		ss = strings.Split(s, ";")
	}
	return ss
}

func TrimSpaceSlice(s []string) (ts []string) {
	for _, v := range s {
		ts = append(ts, strings.TrimSpace(v))
	}
	return
}

func GetDatabaseName(strPath string) (strName string) {
	idx := strings.LastIndex(strPath, "/")
	if idx == -1 {
		return
	}
	return strPath[idx+1:]
}

// 将数据库字段类型转为go语言对应的数据类型
func GetGoColumnType(cmd *CmdFlags, strTableName string, col TableColumn, enableDecimal bool, tinyintAsBool []string) (strGoColType string, isDecimal bool) {

	var bUnsigned bool
	var strColName, strDataType, strColumnType string
	strColName = col.Name
	strDataType = col.DataType
	strColumnType = col.ColumnType

	//log.Debugf("table [%s] column name [%s] type [%s]", strTableName, strColName, strDataType)
	//tinyint type column redeclare as bool
	if len(tinyintAsBool) > 0 && strDataType == DB_COLUMN_TYPE_TINYINT {
		if IsInSlice(strColName, tinyintAsBool) {
			//log.Warnf("table [%s] column [%s] %s redeclare as bool type", strTableName, strColName, strDataType)
			return DB_COLUMN_TYPE_BOOL, false
		}
	}

	if strings.Contains(strColumnType, "unsigned") { //判断字段是否为无符号类型
		bUnsigned = true
	}
	//log.Infof("TableColumn [%+v]", col)
	var ok bool
	if strGoColType, ok = db2goTypes[strDataType]; !ok {
		strGoColType = "string"
		//log.Warnf("table [%v] column [%v] data type [%v] not support yet, set as string type", strTableName, strColName, strDataType)
		return
	}
	if bUnsigned {
		strGoColType = db2goTypesUnsigned[strDataType]
		if strGoColType == "" {
			log.Warnf("data type [%s] column type [%s] have no unsigned type", strDataType, strColumnType)
		}
	}
	switch strDataType {
	case DB_COLUMN_TYPE_DECIMAL:
		if !enableDecimal {
			strGoColType = "float64"
		} else {
			isDecimal = true
			strGoColType = "sqlca.Decimal"
		}
	}

	return ReplaceColumnType(cmd, strTableName, col.Name, strGoColType), isDecimal
}

// 将数据库字段类型转为protobuf对应的数据类型
func GetProtoColumnType(strTableName string, col TableColumn) (strColType string) {
	var ok bool
	var unsigned bool
	var strColName, strDataType, strColumnType string
	strColName = col.Name
	strDataType = col.DataType
	strColumnType = col.ColumnType
	if strings.Contains(strColumnType, "unsigned") { //判断字段是否为无符号类型
		unsigned = true
	}
	if !unsigned {
		if strColType, ok = db2protoTypes[strDataType]; !ok {
			strColType = "string"
			log.Warnf("table [%v] column [%v] data type [%v] not support yet, set as string type", strTableName, strColName, strDataType)
			return
		}
	} else {
		if strColType, ok = db2protoTypesUnsigned[strDataType]; !ok {
			strColType = "string"
			log.Warnf("table [%v] column [%v] data type [%v] not support yet, set as string type", strTableName, strColName, strDataType)
		}
	}
	return strColType
}

func HandleCommentCRLF(table *TableSchema) {
	//write table name in camel case naming
	table.TableComment = ReplaceCRLF(table.TableComment)
	for i, v := range table.Columns {
		table.Columns[i].Comment = ReplaceCRLF(v.Comment)
	}
}
