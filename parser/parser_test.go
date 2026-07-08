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
	fmt.Println("--------------------------------- base ---------------------------------")
	PrintResult(ret1)

	file2 := "files/users.go.work"
	ret2, err := ParseGoFile(file2)
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return
	}
	fmt.Println("--------------------------------- work ---------------------------------")
	PrintResult(ret2)

	var merge string
	merge, err = MergeCode(ret1, ret2, true)
	if err != nil {
		fmt.Printf("merge error: %v\n", err)
		return
	}
	fmt.Println("--------------------------------- merge ---------------------------------")
	fmt.Println(merge)
}
