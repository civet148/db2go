package schema

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
	"os"
	"path/filepath"
	"strings"
)

const(
	TableNamePrefix = "TableName"
)

func ExportTableSchema(cmd *Commander, tables []*TableSchema) (err error) {

	for _, v := range tables {

		_, errStat := os.Stat(cmd.OutDir)
		if errStat != nil && os.IsNotExist(errStat) {

			log.Info("mkdir [%v]", cmd.OutDir)
			if err = os.Mkdir(cmd.OutDir, os.ModeDir); err != nil {
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
			if err = os.Mkdir(v.SchemeDir, os.ModeDir); err != nil {
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

		v.FileName = fmt.Sprintf("%v/%v%v%v.go", v.SchemeDir, strPrefix, v.TableName, strSuffix)
		if err = ExportTableColumns(cmd, v); err != nil {
			return
		}
	}

	return
}

func ExportTableColumns(cmd *Commander, table *TableSchema) (err error) {

	var File *os.File
	File, err = os.OpenFile(table.FileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0)
	if err != nil {
		log.Errorf("open file [%v] error (%v)", table.FileName, err.Error())
		return
	}

	var strHead, strContent string

	//write package name
	strHead += fmt.Sprintf("// Code generated by db2go. DO NOT EDIT.\n")
	strHead += fmt.Sprintf("// https://github.com/civet148/sqlca\n")
	strHead += fmt.Sprintf("\n")
	strHead += fmt.Sprintf("package %v\n\n", cmd.PackageName)

	//write table name in camel case naming
	table.TableNameCamelCase = CamelCaseConvert(table.TableName)
	table.TableComment = ReplaceCRLF(table.TableComment)
	strContent += fmt.Sprintf("const %s%v = \"%v\" //%v \n\n", TableNamePrefix, table.TableNameCamelCase, table.TableName, table.TableComment)

	table.StructName = fmt.Sprintf("%s%s", table.TableNameCamelCase, strings.ToUpper(cmd.Suffix))
	table.StructDAO = fmt.Sprintf("%sDAO", table.TableNameCamelCase)
	for i, v := range table.Columns {
		table.Columns[i].Comment = ReplaceCRLF(v.Comment)
	}
	if haveDecimal(cmd, table, table.Columns, cmd.EnableDecimal) {
		strHead += IMPORT_SQLCA + "\n\n" //根据数据库中是否存在decimal类型决定是否导入sqlca包
	}

	strContent += makeColumnConsts(cmd, table)
	strContent += makeTableStructure(cmd, table)
	strContent += makeObjectMethods(cmd, table)
	strContent += makeTableCreateSQL(cmd, table)
	if _, err = File.WriteString(strHead + strContent); err != nil {
		log.Errorf(err.Error())
		return err
	}
	makeDAO(cmd, table)
	return
}

func haveDecimal(cmd *Commander, table *TableSchema, TableCols []TableColumn, enableDecimal bool) (ok bool) {
	for _, v := range TableCols {
		_, ok = GetGoColumnType(cmd, table.TableName, v, enableDecimal, nil)
		if ok {
			break
		}
	}
	return
}

func makeDAO(cmd *Commander, table *TableSchema) {
	var err error
	var strContent string
	if cmd.DAO == "" {
		return
	}
	var strDir, strDAOFileName string

	if strings.LastIndex(cmd.OutDir, fmt.Sprintf("%v", os.PathSeparator)) == -1 {
		strDir = fmt.Sprintf("%v/%v", cmd.OutDir, cmd.DAO)
	} else {
		strDir = fmt.Sprintf("%v%v", cmd.OutDir, cmd.DAO)
	}
	if _, err = os.Stat(strDir); err != nil {
		if err = os.Mkdir(strDir, os.ModeDir); err != nil {
			log.Errorf("mkdir %s error [%s]", strDir, err.Error())
			return
		}
	}

	var file *os.File
	var fi os.FileInfo

	strDAOFileName = filepath.Join(strDir, table.TableName+".go")
	if fi, err = os.Stat(strDAOFileName); err != nil {
		file, err = os.OpenFile(strDAOFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0)
		if err != nil {
			log.Errorf("open file [%v] error (%v)", strDAOFileName, err.Error())
			return
		}
	} else {
		if fi.IsDir() || fi.Name() != "" {
			//log.Warnf("file %s already exist", strDAOFileName)
			return
		}
	}

	log.Infof("generate [%v]", strDAOFileName)
	strContent += fmt.Sprintf("package %v\n\n", cmd.DAO)
	strContent += fmt.Sprintf(`import (
	"github.com/civet148/sqlca/v2"
	"%s"
)

`, cmd.ImportModels)

	strContent += makeNewMethod(cmd, table)
	strContent += makeOrmMethods(cmd, table)
	if _, err = file.WriteString(strContent); err != nil {
		log.Errorf("export DAO for table [%v] to file [%v] error [%s]", table.TableName, strDAOFileName, err.Error())
		return
	}
}

//make new DAO method
func makeNewMethod(cmd *Commander, table *TableSchema) (strContent string) {
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

func makeModelTableName(cmd *Commander, table *TableSchema) string {
	return fmt.Sprintf("%s.%s%s", cmd.PackageName, TableNamePrefix, table.TableNameCamelCase)
}

func makeModelStructName(cmd *Commander, table *TableSchema) string {
	return fmt.Sprintf("%s.%s", cmd.PackageName, table.StructName)
}

func makeObjectMethods(cmd *Commander, table *TableSchema) (strContent string) {

	for _, v := range table.Columns { //添加结构体成员Get/Set方法

		if IsInSlice(v.Name, cmd.Without) {
			continue
		}
		strColName := CamelCaseConvert(v.Name)
		strColType, _ := GetGoColumnType(cmd, table.TableName, v, cmd.EnableDecimal, cmd.TinyintAsBool)
		strContent += MakeGetter(table.StructName, strColName, strColType)
		strContent += MakeSetter(table.StructName, strColName, strColType)
	}
	return
}

func makeTableCreateSQL(cmd *Commander, table *TableSchema) (strContent string) {
	strContent += "/*\n"
	strContent += table.TableCreateSQL + ";\n"
	strContent += "*/\n"
	return
}

func makeOrmMethods(cmd *Commander, table *TableSchema) (strContent string) {
	strContent += makeOrmInsertMethod(cmd, table)
	strContent += makeOrmUpsertMethod(cmd, table)
	strContent += makeOrmUpdateMethod(cmd, table)
	strContent += makeOrmQueryMethod(cmd, table)
	return
}

func makeOrmInsertMethod(cmd *Commander, table *TableSchema) (strContent string) {
	return fmt.Sprintf(`
//insert into table by data model
func (dao *%v) Insert(do *%s) (lastInsertId int64, err error) {
	return dao.db.Model(do).Table(%s).Insert()
}

`, table.StructDAO, makeModelStructName(cmd, table), makeModelTableName(cmd, table))
}

func makeOrmUpsertMethod(cmd *Commander, table *TableSchema) (strContent string) {
	return fmt.Sprintf(`
//insert if not exist or update columns on duplicate key...
func (dao *%v) Upsert(do *%s, columns...string) (lastInsertId int64, err error) {
	return dao.db.Model(do).Table(%s).Select(columns...).Upsert()
}

`, table.StructDAO, makeModelStructName(cmd, table), makeModelTableName(cmd, table))
}

func makeOrmUpdateMethod(cmd *Commander, table *TableSchema) (strContent string) {
	return fmt.Sprintf(`
//update table set columns where id=xxx
func (dao *%v) Update(do *%s, columns...string) (rows int64, err error) {
	return dao.db.Model(do).Table(%s).Select(columns...).Update()
}

`, table.StructDAO, makeModelStructName(cmd, table), makeModelTableName(cmd, table))
}

func makeOrmQueryMethod(cmd *Commander, table *TableSchema) (strContent string) {
	return fmt.Sprintf(`
func (dao *%v) QueryById(id interface{}, columns...string) (do *%s, err error) {
	if _, err = dao.db.Model(do).Table(%s).Id(id).Select(columns...).Query(); err != nil {
		return nil, err
	}
	return
}

`, table.StructDAO, makeModelStructName(cmd, table), makeModelTableName(cmd, table))
}

func makeColumnConsts(cmd *Commander, table *TableSchema) (strContent string) {
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

func makeTableStructure(cmd *Commander, table *TableSchema) (strContent string) {

	strContent += fmt.Sprintf("type %v struct { \n", table.StructName)
	for _, v := range table.Columns {

		if IsInSlice(v.Name, cmd.Without) {
			continue
		}

		var tagValues []string
		var strColType, strColName string
		strColName = CamelCaseConvert(v.Name)
		strColType, _ = GetGoColumnType(cmd, table.TableName, v, cmd.EnableDecimal, cmd.TinyintAsBool)
		if IsInSlice(v.Name, cmd.ReadOnly) {
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", sqlca.TAG_NAME_SQLCA, sqlca.SQLCA_TAG_VALUE_READ_ONLY))
		}
		for _, t := range cmd.Tags {
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", t, v.Name))
		}
		//添加成员和标签
		strContent += MakeTags(cmd, strColName, strColType, v.Name, v.Comment, strings.Join(tagValues, " "))

		v.GoName = strColName
		v.GoType = strColType
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
