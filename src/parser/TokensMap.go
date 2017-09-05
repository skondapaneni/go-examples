package parser

/* A map of named token list */
type TokensMap struct {
    tokens map[string]*TokenDesc
}

/* Add a token desc to the tokensMap */
func (tm *TokensMap) AddToken(tdesc *TokenDesc) *TokensMap {
    tm.tokens[tdesc.name] = tdesc
    return tm
}

func (tm *TokensMap) GetToken(name string) (td *TokenDesc, ok bool) {
    td, ok = tm.tokens[name]
    return
}

func (tm *TokensMap) Clear() {
    tm.tokens = make(map[string]*TokenDesc)
}

func NewTokensMap() *TokensMap {
    tokensMap := &TokensMap{
        tokens: make(map[string]*TokenDesc),
    }
    return tokensMap
}
