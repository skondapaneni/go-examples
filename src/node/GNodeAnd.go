package node

import (
	"bytes"
	"collections"
	"fmt"
)

type GNodeAnd struct {
	children *collections.ArrayList

	next   Node
	parent Node

	repeatable  bool
	max_repeats int
}

func NewGNodeAnd() *GNodeAnd {
	nodeAnd := &GNodeAnd{
		children:    collections.NewArrayList(),
		repeatable:  false,
		max_repeats: 0,
	}
	return nodeAnd
}

func (n *GNodeAnd) GetName() string {
	return n.ToString()
}

func (n *GNodeAnd) IsTerminal() bool {
	return false
}

func (n *GNodeAnd) IsNullable() bool {

	if (n.children.Length() == 0) { 
	  	return true
	}

	iterator := n.children.Iterator()
	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		if (!v.(Node).IsNullable()) {
			return false
		}
	}

	return true
}

func (n *GNodeAnd) IsOptional() bool {
	return n.IsNullable()
}

func (n *GNodeAnd) IsRepeatable() bool {
	return n.repeatable
}

func (n *GNodeAnd) SetRepeatable(rpt bool, max_repeats int) {
	n.repeatable = rpt
	n.max_repeats = max_repeats
}

func (n *GNodeAnd) GetMaxRepeats() int {
	return n.max_repeats
}

func (n *GNodeAnd) GetSibling() Node {
	return n.next
}

func (n *GNodeAnd) SetSibling(sibling Node) {
	n.next = sibling
}

func (n *GNodeAnd) SetParent(parent Node) {
	n.parent = parent
}

func (self *GNodeAnd) GetParent() Node {
	return self.parent
}

func (n *GNodeAnd) AddNode(node Node) Node {
	last := n.children.GetLast()
	n.children.Add(node)
	if last != nil {
		last.(Node).SetSibling(node)
	}
	node.SetParent(n)
	return n
}

func (n *GNodeAnd) MergeNode(b Node) Node {
	switch b.(type) {
	default:
		return n.AddNode(b)
	case *GNodeAnd:
		target := b.(*GNodeAnd)
		iterator := target.children.Iterator()
		for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
			n.AddNode(v.(Node))
		}
	}
	return n
}

func (n *GNodeAnd) AddChoice(node Node) Node {
	nor := NewGNodeOr()
	nor.AddNode(n)
	nor.AddNode(node)
	return nor
}

func (n *GNodeAnd) AddOptional(node Node) Node {
	if node.IsOptional() {
		return n.Append(node)
	}
	nor := NewGNodeOr()
	nor.AddNode(NodeNull)
	nor.AddNode(node)
	return n.Append(nor)
}

func (n *GNodeAnd) Append(node Node) Node {
	if n.parent != nil {
		return n.parent.(Node).AddNode(node)
	}
	parent := NewGNodeAnd()
	parent.AddNode(n)
	parent.AddNode(node)
	return parent
}

func (n *GNodeAnd) Children() *collections.ArrayList {
	return n.children
}

func (n *GNodeAnd) First() *collections.ArrayList {
	if n.children != nil {
		top := n.children.GetItem(0).(Node)
		fl := top.First()
		if fl.Length() != 0 && fl.Contains(NodeNull) && n.next != nil {
			fl.Remove(NodeNull)
			fl.Join(n.Follow())
		}
		return fl
	}
	return nil
}

func (n *GNodeAnd) FirstFRef(ref Node) *collections.ArrayList {
	if n.children != nil {
		top := n.children.GetItem(0).(Node)
		fl := top.First()
		if fl.Length() != 0 && fl.Contains(NodeNull) && n.next != nil {
			if n.next != ref {
				fl.Remove(NodeNull)
				fl.Join(n.FollowFRef(ref))
			}
		}
		return fl
	}
	return nil
}

func (n *GNodeAnd) Follow() *collections.ArrayList {

	if n.next != nil {
		al := n.next.First()
		return al
	}
	if n.parent != nil {
		return n.parent.Follow()
	}
	return nil
}

func (n *GNodeAnd) FollowFRef(ref Node) *collections.ArrayList {

	if n.next != nil && n.next != ref {
		al := n.next.FirstFRef(ref)
		return al
	}

	if n.parent != nil {
		return n.parent.FollowFRef(ref)
	}

	return nil
}

func (self *GNodeAnd) InOrder(handler NodeHandler) {
	iterator := self.children.Iterator()
	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		v.(Node).InOrder(handler)
	}
}

func (self *GNodeAnd) GetCmdHandler() interface{} {
	return nil
}

func (self *GNodeAnd) SetCmdHandler(handler interface{}) {
	fmt.Println("can't set a  handler for GNodeAnd")
}

func (self *GNodeAnd) Compare(b interface{}) bool {

	switch b.(type) {
	default:
		fmt.Println("GNodeAnd::Unknown compare Type ")
		return false
	case *GNodeAnd:
		to := b.(*GNodeAnd)
		if self.children.Length() != to.children.Length() {
			return false
		}

		iterator := self.children.Iterator()
		for v, ok, i := iterator(); ok; v, ok, i = iterator() {
			subtree_node := to.children.GetItem(i)
			if v != subtree_node && !collections.Compare(v, subtree_node) {
				return false
			}
		}
	}
	return true
}

func (self *GNodeAnd) Union(subtree Node) (Node, error) {

	switch subtree.(type) {
	default:

	case *GNode, *GNodeRef:
		if self.Children().Length() == 0 {
			return subtree, nil
		}

		firstItem := self.Children().GetItem(0)
		if firstItem == subtree || collections.Compare(firstItem, subtree) {
			return self, nil
		}

	case *GNodeAnd:
		if self.Compare(subtree) {
			return self, nil
		}
	}

	nodeOr := NewGNodeOr()
	nodeOr.AddNode(self)
	nodeOr.AddNode(subtree)
	return nodeOr, nil
}

func (n *GNodeAnd) ToString() string {
	var buffer bytes.Buffer
	iterator := n.children.Iterator()
	x := "{ "
	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		buffer.WriteString(x)
		buffer.WriteString(v.(Node).ToString())
		x = " "
	}
	if x != "{ " {
		buffer.WriteString(" }")
		if n.IsRepeatable() {
			buffer.WriteString("*")
		}
	}
	return buffer.String()
}

func (self *GNodeAnd) GetNodeByLabel(label string) *collections.ArrayList {
	result := collections.NewArrayList()
	iterator := self.children.Iterator()
	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		subresult := v.(Node).GetNodeByLabel(label)
		if subresult != nil && subresult.Length() > 0 {
			result.Join(subresult)
		}
	}
	return result
}
