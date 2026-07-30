package mysql

import (
	"fmt"
	"strings"

	"github.com/civet148/db2go/schema"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
)

type ExporterMysql struct {
	Cmd    *schema.CmdFlags
	Engine *sqlca.Engine
}

func init() {
	schema.Register(schema.SCHEME_MYSQL, NewExporterMysql)
}

func NewExporterMysql(cmd *schema.CmdFlags, e *sqlca.Engine) schema.Exporter {
	return &ExporterMysql{
		Cmd:    cmd,
		Engine: e,
	}
}

func (m *ExporterMysql) ExportGo() (err error) {
	return schema.ExportGoSchema(m, m.Cmd)
}

func (m *ExporterMysql) ExportProto() (err error) {
	return schema.ExportProtoSchema(m, m.Cmd)
}

func (m *ExporterMysql) QueryTableSchemas(cmd *schema.CmdFlags) ([]*schema.TableSchema, error) {
	var strQuery string
	var tables []string

	if cmd.Database == "" {
		return nil, fmt.Errorf("no database selected")
	}
	var strDatabaseName = fmt.Sprintf("'%v'", cmd.Database)

	for _, v := range cmd.Tables {
		if v[0] == '-' {
			continue
		}
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

	var schemas []*schema.TableSchema
	_, err := cmd.Engine.Model(&schemas).QueryRaw(strQuery)
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}
	return schemas, nil
}

func (m *ExporterMysql) QueryTableColumns(table *schema.TableSchema) (err error) {
	_, err = m.Cmd.Engine.Model(&table.Columns).QueryRaw("select `TABLE_NAME` as table_name, `COLUMN_NAME` as column_name, `DATA_TYPE` as data_type, `COLUMN_TYPE` as column_type, `EXTRA` as extra,"+
		" `COLUMN_KEY` as column_key, `COLUMN_COMMENT` as column_comment, `IS_NULLABLE` as is_nullable, COLUMN_DEFAULT as column_default, COLUMN_KEY as column_key "+
		" FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = '%v' AND `TABLE_NAME` = '%v' ORDER BY ORDINAL_POSITION ASC", table.SchemeName, table.TableName)
	if err != nil {
		log.Error(err.Error())
		return
	}
	schema.HandleCommentCRLF(table)
	log.Debugf("table [%s] columns %+v", table.TableName, table.Columns)
	return
}

func (m *ExporterMysql) QueryTableIndexes(table *schema.TableSchema) (err error) {
	_, err = m.Cmd.Engine.Model(&table.Indexes).QueryRaw(`SELECT
		    TABLE_SCHEMA AS 'db_name', TABLE_NAME AS 'table_name', INDEX_NAME AS 'index_name', COLUMN_NAME AS 'column_name',
			SEQ_IN_INDEX AS 'seq_in_index', INDEX_TYPE AS 'index_type', NON_UNIQUE AS 'non_unique', INDEX_COMMENT AS 'index_comment'
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'
		ORDER BY INDEX_NAME, SEQ_IN_INDEX`, table.SchemeName, table.TableName)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return nil
}

func (m *ExporterMysql) QueryCreateDatabaseDDL(cmd *schema.CmdFlags) (ddl *schema.CreateDatabaseDDL, err error) {
	_, err = cmd.Engine.Model(&ddl).QueryRaw("SHOW CREATE DATABASE `%s`", cmd.Database)
	if err != nil {
		return nil, log.Error(err.Error())
	}
	return ddl, nil
}

func (m *ExporterMysql) QueryTableCreateSQL(table *schema.TableSchema) error {
	_, err := m.Cmd.Engine.Model(&table.TableName, &table.TableCreateSQL).QueryRaw("SHOW CREATE TABLE `%s`", table.TableName)
	if err != nil {
		return log.Error(err.Error())
	}
	return nil
}
