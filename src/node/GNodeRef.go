package node

import (
	"bytes"
	"collections"
	"fmt"
	"strings"
)

type GNodeRef struct {
	ref    Node
	next   Node
	parent Node

	repeatable  bool
	max_repeats int
	cmdHandler  interface{}
}

func NewGNodeRef(node Node) *GNodeRef {
	gnodeRef := &GNodeRef{
		ref:    node,
		next:   nil,
		parent: nil,
	}
	return gnodeRef
}

func (self *GNodeRef) GetName() string {
	return self.ref.GetName()
}

func (self *GNodeRef) IsTerminal() bool {
	return self.ref.IsTerminal()
}

func (self *GNodeRef) IsNullable() bool {
	return self.ref.IsNullable()
}

func (self *GNodeRef) IsOptional() bool {
	return self.ref.IsOptional()
}

func (self *GNodeRef) IsRepeatable() bool {
	return self.repeatable
}

func (self *GNodeRef) SetRepeatable(rpt bool, max_repeats int) {
	self.repeatable = rpt
	self.max_repeats = max_repeats
}

func (self *GNodeRef) GetMaxRepeats() int {
	return self.ref.GetMaxRepeats()
}

func (self *GNodeRef) AddNode(sibling Node) Node {
	na := NewGNodeAnd()
	na.AddNode(self)
	na.AddNode(sibling)
	return na
}

func (self *GNodeRef) MergeNode(sibling Node) Node {
	return self.AddNode(sibling)
}

func (self *GNodeRef) AddChoice(node Node) Node {
	nor := NewGNodeOr()
	nor.AddNode(self)
	nor.AddNode(node)
	return nor
}

func (self *GNodeRef) Append(node Node) Node {
	if self.parent != nil {
		return self.parent.(Node).AddNode(node)
	}
	parent := NewGNodeAnd()
	parent.AddNode(self)
	parent.AddNode(node)
	return parent
}

func (self *GNodeRef) AddOptional(node Node) Node {
	nor := NewGNodeOr()
	nor.AddNode(NodeNull)
	nor.AddNode(node)
	return self.Append(nor)
}

func (self *GNodeRef) SetSibling(sibling Node) {
	self.next = sibling
}

func (self *GNodeRef) GetSibling() Node {
	return self.next
}

func (self *GNodeRef) SetParent(parent Node) {
	self.parent = parent
}

func (self *GNodeRef) GetParent() Node {
	return self.parent
}

func (self *GNodeRef) Children() *collections.ArrayList {
	return nil
}

func (self *GNodeRef) Compare(b interface{}) bool {
	node := b.(Node)
	fmt.Println("self.GetName(), node.GetName() " + self.GetName() + " " +
		node.GetName())
	if node.IsTerminal() &&
		strings.Compare(self.GetName(), node.GetName()) == 0 {
		return true
	}
	return false
}

func (self *GNodeRef) Union(subtree Node) (Node, error) {

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

func (self *GNodeRef) InOrder(handler NodeHandler) {
	handler(self.ref)
}

func (self *GNodeRef) GetCmdHandler() interface{} {
	return self.cmdHandler
}

func (self *GNodeRef) SetCmdHandler(handler interface{}) {
	self.cmdHandler = handler
}

func (self *GNodeRef) First() *collections.ArrayList {
	
	fl := self.ref.FirstFRef(self)

	if fl.Length() != 0 && fl.Contains(NodeNull) &&
		(self.next != nil || self.parent != nil) {
		fl.Remove(NodeNull)
		fl.Join(self.Follow())
	}
	
	return fl
}

func (self *GNodeRef) Follow() *collections.ArrayList {
		
	if self.next != nil {
		return self.next.First()
	}
	
	if self.parent != nil {
		return self.parent.Follow()
	}
	return nil
}

func (self *GNodeRef) FirstFRef(ref Node) *collections.ArrayList {
	al := collections.NewArrayList()
	if ref == self {
		return al
	}

	al.Add(self.ref.FirstFRef(self))

	if self.ref.IsOptional() {
		al.Join(self.FollowFRef(ref))
	}
	
	return al
}

func (self *GNodeRef) FollowFRef(ref Node) *collections.ArrayList {
		
	if self.next != nil && self.next != ref {
		return self.next.FirstFRef(ref)
	}
	if self.parent != nil {
		return self.parent.FollowFRef(ref)
	}
	return nil
}

func (self *GNodeRef) ToString() string {
	var buffer bytes.Buffer
	buffer.WriteString(self.ref.ToString())
	if self.IsRepeatable() {
		buffer.WriteString("*")
	}
	return buffer.String()
}

func (self *GNodeRef) GetNodeByLabel(label string) *collections.ArrayList {
	if self.GetName() == label {
		result := collections.NewArrayList()
		result.Add(self)
		return result
	}
	return nil
}

