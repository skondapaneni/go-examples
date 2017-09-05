package main

import ( 
    "node" 
    "fmt"
)

func main() {

/*
     * S -> ABCd
     * A -> a | b | null
     * B -> e | f | null
     * C -> i | j
*/

    t1 := node.NewGNode("a")
    t2 := node.NewGNode("b")


    t3 := node.NewGNode("e")
    t4 := node.NewGNode("f")

    t5 := node.NewGNode("i")
    t6 := node.NewGNode("j")

    t7 := node.NewGNode("d")

    A := t1.AddChoice(t2).AddChoice(node.NodeNull)
    B := t3.AddChoice(t4).AddChoice(node.NodeNull)
    C := t5.AddNode(t6)

    A_U_B,_ := A.Union(B)

    S := A_U_B.Append(C).Append(t7)


/*
    switch  p2.(type) {
        default :
            fmt.Println("Unknown")
        case *node.GNodeAnd:
            fmt.Println("GNodeAnd")
        case *node.GNodeOr:
            fmt.Println("GNodeOr")
        case *node.GNode:
            fmt.Println("GNode")
        case *node.GNodeRef:
            fmt.Println("GNodeRef")
    }
*/

    //ref := node.NewGNodeRef(node.NodeNull)

    /*
    p := t1.AddNode(t2).AddChoice(
          t3.AddChoice(t4).AddChoice(node.NodeNull).Append( 
            t5.AddChoice(t6).AddChoice(node.NodeNull) ))
    */
    
    //and.AddNode(node.NodeNull)
    //and.AddNode(ref)

    //fmt.Println("items =%v", p.Children())
    //fmt.Println("nodeNull " + node.NodeNull.ToString())
    //fmt.Println("hasNodeNull =%v", p.Children().Contains(node.NodeNull))

    fmt.Println("A: " + A.ToString() + " %v", A.Children())
    fmt.Println("B: " + B.ToString() + " %v", B.Children())
    fmt.Println("C: " + C.ToString() + " %v", C.Children())
    fmt.Println("S: " + S.ToString())


    u1,_ := A.Union(A)
    fmt.Println("A union A : " + u1.ToString())

    fmt.Println("A union B : " + A_U_B.ToString())

    u3,_ := t2.Union(A_U_B)
    fmt.Println("2 union A_U_B : " + u3.ToString())

    fl := S.First()

    iterator := fl.Iterator();
    for v,ok,_ := iterator(); ok; v,ok,_ = iterator() {
        fmt.Print(v.(node.Node).ToString() + " ")
    }
    fmt.Println()


}
