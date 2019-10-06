package libinfra

import (
	"strings"
)

/*

var elist1 = libinfra.NewEnumList()

var (
	Alpha = elist1.CIota("A")
	Beta  = elist1.CIota("B")
)

type Example struct {
	X Enum
}

func main() {
	fmt.Printf("%+v\n", Example{Alpha})
	fmt.Printf("%+v\n", Example{Beta})
 	fmt.Printf("%+v\n", elist1.Enums())
}

*/

type Enum int

/* A named value is an enum item */
type EnumItem struct {
	name string
	val  Enum
}

/* A list of named values constitute an enumList */
type EnumList struct {
	enums []EnumItem
}

/* A map of named enum list */
var EnumsMap map[string]*EnumList = make(map[string]*EnumList)

func NewEnumList(name string) *EnumList {
	enumList := &EnumList{
		enums: make([]EnumItem, 0),
	}
	EnumsMap[name] = enumList
	return enumList
}

/* Add an enum item (value, name) to an enum List */
func (el *EnumList) addEnum(e Enum, s string) {
	el.enums = append(el.enums, EnumItem{name: s, val: e})
}

/* Return the name of the enum item */
func (el *EnumList) String(e Enum) string {
	return el.enums[int(e)].name
}

/* add an enum to an enumlist with only a name, the value is computed */
func (el *EnumList) CIota(s string) Enum {
	ei := EnumItem{name: s, val: Enum(len(el.enums))}
	el.enums = append(el.enums, ei)
	return ei.val
}

func (el *EnumList) Contains(s string) bool {
	for _, v := range el.enums {
		if strings.EqualFold(v.name, s) { // case insensitive
			return true
		}
	}
	return false
}

func (el *EnumList) HasPrefix(token string) bool {
	for _, v := range el.enums {
		if strings.HasPrefix(v.name, token) { // case insensitive
			return true
		}
	}
	return false
}

func (el *EnumList) Enums() []EnumItem {
	return el.enums
}
