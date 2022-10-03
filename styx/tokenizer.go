package styx

import (
	"log"
	"os"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	Token_Unknown TokenType = iota
	Token_Identifier

	Token_Semicolon
	Token_Colon
	Token_Comma

	Token_CommentLine
	Token_CommentBlock
	Token_Whitespace

	Token_ParentheticalOpen
	Token_ParentheticalClose
	Token_BraceOpen
	Token_BraceClose

	Token_FeedRight
	Token_FeedLeft

	Token_StyxDirective
	Token_EndOfFile
)

type Token struct {
	Type TokenType
	Str  string

	Line uint64
}

type Tokenizer struct {
	file_data string

	offset      int
	next_offset int

	line_at uint64
}

func TokenizeFile(file string) Tokenizer {
	var file_data, err = os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return Tokenizer{file_data: string(file_data), line_at: 1}
}

func (tokens *Tokenizer) token_inc_comment_line() {
	for tokens.next_offset < len(tokens.file_data) &&
		tokens.file_data[tokens.next_offset] != '\n' {
		tokens.next_offset++
	}
	tokens.next_offset++
	tokens.line_at++
}

func (tokens *Tokenizer) token_inc_comment_block() {
	tokens.next_offset++
	for tokens.next_offset < len(tokens.file_data) {
		tokens.next_offset++
		if tokens.file_data[tokens.next_offset] == '\n' {
			tokens.line_at++
		}
		if tokens.file_data[tokens.next_offset] == '*' &&
			tokens.file_data[tokens.next_offset+1] == '/' {
			tokens.line_at += 2
			return
		}
	}
}

func (tokens *Tokenizer) token_inc_whitespace() {
	var whitespace_ch = [...]byte{' ', '\n', '\r', '\t'}
	if tokens.file_data[tokens.next_offset] == '\n' {
		tokens.line_at++
	}

	for tokens.next_offset < len(tokens.file_data) {
		tokens.next_offset++
		for _, ch := range whitespace_ch {
			if tokens.file_data[tokens.next_offset] == ch {
				return
			}
		}
	}

	// What exactly does this do again?
	if tokens.next_offset == '\n' {
		tokens.line_at++
	}
}

func (tokens *Tokenizer) token_inc_def() {
	var delimiters = [...]byte{
		' ', '\n', '\r', '\t',
		'(', ')', '{', '}',
		';', ':', ',',
	}

	for tokens.next_offset < len(tokens.file_data) {
		tokens.next_offset++
		for _, ch := range delimiters {
			if tokens.file_data[tokens.next_offset] == ch {
				return
			}
		}
	}
}

func (tokens *Tokenizer) token_get(token_type TokenType) (result string) {
	tokens.next_offset = tokens.offset

	switch token_type {
	case Token_Whitespace:
		tokens.token_inc_whitespace()
	case Token_CommentLine:
		tokens.token_inc_comment_line()
	case Token_CommentBlock:
		tokens.token_inc_comment_block()
	case Token_BraceOpen, Token_BraceClose, Token_ParentheticalOpen, Token_ParentheticalClose:
		tokens.next_offset++
	default:
		tokens.token_inc_def()
	}

	result = tokens.file_data[tokens.offset:tokens.next_offset]
	return
}

func (tokens *Tokenizer) GetAt() (tok Token) {
	tok.Line = tokens.line_at

	if tokens.offset >= len(tokens.file_data) {
		tok.Type = Token_EndOfFile
		return
	}

	switch tokens.file_data[tokens.offset] {
	case '@':
		tok.Type = Token_StyxDirective

	case ' ', '\t', '\r', '\n':
		tok.Type = Token_Whitespace

	case '/':
		if tokens.file_data[tokens.offset+1] == '/' {
			tok.Type = Token_CommentLine
		} else if tokens.file_data[tokens.offset+1] == '*' {
			tok.Type = Token_CommentBlock
		} else {
			tok.Type = Token_Identifier
		}

	case ';':
		tok.Type = Token_Semicolon
	case ':':
		tok.Type = Token_Colon
	case ',':
		tok.Type = Token_Comma

	case '(':
		tok.Type = Token_ParentheticalOpen
	case ')':
		tok.Type = Token_ParentheticalClose
	case '{':
		tok.Type = Token_BraceOpen
	case '}':
		tok.Type = Token_BraceClose

	case '<':
		if tokens.file_data[tokens.offset+1] == '-' {
			tok.Type = Token_FeedLeft
		} else {
			tok.Type = Token_Identifier
		}
	case '-':
		if tokens.file_data[tokens.offset+1] == '>' {
			tok.Type = Token_FeedRight
		} else {
			tok.Type = Token_Identifier
		}
	}

	if tokens.next_offset > tokens.offset {
		tok.Str = tokens.file_data[tokens.offset:tokens.next_offset]
	} else {
		tok.Str = tokens.token_get(tok.Type)
	}

	return
}

func (tokens *Tokenizer) IncAll() Token {
	tokens.offset = tokens.next_offset
	return tokens.GetAt()
}

func (tokens *Tokenizer) IncNoWhitespace() (tok Token) {
	for tok = tokens.IncAll(); tok.Type != Token_EndOfFile; tok = tokens.IncAll() {
		if tok.Type != Token_Whitespace && tok.Type != Token_CommentLine && tok.Type != Token_CommentBlock {
			break
		}
	}
	return
}

func (tokens *Tokenizer) IncNoComment() (tok Token) {
	for tok = tokens.IncAll(); tok.Type != Token_EndOfFile; tok = tokens.IncAll() {
		if tok.Type != Token_CommentLine && tok.Type != Token_CommentBlock {
			break
		}
	}
	return
}
