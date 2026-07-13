package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/civet148/log"
)

const (
	TableNamePrefix = "TableName"
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
	for _, v := range tables {
		if IsInSlice(v.TableName, cmd.ExcludeTables) {
			continue
		}
		v.InitPackage(cmd)

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

	if haveDecimal(cmd, table, table.Columns, cmd.EnableDecimal) {
		strHead += cmd.ImportVer + "\n\n" //根据数据库中是否存在decimal类型决定是否导入sqlca包
	}

	strContent += makeColumnConsts(cmd, table)
	strContent += makeTableStructure(cmd, table)
	strContent += makeObjectMethods(cmd, table)
	//strContent += makeTableCreateSQL(cmd, table)

	return writeToFile(table.OutFilePath, strHead+strContent, false)
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
		_, ok = GetGoColumnType(cmd, table, v, enableDecimal)
		if ok {
			break
		}
	}
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
		strColType, _ := GetGoColumnType(cmd, table, v, cmd.EnableDecimal)
		strContent += MakeGetter(table.StructName, strColName, strColType)
	}
	strContent += "\n"
	for _, v := range table.Columns { //添加结构体成员Setter方法
		if IsInSlice(v.Name, cmd.Without) {
			continue
		}
		strColName := BigCamelCase(v.Name)
		strColType, _ := GetGoColumnType(cmd, table, v, cmd.EnableDecimal)
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

// makeTableStructure 根据表结构生成Go语言的struct结构体定义
// 参数:
//   - cmd: 命令标志，包含各种配置选项
//   - table: 表结构信息，包含列名、类型等详细信息
//
// 返回值:
//   - strContent: 生成的Go结构体定义字符串
func makeTableStructure(cmd *CmdFlags, table *TableSchema) (strContent string) {

	// 添加结构体类型定义开始部分，使用表名作为结构体名称
	strContent += fmt.Sprintf("type %v struct { \n", table.StructName)

	// 遍历表的每一列
	for _, col := range table.Columns {
		// 跳过在排除列表中的列
		if IsInSlice(col.Name, cmd.Without) {
			continue
		}

		var tagValues []string            // 存储标签值的切片
		var strColType, strColName string // 列类型和列名的Go语言表示
		// 将列名转换为大驼峰命名
		strColName = BigCamelCase(col.Name)
		// 获取列的Go语言类型
		strColType, _ = GetGoColumnType(cmd, table, col, cmd.EnableDecimal)
		// 处理额外的标签
		for _, t := range cmd.ExtraTags {
			tv := col.Name
			// 特殊处理bson标签中的id字段
			if t == "bson" && tv == "id" {
				tv = "_id"
			} else if t == "gorm" {
				// 处理主键
				if col.IsPrimaryKey() {
					tv = fmt.Sprintf("column:%s;primaryKey;autoIncrement;", tv)
				} else {
					// 处理创建时间和更新时间
					if col.IsCreateTime() {
						tv = fmt.Sprintf("column:%s;type:%s;autoCreateTime;", tv, col.ColumnType)
					} else if col.IsUpdateTime() {
						tv = fmt.Sprintf("column:%s;type:%s;autoUpdateTime;", tv, col.ColumnType)
					} else {
						tv = fmt.Sprintf("column:%s;type:%s;", tv, col.ColumnType)
					}
					// 添加索引信息
					index, ok := table.GetGormIndexes(col.Name)
					if ok {
						tv += fmt.Sprintf("%s;", index)
					}
					// 添加默认值
					if col.ColumnDefault != "" {
						tv += fmt.Sprintf("default:%s;", col.ColumnDefault)
					} else {
						tv += fmt.Sprintf("default:null;")
					}
					// 添加列注释
					if col.Comment != "" {
						tv += fmt.Sprintf("comment:%s;", handleColumnComment(col.Comment))
					}
				}
			}
			// 将标签值添加到切片中
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", t, tv))
		}
		// 处理自定义标签类型
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
		// 处理只读列和可空列
		if IsInSlice(col.Name, cmd.ReadOnly) {
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", TAG_NAME_SQLCA, TAG_VALUE_READ_ONLY))
		} else if col.IsNullable == "YES" {
			tagValues = append(tagValues, fmt.Sprintf("%v:\"%v\"", TAG_NAME_SQLCA, TAG_VALUE_IS_NULL))
		}
		//添加成员和标签
		strContent += MakeTags(cmd, strColName, strColType, col.Name, col.Comment, strings.Join(tagValues, " "))

		col.GoName = strColName
	}

	strContent += "}\n\n"

	return
}
