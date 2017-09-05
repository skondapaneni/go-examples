package collections

import (
	"log"
)

/*
 * A Stack is a growable list
 *
 */
type Stack struct {
	items []interface{}
}

type InvalidOpError struct {
	ErrorString string
}

var (
	ErrStackEmpty = &InvalidOpError{"Stack is Empty"}
)

func (err *InvalidOpError) Error() string { return err.ErrorString }

func NewStack() *Stack {
	return &Stack{
		items: make([]interface{}, 0),
	}
}

func (self *Stack) Push(l interface{}) {
	self.items = append(self.items, l)
}

func (self *Stack) Pop() (interface{}, error) {
	var index = len(self.items)
	var item interface{}
	if index != 0 {
		item, self.items = self.items[index-1], self.items[:index-1]
		return item, nil
	}
	return nil, ErrStackEmpty
}

func (self *Stack) hasItem(l interface{}) bool {
	for _, v := range self.items {
		if v == l {
			return true
		}
	}
	return false
}

func (self *Stack) Clear() {
	self.items = make([]interface{}, 0)
}

func (self *Stack) Length() int {
    return (len(self.items))
}

func (self *Stack) IsEmpty() bool {
	return (len(self.items) == 0)
}

func (self *Stack) GetItems() []interface{} {
	return self.items
}

func (self *Stack) GetItem(index int) interface{} {
	return self.items[index]
}

func (self *Stack) Peek() interface{} {
	var index = len(self.items)
	if (index == 0) {
		return nil
	}
	return self.items[index-1]
}

func (self *Stack) PrintItems() {
	for i, v := range self.items {
		log.Printf("[%d] [%v]", i, v)
	}
}
