package parser

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	TOK_UNKNOWN Token = iota
	TOK_EOF
	TOK_WS

	// Literals
	TOK_IDENT // main

	// Misc characters
	TOK_ASTERISK // *
	TOK_COMMA    // ,

	TOK_COMMAND_SEP
	TOK_COMMAND_CONT
	TOK_GROUP_OPEN     // {
	TOK_GROUP_CLOSE    // }
	TOK_EXP_OPEN       // (
	TOK_EXP_CLOSE      // )
	TOK_OPTIONAL_OPEN  // [
	TOK_OPTIONAL_CLOSE // ]
	TOK_SELECT         // |
	TOK_CONCAT         // . or &
        TOK_AND            // &&
        TOK_OR             // ||
	TOK_PLUS           // +
	TOK_QUESTION       // ?
	TOK_PIPE           // |
	TOK_EOC            // ;
)

var PrecedenceMap = map[string]int{
	"(": 1,
	"{": 1,
	"[": 1,
	"|": 2, //alternate
	".": 3, //concatenate
        "||":2, // or
        "&&":3, // and
	"?": 4, // zero or one
	"*": 4, // zero or more
	"+": 4, // one or more
	"^": 5, // complement
	// else 6
}

func (t Token) ToString() string {
	switch t {
	case TOK_UNKNOWN:
		return "tok_unknown"
	case TOK_EOF:
		return "eof"
	case TOK_WS:
		return "ws"
	case TOK_IDENT:
		return "ident"
	case TOK_ASTERISK:
		return "*"
	case TOK_COMMA:
		return ","
	case TOK_COMMAND_SEP:
		return "SEP"
	case TOK_COMMAND_CONT:
		return "CONT"
	case TOK_GROUP_OPEN:
		return "{"
	case TOK_EXP_OPEN:
		return "("
	case TOK_OPTIONAL_OPEN:
		return "["
	case TOK_OPTIONAL_CLOSE:
		return "]"
	case TOK_SELECT:
		return "|"
	case TOK_CONCAT:
		return "."
	case TOK_OR:
		return "||"
	case TOK_AND:
		return "&&"
	case TOK_PLUS:
		return "+"
	case TOK_QUESTION:
		return "?"
	case TOK_PIPE:
		return "|"
	case TOK_EOC:
		return ";"
	default:
		return "unknown"
	}
}
