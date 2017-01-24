package collections

import (
    "fmt"
    "reflect"
    "encoding/json"
    "io"
)

/* Tree Node */
type Node struct {
    Value    interface{} 
    Children []*Node  
    parent   *Node 
    is_dirty bool
}

func NewNode(value  interface{}) *Node {
    return &Node{
        Value:    value,
        Children: make([]*Node, 0),
        parent:   nil,
    }
}

func (n *Node) AddChild(node *Node) {
    n.Children = append(n.Children, node)
    n.is_dirty = true
    node.parent = n
}

func (n *Node) GetChildren() []*Node {
    return n.Children
}

func (n *Node) GetChildrenCount() int {
    if (n.Children != nil) {
        return len(n.Children)
    }
    return 0
}

func (n *Node) GetChildAt(childIndex int) *Node {
    child := n.Children[childIndex]
    return child
}

func (n *Node) Compare(node *Node) bool {
    switch bv := node.Value.(type) {
    default:
        return reflect.DeepEqual(n.Value, node.Value)
    case Comparable:
        switch av := n.Value.(type) {
        default:
            return reflect.DeepEqual(n.Value, node.Value)
        case Comparable:
            return av.Compare(bv)
        }
    }
}

func (n *Node) FindChild(node *Node) *Node {
    for _, v := range n.Children {
        if v.Compare(node) {
            return v
        }
    }
    return nil
}

func (n *Node) removeChildAtIndex(index int) {
    n.Children = append(n.Children[:index], n.Children[index+1:]...)
}

func (n *Node) DeleteChild(node *Node) *Node {
    for i, v := range n.Children {
        if v.Compare(node) {
            n.removeChildAtIndex(i)
            n.is_dirty = true
            return v
        }
    }
    return nil
}

func (n *Node) GetParent() *Node {
    return n.parent
}

func (n *Node) GetValue() interface{} {
    return n.Value
}

func (n *Node) String() string {
    return fmt.Sprint(n.Value)
}

func (n *Node) Write(w io.Writer) {
    b, error := json.Marshal(*n)
    if (error != nil) {
        fmt.Println("error %v\n", error)
    }
    w.Write(b)
}

