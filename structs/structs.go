package structs

import (
	"fmt"
	"github.com/civet148/db2go/cmder"
	"github.com/civet148/db2go/schema"
	"strings"
)

const (
	STRUCT_PROMPT    = "STRUCT> "
	STRUCT_DELIMITER = '}'
)

type GoMember struct {
	Name string
	Type string
	Tags string
}

type GoStruct struct {
	Name     string
	Members  []GoMember
	Contents []string
	cmd      *schema.Commander
	r        *cmder.CmdReader
}

func ExportStruct(cmd *schema.Commander) {

	var r = cmder.Prompt(STRUCT_PROMPT, STRUCT_DELIMITER)
	fmt.Println("")
	fmt.Println("//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println("")

	if gs, err := NewGoStruct(cmd, r); err != nil {
		fmt.Printf("fmt.Sscanf error %s\n", err)
		return
	} else {
		if err = gs.parseMembers(); err != nil {
			fmt.Printf("parse struct members error %s\n", err)
			return
		}
		if err = gs.generateMethods(); err != nil {
			fmt.Printf("generate struct methods error %s\n", err)
			return
		}
	}
}

func NewGoStruct(cmd *schema.Commander, r *cmder.CmdReader) (gs *GoStruct, err error) {

	if len(r.TrimInputs) == 0 {
		fmt.Printf("nil input\n")
		return nil, fmt.Errorf("nil input")
	}

	var strInput = r.TrimInputs[0]
	var strPrefix, strStructName, strSuffix string
	if _, err = fmt.Sscanf(strInput, "%s %s %s", &strPrefix, &strStructName, &strSuffix); err != nil {
		fmt.Printf("fmt.Sscanf error %s\n", err)
		return
	}
	if strStructName == "" {
		panic("struct syntax no illegal, go struct must begin with [type xxx struct {]")
	}
	contents := r.TrimInputs[1:]
	if len(contents) == 0 {
		panic("struct syntax no illegal, go struct have no any member")
	}
	gs = &GoStruct{
		Name:     strStructName,
		Contents: contents,
		cmd:      cmd,
		r:        r,
	}
	return
}

func (g *GoStruct) getShortName() string {
	var strLowerName = strings.ToLower(g.Name)
	return strLowerName[:1]
}

func (g *GoStruct) parseMembers() (err error) {
	for _, v := range g.Contents {

		var strName, strType, strTags string
		fmt.Sscanf(v, "%s %s %s", &strName, &strType, &strTags)
		if strName == "" || strType == "" {
			continue
		}
		g.Members = append(g.Members, GoMember{
			Name: strings.TrimSpace(strName),
			Type: strings.TrimSpace(strType),
			Tags: strings.TrimSpace(strTags),
		})
	}
	return
}

func (g *GoStruct) generateMethods() (err error) {
	var strCode string
	strCode += g.generatePackage()
	strCode += g.generateImports()
	strCode += g.generateStruct()
	strCode += g.generateString()
	strCode += g.generateGoString()
	for _, v := range g.Members {
		strCode += g.generateGetter(v)
		strCode += g.generateSetter(v)
	}

	fmt.Println(strCode)
	return
}

func (g *GoStruct) generatePackage() (strCode string) {
	if g.cmd.PackageName != "" {
		strCode = fmt.Sprintf("package %s\n\n", g.cmd.PackageName)
	}
	return
}

func (g *GoStruct) generateImports() (strCode string) {
	return fmt.Sprintf("import \"encoding/json\"\n\n")
}

func (g *GoStruct) generateStruct() (strCode string) {
	return fmt.Sprintf("%s\n\n", g.r.RawInput)
}

func (g *GoStruct) generateString() (strCode string) {
	var strLogic string
	strLogic = fmt.Sprintf(`    data, _ := json.Marshal(%s)
    return string(data)
`, g.getShortName())
	return schema.GenerateMethodDeclare(g.getShortName(), g.Name, schema.METHOD_NAME_STRING, schema.METHOD_ARGS_NULL, "string", strLogic)
}

func (g *GoStruct) generateGoString() (strCode string) {
	var strLogic string
	strLogic = fmt.Sprintf(`    return %s.String()
`, g.getShortName())
	return schema.GenerateMethodDeclare(g.getShortName(), g.Name, schema.METHOD_NAME_GOSTRING, schema.METHOD_ARGS_NULL, "string", strLogic)
}

func (g *GoStruct) generateGetter(m GoMember) (strCode string) {

	var strLogic string
	strLogic = fmt.Sprintf(`    return %s.%v
`, g.getShortName(), m.Name)

	var strMethodName = schema.METHOD_NAME_GET + schema.CamelCaseConvert(m.Name)
	return schema.GenerateMethodDeclare(g.getShortName(), g.Name, strMethodName, schema.METHOD_ARGS_NULL, m.Type, strLogic)
}

func (g *GoStruct) generateSetter(m GoMember) (strCode string) {

	var strLogic string
	strLogic = fmt.Sprintf(`    %s.%v = v
`, g.getShortName(), m.Name)
	var strArgs = fmt.Sprintf("v %s", m.Type)
	var strMethodName = schema.METHOD_NAME_SET + schema.CamelCaseConvert(m.Name)
	return schema.GenerateMethodDeclare(g.getShortName(), g.Name, strMethodName, strArgs, schema.METHOD_ARGS_NULL, strLogic)
}
