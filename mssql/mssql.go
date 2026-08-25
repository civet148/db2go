package mssql

import (
	"fmt"
	"strings"

	"github.com/civet148/db2go/schema"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
)

type ExporterMssql struct {
	Cmd    *schema.CmdFlags
	Engine *sqlca.Engine
}

func init() {
	schema.Register(schema.SCHEME_MSSQL, NewExporterMssql)
}

func NewExporterMssql(cmd *schema.CmdFlags, e *sqlca.Engine) schema.Exporter {
	return &ExporterMssql{
		Cmd:    cmd,
		Engine: e,
	}
}

func (m *ExporterMssql) GetCmd() *schema.CmdFlags {
	return m.Cmd
}

func (m *ExporterMssql) ExportGo() (err error) {
	return schema.ExportGoSchema(m)
}

func (m *ExporterMssql) ExportProto() (err error) {
	return schema.ExportProtoSchema(m)
}

func (m *ExporterMssql) QueryTableSchemas() ([]*schema.TableSchema, error) {
	var strQuery string
	var tables []string
	cmd := m.Cmd
	if cmd.Database == "" {
		return nil, fmt.Errorf("no database selected")
	}

	if len(cmd.Tables) == 0 {
		var rows int64
		var err error
		strQuery = fmt.Sprintf(`SELECT table_name FROM INFORMATION_SCHEMA.TABLES`)
		if rows, err = cmd.Engine.Model(&cmd.Tables).QueryRaw(strQuery); err != nil {
			log.Errorf(err.Error())
			return nil, err
		}
		if rows == 0 {
			return nil, fmt.Errorf("no table in database [%v]", cmd.Database)
		}
	}

	for _, v := range cmd.Tables {
		if v[0] == '-' {
			continue
		}
		tables = append(tables, fmt.Sprintf("'%v'", v))
	}

	strQuery = fmt.Sprintf(
		`SELECT '%v' as table_schema, A.name as table_name, C.value as table_comment FROM sys.tables A 
                LEFT JOIN sys.extended_properties C ON C.major_id = A.object_id  and minor_id=0 WHERE A.name in (%v)`,
		cmd.Database, strings.Join(tables, ","))
	var schemas []*schema.TableSchema
	_, err := cmd.Engine.Model(&schemas).QueryRaw(strQuery)
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}
	return schemas, nil
}

func (m *ExporterMssql) QueryTableColumns(table *schema.TableSchema) (err error) {
	_, err = m.Cmd.Engine.Model(&table.Columns).QueryRaw(`SELECT table_name, column_name, data_type FROM INFORMATION_SCHEMA.COLUMNS 
                                                        WHERE table_catalog='test' and table_name in ('%v') order by table_name,ordinal_position`, table.TableName)

	if err != nil {
		log.Error(err.Error())
		return
	}
	schema.HandleCommentCRLF(table)
	return schema.ConvertMssqlColumnType(table)
}

func (m *ExporterMssql) QueryTableIndexes(table *schema.TableSchema) (err error) {
	return nil
}

func (m *ExporterMssql) QueryCreateDatabaseDDL() (*schema.CreateDatabaseDDL, error) {
	return nil, nil
}

func (m *ExporterMssql) QueryTableCreateSQL(table *schema.TableSchema) error {
	return nil
}
