package cli

import (
	"bytes"
	"collections"
	"fmt"
	"node"
        "parser"
	"strings"
	"text/scanner"
	"sort"
)

type Cmd struct {
	pc       *parser.ParseContext
	dpt       node.Node
}

var cmds = make(map[string]*Cmd)
var cmd_g *Cmd = nil

func NewCmd() *Cmd {
	return &Cmd{
		pc:  parser.NewParseContext(),
		dpt: nil,
	}
}

type cmdHandler func(args map[string]*collections.ArrayList)

func (cmd *Cmd) parseSyntax(line string, handler cmdHandler) error {
	err := cmd.pc.CreatePostFix(line)
	if err != nil {
		return err
	}
	cmd.dpt, err = cmd.pc.CreateDAG(handler)
	fmt.Printf("node = %v\n", cmd.dpt.ToString())
	return nil
}

func (cmd *Cmd) ParseSyntax(line string, handler cmdHandler) error {
	return cmd.parseSyntax(line, handler)
}

func (self *Cmd) setTokenType(name string, ttype parser.ArgType) {
	td, ok := self.pc.GetToken(name)
	if ok {
		td.SetTokenType(ttype)
	}
}

func (self *Cmd) defineKeyword(name string) {
	self.setTokenType(name, parser.ARG_TYPE_KEYWORD)
}

func (self *Cmd) defineString(name string) {
	self.setTokenType(name, parser.ARG_TYPE_STRING)
}

func (self *Cmd) defineQuotedString(name string) {
	self.setTokenType(name, parser.ARG_TYPE_QUOTED_STRING)
}

func (self *Cmd) defineInteger(name string) {
	self.setTokenType(name, parser.ARG_TYPE_INTEGER)
}

func (self *Cmd) defineHex(name string) {
	self.setTokenType(name, parser.ARG_TYPE_HEX)
}

func (self *Cmd) defineUInteger(name string) {
	self.setTokenType(name, parser.ARG_TYPE_UINTEGER)
}

func (self *Cmd) defineInterfaceType(name string) {
	self.setTokenType(name, parser.ARG_TYPE_INTERFACE)
}

func syntax_processor(args map[string]*collections.ArrayList) {
	var buffer bytes.Buffer
	if v, found := args["<syntax>"]; found {
		iterator := v.Iterator()
		x := ""
		for s, ok, _ := iterator(); ok; s, ok, _ = iterator() {
			buffer.WriteString(x)
			buffer.WriteString(s.(string))
			x = " "
		}
		cmd_line := buffer.String()
		fmt.Println("syntax_processor - cmd_line = " + cmd_line)
		cmd_g.parseSyntax(cmd_line, nil)
	}
}

func keyword_processor(args map[string]*collections.ArrayList) {
	if v, found := args["<label>"]; found {
		fmt.Println("label = ", v.GetItems())
		cmd_g.defineKeyword(args["<label>"].GetItem(0).(string))
	}

	if v, found := args["<help>"]; found {
		fmt.Println("help = ", v.GetItems())
	}
}

func cmd_processor(args map[string]*collections.ArrayList) {
	if v, found := args["<cmd_name>"]; found {
		fmt.Println("cmd_name = ", v.GetItems())
	}
	cmd_g := NewCmd()
	cmds[ args["<cmd_name>"].GetItem(0).(string) ] = cmd_g
}

func end_processor(args map[string]*collections.ArrayList) {
	fmt.Println("end_processor")
	cmd_g = nil
}

func PrintCmdTree(list *collections.ArrayList) {
	iterator := list.Iterator()
	x := "{"
	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		fmt.Printf("%v", x)
		fmt.Printf("%v ", v.(node.Node).ToString())
		x = ","
	}
	fmt.Println("}")
}

func getFirst(list *collections.ArrayList) *collections.ArrayList {
	iterator := list.Iterator()
	nextList := collections.NewArrayList()

	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		nextList.Join(v.(node.Node).First())
	}
	return nextList
}


func sortNodeList(nodelist *collections.ArrayList) *collections.ArrayList {
	if (nodelist == nil || nodelist.Length() < 1) { 
	    return nodelist
	}
	
	gs := node.NewGNodeSorter(nodelist)
	sort.Sort(gs)
	return nodelist
}

func GetFollow(list *collections.ArrayList) *collections.ArrayList {
	iterator := list.Iterator()
	nextList := collections.NewArrayList()

	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		nextList.Join(v.(node.Node).Follow())
	}
	return nextList
}

func nodeListToString(nodelist *collections.ArrayList) string {
	var buffer bytes.Buffer
	iterator := nodelist.Iterator()
	x := "{ "
	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		buffer.WriteString(x)
		buffer.WriteString(v.(node.Node).ToString())
		x = " "
	}
	if x != "{ " {
		buffer.WriteString(" }")
	}
	return buffer.String()
}

func findMatch(list *collections.ArrayList, argList map[*parser.TokenDesc][]*parser.Arg,
	pos scanner.Position, tok rune, lit string) *collections.ArrayList {

	iterator := list.Iterator()
	matchingList := collections.NewArrayList()

	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		flist := v.(node.Node).First()
		fiter := flist.Iterator()
		for fi, fok, _ := fiter(); fok; fi, fok, _ = fiter() {
			gn := fi.(*node.GNode)
			fmt.Printf("node -- %s %d %s \n", gn.ToString(), tok, lit)
			tdesc := gn.GetResource().(*parser.TokenDesc)
			m, tdesc, arg := tdesc.Match(pos, tok, lit)
			if m == true {
				if prev, ok := argList[tdesc]; ok {
					if !gn.IsRepeatable() {
						continue
					}
					argList[tdesc] = append(prev, arg)
					// ToDo: Check if Repeat count is ok..
				} else {
					argslice := make([]*parser.Arg, 1)
					argslice[0] = arg
					argList[tdesc] = argslice
				}
				matchingList.Append(fi)
			}
		}
	}
	return matchingList
}

func (self *Cmd) Complete(line string) *collections.ArrayList {
	// Initialize the scanner.
	var s scanner.Scanner
	s.Filename = ""
	s.Init(strings.NewReader(line))
	s.Mode = (scanner.GoTokens & ^scanner.ScanFloats)
	s.Mode |= scanner.ScanInts

	arg_list := make(map[*parser.TokenDesc][]*parser.Arg)

	fmt.Printf("parse line ...... %s\n", line)
	first_list := self.dpt.First()

	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		tok := s.Scan()
		if tok == scanner.EOF {
			if first_list != nil &&
				first_list.Contains(node.NodeCR) {
				fmt.Println("line matched successfully!!")
				nodeCR := first_list.Find(node.NodeCR)
				fmt.Printf("handler %+v, %+v\n", nodeCR, nodeCR.(node.Node).GetCmdHandler())
				for k, v := range arg_list {
					fmt.Printf("arg-value %+v, %+v\n", k, v[0])
				}
			}
			break
		}
		//fmt.Printf("%s\t%s\t%s\n", s.Pos(), s.TokenText(), scanner.TokenString(tok))
		matching_list := findMatch(first_list, arg_list, s.Pos(), tok, s.TokenText())
		if matching_list.IsEmpty() {
			fmt.Println("syntax error at ", s.Pos(), s.TokenText(), scanner.TokenString(tok))
			return nil
		}
		sortNodeList(matching_list) // sort based on precedence of token type
		
		//fmt.Printf("matching list %+v\n", matching_list)
		first_list = GetFollow(matching_list)
		
		fmt.Printf("first_list %v\n", nodeListToString(first_list))
//		for i, v := range first_list.GetItems() {
//			fmt.Printf("first_list[%d]=%+v, sv=%s\n", i, v, v.(node.Node).ToString())
//		}
	}
	return first_list
}

func (cmd *Cmd) Init() {

	cmd.parseSyntax("command <cmd_name>", cmd_processor)
	cmd.defineKeyword("command")
	cmd.defineString("<cmd_name>")
	directed_parse_tree := cmd.dpt

	cmd.parseSyntax("syntax <cmd_syntax>+", syntax_processor)
	cmd.defineKeyword("syntax")
	cmd.defineQuotedString("<cmd_syntax>")
	directed_parse_tree, _ = directed_parse_tree.Union(cmd.dpt)

	cmd.parseSyntax("keyword <label> <help>", keyword_processor)
	cmd.defineKeyword("keyword")
	cmd.defineString("<label>")
	cmd.defineQuotedString("<help>")
	directed_parse_tree, _ = directed_parse_tree.Union(cmd.dpt)

	cmd.parseSyntax("end", end_processor)
	cmd.defineKeyword("end")
	directed_parse_tree, _ = directed_parse_tree.Union(cmd.dpt)

	fmt.Println("dpt : " + directed_parse_tree.ToString())
	cmd.dpt = directed_parse_tree
}

/*
func (cmd *Cmd) addToken(name string, tok rune, tokenType TokenType) *parser.TokenDesc {
    tokenDesc := NewTokenDesc(name, tok, tokenType)
    cmd.tokens[name] = tokenDesc
    return tokenDesc
}
*/
