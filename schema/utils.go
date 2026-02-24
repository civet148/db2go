package schema

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/civet148/log"
)

type FieldStyle int

const (
	FieldStyle_Default    FieldStyle = iota // 默认跟数据库字段一致
	FieldStyle_SmallCamel                   // 小驼峰
	FieldStyle_BigCamel                     // 大驼峰
)

func writeToFile(strOutputPath, strBody string) (err error) {
	var file *os.File
	strBody += fmt.Sprintf("\n\n%s\n\n", CustomizeCodeTip)
	defer func() {
		defer file.Close()
		_ = command("gofmt", "-w", strOutputPath) //格式化本地文件
	}()

	// 文件不存在或本地没有git，以创建并覆盖方式生成新的文件
	file, err = os.OpenFile(strOutputPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return log.Errorf(err.Error())
	}
	_, err = file.WriteString(strBody)
	if err != nil {
		return log.Errorf("write file [%v] error (%v)", strOutputPath, err.Error())
	}
	return nil
}

// hasGit 检查系统是否安装git并且本地存在git仓库
func hasGit() bool {
	_, err := exec.LookPath("git")
	if err != nil {
		return false
	}
	err = command("git", "status")
	if err != nil {
		return false
	}
	return true
}

func command(name string, args ...string) (err error) {
	var prints []string
	prints = append(prints, name)
	prints = append(prints, args...)
	msg := strings.Join(prints, " ")
	log.Infof("[%v]", msg)
	out, err := exec.Command(name, args...).CombinedOutput()
	var strOutput = string(out)
	if err != nil {
		log.Printf(strOutput)
		return err
	}
	if len(strOutput) > 0 {
		fmt.Println(strOutput)
	}
	return nil
}

func gitCheckout() (err error) {
	_ = command("sh", "-c", "git branch -D db2go 2>/dev/null")
	err = command("sh", "-c", "git checkout -b db2go 2>/dev/null || git checkout db2go")
	if err != nil {
		return log.Errorf("git checkout db2go branch error: %v", err.Error())
	}
	return nil
}

func gitCommit() (err error) {
	var now = time.Now().Format(time.DateTime)
	var commitMsg = fmt.Sprintf("db2go export database models at %s", now)
	_ = command("sh", "-c", fmt.Sprintf("git add -A && git commit -am '%s' 2>/dev/null", commitMsg))
	return nil
}

func gitCheckoutBack() (err error) {
	return command("git", "checkout", "-")
}

func gitMerge() (err error) {
	return command("git", "merge", "db2go")
}

func gitReset() (err error) {
	return command("git", "reset", "--hard", "HEAD")
}

func gitCommitAndMerge() (err error) {
	if err = gitCommit(); err != nil {
		_ = gitReset()        //回滚本地代码
		_ = gitCheckoutBack() //切换到上个分支
		return err
	}
	if err = gitCheckoutBack(); err != nil {
		return err
	}
	if err = gitMerge(); err != nil {
		return err
	}
	return nil
}

// hasUnstagedChanges 检测执行git status命令行是否存在未暂存的代码变更
// 如果没有git或不在git仓库中，直接返回false
func hasUnstagedChanges() (ok bool, err error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		// 如果命令执行失败（例如没有安装git或不在git仓库中），返回false
		return false, err
	}

	// 解析输出，如果存在非空输出，说明有未暂存的变更
	statusOutput := string(output)
	if strings.TrimSpace(statusOutput) == "" {
		return false, nil
	}

	return true, nil
}

func isFileExist(strFilePath string) bool {
	_, err := os.Stat(strFilePath)
	return err == nil
}

func FieldStyleFromString(s string) FieldStyle {
	switch strings.ToLower(s) {
	case "small_camel", "smallcamel", "small", "s":
		return FieldStyle_SmallCamel
	case "big_camel", "bigcamel", "big", "b":
		return FieldStyle_BigCamel
	}
	return FieldStyle_Default
}

// ConvertFieldStyle 数据库字段名风格转换
func ConvertFieldStyle(strColName string, style FieldStyle) string {
	if style == FieldStyle_Default || strColName == "" {
		return strColName
	}

	// 分割单词（支持蛇形命名法和连字符命名法）
	var words []string
	wordStart := 0

	for i, ch := range strColName {
		if ch == '_' || ch == '-' {
			if i > wordStart {
				words = append(words, strColName[wordStart:i])
			}
			wordStart = i + 1
		} else if i == len(strColName)-1 {
			words = append(words, strColName[wordStart:])
		}
	}

	// 如果没有分隔符，直接返回原字符串
	if len(words) == 0 {
		return strColName
	}

	// 根据风格转换
	var result string
	switch style {
	case FieldStyle_SmallCamel:
		// 小驼峰：第一个单词全小写，后续单词首字母大写
		for i, word := range words {
			if i == 0 {
				result += strings.ToLower(word)
			} else {
				result += capitalizeFirst(word)
			}
		}

	case FieldStyle_BigCamel:
		// 大驼峰：所有单词首字母大写
		for _, word := range words {
			result += capitalizeFirst(word)
		}
	}

	return result
}

// capitalizeFirst 将单词首字母大写，其他字母小写
func capitalizeFirst(word string) string {
	if word == "" {
		return word
	}

	// 处理全大写的情况（如ID、URL等缩写）
	if isAllUpper(word) {
		return word
	}

	// 首字母大写，其他字母小写
	return strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
}

// isAllUpper 检查字符串是否全大写
func isAllUpper(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return s != ""
}

func BigCamelCase(strIn string) (strOut string) {
	var idxUnderLine = int(-1)
	for i, v := range strIn {
		strChr := string(v)
		if i == 0 {
			strOut += strings.ToUpper(strChr)
		} else {
			if v == '_' {
				idxUnderLine = i //ignore
			} else {
				if i == idxUnderLine+1 {
					strOut += strings.ToUpper(strChr)
				} else {
					strOut += strChr
				}
			}
		}
	}
	return strOut
}

func SmallCamelCase(strIn string) (strOut string) {
	var idxUnderLine = int(-1)
	for i, v := range strIn {
		strChr := string(v)
		if i == 0 {
			strOut += strings.ToLower(strChr)
		} else {
			if v == '_' {
				idxUnderLine = i //ignore
			} else {
				if i == idxUnderLine+1 {
					strOut += strings.ToUpper(strChr)
				} else {
					strOut += strChr
				}
			}
		}
	}
	return
}

func Split(s string) (ss []string) {
	if strings.Contains(s, ",") {
		ss = strings.Split(s, ",")
	} else {
		ss = strings.Split(s, ";")
	}
	return ss
}

func TrimSpaceSlice(s []string) (ts []string) {
	for _, v := range s {
		ts = append(ts, strings.TrimSpace(v))
	}
	return
}

func GetDatabaseName(strPath string) (strName string) {
	idx := strings.LastIndex(strPath, "/")
	if idx == -1 {
		return
	}
	return strPath[idx+1:]
}

func MakeDir(strDir string) (err error) {
	if !isFileExist(strDir) {
		err = os.MkdirAll(strDir, os.ModePerm)
		if err != nil {
			return log.Errorf("make directory error: %v", err.Error())
		}
		log.Info("make directory [%v] successful", strDir)
	}
	return nil
}

func handleColumnComment(comment string) string {
	comment = strings.ReplaceAll(comment, "\r", "")  //回车
	comment = strings.ReplaceAll(comment, "\n", "")  //换行
	comment = strings.ReplaceAll(comment, ":", "：")  //英文冒号替换成中文冒号
	comment = strings.ReplaceAll(comment, ";", "；")  //英文分号替换成中文分号
	comment = strings.ReplaceAll(comment, "\"", "‘") //英文双引号替换中文单引号
	comment = strings.ReplaceAll(comment, "'", "")   //英文单引号替换中文单引号
	return comment
}

// 特殊复数→单数映射（覆盖所有不规则/特殊情况）
var pluralToSingular = map[string]string{
	"boxes":      "box",
	"bases":      "base",
	"cases":      "case",
	"children":   "child",
	"feet":       "foot",
	"geese":      "goose",
	"teeth":      "tooth",
	"mice":       "mouse",
	"men":        "man",
	"women":      "woman",
	"people":     "person",
	"knives":     "knife",
	"lives":      "life",
	"wives":      "wife",
	"leaves":     "leaf",
	"loaves":     "loaf",
	"potatoes":   "potato",
	"tomatoes":   "tomato",
	"phenomena":  "phenomenon",
	"profiles":   "profile",
	"data":       "data",
	"datas":      "data",
	"media":      "media",
	"quizzes":    "quiz",
	"analyses":   "analysis",
	"cities":     "city",
	"categories": "category",
	"classes":    "class",
	"buses":      "bus",
	"roles":      "role",
	"phases":     "phase",
	"houses":     "house",
	"radius":     "radius",
}

// 不可数/单复数同形单词
var uncountableWords = map[string]bool{
	"news": true,
	"is":   true,
	"as":   true,
}

// 驼峰转下划线正则
var camelRegex = regexp.MustCompile(`([a-z0-9])([A-Z])`)

// TableNameToStructName 表名转结构体名（终极极简版）
func TableNameToStructName(tableName string) string {
	if tableName == "" {
		return ""
	}

	// 步骤1：移除前缀
	lowerName := strings.ToLower(tableName)
	cleaned := tableName
	switch {
	case strings.HasPrefix(lowerName, "tbl_"):
		cleaned = tableName[4:]
	case strings.HasPrefix(lowerName, "tb_"):
		cleaned = tableName[3:]
	case strings.HasPrefix(lowerName, "t_"):
		cleaned = tableName[2:]
	}

	// 步骤2：标准化（驼峰转下划线+小写）
	standardized := camelRegex.ReplaceAllString(cleaned, `${1}_${2}`)
	words := strings.Split(strings.ToLower(standardized), "_")
	// 过滤空单词
	filtered := make([]string, 0, len(words))
	for _, w := range words {
		if w != "" {
			filtered = append(filtered, w)
		}
	}
	if len(filtered) == 0 {
		return ""
	}

	// 步骤3：处理最后一个单词的复数（核心修复：直接用字符串函数，不用rune）
	lastWord := filtered[len(filtered)-1]
	// 优先特殊映射
	if singular, ok := pluralToSingular[lastWord]; ok {
		filtered[len(filtered)-1] = singular
	} else if !uncountableWords[lastWord] {
		// 极简复数规则：只处理最常见场景，避免截取错误
		switch {
		case strings.HasSuffix(lastWord, "ies"):
			filtered[len(filtered)-1] = strings.TrimSuffix(lastWord, "ies") + "y"
		case strings.HasSuffix(lastWord, "es"):
			filtered[len(filtered)-1] = strings.TrimSuffix(lastWord, "es")
		case strings.HasSuffix(lastWord, "s"):
			filtered[len(filtered)-1] = strings.TrimSuffix(lastWord, "s")
		}
	}

	// 步骤4：转大驼峰
	var result strings.Builder
	for _, w := range filtered {
		result.WriteString(strings.ToUpper(w[:1]) + w[1:])
	}
	return result.String()
}
