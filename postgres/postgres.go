package postgres

import (
	"fmt"
	"strings"

	"github.com/civet148/db2go/schema"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
)

type ExporterPostgres struct {
	Cmd    *schema.CmdFlags
	Engine *sqlca.Engine
}

func init() {
	schema.Register(schema.SCHEME_POSTGRES, NewExporterPostgres)
}

func NewExporterPostgres(cmd *schema.CmdFlags, e *sqlca.Engine) schema.Exporter {
	return &ExporterPostgres{
		Cmd:    cmd,
		Engine: e,
	}
}

func (m *ExporterPostgres) GetCmd() *schema.CmdFlags {
	return m.Cmd
}

func (m *ExporterPostgres) ExportGo() (err error) {
	return schema.ExportGoSchema(m)
}

func (m *ExporterPostgres) ExportProto() (err error) {
	return schema.ExportProtoSchema(m)
}

func (m *ExporterPostgres) QueryTableSchemas() ([]*schema.TableSchema, error) {
	var strQuery string
	var tables []string
	cmd := m.GetCmd()
	if cmd.Database == "" {
		return nil, fmt.Errorf("no database selected")
	}
	db := m.Cmd.Engine
	if len(cmd.Tables) == 0 {
		var rows int64
		var err error
		strQuery = `SELECT relname AS table_name FROM pg_class C WHERE relkind = 'r' AND relname NOT LIKE 'pg_%%' AND relname NOT LIKE 'sql_%%' ORDER BY relname`
		if rows, err = db.Model(&cmd.Tables).QueryRaw(strQuery); err != nil {
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
	_, err := db.Model(&schemas).QueryRaw(strQuery)
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}
	return schemas, nil
}

func (m *ExporterPostgres) QueryTableColumns(table *schema.TableSchema) (err error) {
	db := m.Cmd.Engine
	_, err = db.Model(&table.Columns).QueryRaw(`SELECT C.relname as table_name, A.attname AS column_name, format_type(A.atttypid,A.atttypmod) AS data_type,
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

func (m *ExporterPostgres) QueryTableIndexes(table *schema.TableSchema) (err error) {
	db := m.Cmd.Engine
	_, err = db.Model(&table.Indexes).QueryRaw(`SELECT
	n.nspname AS db_name,
	t.relname AS table_name,
	i.relname AS index_name,
	a.attname AS column_name,
	row_number() OVER (PARTITION BY i.relname ORDER BY array_position(ix.indkey, a.attnum)) AS seq_in_index,
	am.amname AS index_type,
	NOT ix.indisunique AS non_unique,
	obj_description(i.oid, 'pg_class') AS index_comment
FROM pg_index ix
JOIN pg_class t ON t.oid = ix.indrelid
JOIN pg_class i ON i.oid = ix.indexrelid
JOIN pg_namespace n ON n.oid = t.relnamespace
JOIN pg_am am ON i.relam = am.oid
CROSS JOIN LATERAL unnest(ix.indkey) WITH ORDINALITY AS k(attnum, ord)
JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = k.attnum AND a.attnum > 0
WHERE n.nspname = '%s' AND t.relname = '%s' AND NOT ix.indisprimary  
ORDER BY index_name, k.ord;`, table.SchemeName, table.TableName)
	if err != nil {
		log.Error(err.Error())
		return
	}
	return nil
}

func (m *ExporterPostgres) QueryCreateDatabaseDDL() (*schema.CreateDatabaseDDL, error) {
	return nil, nil
}

func (m *ExporterPostgres) QueryTableCreateSQL(table *schema.TableSchema) error {
	return nil
}
