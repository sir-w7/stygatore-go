package styx

import "fmt"

type Symbol interface {
	Print()
}

type Declaration struct {
	Identifier string
	Line       uint64

	Parameters []string
	Definition []Token
}

type Reference struct {
	Identifier string
	Line       uint64

	Args    []string
	GenName string
}

type SymbolTable struct {
	Declarations map[string]Declaration
	References   []Reference
}

func NewDeclaration(tokens *Tokenizer, identifier string, line uint64) (decl Declaration) {
	decl.Identifier = identifier
	decl.Line = line

	var tok = tokens.IncNoWhitespace()
	if tok.Type == Token_ParentheticalOpen {
		tok = tokens.IncNoWhitespace()
		for tok.Type != Token_ParentheticalClose {
			if tok.Type != Token_Comma {
				decl.Parameters = append(decl.Parameters, tok.Str)
			}
			tok = tokens.IncNoWhitespace()
		}
	} else {
		decl.Parameters = append(decl.Parameters, tok.Str)
	}

	var prev_state = NewTokenizerState(tokens)
	defer prev_state.Restore()

	tok = tokens.IncNoWhitespace()
	for !tok.KnownStyxDirective() && tok.Type != Token_EndOfFile {
		decl.Definition = append(decl.Definition, tok)
		prev_state = NewTokenizerState(tokens)

		tok = tokens.IncNoComment()
	}

	var last_idx = len(decl.Definition) - 1
	for decl.Definition[last_idx].Type == Token_Whitespace {
		last_idx--
	}

	decl.Definition = decl.Definition[:last_idx+1]

	return
}

func NewReference(tokens *Tokenizer, identifier string, line uint64) (ref Reference) {
	ref.Identifier = identifier
	ref.Line = line

	for tok := tokens.IncNoWhitespace(); tok.Type != Token_Colon; tok = tokens.IncNoWhitespace() {
		if tok.Type == Token_Comma {
			continue
		}

		ref.Args = append(ref.Args, tok.Str)
	}

	var tok = tokens.IncNoWhitespace()
	ref.GenName = tok.Str

	return
}

func (decl Declaration) Print() {
	if len(decl.Identifier) == 0 {
		fmt.Println("sym: (null)")
		return
	}

	fmt.Println("type: Declaration")
	fmt.Printf("%+v\n", decl)
}

func (ref Reference) Print() {
	if len(ref.Identifier) == 0 {
		fmt.Println("sym: (null)")
		return
	}
	fmt.Println("type: Reference")
	fmt.Printf("%+v\n", ref)
}

func ParseNext(tokens *Tokenizer) Symbol {
	var tok = tokens.IncNoWhitespace()

	var tok_identifier = tok.Str
	var tok_line = tok.Line

	tok = tokens.IncNoWhitespace()
	if tok.Type == Token_FeedLeft {
		var decl = NewDeclaration(tokens, tok_identifier, tok_line)
		return decl
	} else {
		var ref = NewReference(tokens, tok_identifier, tok_line)
		return ref
	}
	// error handling
}

func NewSymbolTable() (sym_table SymbolTable) {
	sym_table.Declarations = make(map[string]Declaration)
	return
}

func (sym_table *SymbolTable) Push(sym Symbol) {
	decl, ok := sym.(Declaration)
	if ok {
		sym_table.Declarations[decl.Identifier] = decl
		return
	}

	ref, ok := sym.(Reference)
	if ok {
		sym_table.References = append(sym_table.References, ref)
		return
	}
}

func (sym_table *SymbolTable) Lookup(identifier string) Declaration {
	return sym_table.Declarations[identifier]
}
