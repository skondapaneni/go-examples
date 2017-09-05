package parser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

//The Go language defines the word rune as an alias for the type int32,
//so programs can be clear when an integer value represents a code point.

// Scanner represents a lexical scanner.
type Scanner struct {
	r          *bufio.Reader
	separators string
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader, sep string) *Scanner {
	return &Scanner{r: bufio.NewReader(r),
		separators: sep}
}

func NewScannerFromExp(exp string, sep string) *Scanner {
	return &Scanner{r: bufio.NewReader(strings.NewReader(exp)),
		separators: sep}
}

func (s *Scanner) hasMoreTokens() bool {
	ch := s.read()

	if ch == eof {
		return false
	}

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespace(ch) {
		s.unread()
		s.scanWhitespace()
		return s.hasMoreTokens()
	}

	s.unread()
	return true
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	// If we see a digit then consume as a number.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) || ch == '<' {
		s.unread()
		return s.scanIdent() // scan an identifier
	}

	if ch == eof {
		return TOK_EOF, ""
	}

	if s.separators == "" || strings.IndexRune(s.separators, ch) == -1 {
		return TOK_UNKNOWN, string(ch)
	}

	// Otherwise read the individual character.
	switch ch {
	case '*':
		return TOK_ASTERISK, string(ch)
	case ',':
		return TOK_COMMA, string(ch)
	case '{':
		return TOK_GROUP_OPEN, string(ch)
	case '}':
		return TOK_GROUP_CLOSE, string(ch)
	case '(':
		return TOK_EXP_OPEN, string(ch)
	case ')':
		return TOK_EXP_CLOSE, string(ch)
	case '[':
		return TOK_OPTIONAL_OPEN, string(ch)
	case ']':
		return TOK_OPTIONAL_CLOSE, string(ch)
	case '+':
		return TOK_PLUS, string(ch)
	case '.':
		return TOK_CONCAT, string(ch)
	case '|':
                next := s.read() 
                if (next != '|') {
                    s.unread()
		    return TOK_SELECT, string(ch)
                } 
 		return TOK_OR, "||"
        case  '&':
                next := s.read() 
                if (next != '&') {
                    s.unread()
		    return TOK_CONCAT, string(ch)
                } 
 		return TOK_AND, "&&"

	case ';':
		return TOK_EOC, string(ch)
	}

	return TOK_UNKNOWN, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			s.unread()
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return TOK_WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' && ch != '<' && ch != '>' {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	// Otherwise return as a regular identifier.
	return TOK_IDENT, buf.String()
}

/*
func (s *Scanner) scanMantissa(base int) {
    for digitVal(s.ch) < base {
        s.next()
    }
}

func (s *Scanner) scanNumber(seenDecimalPoint bool) (Token, string) {
    // digitVal(s.ch) < 10
    offs := s.offset
    tok := token.INT

    if seenDecimalPoint {
        offs--
        tok = token.FLOAT
        s.scanMantissa(10)
        goto exponent
    }

    if s.ch == '0' {
        // int or float
        offs := s.offset
        s.next()
        if s.ch == 'x' || s.ch == 'X' {
            // hexadecimal int
            s.next()
            s.scanMantissa(16)
            if s.offset-offs <= 2 {
                // only scanned "0x" or "0X"
                s.error(offs, "illegal hexadecimal number")
            }
        } else {
            // octal int or float
            seenDecimalDigit := false
            s.scanMantissa(8)
            if s.ch == '8' || s.ch == '9' {
                // illegal octal int or float
                seenDecimalDigit = true
                s.scanMantissa(10)
            }
            if s.ch == '.' || s.ch == 'e' || s.ch == 'E' || s.ch == 'i' {
                goto fraction
            }
            // octal int
            if seenDecimalDigit {
                s.error(offs, "illegal octal number")
            }
        }
        goto exit
    }

    // decimal int or float
    s.scanMantissa(10)

  fraction:
    if s.ch == '.' {
        tok = token.FLOAT
        s.next()
        s.scanMantissa(10)
    }

  exponent:
    if s.ch == 'e' || s.ch == 'E' {
        tok = token.FLOAT
        s.next()
        if s.ch == '-' || s.ch == '+' {
            s.next()
        }
        if digitVal(s.ch) < 10 {
            s.scanMantissa(10)
        } else {
            s.error(offs, "illegal floating-point exponent")
        }
    }

    if s.ch == 'i' {
        tok = token.IMAG
        s.next()
    }

  exit:
    return tok, string(s.src[offs:s.offset])
}
*/

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}

// isWhitespace returns true if the rune is a space, tab, or newline.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
