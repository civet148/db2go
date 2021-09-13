package cmder

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CmdReader struct {
	RawInput   string
	TrimInputs []string
}

func readStdin(end byte) (r *CmdReader) {
	var err error
	r = &CmdReader{}
	reader := bufio.NewReader(os.Stdin)
	if r.RawInput, err = reader.ReadString(end); err != nil {
		fmt.Printf("ReadString error %s", err.Error())
		return nil
	}

	ss := strings.Split(r.RawInput, "\n")
	for _, v := range ss {
		s := strings.TrimSpace(v)
		if s == "" {
			continue //remove blank line
		}
		r.TrimInputs = append(r.TrimInputs, s)
	}
	return
}

//strPrompt 命令行提示符(prompt content)
//end 结束字符串(delimiter character)
func Prompt(strPrompt string, end byte) (r *CmdReader) {
	fmt.Printf("%s", strPrompt)
	return readStdin(end)
}
