package schema

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed templates/proto_file.tmpl
var protoFileTmpl string

type ProtoFieldTmplData struct {
	Type    string
	Name    string
	Ordinal int
	Comment string
}

type ProtoMessageTmplData struct {
	Name   string
	Fields []ProtoFieldTmplData
}

type ProtoFileTmplData struct {
	PackageName   string
	HasGogoImport bool
	GogoOptions   []string
	ProtoOptions  map[string]string
	Messages      []ProtoMessageTmplData
}

func MakeProtoHead(cmd *CmdFlags) string {
	data := ProtoFileTmplData{
		PackageName:   cmd.PackageName,
		HasGogoImport: len(cmd.GogoOptions) > 0,
		GogoOptions:   cmd.GogoOptions,
		ProtoOptions:  cmd.ProtoOptions,
	}

	tmpl := template.Must(template.New("proto_head").Parse(`syntax = "proto3";
package {{.PackageName}};

{{if .HasGogoImport}}
import "github.com/gogo/protobuf/gogoproto/gogo.proto";
{{end -}}
{{- range $o := .GogoOptions}}
option {{$o}};
{{end -}}
{{- range $k, $v := .ProtoOptions}}
option {{$k}}="{{$v}}";
{{end}}

`))
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func MakeProtoBody(cmd *CmdFlags, table *TableSchema) string {
	strTableName := TableNameToStructName(table.TableName)
	var fields []ProtoFieldTmplData
	for i, v := range table.Columns {
		if IsInSlice(v.Name, cmd.Without) {
			continue
		}
		no := i + 1
		strColName := ConvertFieldStyle(v.Name, cmd.FieldStyle)
		strColType := GetProtoColumnType(table.TableName, v)
		fields = append(fields, ProtoFieldTmplData{
			Type:    strColType,
			Name:    strColName,
			Ordinal: no,
			Comment: v.Comment,
		})
	}

	data := ProtoMessageTmplData{
		Name:   strTableName,
		Fields: fields,
	}

	tmpl := template.Must(template.New("proto_message").Parse(`message {{.Name}} {
{{- range $f := .Fields}}
	{{$f.Type}} {{$f.Name}} = {{$f.Ordinal}}; //{{$f.Comment}}
{{- end}}
}

`))
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func ExportProtoSchema(provider SchemaProvider, cmd *CmdFlags) error {
	schemas, err := provider.QueryTableSchemas(cmd)
	if err != nil {
		return logErrorf(err.Error())
	}

	for _, v := range schemas {
		if err = provider.QueryTableColumns(v); err != nil {
			return logError(err.Error())
		}

		msgs := []ProtoMessageTmplData{buildProtoMessageData(cmd, v)}
		data := buildProtoFileData(cmd, msgs)
		tmpl := template.Must(template.New("proto_file").Parse(protoFileTmpl))
		var buf strings.Builder
		if err = tmpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("execute proto template error: %w", err)
		}

		file, err := CreateOutputFile(cmd, v, "proto", false)
		if err != nil {
			return err
		}
		_, err = file.WriteString(buf.String())
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func buildProtoFileData(cmd *CmdFlags, msgs []ProtoMessageTmplData) ProtoFileTmplData {
	opts := make(map[string]string, len(cmd.ProtoOptions))
	for k, v := range cmd.ProtoOptions {
		opts[k] = v
	}
	return ProtoFileTmplData{
		PackageName:   cmd.PackageName,
		HasGogoImport: len(cmd.GogoOptions) > 0,
		GogoOptions:   cmd.GogoOptions,
		ProtoOptions:  opts,
		Messages:      msgs,
	}
}

func buildProtoMessageData(cmd *CmdFlags, table *TableSchema) ProtoMessageTmplData {
	strTableName := TableNameToStructName(table.TableName)
	var fields []ProtoFieldTmplData
	for i, v := range table.Columns {
		if IsInSlice(v.Name, cmd.Without) {
			continue
		}
		no := i + 1
		strColName := ConvertFieldStyle(v.Name, cmd.FieldStyle)
		strColType := GetProtoColumnType(table.TableName, v)
		fields = append(fields, ProtoFieldTmplData{
			Type:    strColType,
			Name:    strColName,
			Ordinal: no,
			Comment: v.Comment,
		})
	}
	return ProtoMessageTmplData{
		Name:   strTableName,
		Fields: fields,
	}
}

func logError(msg string) error {
	return fmt.Errorf(msg)
}

func logErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
