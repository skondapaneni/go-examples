package node

type Resource interface {
    SetName(name string)
    GetName() string

    SetType(rt int)
    GetType() int

//    Eval(input interface{}, output interface{}) (status bool, err error)
}
