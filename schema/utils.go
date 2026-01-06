package schema

import (
	"fmt"
	"os"
	"os/exec"
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
		return log.Errorf(err)
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
	err = command("git", "stash")
	if err != nil {
		return log.Errorf("git stash error: %v", err.Error())
	}
	defer func() {
		if err != nil {
			_ = gitStashPop() //命令行执行错误,提前恢复本地变更代码
		}
	}()

	err = command("sh", "-c", "git checkout -b db2go 2>/dev/null || git checkout db2go")
	if err != nil {
		return log.Errorf("git checkout db2go branch error: %v", err.Error())
	}
	return nil
}

func gitCommit() (err error) {
	var now = time.Now().Format(time.DateTime)
	var commitMsg = fmt.Sprintf("db2go export data models at %s", now)
	_ = command("sh", "-c", fmt.Sprintf("git add -A && git commit -am '%s' 2>/dev/null", commitMsg))
	return nil
}

func gitStashPop() (err error) {
	_ = command("sh", "-c", "git stash pop 2>/dev/null")
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
	defer func() {
		if err != nil {
			_ = gitReset()        //回滚本地代码
			_ = gitCheckoutBack() //回到上一个分支
			_ = gitStashPop()     //恢复本地变更代码
		}
	}()
	if err = gitCommit(); err != nil {
		return err
	}
	if err = gitCheckoutBack(); err != nil {
		return err
	}
	if err = gitMerge(); err != nil {
		return err
	}
	if err = gitStashPop(); err != nil { //恢复本地变更代码
		return err
	}
	return nil
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
