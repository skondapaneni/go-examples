package main

import (
   "cli"
   "strings"
   "testing"
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

func TestDAG(t *testing.T) {
   pc := cli.NewParseContext()
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
