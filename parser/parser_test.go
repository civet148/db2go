package parser

import (
	"fmt"
	"testing"
)

func TestParseGoFile(t *testing.T) {
	file1 := "files/users.go.db2go"
	ret1, err := ParseGoFile(file1)
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return
	}
	PrintResult(ret1)
	file2 := "files/users.go.work"
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	ret2, err := ParseGoFile(file2)
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return
	}
	PrintResult(ret2)

	code1 := "func (do *User) SetId(v uint64)           { do.Id = v }"
	code2 := `func (do *User) SetId(v uint64) {
		do.Id = v 
	}`
	hash1 := CodeHash(code1)
	hash2 := CodeHash(code2)
	fmt.Println("hash1:", hash1)
	fmt.Println("hash2:", hash2)

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	MergeCode(ret1, ret2)
}
