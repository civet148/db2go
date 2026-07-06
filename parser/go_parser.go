package parser

import (
	"crypto/md5"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
)

// safeGetSubMatch 安全获取正则分组，防止索引越界panic
func safeGetSubMatch(subMatches []string, idx int) string {
	if idx < 0 || idx >= len(subMatches) {
		return ""
	}
	return subMatches[idx]
}

// SplitCodeAndComment 修复版：精准识别行尾//单行注释、行内/* */块注释，忽略字符串内注释符号
func SplitCodeAndComment(raw string) (code string, comment string) {
	line := raw
	var allComments []string

	// 1. 先提取行尾 // 单行注释（核心修复）
	// 规则：匹配不在双引号/单引号内部的尾部//
	reLineTailComment := regexp.MustCompile(`^(.*?)(?:"[^"]*"|'[^']*')*?(\/\/.*)$`)
	lineMatch := reLineTailComment.FindStringSubmatch(line)
	if len(lineMatch) > 0 {
		codePart := lineMatch[1]
		commentPart := lineMatch[2]
		if commentPart != "" {
			allComments = append(allComments, commentPart)
			line = codePart
		}
	}

	// 2. 提取行内 /* 块注释
	reBlockComment := regexp.MustCompile(`/\*.*?\*/`)
	blockMatches := reBlockComment.FindAllStringIndex(line, -1)
	if len(blockMatches) > 0 {
		// 从后往前替换，避免索引偏移
		for i := len(blockMatches) - 1; i >= 0; i-- {
			start, end := blockMatches[i][0], blockMatches[i][1]
			commentPart := line[start:end]
			allComments = append(allComments, commentPart)
			line = line[:start] + line[end:]
		}
	}

	// 清理纯代码末尾空格制表符
	code = strings.TrimRight(line, " \t")
	// 拼接所有注释
	var commentBuf strings.Builder
	for _, c := range allComments {
		if c != "" {
			commentBuf.WriteString(c)
			commentBuf.WriteString(" ")
		}
	}
	comment = strings.TrimSpace(commentBuf.String())
	return code, comment
}

// CodeLine 单行拆分结构：代码、注释完全分离
type CodeLine struct {
	Raw      string // 原始完整行（输出文件还原格式）
	Code     string // 去注释纯代码，差异对比专用
	Comment  string // 该行全部注释，对比忽略
	Key      string // 代码行对应的key，用于标识具体元素
	Disabled bool   // 已禁用(生成合并代码时忽略)
}

func (cl CodeLine) GetHash() string {
	return CodeHash(cl.Code)
}

func (cl CodeLine) GetRaw() string {
	return cl.Raw
}

func (cl CodeLine) GetCode() string {
	return cl.Code
}

func (cl CodeLine) GetComment() string {
	return cl.Comment
}

func (cl CodeLine) GetKey() string {
	return cl.Key
}

func (cl CodeLine) String() string {
	return cl.Raw
}

func (cl CodeLine) IsTypeStart() bool {
	code := strings.TrimSpace(cl.Code)
	if strings.HasPrefix(code, typeDeclarePrefix) && strings.HasSuffix(code, "{") {
		return true
	}
	return false
}

func (cl CodeLine) IsTypeEnd() bool {
	code := strings.TrimSpace(cl.Code)
	if code == "}" {
		return true
	}
	return false
}

// CodeBlock 通用代码块：import/var/const/func 统一结构
type CodeBlock struct {
	StartLine int         // 代码块起始行号
	Lines     []*CodeLine // 块内分行数据
}

func (cb CodeBlock) GetHash() string {
	return CodeHash(cb.GetCode())
}

func (cb CodeBlock) IsEmpty() bool {
	var count int
	for _, l := range cb.Lines {
		if !l.Disabled {
			count++
		}
	}
	if count == 0 {
		return true
	}
	hash := cb.Lines[0].GetHash()
	if count == 1 {
		if hash == CodeHash(importPkgPrefix) {
			return true
		}
		if hash == CodeHash(constDeclarePrefix) {
			return true
		}
		if hash == CodeHash(varDeclarePrefix) {
			return true
		}
	} else if count == 2 {
		if hash == CodeHash(importPkgMultiPrefix) {
			return true
		}
		if hash == CodeHash(constDeclareMultiPrefix) {
			return true
		}
		if hash == CodeHash(varDeclareMultiPrefix) {
			return true
		}
	}
	return false
}

func (cb CodeBlock) GetRaw() string {
	var code string
	for _, lb := range cb.Lines {
		if lb.Disabled {
			continue
		}
		code += lb.Raw + "\n"
	}
	return code
}

func (cb CodeBlock) String() string {
	var code string
	for _, lb := range cb.Lines {
		if lb.Disabled {
			continue
		}
		code += lb.Raw + "\n"
	}
	return code
}

func (cb CodeBlock) GetCode() string {
	var code string
	for _, lb := range cb.Lines {
		if lb.Disabled {
			continue
		}
		code += lb.Code + "\n"
	}
	return code
}

func (cb CodeBlock) GetKey() string {
	for _, lb := range cb.Lines {
		if lb.Disabled {
			continue
		}
		if lb.Key != "" {
			return lb.Key
		}
	}
	return ""
}

// TypeInfo 单个type定义存储结构
type TypeInfo struct {
	StartLine int          // type起始行
	Lines     []*CodeLine  // 类型完整分行代码
	Methods   []*CodeBlock // 该类型的方法列表
}

func (ti *TypeInfo) InsertField(field *CodeLine) {
	ti.Lines = insertBeforeLast(ti.Lines, field)
}

func (ti *TypeInfo) InsertMethod(method *CodeBlock) {
	ti.Methods = insertBeforeLast(ti.Methods, method)
}

func (ti *TypeInfo) String() string {
	var code string
	for _, l := range ti.Lines {
		if l.Disabled {
			continue
		}
		code += l.Raw + "\n"
	}
	for _, m := range ti.Methods {
		code += m.String() + "\n"
	}
	return code
}

// GoFileParseResult 文件完整解析结果，用于代码合并对比
type GoFileParseResult struct {
	PackageName string               // 包名
	Imports     []*CodeBlock         // import块列表
	Vars        []*CodeBlock         // 顶层var块
	Consts      []*CodeBlock         // 顶层const块
	Functions   []*CodeBlock         // 顶层函数块（不含方法）
	Types       map[string]*TypeInfo // key=类型名，value=类型详情
}

// extractKey 提取单行key
func extractKey(lineCode string, contextType string, typeName string) string {
	trimmed := strings.TrimSpace(lineCode)

	// 如果是空行、括号或只有注释的行，返回空key
	if trimmed == "" || trimmed == "{" || trimmed == "}" {
		return ""
	}

	// 移除行尾注释（但保留代码部分）
	codeOnly := trimmed
	if idx := strings.Index(codeOnly, "//"); idx >= 0 {
		codeOnly = strings.TrimSpace(codeOnly[:idx])
	}
	if idx := strings.Index(codeOnly, "/*"); idx >= 0 {
		codeOnly = strings.TrimSpace(codeOnly[:idx])
	}

	// 处理 import
	if contextType == "import" {
		// 如果是单独的 import 关键字行，返回空
		if trimmed == "import" || trimmed == "import (" {
			return ""
		}
		reImport := regexp.MustCompile(`^import\s+(?:([a-zA-Z_][a-zA-Z0-9_]*)\s+)?["']([^"']+)["']`)
		matches := reImport.FindStringSubmatch(trimmed)
		if len(matches) > 0 {
			if matches[1] != "" {
				return "Import." + matches[1]
			}
			return "Import." + matches[2]
		}
		return ""
	}

	// 处理 var/const 声明（同时支持单行和多行）
	if contextType == "var" || contextType == "const" {
		// 如果是单独的 var/const 关键字行（多行块的第一行），返回空
		if trimmed == "var" || trimmed == "var (" || trimmed == "const" || trimmed == "const (" {
			return ""
		}

		// 匹配 var/const 标识符（包括单行和多行声明）
		// 单行: const a = 1
		// 多行: a = 1 或 a int
		reIdent := regexp.MustCompile(`^(?:var|const)?\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:,|=|:|\(|\)|\s)`)
		matches := reIdent.FindStringSubmatch(codeOnly)
		if len(matches) > 1 {
			// 如果是 var a = 1 或 const a = 1 这种单行声明
			if strings.HasPrefix(trimmed, "var") || strings.HasPrefix(trimmed, "const") {
				if contextType == "var" {
					return "Var." + matches[1]
				}
				return "Const." + matches[1]
			}
			// 多行声明中的变量/常量
			if contextType == "var" {
				return "Var." + matches[1]
			}
			return "Const." + matches[1]
		}

		// 如果是简单的赋值声明
		reAssign := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s*[\+=]`)
		matchesAssign := reAssign.FindStringSubmatch(codeOnly)
		if len(matchesAssign) > 1 {
			if contextType == "var" {
				return "Var." + matchesAssign[1]
			}
			return "Const." + matchesAssign[1]
		}
		return ""
	}

	// 处理 type 结构体字段
	if contextType == "type_field" && typeName != "" {
		// 匹配结构体字段名
		reField := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\s+`)
		matches := reField.FindStringSubmatch(trimmed)
		if len(matches) > 1 && !strings.Contains(trimmed, "func") && !strings.Contains(trimmed, "interface") && !strings.Contains(trimmed, "struct") {
			return typeName + "." + matches[1]
		}
		return ""
	}

	// 处理函数
	if contextType == "func" {
		// 匹配函数名，包括接收者
		reFunc := regexp.MustCompile(`^func\s+(?:\(([^)]*)\)\s+)?([a-zA-Z_][a-zA-Z0-9_]*)`)
		matches := reFunc.FindStringSubmatch(trimmed)
		if len(matches) > 2 {
			if matches[1] != "" {
				// 有接收者：方法
				receiver := strings.TrimSpace(matches[1])
				// 提取接收者类型名
				receiverParts := strings.Fields(receiver)
				if len(receiverParts) > 0 {
					recvType := receiverParts[len(receiverParts)-1]
					// 移除指针符号
					recvType = strings.TrimPrefix(recvType, "*")
					return "Func." + recvType + "." + matches[2]
				}
				return "Func." + matches[2]
			}
			// 普通函数
			return "Func." + matches[2]
		}
		return ""
	}

	return ""
}

// extractReceiverType 提取函数接收者类型
func extractReceiverType(funcDecl *ast.FuncDecl) string {
	if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
		return ""
	}

	recv := funcDecl.Recv.List[0]
	switch t := recv.Type.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	case *ast.IndexExpr:
		// 泛型类型
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

// ParseGoFile 解析go文件入口
func ParseGoFile(filePath string) (*GoFileParseResult, error) {
	srcBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file err: %w", err)
	}
	fullSrc := string(srcBytes)
	srcRawLines := strings.Split(fullSrc, "\n")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, srcBytes, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("ast parse err: %w", err)
	}

	res := &GoFileParseResult{
		PackageName: node.Name.Name,
		Types:       make(map[string]*TypeInfo),
	}

	getLineNum := func(pos token.Pos) int {
		return fset.Position(pos).Line
	}

	// 先收集所有type定义，初始化Types map
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					tName := typeSpec.Name.Name
					if _, exists := res.Types[tName]; !exists {
						res.Types[tName] = &TypeInfo{
							StartLine: 0,
							Lines:     []*CodeLine{},
							Methods:   []*CodeBlock{},
						}
					}
				}
			}
		}
	}

	// 遍历顶层所有声明
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			startL := getLineNum(d.Pos())
			endL := getLineNum(d.End())
			var codeLines []*CodeLine

			// 根据不同的声明类型处理
			switch d.Tok {
			case token.IMPORT:
				// 处理 import
				for i := startL - 1; i < endL; i++ {
					if i < 0 || i >= len(srcRawLines) {
						continue
					}
					raw := srcRawLines[i]
					c, cm := SplitCodeAndComment(raw)
					key := extractKey(c, "import", "")
					codeLines = append(codeLines, &CodeLine{
						Raw:     raw,
						Code:    c,
						Comment: cm,
						Key:     key,
					})
				}
				res.Imports = append(res.Imports, &CodeBlock{
					StartLine: startL,
					Lines:     codeLines,
				})

			case token.VAR:
				// 处理 var 块，逐行提取 key
				for i := startL - 1; i < endL; i++ {
					if i < 0 || i >= len(srcRawLines) {
						continue
					}
					raw := srcRawLines[i]
					c, cm := SplitCodeAndComment(raw)
					key := extractKey(c, "var", "")
					codeLines = append(codeLines, &CodeLine{
						Raw:     raw,
						Code:    c,
						Comment: cm,
						Key:     key,
					})
				}
				res.Vars = append(res.Vars, &CodeBlock{
					StartLine: startL,
					Lines:     codeLines,
				})

			case token.CONST:
				// 处理 const 块，逐行提取 key
				for i := startL - 1; i < endL; i++ {
					if i < 0 || i >= len(srcRawLines) {
						continue
					}
					raw := srcRawLines[i]
					c, cm := SplitCodeAndComment(raw)
					key := extractKey(c, "const", "")
					codeLines = append(codeLines, &CodeLine{
						Raw:     raw,
						Code:    c,
						Comment: cm,
						Key:     key,
					})
				}
				res.Consts = append(res.Consts, &CodeBlock{
					StartLine: startL,
					Lines:     codeLines,
				})

			case token.TYPE:
				// 处理 type 定义
				for _, spec := range d.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					tName := typeSpec.Name.Name
					tStart := getLineNum(spec.Pos())
					tEnd := getLineNum(spec.End())

					var typeLines []*CodeLine

					// 先处理 type 声明行
					for i := tStart - 1; i < tEnd; i++ {
						if i < 0 || i >= len(srcRawLines) {
							continue
						}
						raw := srcRawLines[i]
						c, cm := SplitCodeAndComment(raw)
						// 第一行可能是 type xxx struct/interface
						key := extractKey(c, "type", tName)
						typeLines = append(typeLines, &CodeLine{
							Raw:     raw,
							Code:    c,
							Comment: cm,
							Key:     key,
						})
					}

					// 如果是结构体，处理字段
					if structType, ok := typeSpec.Type.(*ast.StructType); ok && structType.Fields != nil {
						// 获取结构体字段的起止行
						fieldStart := getLineNum(structType.Fields.Pos())
						fieldEnd := getLineNum(structType.Fields.End())

						// 创建字段行到typeLines的映射
						for i := fieldStart - 1; i < fieldEnd; i++ {
							if i < 0 || i >= len(srcRawLines) {
								continue
							}
							raw := srcRawLines[i]
							c, _ := SplitCodeAndComment(raw)
							key := extractKey(c, "type_field", tName)
							if key != "" {
								// 找到对应的行并更新key
								for j := range typeLines {
									if i == tStart-1+j {
										typeLines[j].Key = key
										break
									}
								}
							}
						}
					}

					// 更新或创建TypeInfo
					if info, exists := res.Types[tName]; exists {
						info.StartLine = tStart
						info.Lines = typeLines
						res.Types[tName] = info
					} else {
						res.Types[tName] = &TypeInfo{
							StartLine: tStart,
							Lines:     typeLines,
							Methods:   []*CodeBlock{},
						}
					}
				}
			}

		case *ast.FuncDecl:
			if d.Name == nil {
				continue
			}

			// 检查是否有接收者（是否是方法）
			receiverType := extractReceiverType(d)
			startL := getLineNum(d.Pos())
			endL := getLineNum(d.End())
			var codeLines []*CodeLine

			// 处理函数声明
			for i := startL - 1; i < endL; i++ {
				if i < 0 || i >= len(srcRawLines) {
					continue
				}
				raw := srcRawLines[i]
				c, cm := SplitCodeAndComment(raw)
				key := extractKey(c, "func", "")
				codeLines = append(codeLines, &CodeLine{
					Raw:     raw,
					Code:    c,
					Comment: cm,
					Key:     key,
				})
			}

			methodBlock := &CodeBlock{
				StartLine: startL,
				Lines:     codeLines,
			}

			if receiverType != "" {
				// 这是方法，添加到对应类型的Methods中
				if info, exists := res.Types[receiverType]; exists {
					info.Methods = append(info.Methods, methodBlock)
					res.Types[receiverType] = info
				} else {
					// 如果类型不存在（可能是未定义的类型或接口），但为了安全，我们仍然创建
					// 这种情况通常不会发生，因为类型应该已经定义
					res.Types[receiverType] = &TypeInfo{
						StartLine: 0,
						Lines:     []*CodeLine{},
						Methods:   []*CodeBlock{methodBlock},
					}
				}
			} else {
				// 普通函数，添加到Functions中
				res.Functions = append(res.Functions, methodBlock)
			}
		}
	}

	return res, nil
}

// PrintResult 调试打印所有解析内容
func PrintResult(r *GoFileParseResult) {
	fmt.Printf("===== Package Name: %s =====\n\n", r.PackageName)

	fmt.Println("======== Imports ========")
	for idx, blk := range r.Imports {
		fmt.Printf("Block %d LineStart:%d\n", idx+1, blk.StartLine)
		for _, cl := range blk.Lines {
			fmt.Printf("CODE: %q | COMMENT: %q | KEY: %s\n", cl.Code, cl.Comment, cl.Key)
		}
		fmt.Println("------------------------")
	}

	fmt.Println("\n======== Vars ========")
	for idx, blk := range r.Vars {
		fmt.Printf("Block %d LineStart:%d\n", idx+1, blk.StartLine)
		for _, cl := range blk.Lines {
			fmt.Printf("CODE: %q | COMMENT: %q | KEY: %s\n", cl.Code, cl.Comment, cl.Key)
		}
		fmt.Println("------------------------")
	}

	fmt.Println("\n======== Consts ========")
	for idx, blk := range r.Consts {
		fmt.Printf("Block %d LineStart:%d\n", idx+1, blk.StartLine)
		for _, cl := range blk.Lines {
			fmt.Printf("CODE: %q | COMMENT: %q | KEY: %s\n", cl.Code, cl.Comment, cl.Key)
		}
		fmt.Println("------------------------")
	}

	fmt.Println("\n======== Functions ========")
	for idx, blk := range r.Functions {
		fmt.Printf("Block %d LineStart:%d\n", idx+1, blk.StartLine)
		for _, cl := range blk.Lines {
			fmt.Printf("CODE: %q | COMMENT: %q | KEY: %s\n", cl.Code, cl.Comment, cl.Key)
		}
		fmt.Println("------------------------")
	}

	fmt.Println("\n======== Types ========")
	for name, info := range r.Types {
		fmt.Printf("Type[%s] LineStart:%d\n", name, info.StartLine)
		fmt.Println("--- Fields ---")
		for _, cl := range info.Lines {
			fmt.Printf("CODE: %q | COMMENT: %q | KEY: %s\n", cl.Code, cl.Comment, cl.Key)
		}
		if len(info.Methods) > 0 {
			fmt.Println("--- Methods ---")
			for idx, method := range info.Methods {
				fmt.Printf("Method %d LineStart:%d\n", idx+1, method.StartLine)
				for _, cl := range method.Lines {
					fmt.Printf("  CODE: %q | COMMENT: %q | KEY: %s\n", cl.Code, cl.Comment, cl.Key)
				}
			}
		}
		fmt.Println("------------------------")
	}
}

// DiffCodeLine 对比两行逻辑代码，注释不参与差异判断
func DiffCodeLine(a, b CodeLine) bool {
	return a.Code != b.Code
}

func CodeHash(code string) string {
	for k, v := range replaceChars {
		code = strings.Replace(code, k, v, -1)
	}
	hash := md5.Sum([]byte(code))
	return fmt.Sprintf("%x", hash)
}

// insertBeforeLast 在切片最后一个元素前面插入元素
// T any 支持所有类型：基础类型、结构体、自定义类型
func insertBeforeLast[T any](s []T, elems ...T) []T {
	// 空切片直接返回带插入元素的切片
	if len(s) == 0 {
		return elems
	}
	// 最后一个元素下标
	lastIdx := len(s) - 1
	// 预分配容量，减少扩容
	newSlice := make([]T, 0, len(s)+1)
	// 前半段：0 ~ lastIdx（不包含最后一个元素）
	newSlice = append(newSlice, s[:lastIdx]...)
	// 插入新元素
	newSlice = append(newSlice, elems...)
	// 拼接原来最后的元素，保证它永远在末尾
	newSlice = append(newSlice, s[lastIdx:]...)
	return newSlice
}
