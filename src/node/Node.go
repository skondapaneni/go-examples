package node

import (
    "collections"
    "libinfra"
)

/**
 * Parse trees are an in-memory representation of the input with a structure that conforms to the grammar.

 * The advantages of using parse trees instead of semantic actions:
 *    - You can make multiple passes over the data without having to re-parse the input.
 *    - You can perform transformations on the tree.
 *    - You can evaluate things in any order you want, whereas with attribute schemes you have to process in a begin to end fashion.
 *    - You do not have to worry about backtracking and action side effects that may occur with an ambiguous grammar.
 * 
 **/
type NodeHandler func(n Node)

type RunTimeException struct {
   msg  string
}

type NodeExpressionError struct {
    errno libinfra.Enum
}

var errorList = libinfra.NewEnumList("NodeExpressionError")

var (
    ERR_AMBIGUOUS = errorList.CIota("Ambiguous Expression")
    //BOOL_OR  = errorList.CIota("OR")
)

func (e *NodeExpressionError) Error() string { 
    return errorList.String(e.errno)
}

type Node interface {
    GetName() string
    IsTerminal() bool /* Is the Node a Terminating node in the grammar */
    IsNullable() bool /* Is the node Nullable or skippable for
                         matching in the grammar */
    IsOptional() bool   /* Is the node Optional in the grammar */
    IsRepeatable() bool /* Does the node repeat itslef in the grammar */
    SetRepeatable(rpt bool, max_repeats int)
    GetMaxRepeats() int
    

    AddNode(sibling Node) Node
    MergeNode(sibling Node) Node
    AddChoice(node Node) Node
    AddOptional(node Node) Node
    Append(node Node) Node

    Union(subtree Node) (Node, error)

    SetSibling(sibling Node)
    GetSibling() Node

    SetParent(parent Node)
    GetParent() Node
    
    GetCmdHandler() interface{}
    SetCmdHandler(handler interface{})

    InOrder(handler NodeHandler)
    Children() *collections.ArrayList
    
    First() *collections.ArrayList
    Follow() *collections.ArrayList
    
    FirstFRef(ref Node) *collections.ArrayList
    FollowFRef(ref Node) *collections.ArrayList

    Compare(b interface{}) bool
    GetNodeByLabel(label string) *collections.ArrayList

    //RemoveNode(node Node) Node
    //Help() string
    //EnumRef() string

    ToString() string
}

func Optional(node Node) Node {
	nor := NewGNodeOr()
	nor.AddNode(NodeNull)
	nor.AddNode(node)
	return nor
}
