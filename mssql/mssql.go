package mssql

import "C"
import (
	"fmt"
	"github.com/civet148/db2go/schema"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
	"os"
	"strings"
)

/*
-- 查询所有数据表和注释
SELECT table_name FROM INFORMATION_SCHEMA.TABLES
SELECT A.name as table_name, C.value as table_comment FROM sys.tables A left JOIN sys.extended_properties C ON C.major_id = A.object_id  and minor_id=0 WHERE A.name = 'classes'

-- 查询某些表字段名、类型和注释
SELECT A.name as table_name, C.value as table_comment FROM sys.tables A
LEFT JOIN sys.extended_properties C ON C.major_id = A.object_id  and minor_id=0 WHERE A.name in ('users')
*/

type ExporterMssql struct {
	Cmd     *schema.Commander
	Engine  *sqlca.Engine
	Schemas []*schema.TableSchema
}

func init() {
	schema.Register(schema.SCHEME_MSSQL, NewExporterMssql)
}

func NewExporterMssql(cmd *schema.Commander, e *sqlca.Engine) schema.Exporter {

	return &ExporterMssql{
		Cmd:    cmd,
		Engine: e,
	}
}

func (m *ExporterMssql) ExportGo() (err error) {
	var cmd = m.Cmd
	var schemas = m.Schemas
	//var tableNames []string

	if cmd.Database == "" {
		err = fmt.Errorf("no database selected")
		log.Error(err.Error())
		return
	}
	//var strDatabaseName = fmt.Sprintf("'%v'", cmd.Database)
	log.Infof("ready to export tables [%v]", cmd.Tables)

	if schemas, err = m.queryTableSchemas(); err != nil {
		log.Errorf("query tables error [%s]", err.Error())
		return
	}
	for _, v := range schemas {
		if err = m.queryTableColumns(v); err != nil {
			log.Error(err.Error())
			return
		}
	}
	return schema.ExportTableSchema(cmd, schemas)
}

func (m *ExporterMssql) ExportProto() (err error) {
	var cmd = m.Cmd
	var schemas = m.Schemas
	if schemas, err = m.queryTableSchemas(); err != nil {
		log.Errorf(err.Error())
		return
	}

	var file *os.File
	strHead := schema.MakeProtoHead(cmd)
	for i, v := range schemas {
		if err = m.queryTableColumns(v); err != nil {
			log.Error(err.Error())
			return
		}

		var append bool
		if i > 0 && cmd.OneFile {
			append = true
		}

		strBody := schema.MakeProtoBody(cmd, v)

		if file, err = schema.CreateOutputFile(cmd, v, "proto", append); err != nil {
			log.Error(err.Error())
			return
		}

		if i == 0 {
			file.WriteString(strHead)
		} else if !cmd.OneFile {
			file.WriteString(strHead)
		}
		file.WriteString(strBody)
	}
	file.Close()
	return
}

//查询当前库下所有表名
func (m *ExporterMssql) queryTableNames() (rows int64, err error) {
	var e = m.Engine
	var cmd = m.Cmd
	strQuery := fmt.Sprintf(`SELECT table_name FROM INFORMATION_SCHEMA.TABLES`)
	if rows, err = e.Model(&cmd.Tables).QueryRaw(strQuery); err != nil {
		log.Errorf(err.Error())
		return
	}
	return
}

//查询表和注释、引擎等等基本信息
func (m *ExporterMssql) queryTableSchemas() (schemas []*schema.TableSchema, err error) {

	var cmd = m.Cmd
	var e = m.Engine
	var strQuery string
	var tables []string

	if cmd.Database == "" {
		err = fmt.Errorf("no database selected")
		log.Error(err.Error())
		return
	}

	if len(cmd.Tables) == 0 {
		var rows int64
		if rows, err = m.queryTableNames(); err != nil {
			log.Errorf("query table names from database [%v] err [%+v]", cmd.Database, err.Error())
			return
		}
		if rows == 0 {
			err = fmt.Errorf("no table in database [%v]", cmd.Database)
			log.Errorf(err.Error())
			return
		}
	}

	for _, v := range cmd.Tables {
		tables = append(tables, fmt.Sprintf("'%v'", v))
	}

	log.Infof("ready to export tables %v", tables)

	strQuery = fmt.Sprintf(
		`SELECT '%v' as table_schema, A.name as table_name, C.value as table_comment FROM sys.tables A 
                LEFT JOIN sys.extended_properties C ON C.major_id = A.object_id  and minor_id=0 WHERE A.name in (%v)`,
		cmd.Database, strings.Join(tables, ","))
	_, err = e.Model(&schemas).QueryRaw(strQuery)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	return
}

func (m *ExporterMssql) queryTableColumns(table *schema.TableSchema) (err error) {

	var e = m.Engine
	_, err = e.Model(&table.Columns).QueryRaw(`SELECT table_name, column_name, data_type FROM INFORMATION_SCHEMA.COLUMNS 
                                                        WHERE table_catalog='test' and table_name in ('%v') order by table_name,ordinal_position`, table.TableName)

	if err != nil {
		log.Error(err.Error())
		return
	}
	schema.HandleCommentCRLF(table)
	return schema.ConvertMssqlColumnType(table) //转换Mssql数据库字段类型为MYSQL映射的类型
}
