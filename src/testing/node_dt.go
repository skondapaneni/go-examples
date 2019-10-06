package main

import (
	"bytes"
	"parser"
	"collections"
	"fmt"
	"node"
	"strings"
	"text/scanner"
)

type testResult struct {
	input  string
	output string
}

var exp1 = "a [b | c] d"
var tests_follow1 = []testResult{
	{"a", "b c d"},
	{"b", "d"},
	{"c", "d"},
	{"d", "T:NodeCR"},
}

var exp2 = "show [a | b | c]+ z "
var tests_follow2 = []testResult{
	{"show", "a b c z"},
	{"a", "a b c z"},
	{"b", "a b c z"},
	{"c", "a b c z"},
	{"z", "T:NodeCR"},
}

var exp3 = "show a  b  c* z "
var tests_follow3 = []testResult{
	{"show", "a"},
	{"a", "b"},
	{"b", "c z"},
	{"c", "c z"},
	{"z", "T:NodeCR"},
}

var exp4 = "show [a | b | c]* z "
var tests_follow4 = []testResult{
	{"show", "a b c z"},
	{"a", "a b c z"},
	{"b", "a b c z"},
	{"c", "a b c z"},
	{"z", "T:NodeCR"},
}

var exp5 = "show {a | b | c}+ z "
var tests_follow5 = []testResult{
	{"show", "a b c"},
	{"a", "a b c z"},
	{"b", "a b c z"},
	{"c", "a b c z"},
	{"z", "T:NodeCR"},
}

var exp6 = "show {a {b | d}+ c}+ z"
var tests_follow6 = []testResult{
	{"show", "a"},
	{"a", "b d"},
	{"b", "b d c"},
	{"d", "b d c"},
	{"c", "a z"},
	{"z", "T:NodeCR"},
}

var exp7 = "show [a {b | d}+ c]+ z"
var tests_follow7 = []testResult{
	{"show", "a z"},
	{"a", "b d"},
	{"b", "b d c"},
	{"d", "b d c"},
	{"c", "a z"},
	{"z", "T:NodeCR"},
}

var exp8 = "show [a b c] [d]* z"
var tests_follow8 = []testResult{
	{"show", "a d z"},
	{"a", "b"},
	{"b", "c"},
	{"c", "d z"},
	{"d", "d z"},
	{"z", "T:NodeCR"},
}

var exp9 = "show [a]+ b c"
var tests_follow9 = []testResult{
	{"show", "a, b"},
	{"a", "a, b"},
	{"b", "c"},
	{"c", "T:NodeCR"},
}

var exp10 = "show a+ b c"
var tests_follow10 = []testResult{
	{"show", "a"},
	{"a", "a, b"},
	{"b", "c"},
	{"c", "T:NodeCR"},
}

var exp11 = "show {[a | b | d]+ c}+ z"
var tests_follow11 = []testResult{
	{"show", "a b d c"},
	{"a", "a b d c"},
	{"b", "a b d c"},
	{"d", "a b d c"},
	{"c", "a b d c z"},
	{"z", "T:NodeCR"},
}

var exp12 = "show [[a | b | d]+ c]+ z"
var tests_follow12 = []testResult{
	{"show", "a b d c z"},
	{"a", "a b d c"},
	{"b", "a b d c"},
	{"d", "a b d c"},
	{"c", "a b d c z"},
	{"z", "T:NodeCR"},
}

func NodeListToString(nodelist *collections.ArrayList) string {
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

func matchFollow(followList *collections.ArrayList,
	expected string) bool {
	var s scanner.Scanner
	var tok rune
	s.Filename = ""
	s.Init(strings.NewReader(expected))
	for tok != scanner.EOF {
		tok = s.Scan()
		iterator := followList.Iterator()
		for v, ok, i := iterator(); ok; v, ok, i = iterator() {
			if v.(node.Node).GetName() == s.TokenText() {
				followList.RemoveIndex(i)
				break
			}
		}
	}

	if followList.Length() == 0 {
		return true
	}
	return false
}

func testFollow(exp string, tests_follow []testResult) {
	pc := parser.NewParseContext()
	pc.PrintPostFix(exp)
	p, _ := pc.CreateDAG(nil)

	for _, pair := range tests_follow {
		nodeList := p.GetNodeByLabel(pair.input)
		if nodeList == nil || nodeList.Length() == 0 {
			fmt.Errorf("%s %v %s %v %s %v",
				"follow list for ", pair.input,
				"expected", pair.output,
				"got", nil,
			)
		}

		fmt.Printf("nodeList for label %s, %s\n", pair.input, NodeListToString(nodeList))
		followList := parser.GetFollow(nodeList)
		fmt.Printf("expected follow for %s : %s\n", pair.input, pair.output)
		fmt.Printf("Follow %s\n", NodeListToString(followList))
		fmt.Printf("--------------- \n")

		if !matchFollow(followList, pair.output) {
			fmt.Errorf("%s %v %s %v %s %v",
				"postfix for ", pair.input,
				"expected", pair.output,
				"got", followList,
			)
		}
	}
}

func main() {
//	fmt.Println("==========Test1")
//	testFollow(exp1, tests_follow1)
//	fmt.Println("==========Test2")
//	testFollow(exp2, tests_follow2)
//	fmt.Println("==========Test3")
//	testFollow(exp3, tests_follow3)
//	fmt.Println("==========Test4")
//	testFollow(exp4, tests_follow4)
//	fmt.Println("==========Test5")
//	testFollow(exp5, tests_follow5)
	fmt.Println("==========Test6")
	testFollow(exp6, tests_follow6)
//	fmt.Println("==========Test7")
//	testFollow(exp7, tests_follow7)
//	fmt.Println("==========Test8")
//	testFollow(exp8, tests_follow8)
//	fmt.Println("==========Test9")
//	testFollow(exp9, tests_follow9)
//	fmt.Println("==========Test10")
//	testFollow(exp10, tests_follow10)
	fmt.Println("==========Test11")
	testFollow(exp11, tests_follow11)

	fmt.Println("==========Test12")
	testFollow(exp12, tests_follow12)

}
