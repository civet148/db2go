package mysql

import (
	"fmt"
	"os"
	"strings"

	"github.com/civet148/db2go/schema"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
)

type ExporterMysql struct {
	Cmd     *schema.CmdFlags
	Engine  *sqlca.Engine
	Schemas []*schema.TableSchema
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
	var cmd = m.Cmd
	var e = m.Engine
	var schemas = m.Schemas
	//var tableNames []string

	if cmd.Database == "" {
		err = fmt.Errorf("no database selected")
		return log.Error(err.Error())
	}
	//var strDatabaseName = fmt.Sprintf("'%v'", cmd.Database)
	log.Infof("ready to export tables %v", cmd.Tables)
	var ddl *schema.CreateDatabaseDDL
	ddl, err = m.queryCreateDatabaseDDL(cmd, e)
	if err != nil {
		return err
	}
	if schemas, err = m.queryTableSchemas(cmd, e); err != nil {

		return log.Errorf("query tables error [%s]", err.Error())
	}
	for _, v := range schemas {
		if err = m.queryTableColumns(v); err != nil {
			return log.Error(err.Error())
		}
		if err = m.queryTableIndexes(v); err != nil {
			return log.Error(err.Error())
		}
		if err = m.queryTableCreateStructure(v); err != nil {
			return log.Error(err.Error())
		}
	}
	err = schema.ExportTableSchema(cmd, schemas)
	if err != nil {
		return err
	}
	if cmd.ExportDDL != "" {
		err = schema.ExportToSqlFile(cmd, ddl, schemas)
		if err != nil {
			log.Warnf("export to file [%s] error [%s]", cmd.ExportDDL, err.Error())
		}
	}
	return nil
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

func (m *ExporterMysql) queryCreateDatabaseDDL(cmd *schema.CmdFlags, e *sqlca.Engine) (ddl *schema.CreateDatabaseDDL, err error) {
	_, err = e.Model(&ddl).QueryRaw("SHOW CREATE DATABASE `%s`", cmd.Database)
	if err != nil {
		return nil, log.Error(err.Error())
	}
	return ddl, nil
}

func (m *ExporterMysql) queryTableSchemas(cmd *schema.CmdFlags, e *sqlca.Engine) (schemas []*schema.TableSchema, err error) {

	var strQuery string
	var tables []string

	if cmd.Database == "" {
		err = fmt.Errorf("no database selected")
		log.Error(err.Error())
		return
	}
	var strDatabaseName = fmt.Sprintf("'%v'", cmd.Database)

	log.Infof("ready to export tables %v", cmd.Tables)
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
		 SELECT
				`TABLE_NAME` as table_name, `COLUMN_NAME` as column_name, `DATA_TYPE` as data_type, `COLUMN_TYPE` as column_type, `EXTRA` as extra,
				`COLUMN_KEY` as column_key, `COLUMN_COMMENT` as column_comment, `IS_NULLABLE` as is_nullable, COLUMN_DEFAULT as column_default, COLUMN_KEY as column_key
		 FROM `INFORMATION_SCHEMA`.`COLUMNS`
		 WHERE `TABLE_SCHEMA` = 'test' AND `TABLE_NAME` = 'users' ORDER BY ORDINAL_POSITION ASC
	*/
	var e = m.Engine
	_, err = e.Model(&table.Columns).QueryRaw("select `TABLE_NAME` as table_name, `COLUMN_NAME` as column_name, `DATA_TYPE` as data_type, `COLUMN_TYPE` as column_type, `EXTRA` as extra,"+
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

func (m *ExporterMysql) queryTableCreateStructure(table *schema.TableSchema) (err error) {
	if _, err = m.Engine.Model(&table.TableName, &table.TableCreateSQL).QueryRaw("SHOW CREATE TABLE `%s`", table.TableName); err != nil {
		log.Error(err.Error())
		return
	}
	return
}

func (m *ExporterMysql) queryTableIndexes(table *schema.TableSchema) (err error) {
	/*
		SELECT
		    TABLE_SCHEMA AS 'db_name', TABLE_NAME AS 'table_name', INDEX_NAME AS 'index_name', COLUMN_NAME AS 'column_name',
			SEQ_IN_INDEX AS 'seq_in_index', INDEX_TYPE AS 'index_type', NON_UNIQUE AS 'non_unique', INDEX_COMMENT AS 'index_comment'
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = 'test' AND TABLE_NAME = 'inventory_out'
		ORDER BY INDEX_NAME, SEQ_IN_INDEX;
	*/
	var e = m.Engine
	_, err = e.Model(&table.Indexes).QueryRaw(`SELECT
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
