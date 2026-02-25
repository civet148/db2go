package postgres

import (
	"fmt"
	"os"
	"strings"

	"github.com/civet148/db2go/schema"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
)

/*
-- 查询所有数据表和注释
SELECT
	relname AS table_name,
	CAST ( obj_description ( relfilenode, 'pg_class' ) AS VARCHAR ) AS table_comment
FROM
	pg_class C
WHERE
	relkind = 'r'
	AND relname NOT LIKE'pg_%'
	AND relname NOT LIKE'sql_%'
ORDER BY
	relname

-- 查询某些表字段名、类型和注释
SELECT
  C.relname as table_name,
	A.attname AS column_name,
	format_type ( A.atttypid, A.atttypmod ) AS data_type,
	col_description ( A.attrelid, A.attnum ) AS column_comment
FROM
	pg_class AS C,
	pg_attribute AS A
WHERE
	C.relname in ('users','classes')
	AND A.attrelid = C.oid
	AND A.attnum > 0
ORDER BY C.relname,A.attnum
*/

type ExporterPostgres struct {
	Cmd     *schema.CmdFlags
	Engine  *sqlca.Engine
	Schemas []*schema.TableSchema
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

func (m *ExporterPostgres) ExportGo() (err error) {
	var cmd = m.Cmd
	var schemas = m.Schemas
	//var tableNames []string

	if cmd.Database == "" {
		err = fmt.Errorf("no database selected")
		log.Error(err.Error())
		return
	}
	//var strDatabaseName = fmt.Sprintf("'%v'", cmd.Database)

	if schemas, err = m.queryTableSchemas(); err != nil {
		log.Errorf("query tables error [%s]", err.Error())
		return
	}
	for _, v := range schemas {
		if err = m.queryTableColumns(v); err != nil {
			return log.Error(err.Error())
		}
		if err = m.queryTableIndexes(v); err != nil {
			return log.Error(err.Error())
		}
	}
	return schema.ExportTableSchema(cmd, schemas)
}

func (m *ExporterPostgres) ExportProto() (err error) {
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

// 查询当前库下所有表名
func (m *ExporterPostgres) queryTableNames() (rows int64, err error) {
	var e = m.Engine
	var cmd = m.Cmd
	strQuery := `SELECT relname AS table_name FROM pg_class C WHERE relkind = 'r' AND relname NOT LIKE 'pg_%%' AND relname NOT LIKE 'sql_%%' ORDER BY relname`
	if rows, err = e.Model(&cmd.Tables).QueryRaw(strQuery); err != nil {
		log.Errorf(err.Error())
		return
	}
	return
}

// 查询表和注释、引擎等等基本信息
func (m *ExporterPostgres) queryTableSchemas() (schemas []*schema.TableSchema, err error) {

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
		if v[0] == '-' {
			continue
		}
		tables = append(tables, fmt.Sprintf("'%v'", v))
	}

	strQuery = fmt.Sprintf(
		`SELECT '%v' as table_schema, relname AS table_name, CAST ( obj_description ( relfilenode, 'pg_class' ) AS VARCHAR ) AS table_comment
                 FROM pg_class C WHERE relkind = 'r' AND relname in (%v) ORDER BY relname`, cmd.Database, strings.Join(tables, ","))
	_, err = e.Model(&schemas).QueryRaw(strQuery)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	return
}

func (m *ExporterPostgres) queryTableColumns(table *schema.TableSchema) (err error) {

	var e = m.Engine
	_, err = e.Model(&table.Columns).QueryRaw(`SELECT C.relname as table_name, A.attname AS column_name, format_type(A.atttypid,A.atttypmod) AS data_type,
	col_description ( A.attrelid, A.attnum ) AS column_comment FROM pg_class AS C, pg_attribute AS A WHERE	C.relname = '%v' AND A.attrelid = C.oid	AND A.attnum > 0 
    AND format_type(A.atttypid,A.atttypmod) != '-'
    ORDER BY C.relname,A.attnum`, table.TableName)

	if err != nil {
		log.Error(err.Error())
		return
	}
	schema.HandleCommentCRLF(table)
	return schema.ConvertPostgresColumnType(table) //转换postgres数据库字段类型为MYSQL映射的类型
}

/*
SELECT

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
WHERE n.nspname = 'test'

	AND t.relname = 'inventory_out'
	AND NOT ix.indisprimary  -- 排除主键，如果需要主键则移除这个条件

ORDER BY index_name, k.ord;
*/
func (m *ExporterPostgres) queryTableIndexes(table *schema.TableSchema) (err error) {
	var e = m.Engine
	_, err = e.Model(&table.Indexes).QueryRaw(`SELECT

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
