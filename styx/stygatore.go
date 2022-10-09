package styx

import (
	"log"
)

type CompilationSettings struct {
	OutputName     string
	CommentKeyword string
}

func Generate(sym_table *SymbolTable, compile_settings CompilationSettings) (generated string) {
	generated += "// Courtesy of Stygatore\n\n"

	for _, ref := range sym_table.References {
		var template_symbol = sym_table.Lookup(ref.Identifier)
		var param_count = len(template_symbol.Parameters)
		var arg_count = len(ref.Args)

		if param_count != arg_count {
			log.Fatalf("Line %d: Number of arguments does not match number of parameters.\n\tExpected %d arguments, got %d.",
				ref.Line, param_count, arg_count)
		}

		{
			var comment = ref.Identifier
			comment += " -> "
			for i := 0; i < arg_count; i++ {
				comment += ref.Args[i]
				if i < arg_count-1 {
					comment += ", "
				}
			}

			comment += ": "
			comment += ref.GenName

			generated += "// " + comment + "\n"
		}

		for i := 0; i < len(template_symbol.Definition); i++ {
			var tok = template_symbol.Definition[i]
			if tok.Type == Token_StyxDirective {
				var directive = tok.Str[1:]
				if directive == "t_name" {
					generated += ref.GenName
				} else {
					var idx = 0
					for idx < len(template_symbol.Parameters) {
						if template_symbol.Parameters[idx] == directive {
							break
						}
						idx++
					}

					generated += ref.Args[idx]
				}

				continue
			}

			generated += tok.Str
		}

		generated += "\n\n"
	}

	return
}
