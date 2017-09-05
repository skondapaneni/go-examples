package main

import (
	"strings"
	"testing"
)

// Ensure the scanner can scan tokens correctly.
func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok Token
		lit string
	}{
		// Special tokens (EOF, ILLEGAL, WS)
		{s: ``, tok: TOK_EOF},
		{s: `#`, tok: TOK_ILLEGAL, lit: `#`},
		{s: ` `, tok: TOK_WS, lit: " "},
		{s: "\t", tok: TOK_WS, lit: "\t"},
		{s: "\n", tok: TOK_WS, lit: "\n"},

		// Misc characters
		{s: `*`, tok: TOK_ILLEGAL, lit: "*"},

		// Identifiers
		{s: `foo`, tok: TOK_IDENT, lit: `foo`},
		{s: `Zx12_3U_-`, tok: TOK_IDENT, lit: `Zx12_3U_`},

		// Keywords
		//{s: `FROM`, tok: FROM, lit: "FROM"},
		//{s: `SELECT`, tok: SELECT, lit: "SELECT"},
	}

	for i, tt := range tests {
		s := NewScanner(strings.NewReader(tt.s), " \t\n\r\f|&()+*[]{}")
		tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
