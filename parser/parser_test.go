package parser

import (
	"fmt"
	"testing"
)

func TestParseGoFile(t *testing.T) {
	file1 := "users.go.db2go"
	ret, err := ParseGoFile(file1)
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return
	}
	PrintResult(ret)
}
