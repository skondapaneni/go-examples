package parser

/**
 * parseContext.go
 *
 * @author skondapaneni
 *
 * parses a line in regular expression form and
 * creates a cmd expression tree
 *
 */
import (
	"collections"
	"fmt"
	"node"
	"os"
	"runtime"
)

type ParseContext struct {
	postFixList *collections.ArrayList
	opStack     *collections.Stack
	tokensMap   *TokensMap
	nodeCR      *node.GNodeRef
}

func NewParseContext() *ParseContext {
	return &ParseContext{
		postFixList: collections.NewArrayList(),
		opStack:     collections.NewStack(),
		tokensMap:   NewTokensMap(),
		nodeCR:      node.NewGNodeRef(node.NodeCR),
	}
}

/* We do this to save a value of different type's into a stack
 * by converting the actual value into an interface and saving the
 * interface object as a generic type.
 */
type StringInterface interface{}
type TokenInterface interface{}

func ConvertStringToInterface(token string) StringInterface {
	return StringInterface(token)
}

func InterfaceToString(item interface{}) string {
	switch v := item.(type) {
	case string:
		return v
	default:
		fmt.Printf("%T is of a type I don't know how to handle", item)
		pc, file, no, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("%s %s#%d\n", runtime.FuncForPC(pc).Name(), file, no)
		}
		os.Exit(203)
		return ""
	}
}

func ConvertTokenToInterface(token Token) TokenInterface {
	return TokenInterface(token)
}

func InterfaceToToken(item interface{}) Token {
	switch v := item.(type) {
	case Token:
		return v
	default:
		fmt.Printf("%T is of a type I don't know how to handle", item)
		pc, file, no, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("%s %s#%d\n", runtime.FuncForPC(pc).Name(), file, no)
		}
		os.Exit(203)
		return TOK_UNKNOWN
	}
}

func PrecedenceOf(tok string) int {
	i, ok := PrecedenceMap[tok]
	if ok {
		return i
	}
	return 6
}

func (self *ParseContext) AddImplicitConcatOp(last_op Token) {
	if last_op == TOK_GROUP_OPEN ||
		last_op == TOK_EXP_OPEN ||
		last_op == TOK_OPTIONAL_OPEN ||
		last_op == TOK_SELECT ||
		last_op == TOK_CONCAT {
		return
	}

	if !self.postFixList.IsEmpty() || !self.opStack.IsEmpty() {
		for !self.opStack.IsEmpty() &&
			PrecedenceOf(InterfaceToToken(
				self.opStack.Peek()).ToString()) >=
				PrecedenceOf(TOK_CONCAT.ToString()) {
			top, _ := self.opStack.Pop()
			self.postFixList.Append(top)
		}
		self.opStack.Push(ConvertTokenToInterface(TOK_CONCAT))
	}
}

func (self *ParseContext) CreatePostFix(exp string) error {
	scanner := NewScannerFromExp(exp, " \t\n\r\f|&().+*[]{}")
	fmt.Println("line is :", exp)

	self.postFixList.Clear()
	self.opStack.Clear()
	self.tokensMap.Clear()
	self.nodeCR = node.NewGNodeRef(node.NodeCR)

	last_op := TOK_UNKNOWN

	for scanner.hasMoreTokens() {
		nextTok, token := scanner.Scan()

		switch nextTok {
		case TOK_GROUP_OPEN, TOK_EXP_OPEN:
			// implicit AndToken -In A {B | C} == A.[B | C]
			self.AddImplicitConcatOp(last_op)

			// Push '(' as a marker
			self.opStack.Push(ConvertTokenToInterface(nextTok))
			last_op = nextTok

		case TOK_OPTIONAL_OPEN:
			// implicit AndToken -In A [B | C] == A.[B | C]
			self.AddImplicitConcatOp(last_op)

			// Push '[' as a marker
			self.opStack.Push(ConvertTokenToInterface(nextTok))
			last_op = nextTok

		case TOK_OPTIONAL_CLOSE:
			top, _ := self.opStack.Pop()
			topToken := InterfaceToToken(top)
			for topToken != TOK_OPTIONAL_OPEN {
				self.postFixList.Append(top)
				top, _ = self.opStack.Pop()
				topToken = InterfaceToToken(top)
			}
			self.postFixList.Append(ConvertTokenToInterface(TOK_QUESTION))
			last_op = TOK_UNKNOWN

		case TOK_GROUP_CLOSE, TOK_EXP_CLOSE:
			top, _ := self.opStack.Pop()
			topToken := InterfaceToToken(top)

			for topToken != TOK_GROUP_OPEN &&
				topToken != TOK_EXP_OPEN {
				self.postFixList.Append(top)
				top, _ = self.opStack.Pop()
				topToken = InterfaceToToken(top)
			}
			last_op = TOK_UNKNOWN

		case TOK_IDENT:
			self.AddImplicitConcatOp(last_op)
			fmt.Printf("token %v\n", token)
			self.postFixList.Append(ConvertStringToInterface(token))
			last_op = TOK_UNKNOWN

		case TOK_WS:

		default:
			for !self.opStack.IsEmpty() &&
				PrecedenceOf(InterfaceToToken(self.opStack.Peek()).ToString()) >=
					PrecedenceOf(nextTok.ToString()) {
				top, _ := self.opStack.Pop()
				self.postFixList.Append(top)
			}
			last_op = nextTok
			self.opStack.Push(ConvertTokenToInterface(nextTok))
		}
	}

	for !self.opStack.IsEmpty() {
		top, _ := self.opStack.Pop()
		self.postFixList.Append(top)
	}

	return nil
}

func (self *ParseContext) GetNode(n interface{}) node.Node {
	switch x := n.(type) {
	case Token:
		td := &TokenDesc{}
		td.SetName(x.ToString())
		return node.NewGNode(td)
	case string:
		td := &TokenDesc{}
		td.SetName(x)
		self.tokensMap.AddToken(td)
		return node.NewGNode(td)
	case node.Node:
		return x
	default:
		fmt.Println()
		fmt.Println("x is unknown type, return nil ", x)
		return nil
	}
}

func (self *ParseContext) GetToken(name string) (td *TokenDesc, ok bool) {
	return self.tokensMap.GetToken(name)
}

func (self *ParseContext) createBinaryExp(root node.Node, left interface{},
	right interface{}) {
	root.MergeNode(self.GetNode(left))
	root.MergeNode(self.GetNode(right))
}

func (self *ParseContext) CreateDAG(handler interface{}) (node.Node, error) {
	var done = false
	var p node.Node = nil

	for !done {
		iter := self.postFixList.Iterator()
		for v, ok, i := iter(); ok; v, ok, i = iter() {
			switch x := v.(type) {
			case Token:
				if x == TOK_CONCAT {
					p := node.NewGNodeAnd()
					self.createBinaryExp(p,
						self.postFixList.RemoveIndex(i-2), // left
						self.postFixList.RemoveIndex(i-2)) // right
					self.postFixList.ReplaceIndex(p, i-2) // replace '.'
				} else if x == TOK_SELECT {
					p := node.NewGNodeOr()
					self.createBinaryExp(p,
						self.postFixList.RemoveIndex(i-2),
						self.postFixList.RemoveIndex(i-2))
					self.postFixList.ReplaceIndex(p, i-2) // replace '|'
				} else if x == TOK_PLUS {
					p := self.GetNode(self.postFixList.RemoveIndex(i - 1))
					pref := node.NewGNodeRef(p)
					pref.SetRepeatable(true, 100)
					self.postFixList.ReplaceIndex(p.AddOptional(pref), i-1) // replace '+'
				} else if x == TOK_ASTERISK {
					p := self.GetNode(self.postFixList.RemoveIndex(i - 1))
					nor := node.Optional(p)
					nor.SetRepeatable(true, 100)
					self.postFixList.ReplaceIndex(nor, i-1) // replace '*'
				} else if x == TOK_QUESTION {
					p := self.GetNode(self.postFixList.RemoveIndex(i - 1))
					self.postFixList.ReplaceIndex(node.Optional(p), i-1) // replace '?'
				}
				ok = false
			case string:
				p := self.GetNode(v)
				self.postFixList.ReplaceIndex(p, i)
			case node.Node:
			default:
			}
			if ok == false {
				break
			}
		}
		done = (self.postFixList.Length() <= 2)
	}

	p = self.postFixList.GetItem(0).(node.Node).MergeNode(self.nodeCR)
	fmt.Println("p = ", self.postFixList.GetItem(0).(node.Node).ToString())
	self.nodeCR.SetCmdHandler(handler)
	return p, nil
}

func (self *ParseContext) Print() {
	iter := self.postFixList.Iterator()
	for v, ok, _ := iter(); ok; v, ok, _ = iter() {
		switch x := v.(type) {
		case Token:
			fmt.Printf("%v ", x.ToString())
		case string:
			fmt.Printf("%v ", x)
		case node.Node:
			fmt.Printf("%v ", x.ToString())
		default:
			fmt.Printf("%v ", x)
		}
	}
	fmt.Println()
}

func (self *ParseContext) PrintPostFix(exp string) error {
	err := self.CreatePostFix(exp)
	if err != nil {
		return err
	}
	self.Print()
	return nil
}

func GetFollow(list *collections.ArrayList) *collections.ArrayList {
	iterator := list.Iterator()
	nextList := collections.NewArrayList()

	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		nextList.Join(v.(node.Node).Follow())
	}
	return nextList
}
