package mysql

import (
	"fmt"
	"github.com/civet148/db2go/schema"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
	"os"
	"strings"
)

type ExporterMysql struct {
	Cmd     *schema.Commander
	Engine  *sqlca.Engine
	Schemas []*schema.TableSchema
}

func init() {
	schema.Register(schema.SCHEME_MYSQL, NewExporterMysql)
}

func NewExporterMysql(cmd *schema.Commander, e *sqlca.Engine) schema.Exporter {

	return &ExporterMysql{
		Cmd:    cmd,
		Engine: e,
	}
}

func (m *ExporterMysql) ExportGo() (err error) {
	var cmd = m.Cmd
	var e = m.Engine
	var schemas = m.Schemas
	//var tableNames []string

	if cmd.Database == "" {
		err = fmt.Errorf("no database selected")
		log.Error(err.Error())
		return
	}
	//var strDatabaseName = fmt.Sprintf("'%v'", cmd.Database)
	log.Infof("ready to export tables [%v]", cmd.Tables)

	if schemas, err = m.queryTableSchemas(cmd, e); err != nil {
		log.Errorf("query tables error [%s]", err.Error())
		return
	}
	for _, v := range schemas {
		if err = m.queryTableColumns(v); err != nil {
			log.Error(err.Error())
			return
		}
		if err = m.queryTableCreateStructure(v); err != nil {
			log.Error(err.Error())
			return
		}
	}

	return schema.ExportTableSchema(cmd, schemas)
}

func (m *ExporterMysql) ExportProto() (err error) {
	var cmd = m.Cmd
	var e = m.Engine
	var schemas = m.Schemas
	if schemas, err = m.queryTableSchemas(cmd, e); err != nil {
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

func (m *ExporterMysql) queryTableSchemas(cmd *schema.Commander, e *sqlca.Engine) (schemas []*schema.TableSchema, err error) {

	var strQuery string
	var tables []string

	if cmd.Database == "" {
		err = fmt.Errorf("no database selected")
		log.Error(err.Error())
		return
	}
	var strDatabaseName = fmt.Sprintf("'%v'", cmd.Database)

	log.Infof("ready to export tables [%v]", cmd.Tables)
	for _, v := range cmd.Tables {
		tables = append(tables, fmt.Sprintf("'%v'", v))
	}

	if len(tables) == 0 {
		strQuery = fmt.Sprintf("SELECT `TABLE_SCHEMA` as table_schema, `TABLE_NAME` as table_name, `ENGINE` as engine, `TABLE_COMMENT` as table_comment "+
			"FROM `INFORMATION_SCHEMA`.`TABLES` "+
			"where (`ENGINE`='myisam' OR `ENGINE` = 'innodb' OR `ENGINE` = 'tokudb') and `TABLE_SCHEMA` IN (%v) ORDER BY TABLE_SCHEMA",
			strDatabaseName)
	} else {
		strQuery = fmt.Sprintf("SELECT `TABLE_SCHEMA` as table_schema, `TABLE_NAME` as table_name, `ENGINE` as engine, `table_comment` as table_comment "+
			" FROM `INFORMATION_SCHEMA`.`TABLES` "+
			" WHERE (`ENGINE`='myisam' or `ENGINE` = 'innodb' or `ENGINE` = 'tokudb') and `TABLE_SCHEMA` in (%v) AND TABLE_NAME in (%v) ORDER BY TABLE_SCHEMA",
			strDatabaseName, strings.Join(tables, ","))
	}

	_, err = e.Model(&schemas).QueryRaw(strQuery)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	return
}

func (m *ExporterMysql) queryTableColumns(table *schema.TableSchema) (err error) {

	/*
	 SELECT `TABLE_NAME`, `COLUMN_NAME`, `DATA_TYPE`, `EXTRA`, `COLUMN_KEY`, `COLUMN_COMMENT`, `IS_NULLABLE` FROM `INFORMATION_SCHEMA`.`COLUMNS`
	 WHERE `TABLE_SCHEMA` = 'test' AND `TABLE_NAME` = 'users' ORDER BY ORDINAL_POSITION ASC
	*/
	var e = m.Engine
	_, err = e.Model(&table.Columns).QueryRaw("select `TABLE_NAME` as table_name, `COLUMN_NAME` as column_name, `DATA_TYPE` as data_type, `COLUMN_TYPE` as column_type, `EXTRA` as extra,"+
		" `COLUMN_KEY` as column_key, `COLUMN_COMMENT` as column_comment, `IS_NULLABLE` as is_nullable "+
		" FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = '%v' AND `TABLE_NAME` = '%v' ORDER BY ORDINAL_POSITION ASC", table.SchemeName, table.TableName)
	if err != nil {
		log.Error(err.Error())
		return
	}
	schema.HandleCommentCRLF(table)
	log.Debugf("table [%s] columns %+v", table.TableName, table.Columns)
	return
}

func (m *ExporterMysql) queryTableCreateStructure(table *schema.TableSchema) (err error) {
	if _, err = m.Engine.Model(&table.TableName, &table.TableCreateSQL).QueryRaw("SHOW CREATE TABLE %s", table.TableName); err != nil {
		log.Error(err.Error())
		return
	}
	return
}
