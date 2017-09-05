package parser

import (
    "text/scanner"
    "strconv"
    "strings"
    "net"
)

type ArgType  int

const (
    ARG_TYPE_UNSET      ArgType  = 1 + iota
    ARG_TYPE_STRING
    ARG_TYPE_DATA
    ARG_TYPE_QUOTED_STRING    
    ARG_TYPE_INTEGER 
    ARG_TYPE_IPADDR 
    ARG_TYPE_IPPREFIX
    ARG_TYPE_INTERFACE
    ARG_TYPE_ETHERNET
    ARG_TYPE_UINTEGER
    ARG_TYPE_SID
    ARG_TYPE_NID
    ARG_TYPE_LID
    ARG_TYPE_NET
    ARG_TYPE_COMMUNITY
    ARG_TYPE_LONGLONG
    ARG_TYPE_INTERFACE_RANGE
    ARG_TYPE_PASSWORD
    ARG_TYPE_IPV6
    ARG_TYPE_IPV6_PREFIX
    ARG_TYPE_HEX
    ARG_TYPE_KEYWORD
    ARG_TYPE_COMMAND_CONT
    ARG_TYPE_GREP_EXP
    ARG_TYPE_BIT_EXP
    ARG_TYPE_OPERATOR
    ARG_TYPE_ANY 
)

var argTypeList = [...]string {
    "UNSET",
    "STRING",
    "QUOTED_STRING",  
    "INTEGER", 
    "IPADDR", 
    "IPPREFIX",
    "DATA",
    "INTERFACE",
    "ETHERNET",
    "UINTEGER",
    "SID",
    "NID",
    "LID",
    "NET",
    "COMMUNITY",
    "LONGLONG",
    "INTERFACE_RANGE",
    "PASSWORD",
    "IPV6",
    "IPV6_PREFIX",
    "HEX",
    "KEYWORD",
    "COMMAND_CONT",
    "GREP_EXP",
    "BIT_EXP",
    "OPERATOR",
    "ANY",
}

/**
 * func : String() function will return the string form for the
 * arg_t enumeration.
 */
func (arg_t ArgType) String() string {
   return argTypeList[arg_t-1]
}

type Arg struct {
    atype       ArgType
    label       string
    str_rep     string
    value       interface{} 
}

func NewArg(atype ArgType, label string,
    value interface{}, str_val string) *Arg {
    arg := &Arg{
        atype: atype,
        label: label,
        value: value,
        str_rep: str_val,
    }
    return arg
}

func (self *Arg) GetType() ArgType {
    return self.atype
}

func (self *Arg) GetLabel() string {
    return self.label
}

func (self *Arg) GetValue() interface{} {
    return self.value
}

func (self *Arg) GetStringRep() string {
    return self.str_rep
}

func (self *Arg) Compare(to *Arg) bool {
	return (self.atype == to.atype &&
		self.label == to.label)
}

/* A Token Description */
type TokenDesc struct {
    name      string
    tokenType ArgType
}

/**
  * func : NewTokenDesc
  *
  * Constructor for TokenDesc
  *
  * Arg1 : name
  * Arg2 : ttype  TokenType
  */
func NewTokenDesc(name string, tok rune, ttype ArgType) *TokenDesc {
    tokenItem := &TokenDesc{
        name:      name,
        tokenType: ttype,
    }
    return tokenItem
}

func (self *TokenDesc) GetName() string {
    return self.name
}

func (self *TokenDesc) SetName(name string) {
    self.name = name
}

func (self *TokenDesc) SetType(rt int) {
     self.SetTokenType(ArgType(rt))
}

func (self *TokenDesc) GetType() int {
     return int(self.GetTokenType())
}

func (self *TokenDesc) GetTokenType() ArgType {
    return self.tokenType
}

func (self *TokenDesc) SetTokenType(tokenType ArgType) *TokenDesc {
    self.tokenType = tokenType
    return self
}

//type LexInput struct {
//    pos  scanner.Position
//    tok  rune
//    lit  string
//}
//
//type LexOutput struct {
//    tdesc *TokenDesc
//    arg *Arg
//}

//func (self *TokenDesc) Eval(input interface{}, output interface{}) (status bool, err error) {
//
//     status, tdesc, arg := self.Match(input.(*LexInput).pos,
//			input.(*LexInput).tok,
//			input.(*LexInput).lit)
//
//     output.(*LexOutput).tdesc = tdesc
//     output.(*LexOutput).arg = arg
//     return status, nil
//}

/**
  * func : Match a lex token
  */
func (tdesc *TokenDesc) Match(pos scanner.Position, tok rune, lit string) (bool, *TokenDesc, *Arg) {
    switch tdesc.tokenType {
    case ARG_TYPE_KEYWORD:
        if tok == scanner.Ident {
            if (strings.Compare(tdesc.name, lit) == 0) {
            	return true, tdesc, NewArg(tdesc.tokenType, tdesc.name,
				    tdesc.name, tdesc.name)
            }
        }
        
    case ARG_TYPE_OPERATOR:
        if tok == scanner.Char {
            if (strings.Compare(tdesc.name, lit) == 0) {
            	return true, tdesc, NewArg(tdesc.tokenType, tdesc.name,
				    lit, lit)
            }
        }

    case ARG_TYPE_INTEGER :
        if tok == scanner.Int {
        	val, err := strconv.Atoi(lit)
        	if (err == nil) {
				return true, tdesc, NewArg(tdesc.tokenType, tdesc.name,
				    val, lit)
        	}
        }

    case ARG_TYPE_HEX :
        if tok == scanner.Int {
        	val, err := strconv.Atoi(lit)
        	if (err == nil) {
	            return true, tdesc, NewArg(tdesc.tokenType, tdesc.name,
				    val, lit)
        	}
        }

        if (tok == scanner.String) {
            val, err := strconv.ParseInt(lit, 0, 64)
            if (err == nil) {
                return true, tdesc, NewArg(tdesc.tokenType, tdesc.name,
				    val, lit)
            }
        }

    case ARG_TYPE_STRING:
        if tok == scanner.String || tok == scanner.Ident { 
            return true, tdesc, NewArg(tdesc.tokenType, tdesc.name,
				    lit, lit)
        }

    case ARG_TYPE_QUOTED_STRING:
        if tok == scanner.String { 
            return true, tdesc, NewArg(tdesc.tokenType, tdesc.name,
				    lit, lit)
        }

    case ARG_TYPE_IPADDR :
        if tok == scanner.String {
        	ip := net.ParseIP(strings.Trim(lit, "\" "))
            if ip != nil {
                return true, tdesc, NewArg(tdesc.tokenType, tdesc.name,
				   ip, lit)
            }
        }
    }
    return false, nil, nil
}
