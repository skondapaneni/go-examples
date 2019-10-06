package libinfra

var booleanE = NewEnumList("BooleanOp")

var (
	BOOL_AND = booleanE.CIota("AND")
	BOOL_OR  = booleanE.CIota("OR")
	BOOL_UNKNOWN = booleanE.CIota("UNKNOWN")
)


type EnumInterface interface{}

func EnumConvert(e Enum) EnumInterface {
	return EnumInterface(e)
}
