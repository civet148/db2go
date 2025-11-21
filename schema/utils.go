package schema

import (
	"strings"
	"unicode"
)

type FieldStyle int

const (
	FieldStyle_Default    FieldStyle = iota // 默认跟数据库字段一致
	FieldStyle_SmallCamel                   // 小驼峰
	FieldStyle_BigCamel                     // 大驼峰
)

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
