package opengauss

import (
	"fmt"
	"strings"

	"github.com/civet148/db2go/schema"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
)

type ExporterOpenGauss struct {
	Cmd    *schema.CmdFlags
	Engine *sqlca.Engine
}

func init() {
	schema.Register(schema.SCHEME_OPEN_GAUSS, NewExporterOpenGauss)
}

func NewExporterOpenGauss(cmd *schema.CmdFlags, e *sqlca.Engine) schema.Exporter {
	return &ExporterOpenGauss{
		Cmd:    cmd,
		Engine: e,
	}
}

func (m *ExporterOpenGauss) ExportGo() (err error) {
	return schema.ExportGoSchema(m, m.Cmd)
}

func (m *ExporterOpenGauss) ExportProto() (err error) {
	return schema.ExportProtoSchema(m, m.Cmd)
}

func (m *ExporterOpenGauss) QueryTableSchemas(cmd *schema.CmdFlags) ([]*schema.TableSchema, error) {
	var strQuery string
	var tables []string

	if cmd.Database == "" {
		return nil, fmt.Errorf("no database selected")
	}

	if len(cmd.Tables) == 0 {
		var rows int64
		var err error
		strQuery = `SELECT relname AS table_name FROM pg_class C WHERE relkind = 'r' AND relname NOT LIKE 'pg_%%' AND relname NOT LIKE 'sql_%%' ORDER BY relname`
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
		`SELECT '%v' as table_schema, relname AS table_name, CAST ( obj_description ( relfilenode, 'pg_class' ) AS VARCHAR ) AS table_comment
                 FROM pg_class C WHERE relkind = 'r' AND relname in (%v) ORDER BY relname`, cmd.Database, strings.Join(tables, ","))
	var schemas []*schema.TableSchema
	_, err := cmd.Engine.Model(&schemas).QueryRaw(strQuery)
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}
	return schemas, nil
}

func (m *ExporterOpenGauss) QueryTableColumns(table *schema.TableSchema) (err error) {
	_, err = m.Cmd.Engine.Model(&table.Columns).QueryRaw(`SELECT C.relname as table_name, A.attname AS column_name, format_type(A.atttypid,A.atttypmod) AS data_type,
	col_description ( A.attrelid, A.attnum ) AS column_comment FROM pg_class AS C, pg_attribute AS A WHERE	C.relname = '%v' AND A.attrelid = C.oid	AND A.attnum > 0 
    AND format_type(A.atttypid,A.atttypmod) != '-'
    ORDER BY C.relname,A.attnum`, table.TableName)

	if err != nil {
		log.Error(err.Error())
		return
	}
	schema.HandleCommentCRLF(table)
	return schema.ConvertPostgresColumnType(table)
}

func (m *ExporterOpenGauss) QueryTableIndexes(table *schema.TableSchema) (err error) {
	return nil
}

func (m *ExporterOpenGauss) QueryCreateDatabaseDDL(cmd *schema.CmdFlags) (*schema.CreateDatabaseDDL, error) {
	return nil, nil
}

func (m *ExporterOpenGauss) QueryTableCreateSQL(table *schema.TableSchema) error {
	return nil
}
