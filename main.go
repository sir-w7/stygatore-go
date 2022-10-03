package main

import (
	"fmt"
	"os"
	"time"

	"stygatore/styx"
)

func GenerateFile(file string) {
	var file_info = styx.QueryFileInfo(file)
	defer styx.Profile(time.Now(), file_info.Filename)

	var tokens = styx.TokenizeFile(file)
	for tok := tokens.GetAt(); tok.Type != styx.Token_EndOfFile; tok = tokens.IncNoWhitespace() {
		fmt.Printf("%+v\n", tok)
	}
}

func main() {
	var args = os.Args
	if len(args) == 1 {
		fmt.Println("stygatore is a sane, performant metaprogramming tool",
			"for language-agnostic generics with readable diagnostics",
			"for maximum developer productivity.")
		fmt.Println("Usage:", args[0], "[files/directories]")
		return
	}

	var fileList = styx.GetFileList(args[1:], "styx")
	for _, file := range fileList {
		GenerateFile(file)
	}
}
