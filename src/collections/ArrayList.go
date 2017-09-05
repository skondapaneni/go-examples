package collections

import (
        "fmt"
        //"time"
        "sync"
        "reflect"
       // "bytes"
)

/* 
 * A ArrayList is a growable list 
 * Implemented using a go slice
 *
 */
type ArrayList struct {
    items []interface{}
}

func NewArrayList() *ArrayList {
    return &ArrayList{
        items: make ([]interface{}, 0),
    }
}

func (self *ArrayList) Add(l interface{}) {
    self.items = append(self.items, l)
}

func (self *ArrayList) Append(l interface{}) {
    self.items = append(self.items, l)
}

func (self *ArrayList) Join(other *ArrayList) {
	if (other == nil) {
		return
	}
    for _, v := range other.items {
        if (!self.Contains(v)) {
	        self.items = append(self.items, v)
        }
    }
}

func (self *ArrayList) Clear() {
    self.items = make ([]interface{}, 0)
}

func (self *ArrayList) Intersect(other *ArrayList) *ArrayList {
    result := NewArrayList()

    for _, v := range other.items {
        if (self.Contains(v)) {
            result.Add(v)
        }
    }

    if result.Length() > 0 {
        return result
    }

    return nil
}

func (self *ArrayList) Union(other *ArrayList) *ArrayList {

	if (other == nil) {
		return self
	}

    result := NewArrayList()
    for _, v := range self.items {
        result.Add(v)
    }

    for _, v := range other.items {
        if (!self.Contains(v)) {
            result.Add(v)
        }
    }

    if result.Length() > 0 {
        return result
    }

    return nil
}

func (self *ArrayList) Remove(l interface{}) interface{} {
    for k,v := range self.items {
        if (v == l || Compare(v, l)) {
            self.RemoveIndex(k)
            return v
        }
    }
    return nil
}

func (self *ArrayList) RemoveLast() interface{} {
    l := len(self.items)
    if (l == 0) {
        return nil
    }
    last := self.items[l-1]
    self.items = self.items[:l-1]
    return last
}

func Compare(a interface{}, b interface{}) bool {
    switch bv := b.(type) {
    case Comparable:
        switch av := a.(type) {
        default:
            return reflect.DeepEqual(a, b)
        case Comparable:
            return av.Compare(bv)
        }
    default:
        return reflect.DeepEqual(a, b)
    }
}

func (self *ArrayList) Contains(l interface{}) bool {
    for _,v := range self.items {
        if (v == l || Compare(v, l)) {
            return true
        }
    }
    return false
}

func (self *ArrayList) Contains_i(l interface{}) (int, bool) {
    for i,v := range self.items {
        if (v == l || Compare(v, l)) {
            return i,true
        }
    }
    return -1,false
}

/**
 * Func : Find
 * Does a linear search.
 */
func (self *ArrayList) Find(l interface{}) (item interface{}) {
    for _,v := range self.items {
        if (v == l || Compare(v, l)) {
            return v
        }
    }
    return nil
}


func (self *ArrayList) RemoveIndex(index int) interface{} {
    retVal := self.items[index]
    self.items = append(self.items[:index], self.items[index+1:]...)
    return retVal
}

func (self *ArrayList) ReplaceIndex(l interface{}, index int) {
    self.items[index] = l
}

func (self *ArrayList) Push(l interface{}) {
    self.items = append([]interface{} {l}, self.items...)
}

func (self *ArrayList) Pop() interface{} {
    first := self.items[0]
    self.RemoveIndex(0)
    return first;
}

func (self *ArrayList) IsQEmpty() bool {
    return (len(self.items) == 0)
}

func (self *ArrayList) IsEmpty() bool {
    return (len(self.items) == 0)
}

func (self *ArrayList) GetItems() []interface{} {
    return self.items
}

func (self *ArrayList) GetItem(index int) interface{} {
    return self.items[index]
}

func (self *ArrayList) GetLast() interface{} {
    if (len(self.items) == 0) {
        return nil
    }
    return self.items[len(self.items)-1]
}

func (self *ArrayList) Length() int {
    return (len(self.items))
}

func (self *ArrayList) Swap(i, j int) {
	self.items[i], self.items[j] = self.items[j], self.items[i]
}

/**
 * func : Iterator
 *
 * Returns an iterator func, which can be used to iterator over the list
 */
func (self *ArrayList) Iterator() func() (interface{}, bool, int) {
    i := -1
    return func() (interface{}, bool, int) {
        i++
        if i >= len(self.items) {
            return nil, false, i
        }
        return self.items[i], true, i
    }
}


/**
 * An iterator struct for providing an async interface to 
 * iteration
 */
type Iterator struct {
    list *ArrayList 
    iterChan chan interface{}  /* channel on which next item 
                                     is written by a go routine */ 
    nextChan chan bool /* channel on which next 
                              request is sent */
    closeChan chan bool /* channel to send a close 
                             to the iterator go method */
    hasMore bool
}

/**
 * Constructor for Iterator object
 */
func (self *ArrayList) NewIterator()  *Iterator {
     n := &Iterator { 
             list: self, 
             iterChan : make(chan interface{}),
             nextChan : make (chan bool),
             closeChan : make (chan bool),
             hasMore : true,
    }

    if (self.IsQEmpty()) {
        n.hasMore = false
    } 

    n.init()
    return n
}

func (iter *Iterator) init() {
    go func(iter *Iterator) {
        done := false
        iterator := iter.list.Iterator();
        for v,ok,i := iterator(); ok; v,ok,i= iterator() {
            select {
            case <-iter.nextChan:
                iter.hasMore = (i+1 < len(iter.list.items))
                iter.iterChan <- v        
            case <-iter.closeChan:
                done = true
                break
            }
        }
        iter.hasMore = false
        for (!done) {
            select {
            case <-iter.nextChan:
                iter.iterChan <- nil        
            case <-iter.closeChan:
                done = true
                break
            }
        }

        close(iter.nextChan)
        close(iter.closeChan)
        close(iter.iterChan)
        fmt.Println("Closed all channels")
    } (iter)
}

func (iter *Iterator) Next() interface{} {
    if (iter.hasMore) {
        iter.nextChan <- true
        select {
        case v := <- iter.iterChan :
            fmt.Println("Next() : ",  v)
            return v
        }
    }
    return nil
}

func (iter *Iterator) HasNext() bool {
    return iter.hasMore
}

func (iter *Iterator) Close() {
   var wg sync.WaitGroup
   wg.Add(1)

   go func(iter *Iterator) {
      defer wg.Done()
      for _, ok := <- iter.iterChan; ok; _,ok = <- iter.iterChan {
      }
   }(iter)

   iter.closeChan <- true 
   wg.Wait()

}

func (self *ArrayList) Print() {
    iterator := self.Iterator();
    x := "{"
    for v,ok,_ := iterator(); ok; v,ok,_ = iterator() {
        fmt.Printf("%+v", x)
        fmt.Printf("%+v ", v)
        x = " & "
    }
    fmt.Println("}")
}

/*
type Keyword struct {
    name  string 
}

var keyword1 = &Keyword { name : "one" }
var keyword2 = &Keyword { name : "two" }

func (k *Keyword) toString() string {
     return k.name
}

type Node interface {
     toString() string
}

type NodeAnd struct {
      al *ArrayList
}

type NodeOr struct {
      al *ArrayList
}

func NewNodeAnd() *NodeAnd {
    nodeAnd := &NodeAnd {
        al :   NewArrayList(),
    }
    return nodeAnd
}

func (na *NodeAnd) Add(node Node) {
    na.al.Add(node)
}

func (na *NodeAnd) toString() string {
    var buffer bytes.Buffer
    iterator := na.al.Iterator();
    x := "{"
    for v,ok,_ := iterator(); ok; v,ok,_ = iterator() {
        buffer.WriteString(x)
        buffer.WriteString(v.(Node).toString())
        x = " & "
    }
    buffer.WriteString("}")
    return buffer.String()
}

func NewNodeOr() *NodeOr {
    nodeOr := &NodeOr {
        al :   NewArrayList(),
    }
    return nodeOr
}

func (no *NodeOr) Add(node Node) {
    no.al.Add(node)
}

func (no NodeOr) toString() string {
    var buffer bytes.Buffer
    iterator := no.al.Iterator();
    x := "{"
    for v,ok,_ := iterator(); ok; v,ok,_ = iterator() {
        buffer.WriteString(x)
        buffer.WriteString(v.(Node).toString())
        x = " | "
    }
    buffer.WriteString("}")
    return buffer.String()
}
*/


/*
func main() {

    var q *ArrayList = NewArrayList()
    q.Add(1)
    q.Add(2)
    q.Add(3)
    q.Add(4)
    q.Add(5)
    fmt.Println("items =", q.GetItems())
    //Output: items = [1 2 3 4 5]

    q.Add(6)
    q.Add(7)
    fmt.Println("items =", q.GetItems())
    //Output: items = [1 2 3 4 5]
    //items = [1 2 3 4 5 6 7]

    q.RemoveIndex(4)
    q.RemoveIndex(4)
    q.Push(10)
    q.Push(11)
    fmt.Println("items =", q.GetItems())
    //Output: items = [1 2 3 4 5]
    //items = [1 2 3 4 5 6 7]
    //items = [11 10 1 2 3 4 7]

    q.Pop()
    fmt.Println("items =", q.GetItems())
    //Output: items = [1 2 3 4 5]
    //items = [1 2 3 4 5 6 7]
    //items = [11 10 1 2 3 4 7]
    //items = [10 1 2 3 4 7]

    // concurrent modification
    go func() {
        q.Add(100)
        fmt.Println("modified items =", q.GetItems())
    }()

    // Example 1 - using the iterator method
    iterator := q.Iterator();
    for v,ok,_ := iterator(); ok; v,ok,_ = iterator() {
        fmt.Println("@@item =", v)
        time.Sleep(1 * time.Second)
    }

    // concurrent modification
    go func() {
        q.RemoveIndex(4)
	q.Add(90)
	q.Add(91)
        fmt.Println("modified items =", q.GetItems())
    }()

    // Example 2 - more traditional iteration with an iterator object..
    iter := q.NewIterator()
    for (iter.HasNext()) {
        fmt.Println("##item =", iter.Next())
        time.Sleep(1 * time.Second)
    }
    iter.Close()


    na := NewNodeAnd() 
    na.Add(keyword1)
    na.Add(keyword2)

    fmt.Println("na =", na.toString())
}
*/
