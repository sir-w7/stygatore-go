package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"stygatore/styx"
)

var wg sync.WaitGroup

func GenerateFile(file string) {
	defer wg.Done()

	var file_info = styx.QueryFileInfo(file)

	defer styx.Profile(time.Now(), file_info.Filename)

	var compile_settings styx.CompilationSettings
	var tokens = styx.TokenizeFile(file)
	var sym_table = styx.NewSymbolTable()

	for tok := tokens.GetAt(); tok.Type != styx.Token_EndOfFile; tok = tokens.IncNoWhitespace() {
		if tok.Type == styx.Token_StyxDirective {
			switch tok.Str {
			case "@template":
				var symbol = styx.ParseNext(&tokens)
				sym_table.Push(symbol)
			case "@output":
				tok = tokens.IncNoWhitespace()
				compile_settings.OutputName = tok.Str
			}
		}
	}

	var generated = styx.Generate(&sym_table, compile_settings)
	var output_path = ""

	if len(compile_settings.OutputName) != 0 {
		output_path = file_info.WorkingDir + "/" + compile_settings.OutputName
	} else {
		output_path = file_info.WorkingDir + "/" + file_info.BaseName + ".h"
	}

	err := os.WriteFile(output_path, []byte(generated), 0644)

	if err != nil {
		log.Fatal(err)
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
		wg.Add(1)
		go GenerateFile(file)
	}

	wg.Wait()
}
