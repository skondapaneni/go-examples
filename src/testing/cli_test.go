package main

import (
   "strings"
   "testing"
   "fmt"
   "node"
   "collections"
   "text/scanner"
)

type testResult struct {
   input string
   output string
}

var tests = []testResult {
  { "syntax <syntax>+",  "{ syntax <syntax> { T:NodeNull | <syntax>* } T:NodeCR }"}, 
  { "show [A | B ]", "{ show { T:NodeNull | { A | B } } T:NodeCR }"},
  { "show A B C", "{ show A B C T:NodeCR }"},
  { "show D [E F]", "{ show D { T:NodeNull | { E F } } T:NodeCR }"}, 
  { "a.b.c", "{ a b c T:NodeCR }"},
  { "a b | c", "{ { a b } | c | T:NodeCR }"},
  { "a b + c", "{ a b { T:NodeNull | b* } c T:NodeCR }"},
  { "a (b b)+ c", "{ a { b b } { T:NodeNull | { b b }* } c T:NodeCR }" },
}

/*
func TestDAG(t *testing.T) {
   pc := NewParseContext()
   for _, pair := range tests {
       pc.PrintPostFix(pair.input)
       p, _ := pc.CreateDAG(nil)
       v := p.ToString()
       if (strings.Compare(v, pair.output) != 0) {
           t.Error(
             "postfix for ", pair.input,
             "expected", pair.output,
             "got", v,
           )
       }
   }
}
*/

var exp = "a [b | c] d"
var tests_follow = []testResult {
  { "a", "b c d" },
  { "b", "d" },
  { "c", "d" },
  { "d", "T:NodeCR" },
}

func matchFollow(followList *collections.ArrayList, 
   expected string) bool {
	var s scanner.Scanner
        var tok rune
        s.Filename = ""
        s.Init(strings.NewReader(expected))
        for tok != scanner.EOF {
                tok = s.Scan()
                fmt.Println("expected follow tok:", s.TokenText())
	        iterator := followList.Iterator()
	        for v, ok, i := iterator(); ok; v, ok, i = iterator() {
                    if v.(node.Node).GetName() == s.TokenText() {
			followList.RemoveIndex(i)
			break
                    }
		}
        }

	if (followList.Length() == 0) {
	    return true
	}
	return false
}

/*
func TestFollow(t *testing.T) {

    cmd := NewCmd()
    cmd.parseSyntax(exp, nil)
    pc := NewParseContext()
    pc.PrintPostFix(exp)
    p, _ := pc.CreateDAG(nil)

    for _, pair := range tests_follow {
       nodeList := p.GetNodeByLabel(pair.input)
       if (nodeList == nil || nodeList.Length() == 0) {
           t.Error(
             "follow list for ", pair.input,
             "expected", pair.output,
             "got", nil,
           )
       }

       followList := getFollow(nodeList)

       followList.Print()

       if (!matchFollow(followList, pair.output)) {
           t.Error(
             "postfix for ", pair.input,
             "expected", pair.output,
             "got", followList,
           )
       }
   }
}
*/

/*
func main() {
    pc := cli.NewParseContext()

    pc.PrintPostFix("syntax <syntax>+")
    pc.CreateDAG(nil)

    pc.PrintPostFix("show [A | B ]")
    pc.CreateDAG(nil)

    pc.PrintPostFix("show A B C")
    pc.CreateDAG(nil)

    pc.PrintPostFix("show D [E F]")
    pc.CreateDAG(nil)

    pc.PrintPostFix("a.b.c")
    pc.CreateDAG(nil)

    pc.PrintPostFix("a b | c")
    pc.CreateDAG(nil)

    pc.PrintPostFix("a b + c")
    pc.CreateDAG(nil)

    pc.PrintPostFix("a (b b)+ c")
    pc.CreateDAG(nil)
}
*/

func main() {

    cmd := NewCmd()
    cmd.parseSyntax(exp, nil)
    pc := NewParseContext()
    pc.PrintPostFix(exp)
    p, _ := pc.CreateDAG(nil)

    for _, pair := range tests_follow {
       nodeList := p.GetNodeByLabel(pair.input)
       if (nodeList == nil || nodeList.Length() == 0) {
           fmt.Println(&Error(
             "follow list for ", pair.input,
             "expected", pair.output,
             "got", nil,
           ))
       }

       followList := getFollow(nodeList)

       followList.Print()

       if (!matchFollow(followList, pair.output)) {
           fmt.Println( &Error(
             "postfix for ", pair.input,
             "expected", pair.output,
             "got", followList,
           ))
       }
   }
}
