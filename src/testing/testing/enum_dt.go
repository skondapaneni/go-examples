
package main

import (
   "fmt"
   "libinfra"
)

var elist1 = libinfra.NewEnumList("Test")

var (
	Alpha = elist1.CIota("A")
	Beta  = elist1.CIota("B")
	Gamma  = elist1.CIota("C")
)

type Example struct {
	X libinfra.Enum
}

func main() {
	//fmt.Printf("%+v\n", Example{Alpha})
	//fmt.Printf("%+v\n", Example{Beta})
	//fmt.Printf("%+v\n", Example{Gamma})

 	fmt.Printf("%+v\n", elist1.Enums())
}

