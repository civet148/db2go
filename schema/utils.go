package schema

import (
	"fmt"
	"github.com/civet148/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

type FieldStyle int

const (
	FieldStyle_Default    FieldStyle = iota // 默认跟数据库字段一致
	FieldStyle_SmallCamel                   // 小驼峰
	FieldStyle_BigCamel                     // 大驼峰
)

// hasGit 检查系统是否安装Git
func hasGit() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func writeToFile(strOutputPath, strBody string) (err error) {
	var file *os.File
	strBody += fmt.Sprintf("\n\n%s\n\n", CustomizeCodeTip)
	defer func() {
		defer file.Close()
		_ = exec.Command("gofmt", "-w", strOutputPath).Run() //格式化本地文件
	}()

	if !hasGit() || !isFileExist(strOutputPath) {
		// 文件不存在或本地没有git，以创建并覆盖方式生成新的文件
		file, err = os.OpenFile(strOutputPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return log.Errorf(err)
		}
		_, err = file.WriteString(strBody)
		return log.Errorf(err)
	}
	var dir = filepath.Dir(strOutputPath)
	var name = filepath.Base(strOutputPath)
	var datetime = time.Now().Format("20060102150405")

	baseFile := filepath.Join(dir, fmt.Sprintf(".%s-%s", datetime, name))
	blankFile := filepath.Join(dir, fmt.Sprintf(".%s-blank.go", datetime))
	bf, err := os.Create(blankFile)
	if err != nil {
		return log.Errorf(err)
	}
	defer func() {
		_ = bf.Close()
		_ = os.Remove(blankFile)
	}()

	file, err = os.OpenFile(baseFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return log.Errorf(err)
	}
	_, err = file.WriteString(strBody)
	if err != nil {
		return log.Errorf(err)
	}

	defer os.Remove(baseFile)

	_ = exec.Command("gofmt", "-w", baseFile).Run() //格式化新的数据模型副本

	cmd := exec.Command("git", "merge-file", "-q", strOutputPath, blankFile, baseFile)
	if err = cmd.Run(); err != nil {
		log.Errorf("file %s merge conflict occurred", strOutputPath)
		//err = gitMergeFile(strOutputPath)
		//if err != nil {
		//	log.Errorf("file %s merge conflict error: %s", strOutputPath, err.Error())
		//}
	}
	return nil
}

//	func gitMergeFile(strOutputPath string) (err error) {
//		/*
//			git config merge.tool vimdiff
//			git config mergetool.keepBackup false
//			git mergetool --no-prompt --tool=vimdiff output.go
//		*/
//		if hasGit() {
//			//_ = exec.Command("git", "config", "merge.tool", "vimdiff").Run()
//			//_ = exec.Command("git", "config", "mergetool.keepBackup", "false").Run()
//			_ = exec.Command("git", "mergetool", "--no-prompt", "--tool=vimdiff", strOutputPath).Run()
//		}
//		return nil
//	}
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
