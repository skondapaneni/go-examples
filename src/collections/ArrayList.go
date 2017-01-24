package collections

import (
)

/* 
 * A ArrayList is a growable list 
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

func (self *ArrayList) AddItem(l interface{}) {
    self.items = append(self.items, l)
}

func (self *ArrayList) RemoveItem(l interface{}) {
   for k,v := range self.items {
       if (v == l) {
           self.RemoveIndex(k)
           return
       }
   }
}

func (self *ArrayList) RemoveIndex(index int) {
    self.items = append(self.items[:index], self.items[index+1:]...)
}

func (self *ArrayList) IsQEmpty() bool {
    return (len(self.items) == 0)
}

func (self *ArrayList) GetItems() []interface{} {
    return self.items
}

func (self *ArrayList) GetItem(index int) interface{} {
    return self.items[index]
}


/*
func main() {

    var q *ArrayList = NewArrayList()
    q.AddItem(1)
    q.AddItem(2)
    q.AddItem(3)
    q.AddItem(4)
    q.AddItem(5)

    fmt.Println("items =%v", q.GetItems())

    q.AddItem(6)
    q.AddItem(7)

    fmt.Println("items =%v", q.GetItems())

    q.RemoveIndex(4)
    q.RemoveIndex(4)

    fmt.Println("items =%v", q.GetItems())

}
*/
