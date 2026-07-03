package parser

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	importPkgPrefix         = "import"
	importPkgMultiPrefix    = "import ("
	importPkgMultiSuffix    = ")"
	varDeclarePrefix        = "var"
	varDeclareMultiPrefix   = "var ("
	varDeclareMultiSuffix   = ")"
	constDeclarePrefix      = "const"
	constDeclareMultiPrefix = "const ("
	constDeclareMultiSuffix = ")"
	typeDeclarePrefix       = "type"
	typeDeclareStruct       = "struct"
	typeDeclareInterface    = "interface"
)

var replaceChars = map[string]string{
	" ":  "",
	"\t": "",
	"\n": "",
	"\r": "",
}

// MergeCode 合并两个差异代码
// base: 基础代码
// work: 工作代码
// return: merge合并后的代码内容
func MergeCode(base, work *GoFileParseResult) (merge string, err error) {
	if base == nil || work == nil {
		return merge, errors.New("base or work code parse result is nil")
	}
	var importCodes = mergeImportPackages(base, work)
	log.Printf("------------------------import-------------------------------")
	for _, code := range importCodes {
		fmt.Printf("%v\n", code.GetRaw())
	}
	log.Printf("------------------------var----------------------------------")
	var varCodes = mergeVars(base, work)
	for _, code := range varCodes {
		fmt.Printf("%v\n", code.GetRaw())
	}
	log.Printf("------------------------const--------------------------------")
	var constCodes = mergeConsts(base, work)
	for _, code := range constCodes {
		fmt.Printf("%v\n", code.GetRaw())
	}
	log.Printf("------------------------struct--------------------------------")
	var typesCodes = mergeTypes(base, work)
	for _, code := range typesCodes {
		fmt.Printf("%s\n", code)
	}
	log.Printf("------------------------function--------------------------------")
	var funcCodes = mergeFuncs(base, work)
	for _, code := range funcCodes {
		fmt.Printf("%s\n", code)
	}
	return merge, nil
}

func mergeImportPackages(base, work *GoFileParseResult) (codes []*CodeLine) {
	return mergePackageVarConst(base.Imports, work.Imports, importPkgPrefix, importPkgMultiPrefix, importPkgMultiSuffix)
}

func mergeVars(base, work *GoFileParseResult) (codes []*CodeLine) {
	return mergePackageVarConst(base.Vars, work.Vars, varDeclarePrefix, varDeclareMultiPrefix, varDeclareMultiSuffix)
}

func mergeConsts(base, work *GoFileParseResult) (codes []*CodeBlock) {
	var codeHashMap = make(map[string]string)
	var lineHashMap = make(map[string]string)
	for _, cb := range base.Consts {
		key := cb.GetHash()
		codeHashMap[key] = cb.GetCode()
		codes = append(codes, cb)
		for _, line := range cb.Lines {
			lineHashMap[line.GetHash()] = line.GetCode()
		}
	}

	for _, cb := range work.Consts {
		key := cb.GetHash()
		if _, ok := codeHashMap[key]; !ok {
			for _, lc := range cb.Lines {
				var code string
				if lc.Key != "" { //单独声明
					code = strings.Replace(lc.Code, constDeclarePrefix, "", 1)
				} else {
					if strings.Contains(lc.Code, constDeclareMultiPrefix) || strings.Contains(lc.Code, constDeclareMultiSuffix) {
						continue
					}
					code = lc.Code
				}
				if _, ok = lineHashMap[CodeHash(code)]; ok {
					lc.Disabled = true
				}
			}
			if cb.IsEmpty() {
				continue
			}
			codes = append(codes, cb)
		}
	}
	return codes
}

func mergePackageVarConst(baseLineBlocks, workLineBlocks []*CodeBlock, singlePrefix, multiPrefix, multiSuffix string) (codes []*CodeLine) {
	var codeHashMap = make(map[string]string)
	codes = append(codes, &CodeLine{
		Raw:  multiPrefix,
		Code: multiPrefix,
	})
	for _, lb := range baseLineBlocks {
		for _, lc := range lb.Lines {
			var code string
			if lc.Key != "" { //单独声明
				code = strings.Replace(lc.Code, singlePrefix, "", 1)
			} else { // 多行声明
				if strings.Contains(lc.Code, multiPrefix) || strings.Contains(lc.Code, multiSuffix) {
					continue
				}
				code = lc.Code
			}
			codes = append(codes, &CodeLine{
				Raw:     code + " " + lc.Comment,
				Code:    code,
				Comment: lc.Comment,
			})
			codeHashMap[CodeHash(code)] = code
		}
	}
	for _, lb := range workLineBlocks {
		for _, lc := range lb.Lines {
			var code string
			if lc.Key != "" { //单独声明
				code = strings.Replace(lc.Code, singlePrefix, "", 1)
			} else {
				if strings.Contains(lc.Code, multiPrefix) || strings.Contains(lc.Code, multiSuffix) {
					continue
				}
				code = lc.Code
			}
			if _, ok := codeHashMap[CodeHash(code)]; !ok {
				codes = append(codes, &CodeLine{
					Raw:     code + " " + lc.Comment,
					Code:    code,
					Comment: lc.Comment,
				})
			}
		}
	}
	codes = append(codes, &CodeLine{
		Raw:  multiSuffix,
		Code: multiSuffix,
	})
	return codes
}

func mergeTypes(base, work *GoFileParseResult) map[string]*TypeInfo {
	var codeFieldMap = make(map[string]string)
	var codeMethodMap = make(map[string]string)
	for _, bt := range base.Types {
		for _, line := range bt.Lines {
			if line.GetKey() == "" || line.IsTypeStart() || line.IsTypeEnd() {
				continue
			}
			codeFieldMap[line.Key] = line.Code
		}
		for _, method := range bt.Methods {
			codeMethodMap[method.GetKey()] = method.GetCode()
		}
	}
	for k, wt := range work.Types {
		for _, line := range wt.Lines {
			if line.GetKey() == "" || line.IsTypeStart() || line.IsTypeEnd() {
				continue
			}
			if _, ok := codeFieldMap[line.Key]; !ok {
				var bt *TypeInfo
				if bt, ok = base.Types[k]; ok {
					bt.InsertField(line)
				}
			}
		}
		for _, method := range wt.Methods {
			if _, ok := codeMethodMap[method.GetKey()]; !ok {
				var bt *TypeInfo
				if bt, ok = base.Types[k]; ok {
					bt.InsertMethod(method)
				}
			}
		}
	}
	return base.Types
}

func mergeFuncs(base, work *GoFileParseResult) []*CodeBlock {
	var codeFuncMap = make(map[string]string)
	for _, bf := range base.Functions {
		codeFuncMap[bf.GetKey()] = bf.GetCode()
	}
	for _, wf := range work.Functions {
		if _, ok := codeFuncMap[wf.GetKey()]; !ok {
			base.Functions = insertBeforeLast(base.Functions, wf)
		}
	}
	return base.Functions
}
