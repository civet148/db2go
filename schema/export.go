package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
)

type SchemaProvider interface {
	QueryTableSchemas(cmd *CmdFlags) ([]*TableSchema, error)
	QueryTableColumns(table *TableSchema) error
	QueryTableIndexes(table *TableSchema) error
	QueryCreateDatabaseDDL(cmd *CmdFlags) (*CreateDatabaseDDL, error)
}

func NewEngine(cmd *CmdFlags) (*sqlca.Engine, error) {
	if strings.TrimSpace(cmd.SSH) != "" {
		return sqlca.NewEngine(cmd.ConnUrl, sshOption(cmd.SSH))
	}
	return sqlca.NewEngine(cmd.ConnUrl)
}

func sshOption(strSSH string) *sqlca.Options {
	if strSSH == "" {
		return nil
	}
	ssh := sqlca.ParseUrl(strSSH)
	return &sqlca.Options{
		SSH: &sqlca.SSH{
			User:     ssh.User,
			Password: ssh.Password,
			Host:     ssh.Host,
		},
	}
}

func ExportGoSchema(provider SchemaProvider, cmd *CmdFlags) error {
	ddl, err := provider.QueryCreateDatabaseDDL(cmd)
	if err != nil {
		log.Warnf("query create database DDL error [%s]", err.Error())
	}

	schemas, err := provider.QueryTableSchemas(cmd)
	if err != nil {
		return log.Errorf("query tables error [%s]", err.Error())
	}

	for _, v := range schemas {
		if err = provider.QueryTableColumns(v); err != nil {
			return log.Error(err.Error())
		}
		if err = provider.QueryTableIndexes(v); err != nil {
			return log.Error(err.Error())
		}
	}

	if err = ExportTableSchema(cmd, schemas); err != nil {
		return err
	}

	if cmd.ExportDDL != "" {
		if ddl == nil {
			ddl, err = provider.QueryCreateDatabaseDDL(cmd)
			if err != nil {
				log.Warnf("query DDL for export error [%s]", err.Error())
			}
		}
		if ddl != nil {
			if err = ExportToSqlFile(cmd, ddl, schemas); err != nil {
				log.Warnf("export to file [%s] error [%s]", cmd.ExportDDL, err.Error())
			}
		}
	}
	return nil
}

func ExportProtoSchema(provider SchemaProvider, cmd *CmdFlags) error {
	schemas, err := provider.QueryTableSchemas(cmd)
	if err != nil {
		return log.Errorf(err.Error())
	}

	strHead := MakeProtoHead(cmd)
	var file *os.File
	for i, v := range schemas {
		if err = provider.QueryTableColumns(v); err != nil {
			return log.Error(err.Error())
		}

		var appendMode bool
		if i > 0 && cmd.OneFile {
			appendMode = true
		}

		strBody := MakeProtoBody(cmd, v)

		if file, err = CreateOutputFile(cmd, v, "proto", appendMode); err != nil {
			return log.Error(err.Error())
		}

		if i == 0 {
			file.WriteString(strHead)
		} else if !cmd.OneFile {
			file.WriteString(strHead)
		}
		file.WriteString(strBody)
	}
	if file != nil {
		file.Close()
	}
	return nil
}

func ExportToSqlFile(cmd *CmdFlags, ddl *CreateDatabaseDDL, tables []*TableSchema) (err error) {
	if len(tables) == 0 {
		return nil
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
