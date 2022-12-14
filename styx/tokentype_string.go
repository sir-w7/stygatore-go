// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package styx

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Token_Unknown-0]
	_ = x[Token_Identifier-1]
	_ = x[Token_Semicolon-2]
	_ = x[Token_Colon-3]
	_ = x[Token_Comma-4]
	_ = x[Token_CommentLine-5]
	_ = x[Token_CommentBlock-6]
	_ = x[Token_Whitespace-7]
	_ = x[Token_ParentheticalOpen-8]
	_ = x[Token_ParentheticalClose-9]
	_ = x[Token_BraceOpen-10]
	_ = x[Token_BraceClose-11]
	_ = x[Token_FeedRight-12]
	_ = x[Token_FeedLeft-13]
	_ = x[Token_StyxDirective-14]
	_ = x[Token_EndOfFile-15]
}

const _TokenType_name = "Token_UnknownToken_IdentifierToken_SemicolonToken_ColonToken_CommaToken_CommentLineToken_CommentBlockToken_WhitespaceToken_ParentheticalOpenToken_ParentheticalCloseToken_BraceOpenToken_BraceCloseToken_FeedRightToken_FeedLeftToken_StyxDirectiveToken_EndOfFile"

var _TokenType_index = [...]uint16{0, 13, 29, 44, 55, 66, 83, 101, 117, 140, 164, 179, 195, 210, 224, 243, 258}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
