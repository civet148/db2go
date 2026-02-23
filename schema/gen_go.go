package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/civet148/log"
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
	dir := filepath.Dir(cmd.ExportDDL)
	if err = MakeDir(dir); err != nil {
		return err
	}

	var fi *os.File
	fi, err = os.OpenFile(cmd.ExportDDL, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return log.Errorf("open file [%v] error (%v)", cmd.ExportDDL, err.Error())
	}
	_, err = fi.WriteString(strTemplate)
	if err != nil {
		return log.Errorf(err.Error())
	}
	return nil
}

func ExportTableSchema(cmd *CmdFlags, tables []*TableSchema) (err error) {

	var git = hasGit()
	if git {
		var ok bool
		ok, err = hasUnstagedChanges()
		if err != nil {
			return log.Errorf("git status error: %s", err)
		}
		if !ok {
			return log.Errorf("请先暂存/提交本地代码再重试 (Please stash or commit your work code and try it later)")
		}
		if err = gitCheckout(); err != nil {
			return err
		}
	}
	for _, v := range tables {
		if len(v.ImportPackages) == 0 {
			v.ImportPackages = make(map[string]bool)
		}
		if err = MakeDir(cmd.OutDir); err != nil {
			return err
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

		if err = MakeDir(v.SchemeDir); err != nil {
			return err
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

	if git {
		if err = gitCommitAndMerge(); err != nil {
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

	table.StructName = TableNameToStructName(table.TableNameCamelCase)
	table.StructDAO = TableNameToStructName(table.TableNameCamelCase)
	for i, v := range table.Columns {
		table.Columns[i].Comment = ReplaceCRLF(v.Comment)
	}
	//导入固定的系统包(比如time)
	for k, _ := range table.ImportPackages {
		strHead += fmt.Sprintf(`import "%s"`, k)
		strHead += "\n"
		packages[k] = true
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
		_, ok = GetGoColumnType(cmd, table, v, enableDecimal, nil)
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
	if err = MakeDir(strDir); err != nil {
		return err
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
	strContent += MakeTableNameGetter(table.StructName, table.TableName)
	strContent += "\n"
	for _, v := range table.Columns { //添加结构体成员Getter方法
		if IsInSlice(v.Name, cmd.Without) {
			continue
		}
		strColName := BigCamelCase(v.Name)
		strColType, _ := GetGoColumnType(cmd, table, v, cmd.EnableDecimal, cmd.TinyintAsBool)
		strContent += MakeGetter(table.StructName, strColName, strColType)
	}
	strContent += "\n"
	for _, v := range table.Columns { //添加结构体成员Setter方法
		if IsInSlice(v.Name, cmd.Without) {
			continue
		}
		strColName := BigCamelCase(v.Name)
		strColType, _ := GetGoColumnType(cmd, table, v, cmd.EnableDecimal, cmd.TinyintAsBool)
		strContent += MakeSetter(table.StructName, strColName, strColType)
	}
	strContent += "\n"
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

	for _, col := range table.Columns {
		if cmd.IsBaseColumn(col.Name) {
			continue
		}
		if IsInSlice(col.Name, cmd.Without) {
			continue
		}

		var tagValues []string
		var strColType, strColName string
		strColName = BigCamelCase(col.Name)
		strColType, _ = GetGoColumnType(cmd, table, col, cmd.EnableDecimal, cmd.TinyintAsBool)

		for _, t := range cmd.ExtraTags {
			tv := col.Name
			if t == "bson" && tv == "id" {
				tv = "_id"
			} else if t == "gorm" {
				if col.IsPrimaryKey() {
					tv = fmt.Sprintf("column:%s;primaryKey;autoIncrement;", tv)
				} else {
					if col.IsCreateTime() {
						tv = fmt.Sprintf("column:%s;type:%s;autoCreateTime;", tv, col.ColumnType)
					} else if col.IsUpdateTime() {
						tv = fmt.Sprintf("column:%s;type:%s;autoUpdateTime;", tv, col.ColumnType)
					} else {
						tv = fmt.Sprintf("column:%s;type:%s;", tv, col.ColumnType)
					}
					index, ok := table.GetGormIndexes(col.Name)
					if ok {
						tv += fmt.Sprintf("%s;", index)
					}
					if col.ColumnDefault != "" {
						tv += fmt.Sprintf("default:%s;", col.ColumnDefault)
					}
					if col.Comment != "" {
						tv += fmt.Sprintf("comment:%s;", handleColumnComment(col.Comment))
					}
				}
			}
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", t, tv))
		}
		for _, t := range cmd.TagTypes {
			if t.Column != col.Name {
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
		if IsInSlice(col.Name, cmd.ReadOnly) {
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", TAG_NAME_SQLCA, TAG_VALUE_READ_ONLY))
		} else if col.IsNullable == "YES" {
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", TAG_NAME_SQLCA, TAG_VALUE_IS_NULL))
		}
		//添加成员和标签
		strContent += MakeTags(cmd, strColName, strColType, col.Name, col.Comment, strings.Join(tagValues, " "))

		col.GoName = strColName
		col.GoType = strColType
	}

	strContent += "}\n\n"

	return
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
