package schema

import (
	"fmt"
	"github.com/civet148/log"
	"os"
	"path/filepath"
	"strings"
)

const (
	TableNamePrefix  = "TableName"
	CustomizeCodeTip = "////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////"
)

func ExportToSqlFile(cmd *CmdFlags, ddl *CreateDatabaseDDL, tables []*TableSchema) (err error) {
	if len(tables) == 0 {
		return nil //no table found
	}
	var strDatabase = fmt.Sprintf("`%s`", cmd.Database)
	var strTemplate string

	strTemplate += ddl.CreateSQL + ";\n"
	strTemplate += fmt.Sprintf(`USE %s;`, strDatabase)
	strTemplate += "\n\n"
	for _, t := range tables {
		strTemplate += "\n"
		strTemplate += t.TableCreateSQL
		strTemplate += ";\n"
	}
	var fi *os.File
	fi, err = os.OpenFile(cmd.ExportTo, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return log.Errorf("open file [%v] error (%v)", cmd.ExportTo, err.Error())
	}
	_, err = fi.WriteString(strTemplate)
	if err != nil {
		return log.Errorf(err.Error())
	}
	return nil
}

func ExportTableSchema(cmd *CmdFlags, tables []*TableSchema) (err error) {

	for _, v := range tables {

		_, errStat := os.Stat(cmd.OutDir)
		if errStat != nil && os.IsNotExist(errStat) {

			log.Info("mkdir [%v]", cmd.OutDir)
			if err = os.MkdirAll(cmd.OutDir, os.ModeDir); err != nil {
				log.Error("mkdir [%v] error (%v)", cmd.OutDir, err.Error())
				return
			}
		}

		v.OutDir = cmd.OutDir

		if cmd.PackageName == "" {
			//mkdir by output dir + scheme name
			cmd.PackageName = v.SchemeName
			if strings.LastIndex(cmd.OutDir, fmt.Sprintf("%v", os.PathSeparator)) == -1 {
				v.SchemeDir = fmt.Sprintf("%v/%v", cmd.OutDir, cmd.PackageName)
			} else {
				v.SchemeDir = fmt.Sprintf("%v%v", cmd.OutDir, cmd.PackageName)
			}
		} else {
			v.SchemeDir = fmt.Sprintf("%v/%v", cmd.OutDir, cmd.PackageName) //mkdir by package name
		}

		_, errStat = os.Stat(v.SchemeDir)

		if errStat != nil && os.IsNotExist(errStat) {

			log.Info("mkdir [%v]", v.SchemeDir)
			if err = os.MkdirAll(v.SchemeDir, os.ModePerm); err != nil {
				log.Errorf("mkdir path name [%v] error (%v)", v.SchemeDir, err.Error())
				return
			}
		}

		var strPrefix, strSuffix string
		if cmd.Prefix != "" {
			strPrefix = fmt.Sprintf("%v_", cmd.Prefix)
		}
		if cmd.Suffix != "" {
			strSuffix = fmt.Sprintf("_%v", cmd.Suffix)
		}

		v.OutFilePath = fmt.Sprintf("%v/%v%v%v.go", v.SchemeDir, strPrefix, v.TableName, strSuffix)
		if err = exportModels(cmd, v); err != nil {
			return err
		}
		if err = exportDAO(cmd, v); err != nil {
			return err
		}
	}
	return nil
}

func exportModels(cmd *CmdFlags, table *TableSchema) (err error) {

	var strHead, strContent string
	var packages = make(map[string]bool)
	//write package name
	strHead += fmt.Sprintf("package %v\n\n", cmd.PackageName)

	//write table name in camel case naming
	table.TableNameCamelCase = BigCamelCase(table.TableName)
	table.TableComment = ReplaceCRLF(table.TableComment)
	strContent += fmt.Sprintf("const %s%v = \"%v\" //%v \n\n", TableNamePrefix, table.TableNameCamelCase, fmt.Sprintf("%s", table.TableName), table.TableComment)

	table.StructName = fmt.Sprintf("%s%s", table.TableNameCamelCase, strings.ToUpper(cmd.Suffix))
	table.StructDAO = fmt.Sprintf("%sDAO", table.TableNameCamelCase)
	for i, v := range table.Columns {
		table.Columns[i].Comment = ReplaceCRLF(v.Comment)
	}
	//检查此表是否存在需要导入第三方包类型的字段
	spectTypes := getImportSpecTypes(cmd, table)

	for _, st := range spectTypes {
		if table.TableName != st.Table && st.Table != TABLE_ALL {
			continue
		}

		for k, v := range st.Package {
			if ok := packages[v]; ok {
				continue //package already exist
			}
			strHead += fmt.Sprintf(`import %s "%s"`, k, v)
			strHead += "\n"
			packages[v] = true
		}
	}
	if cmd.BaseModel != nil && len(cmd.BaseModel.Package) > 0 {
		for k, v := range cmd.BaseModel.Package {
			if ok := packages[v]; ok {
				continue //package already exist
			}
			strHead += fmt.Sprintf(`import %s "%s"`, k, v)
			strHead += "\n"
			packages[v] = true
		}
	}
	if haveDecimal(cmd, table, table.Columns, cmd.EnableDecimal) {
		strHead += cmd.ImportVer + "\n\n" //根据数据库中是否存在decimal类型决定是否导入sqlca包
	}

	strContent += makeColumnConsts(cmd, table)
	strContent += makeTableStructure(cmd, table)
	strContent += makeObjectMethods(cmd, table)
	//strContent += makeTableCreateSQL(cmd, table)

	return writeToFile(table.OutFilePath, strHead+strContent)
}

func getImportSpecTypes(cmd *CmdFlags, table *TableSchema) (specTypes []*SpecType) {
	for _, col := range table.Columns {
		for _, st := range cmd.SpecTypes {
			if col.Name == st.Column {
				specTypes = append(specTypes, st)
			}
		}
	}
	return
}

func haveDecimal(cmd *CmdFlags, table *TableSchema, TableCols []TableColumn, enableDecimal bool) (ok bool) {
	for _, v := range TableCols {
		_, ok = GetGoColumnType(cmd, table.TableName, v, enableDecimal, nil)
		if ok {
			break
		}
	}
	return
}

func exportDAO(cmd *CmdFlags, table *TableSchema) (err error) {
	var strContent string
	if cmd.DAO == "" {
		return nil
	}
	var strDir, strOutputFilePath string

	if strings.LastIndex(cmd.OutDir, fmt.Sprintf("%v", os.PathSeparator)) == -1 {
		strDir = fmt.Sprintf("%v/%v", cmd.OutDir, cmd.DAO)
	} else {
		strDir = fmt.Sprintf("%v%v", cmd.OutDir, cmd.DAO)
	}
	if _, err = os.Stat(strDir); err != nil {
		if err = os.MkdirAll(strDir, os.ModePerm); err != nil {
			return log.Errorf("mkdir %s error [%s]", strDir, err.Error())
		}
	}

	var fi os.FileInfo

	strOutputFilePath = filepath.Join(strDir, table.TableName+".go")
	if fi, err = os.Stat(strOutputFilePath); err == nil {
		if fi.IsDir() {
			return log.Errorf("dao file [%v] is dir", strOutputFilePath)
		}
	}

	log.Infof("generate [%v]", strOutputFilePath)
	strContent += fmt.Sprintf("package %v\n\n", cmd.DAO)
	strContent += fmt.Sprintf(`import (
    "fmt"
	"%s"
	"%s"
)

`, cmd.SqlcaPkg, cmd.ImportModels)

	strContent += makeNewMethod(cmd, table)
	strContent += makeOrmMethods(cmd, table)

	if err = writeToFile(strOutputFilePath, strContent); err != nil {
		return log.Errorf("export dao for table [%v] to file [%v] error [%s]", table.TableName, strOutputFilePath, err.Error())
	}

	return nil
}

// make new DAO method
func makeNewMethod(cmd *CmdFlags, table *TableSchema) (strContent string) {
	strContent += fmt.Sprintf(`
type %s struct {
	db *sqlca.Engine
}

`, table.StructDAO)

	strContent += fmt.Sprintf(`
func New%v(db *sqlca.Engine) *%v {
	return &%v{
		db: db,
	}
}

`, table.StructDAO, table.StructDAO, table.StructDAO)
	return
}

func makeModelTableName(cmd *CmdFlags, table *TableSchema) string {
	return fmt.Sprintf("%s.%s%s", cmd.PackageName, TableNamePrefix, table.TableNameCamelCase)
}

func makeModelStructName(cmd *CmdFlags, table *TableSchema) string {
	return fmt.Sprintf("%s.%s", cmd.PackageName, table.StructName)
}

func makeObjectMethods(cmd *CmdFlags, table *TableSchema) (strContent string) {

	for _, v := range table.Columns { //添加结构体成员Get/Set方法

		if IsInSlice(v.Name, cmd.Without) {
			continue
		}
		strColName := BigCamelCase(v.Name)
		strColType, _ := GetGoColumnType(cmd, table.TableName, v, cmd.EnableDecimal, cmd.TinyintAsBool)
		strContent += MakeGetter(table.StructName, strColName, strColType)
		strContent += MakeSetter(table.StructName, strColName, strColType)
	}
	return
}

func makeTableCreateSQL(cmd *CmdFlags, table *TableSchema) (strContent string) {
	strContent += "/*\n"
	strContent += table.TableCreateSQL + ";\n"
	strContent += "*/\n"
	return
}

func makeOrmMethods(cmd *CmdFlags, table *TableSchema) (strContent string) {
	strContent += makeOrmInsertMethod(cmd, table)
	strContent += makeOrmUpsertMethod(cmd, table)
	strContent += makeOrmUpdateMethod(cmd, table)
	strContent += makeOrmQueryByIdMethod(cmd, table)
	strContent += makeOrmQueryByConditionMethod(cmd, table)
	return
}

func makeTableBaseModel(cmd *CmdFlags) (strContent string) {
	var baseModel = cmd.BaseModel
	if baseModel != nil {
		modelType := baseModel.Type
		strContent = fmt.Sprintf("%s\n", modelType)
		for k := range baseModel.Package {
			if k != "" {
				strContent = fmt.Sprintf("%s.%s\n", k, modelType)
			}
		}
	}
	return strContent
}

func makeOrmInsertMethod(cmd *CmdFlags, table *TableSchema) (strContent string) {
	return fmt.Sprintf(`
//insert into table by data model
func (dao *%v) Insert(do *%s) (lastInsertId, rowsAffected int64, err error) {
	return dao.db.Model(&do).Table(%s).Insert()
}

`, table.StructDAO, makeModelStructName(cmd, table), makeModelTableName(cmd, table))
}

func makeOrmUpsertMethod(cmd *CmdFlags, table *TableSchema) (strContent string) {
	return fmt.Sprintf(`
//insert if not exist or update columns on duplicate key...
func (dao *%v) Upsert(do *%s, columns...string) (lastInsertId int64, err error) {
    if len(columns) == 0 {
        return 0, fmt.Errorf("no columns to update")
    }
	return dao.db.Model(&do).Table(%s).Select(columns...).Upsert()
}

`, table.StructDAO, makeModelStructName(cmd, table), makeModelTableName(cmd, table))
}

func makeOrmUpdateMethod(cmd *CmdFlags, table *TableSchema) (strContent string) {
	return fmt.Sprintf(`
//update table set columns where id=xxx
func (dao *%v) Update(do *%s, columns...string) (rows int64, err error) {
    if len(columns) == 0 {
        return 0, fmt.Errorf("no columns to update")
    }
	return dao.db.Model(&do).Table(%s).Select(columns...).Update()
}

`, table.StructDAO, makeModelStructName(cmd, table), makeModelTableName(cmd, table))
}

func makeOrmQueryByIdMethod(cmd *CmdFlags, table *TableSchema) (strContent string) {
	return fmt.Sprintf(`
//query records by id
func (dao *%v) QueryById(id interface{}, columns...string) (do *%s, err error) {
	if _, err = dao.db.Model(&do).Table(%s).Id(id).Select(columns...).Query(); err != nil {
		return nil, err
	}
	return
}

`, table.StructDAO, makeModelStructName(cmd, table), makeModelTableName(cmd, table))
}

func makeOrmQueryByConditionMethod(cmd *CmdFlags, table *TableSchema) (strContent string) {
	return fmt.Sprintf(`
//query records by conditions
func (dao *%v) QueryByCondition(conditions map[string]interface{}, columns...string) (dos []*%s, err error) {
    if len(conditions) == 0 {
        return nil, fmt.Errorf("condition must not be empty")
    }
    e := dao.db.Model(&dos).Table(%s).Select(columns...)
    for k, v := range conditions {
        e.Eq(k, v)
    }
	if _, err = e.Query(); err != nil {
		return nil, err
	}
	return
}

`, table.StructDAO, makeModelStructName(cmd, table), makeModelTableName(cmd, table))
}

func makeColumnConsts(cmd *CmdFlags, table *TableSchema) (strContent string) {
	var strUpperTableName string
	var strUpperColumnName string
	strUpperTableName = strings.ToUpper(table.TableName)

	strContent += fmt.Sprintf("const (\n")
	for _, v := range table.Columns {
		strUpperColumnName = strings.ToUpper(v.Name)
		strContent += fmt.Sprintf("%s_%s_%s = \"%s\"\n", strUpperTableName, "COLUMN", strUpperColumnName, v.Name)
	}
	strContent += fmt.Sprintf(")\n\n")
	return
}

func makeTableStructure(cmd *CmdFlags, table *TableSchema) (strContent string) {

	strContent += fmt.Sprintf("type %v struct { \n", table.StructName)

	strContent += makeTableBaseModel(cmd) //base model

	for _, v := range table.Columns {
		if cmd.IsBaseColumn(v.Name) {
			continue
		}
		if IsInSlice(v.Name, cmd.Without) {
			continue
		}

		var tagValues []string
		var strColType, strColName string
		strColName = BigCamelCase(v.Name)
		strColType, _ = GetGoColumnType(cmd, table.TableName, v, cmd.EnableDecimal, cmd.TinyintAsBool)

		for _, t := range cmd.ExtraTags {
			tv := v.Name
			if t == "bson" && tv == "id" {
				tv = "_id"
			} else if t == "gorm" {
				if isPrimartyKey(tv) {
					tv = fmt.Sprintf("column:%s;primaryKey;autoIncrement", tv)
				} else if isCreateTime(tv) {
					tv = fmt.Sprintf("column:%s;default:CURRENT_TIMESTAMP;autoCreateTime", tv)
				} else if isUpdateTime(tv) {
					tv = fmt.Sprintf("column:%s;default:CURRENT_TIMESTAMP;autoUpdateTime", tv)
				} else {
					tv = fmt.Sprintf("column:%s", tv)
				}
			}
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", t, tv))
		}
		for _, t := range cmd.TagTypes {
			if t.Column != v.Name {
				continue
			}
			if t.Table == table.TableName || t.Table == TABLE_ALL {
				tv := t.TagValue
				if strings.Contains(tv, "\"") {
					tv = strings.ReplaceAll(tv, "\"", "")
				}
				tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", t.TagName, tv))
			}
		}
		if IsInSlice(v.Name, cmd.ReadOnly) {
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", TAG_NAME_SQLCA, TAG_VALUE_READ_ONLY))
		} else if v.IsNullable == "YES" {
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", TAG_NAME_SQLCA, TAG_VALUE_IS_NULL))
		}
		//添加成员和标签
		strContent += MakeTags(cmd, strColName, strColType, v.Name, v.Comment, strings.Join(tagValues, " "))

		v.GoName = strColName
		v.GoType = strColType
	}

	strContent += "}\n\n"

	return
}

func isPrimartyKey(tv string) bool {
	return tv == "id" || tv == "uid"
}

func isCreateTime(tv string) bool {
	return tv == "create_time" || tv == "create_at" || tv == "created_time" || tv == "created_at"
}

func isUpdateTime(tv string) bool {
	return tv == "update_time" || tv == "update_at" || tv == "updated_time" || tv == "updated_at"
}

func isDeleteTime(tv string) bool {
	return tv == "delete_time" || tv == "delete_at" || tv == "deleted_time" || tv == "deleted_at"
}

func GenerateMethodDeclare(strShortName, strStructName, strMethodName, strArgs, strReturn, strLogic string) (strFunc string) {
	if strReturn == "" {
		strFunc = fmt.Sprintf("func (%s *%s) %s(%s) {\n", strShortName, strStructName, strMethodName, strArgs)
	} else {
		strFunc = fmt.Sprintf("func (%s *%s) %s(%s) %s {\n", strShortName, strStructName, strMethodName, strArgs, strReturn)
	}
	strFunc += strLogic
	strFunc += fmt.Sprintf("}\n\n")
	return
}
