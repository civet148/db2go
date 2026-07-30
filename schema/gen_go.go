package schema

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"
)

const (
	TableNamePrefix = "TableName"
)

//go:embed templates/go_model.tmpl
var goModelTmpl string

type ColumnTmplData struct {
	ConstName  string
	ColumnName string
	GoName     string
	GoType     string
	Tags       string
	Comment    string
}

type GoFileTmplData struct {
	PackageName    string
	TableNameConst string
	TableName      string
	TableComment   string
	StructName     string
	ImportPackages []string
	ImportVer      string
	Columns        []ColumnTmplData
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
			cmd.PackageName = v.SchemeName
			if strings.LastIndex(cmd.OutDir, fmt.Sprintf("%v", os.PathSeparator)) == -1 {
				v.SchemeDir = fmt.Sprintf("%v/%v", cmd.OutDir, cmd.PackageName)
			} else {
				v.SchemeDir = fmt.Sprintf("%v%v", cmd.OutDir, cmd.PackageName)
			}
		} else {
			v.SchemeDir = fmt.Sprintf("%v/%v", cmd.OutDir, cmd.PackageName)
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

func exportModels(cmd *CmdFlags, table *TableSchema) error {
	table.TableNameCamelCase = BigCamelCase(table.TableName)
	table.TableComment = ReplaceCRLF(table.TableComment)
	table.StructName = TableNameToStructName(table.TableNameCamelCase)
	table.StructDAO = table.StructName
	for i, v := range table.Columns {
		table.Columns[i].Comment = ReplaceCRLF(v.Comment)
	}

	var importPkgs []string
	for k := range table.ImportPackages {
		importPkgs = append(importPkgs, k)
	}

	specTypes := getImportSpecTypes(cmd, table)
	for _, st := range specTypes {
		if table.TableName != st.Table && st.Table != TABLE_ALL {
			continue
		}
		for _, v := range st.Package {
			if !contains(importPkgs, v) {
				importPkgs = append(importPkgs, v)
			}
		}
	}

	var importVer string
	if haveDecimal(cmd, table, table.Columns, cmd.EnableDecimal) {
		importVer = cmd.ImportVer
	}

	var cols []ColumnTmplData
	for _, col := range table.Columns {
		if IsInSlice(col.Name, cmd.Without) {
			continue
		}
		cols = append(cols, buildColumnTmplData(cmd, table, col))
	}

	data := GoFileTmplData{
		PackageName:    cmd.PackageName,
		TableNameConst: fmt.Sprintf("%s%v", TableNamePrefix, table.TableNameCamelCase),
		TableName:      table.TableName,
		TableComment:   table.TableComment,
		StructName:     table.StructName,
		ImportPackages: importPkgs,
		ImportVer:      importVer,
		Columns:        cols,
	}

	tmpl := template.Must(template.New("go_model").Parse(goModelTmpl))
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute go template error: %w", err)
	}

	return writeToFile(table.OutFilePath, buf.String(), false)
}

func buildColumnTmplData(cmd *CmdFlags, table *TableSchema, col TableColumn) ColumnTmplData {
	strColName := BigCamelCase(col.Name)
	strColType, _ := GetGoColumnType(cmd, table, col, cmd.EnableDecimal)
	tags := buildColumnTags(cmd, table, col)

	col.GoName = strColName

	return ColumnTmplData{
		ConstName:  fmt.Sprintf("%s%s_%s", BigCamelCase(table.TableName), "Column", strColName),
		ColumnName: col.Name,
		GoName:     strColName,
		GoType:     strColType,
		Tags:       tags,
		Comment:    col.Comment,
	}
}

func buildColumnTags(cmd *CmdFlags, table *TableSchema, col TableColumn) string {
	var tagValues []string

	var strJsonValue = col.Name
	if cmd.JsonStyle == JSON_STYLE_SMALL_CAMELCASE {
		strJsonValue = SmallCamelCase(col.Name)
	} else if cmd.JsonStyle == JSON_STYLE_BIG_CAMELCASE {
		strJsonValue = BigCamelCase(col.Name)
	}
	tagValues = append(tagValues, fmt.Sprintf(`json:"%v"`, strJsonValue))
	tagValues = append(tagValues, fmt.Sprintf(`db:"%v"`, col.Name))

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
				} else {
					tv += "default:null;"
				}
				if col.Comment != "" {
					tv += fmt.Sprintf("comment:%s;", handleColumnComment(col.Comment))
				}
			}
		}
		tagValues = append(tagValues, fmt.Sprintf(`%v:"%v"`, t, tv))
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
			tagValues = append(tagValues, fmt.Sprintf(`%v:"%v"`, t.TagName, tv))
		}
	}

	if IsInSlice(col.Name, cmd.ReadOnly) {
		tagValues = append(tagValues, fmt.Sprintf(`%v:"%v"`, TAG_NAME_SQLCA, TAG_VALUE_READ_ONLY))
	} else if col.IsNullable == "YES" {
		tagValues = append(tagValues, fmt.Sprintf(`%v:"%v"`, TAG_NAME_SQLCA, TAG_VALUE_IS_NULL))
	}

	return strings.Join(tagValues, " ")
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

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
