package node

import (
	"collections"
	"fmt"
	"strings"
)

/* A Recursive Descent Parser Grammar Terminal Node */
type GNode struct {
	next   Node /* It's sibling */
	parent Node /* It's parent if part of a container */

	repeatable  bool /* Is it a repeatable production */
	max_repeats int  /* max number of repeats */

	resource   Resource    /* A resource desc which can represent anything like a lex token, task */
	cmdHandler interface{} /* A cmd handler that will be executed on matching an input and
	   this node is the terminating node */
}

/* Declare const variables */
var NodeNull *GNode = &GNode{}
var NodeCR *GNode = &GNode{}

func Less(a *GNode, b *GNode) bool {
	return (a.LessThanByTokenType(b))
}

// GNodeSorter joins a By function and a slice of GNodes to be sorted.
type GNodeSorter struct {
	gnodes *collections.ArrayList
}

func NewGNodeSorter(nodelist *collections.ArrayList) *GNodeSorter {
	gnode_sorter := &GNodeSorter{
		gnodes: nodelist,
	}
	return gnode_sorter
}

// Len is part of sort.Interface.
func (s *GNodeSorter) Len() int {
	return s.gnodes.Length()
}

// Swap is part of sort.Interface.
func (s *GNodeSorter) Swap(i, j int) {
	s.gnodes.Swap(i, j)
}

// Less is part of sort.Interface.
// It is implemented by calling the "by" closure in the sorter.
func (s *GNodeSorter) Less(i, j int) bool {
	return Less(s.gnodes.GetItem(i).(*GNode),
		s.gnodes.GetItem(i).(*GNode))
}

/**
 * Func: NewGNode
 * Desc:
 *   A constructor for a terminal Node. Takes a tokenDesc which
 *   knows how to match a lex token.
 */
func NewGNode(resource Resource) *GNode {
	gnode := &GNode{
		resource: resource,
	}
	return gnode
}

/**
 * Func: GetName
 * Desc:
 *  Returns the name of the Terminal Node
 */
func (self *GNode) GetName() string {
	return self.ToString()
}

/**
 * Func : IsTermial
 * Desc:
 *    Returns whether the Grammar Node is a Terminal.
 *    By default all GNode's are terminal's.
 */
func (self *GNode) IsTerminal() bool {
	return true
}

/**
 * Func : IsNullable
 * Desc:
 *    Returns whether the Grammar Node is a nullable node which means, the production
 * can be skipped to match it's following node..
 */
func (self *GNode) IsNullable() bool {
	if self == NodeNull {
		return true
	}
	return false
}

/**
 * Func : IsOptional
 * Desc:
 *    Returns whether the Grammar Node is optional which means, the production
 * can be skipped to match it's following node..
 */
func (self *GNode) IsOptional() bool {
	return self.IsNullable()
}

/**
 * Func : IsRepeatable
 * Desc:
 *    Returns whether the Grammar Node can repeat itself for matching a production.
 */
func (self *GNode) IsRepeatable() bool {
	return self.repeatable
}

/**
 * Func : SetRepeatable
 * Desc:
 *    Set's whether the Grammar Node can repeat itself for matching a production.
 */
func (self *GNode) SetRepeatable(rpt bool, max_repeats int) {
	self.repeatable = rpt
	self.max_repeats = max_repeats
}

/**
 * Func : GetMaxRepeats
 * Desc:
 *    Returns the number of times the Grammar Node can repeat itself
 * for matching input productions.
 */
func (self *GNode) GetMaxRepeats() int {
	return self.max_repeats
}

/**
 * Func: AddNode
 *   Adds a node as a sibling with a new container node as parent.
 */
func (self *GNode) AddNode(sibling Node) Node {
	na := NewGNodeAnd()
	na.AddNode(self)
	na.AddNode(sibling)
	return na
}

/**
 * Func: MergeNode
 *   Merges 2 production Trees together.
 */
func (self *GNode) MergeNode(sibling Node) Node {
	return self.AddNode(sibling)
}

/**
 * func: AddChoice
 *   Adds a node as an alternate node with a new container node as parent.
 */
func (self *GNode) AddChoice(node Node) Node {
	nor := NewGNodeOr()
	nor.AddNode(self)
	nor.AddNode(node)
	return nor
}

/**
 * func: AddOptional
 *   Adds a new optional node as sibling node.
 */
func (self *GNode) AddOptional(node Node) Node {
	nor := NewGNodeOr()
	nor.AddNode(NodeNull)
	nor.AddNode(node)
	return self.Append(nor)
}

func (self *GNode) Append(node Node) Node {
	if self.parent != nil {
		return self.parent.(Node).AddNode(node)
	}
	parent := NewGNodeAnd()
	parent.AddNode(self)
	parent.AddNode(node)
	return parent
}

func (self *GNode) SetSibling(sibling Node) {
	self.next = sibling
}

func (self *GNode) GetSibling() Node {
	return self.next
}

func (self *GNode) SetParent(parent Node) {
	self.parent = parent
}

func (self *GNode) GetParent() Node {
	return self.parent
}

func (self GNode) Children() *collections.ArrayList {
	return nil
}

func (self *GNode) InOrder(handler NodeHandler) {
	handler(self)
}

func (self *GNode) GetCmdHandler() interface{} {
	return self.cmdHandler
}

func (self *GNode) SetCmdHandler(handler interface{}) {
	self.cmdHandler = handler
}

func (self *GNode) GetResource() Resource {
	return self.resource
}

func (self *GNode) Compare(b interface{}) bool {
	node := b.(Node)
	if self.IsTerminal() && node.IsTerminal() &&
		strings.Compare(self.GetName(), node.GetName()) == 0 {
		if self.GetCmdHandler() != node.GetCmdHandler() {
			return false
		}
		return true
	}
	return false
}

/*
func (self *GNode) Match(pos scanner.Position, tok rune, lit string) (bool,
	*TokenDesc, *Arg) {
	if self.tokenDesc != nil {
		return self.tokenDesc.Match(pos, tok, lit)
	}
	return false, nil, nil
}
*/

//func (self *GNode) Eval(input interface{}, output interface{}) (status bool, error interface{}) {
//	if self.resource != nil {
//		return self.resource.Eval(input, output)
//	}
//	return false, nil
//}

func (self *GNode) Union(subtree Node) (Node, error) {

	switch subtree.(type) {
	default:

	case *GNode, *GNodeRef:
		if self.Compare(subtree) {
			return self, nil
		}

	case *GNodeAnd, *GNodeOr:
		return subtree.Union(self)
	}

	nodeOr := NewGNodeOr()
	nodeOr.AddNode(self)
	nodeOr.AddNode(subtree)
	return nodeOr, nil
}

func (self *GNode) First() *collections.ArrayList {
	fl := collections.NewArrayList()
	fl.Add(self)
	return fl
}

func (self *GNode) Follow() *collections.ArrayList {

	if self.next != nil {
		return self.next.First()
	}

	if self.parent != nil {
		return self.parent.Follow()
	}
	return nil
}

func (n *GNode) FirstFRef(ref Node) *collections.ArrayList {
	fl := collections.NewArrayList()
	fl.Add(n)
	return fl
}

func (self *GNode) FollowFRef(ref Node) *collections.ArrayList {

	fmt.Printf("self => %s repeatable %t\n", self.ToString(), self.IsRepeatable())
	fmt.Printf("ref => %s\n", ref.ToString())

	if self.IsRepeatable() {
		fl := self.First()
		if self.IsOptional() {
			if self.next != nil && self.next != ref {
				fl.Join(self.next.FirstFRef(ref))
			}
			if self.parent != nil {
				fl.Join(self.parent.FollowFRef(ref))
			}
		}
		return fl
	}

	if self.next != nil && self.next != ref {
		fmt.Printf("self.next => %s\n", self.next.ToString())

		return self.next.FirstFRef(ref)
	}

	if self.parent != nil {
		fmt.Printf("self.parent => %s\n", self.parent.ToString())

		return self.parent.FollowFRef(ref)
	}
	return nil
}

func (self *GNode) ToString() string {
	if self == NodeNull {
		return "T:NodeNull"
	}
	if self == NodeCR {
		return "T:NodeCR"
	}
	return self.resource.GetName()
}

func (self *GNode) LessThanByTokenType(to *GNode) bool {
	if self.resource != nil && to.resource != nil {
		return self.resource.GetType() <
			to.resource.GetType()
	}
	return false
}

func (self *GNode) GetNodeByLabel(label string) *collections.ArrayList {
	if self.GetName() == label {
		result := collections.NewArrayList()
		result.Add(self)
		return result
	}
	return nil
}
