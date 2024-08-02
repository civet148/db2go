package schema

import (
	"fmt"
	"strings"
)

func MakeProtoHead(cmd *Commander) (strContent string) {

	strContent += "syntax = \"proto3\";\n"
	strContent += fmt.Sprintf("package %v;\n\n", cmd.PackageName)

	if len(cmd.GogoOptions) > 0 {
		strContent += IMPORT_GOGO_PROTO + "\n"
		//strContent += IMPORT_GOOGOLE_PROTOBUF + "\n"
	}
	strContent += "\n"
	for _, v := range cmd.GogoOptions {
		strContent += fmt.Sprintf("option %v;\n", v)
	}
	for k, v := range cmd.ProtoOptions {
		if strings.Contains(v, "\"") {
			strContent += fmt.Sprintf("option %v=%v;\n", k, v)
		} else {
			strContent += fmt.Sprintf("option %v=\"%v\";\n", k, v)
		}
	}
	strContent += "\n"
	return
}

func MakeProtoBody(cmd *Commander, table *TableSchema) (strContent string) {

	strTableName := BigCamelCase(table.TableName)
	strContent += fmt.Sprintf("message %v {\n", strTableName)
	for i, v := range table.Columns {

		if IsInSlice(v.Name, cmd.Without) {
			continue
		}
		no := i + 1
		strColName := v.Name
		strColType := GetProtoColumnType(table.TableName, v.Name, v.DataType)
		strContent += fmt.Sprintf("	%-10s %-22s = %-2d; //%v\n", strColType, strColName, no, v.Comment)
	}
	strContent += "}\n\n"
	return
}

