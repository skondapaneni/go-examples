package node

import (
	"bytes"
	"collections"
	"fmt"
	//	"runtime"
)

type GNodeOr struct {
	children *collections.ArrayList

	next   Node
	parent Node

	repeatable  bool
	max_repeats int
}

func NewGNodeOr() *GNodeOr {
	nodeOr := &GNodeOr{
		children:    collections.NewArrayList(),
		repeatable:  false,
		max_repeats: 0,
	}
	return nodeOr
}

func (self *GNodeOr) GetName() string {
	return self.ToString()
}

func (self *GNodeOr) IsTerminal() bool {
	return false
}

func (self *GNodeOr) IsNullable() bool {
	iterator := self.children.Iterator()
	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
                if (v.(Node).IsNullable()) {
		    return true
		}
	}
	return false
}

func (self *GNodeOr) IsOptional() bool {
	return self.IsNullable()
}

func (self *GNodeOr) IsRepeatable() bool {
	return self.repeatable
}

func (self *GNodeOr) SetRepeatable(rpt bool, max_repeats int) {
	self.repeatable = rpt
	self.max_repeats = max_repeats
}

func (self *GNodeOr) GetMaxRepeats() int {
	return self.max_repeats
}

func (self *GNodeOr) GetSibling() Node {
	return self.next
}

func (self *GNodeOr) SetSibling(sibling Node) {
	self.next = sibling
}

func (self *GNodeOr) SetParent(parent Node) {
	self.parent = parent
}

func (self *GNodeOr) GetParent() Node {
	return self.parent
}

func (self *GNodeOr) AddNode(node Node) Node {
	self.children.Add(node)
	node.SetParent(self)
	return self
}

func (n *GNodeOr) MergeNode(b Node) Node {
	switch b.(type) {
	default:
		return n.AddNode(b)
	case *GNodeOr:
		target := b.(*GNodeOr)
		iterator := target.children.Iterator()
		for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
			n.AddNode(v.(Node))
		}
	}
	return n
}

func (self *GNodeOr) AddChoice(node Node) Node {
	self.AddNode(node)
	return self
}

func (self *GNodeOr) AddOptional(node Node) Node {
	if node.IsOptional() {
		return self.Append(node)
	}
	nor := NewGNodeOr()
	nor.AddNode(NodeNull)
	nor.AddNode(node)
	return self.Append(nor)
}

func (self *GNodeOr) Append(node Node) Node {
	if self.parent != nil {
		return self.parent.AddNode(node)
	}
	parent := NewGNodeAnd()
	parent.AddNode(self)
	parent.AddNode(node)
	return parent
}

func (self *GNodeOr) Children() *collections.ArrayList {
	return self.children
}

func (self *GNodeOr) First() *collections.ArrayList {

	iterator := self.children.Iterator()
	fl := collections.NewArrayList()

	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		fl.Join(v.(Node).First())
	}

	if fl.Length() != 0 && fl.Contains(NodeNull) &&
		(self.next != nil || self.parent != nil) {
		fl.Remove(NodeNull)
		fl.Join(self.Follow())
	}

	return fl
}

/**
 * Func : Follow
 * Desc:
 *   To get the follow set of Terminal Nodes from this Node
 */
func (self *GNodeOr) Follow() *collections.ArrayList {

	if self.next != nil {
		return self.next.First()
	}

	if self.parent != nil {
		return self.parent.Follow()
	}

	return nil
}

func (self *GNodeOr) FirstFRef(ref Node) *collections.ArrayList {

//	fmt.Printf("self => %s repeatable %t\n", self.ToString(), self.IsRepeatable())
//	fmt.Printf("ref => %s\n", ref.ToString())
//	fmt.Printf("ref == self.next => %t\n", (self.next == ref))

	iterator := self.children.Iterator()
	fl := collections.NewArrayList()

	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		fl.Join(v.(Node).FirstFRef(ref))
	}

	if fl.Length() != 0 && fl.Contains(NodeNull) &&
		(self.next != nil || self.parent != nil) {

		if self.next != ref {
			fl.Remove(NodeNull)
			fl.Join(self.FollowFRef(ref))
		}
	}

	return fl
}

func (self *GNodeOr) FollowFRef(ref Node) *collections.ArrayList {

//	fmt.Printf("follow self => %s repeatable %t\n", self.ToString(), self.IsRepeatable())
//	fmt.Printf("follow ref => %s\n", ref.ToString())
//	fmt.Printf("follow ref == self.next => %t\n", (self.next == ref))

	if self.next != nil && self.next != ref {
		return self.next.FirstFRef(ref)
	}

	if self.parent != nil {
		return self.parent.FollowFRef(ref)
	}
	return nil
}

func (self *GNodeOr) Compare(b interface{}) bool {
	switch b.(type) {
	default:
		fmt.Println("GNodeOr::Unknown compare Type ")
		return false

	case *GNodeOr:
		to := b.(*GNodeOr)
		if self.children.Length() != to.children.Length() {
			return false
		}

		iterator := self.children.Iterator()
		for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
			if _, to_ok := to.Children().Contains_i(v); to_ok {
				continue
			}
			return false
		}
	}
	return true
}

func (self *GNodeOr) Union(subtree Node) (Node, error) {

	switch subtree.(type) {
	default:

	case *GNode, *GNodeRef:
		if self.Children().Length() == 0 {
			return subtree, nil
		}

		if _, ok := self.children.Contains_i(subtree); ok {
			return self, nil
		}

	case *GNodeOr:
		if self.Compare(subtree) {
			return self, nil
		}

		if self.next == nil && subtree.GetSibling() == nil {
			iterator := subtree.Children().Iterator()
			for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
				if _, to_ok := self.Children().Contains_i(v); to_ok {
					continue
				}
				self.AddNode(v.(Node))
			}
			return self, nil
		}
	}

	self.AddNode(subtree)
	return self, nil
}

func (self *GNodeOr) InOrder(handler NodeHandler) {
	iterator := self.children.Iterator()
	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		v.(Node).InOrder(handler)
	}
}

func (self *GNodeOr) GetCmdHandler() interface{} {
	return nil
}

func (self *GNodeOr) SetCmdHandler(handler interface{}) {
	fmt.Println("can't set a  handler for GNodeOr")
}

func (self *GNodeOr) ToString() string {
	var buffer bytes.Buffer
	iterator := self.children.Iterator()
	x := "{ "
	for v, ok, _ := iterator(); ok; v, ok, _ = iterator() {
		buffer.WriteString(x)
		buffer.WriteString(v.(Node).ToString())
		x = " | "
	}
	if x != "{ " {
		buffer.WriteString(" }")
		if self.IsRepeatable() {
			buffer.WriteString("*")
		}
	}
	return buffer.String()
}

func (self *GNodeOr) GetNodeByLabel(label string) *collections.ArrayList {
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
