package cli

import (
	"testing"
        "go/scanner"
        "go/token"
)

// Ensure the scanner can scan tokens correctly.
func TestGoScanner_Scan(t *testing.T) {
	var tests = []struct {
		s   string
		tok  token.Token
		lit string
	}{
		// Special tokens (EOF, ILLEGAL, WS)
		{s: ``, tok: token.EOF},
		{s: `#`, tok: token.ILLEGAL, lit: `#`},

		// Misc characters
		{s: `*`, tok: token.MUL, lit: "*"},
		{s: `]`, tok: token.RBRACK, lit: "]"},

		// Identifiers
		{s: `foo`, tok: token.IDENT, lit: `foo`},
		{s: `Zx12_3U_-`, tok: token.IDENT, lit: `Zx12_3U_`},

		// Keywords
		//{s: `FROM`, tok: FROM, lit: "FROM"},
		//{s: `SELECT`, tok: SELECT, lit: "SELECT"},
	}

	for i, tt := range tests {
                var s scanner.Scanner
                fset := token.NewFileSet()  // positions are relative to fset
                file := fset.AddFile("", fset.Base(), len(tt.s)) //register input "file"
                s.Init(file, []byte(tt.s), nil /* no error handler */, scanner.ScanComments)

		_, tok, lit := s.Scan()
		if tt.tok != tok {
			t.Errorf("%d. %q token mismatch: exp=%q got=%q <%q>", i, tt.s, tt.tok, tok, lit)
		} else if tt.lit != lit {
			t.Errorf("%d. %q literal mismatch: exp=%q got=%q", i, tt.s, tt.lit, lit)
		}
	}
}
