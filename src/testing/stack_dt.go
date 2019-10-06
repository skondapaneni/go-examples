package main

import ( 
    "collections" 
    "fmt"
)

func main() {

    var q = collections.NewStack()
    q.Push(1)
    q.Push(2)
    q.Push(3)
    q.Push(4)
    q.Push(5)

    fmt.Println("items =%v", q.GetItems())

    q.Pop()
    q.Pop()

    fmt.Println("items =%v", q.GetItems())

    q.Pop()
    q.Pop()

    fmt.Println("items =%v", q.GetItems())

    q.Pop()
    q.Pop()

    fmt.Println("items =%v", q.GetItems())

    _, err := q.Pop()
    fmt.Println("items =%v %v", q.GetItems(), err)

}
