package datatype

type DType int

const (
	Undefined DType = iota
	Type
	Null
	Num
	Bool
	Int
	Ratio
	Float
	Complex
	Byte
	Container
	List
	Tuple
	Dict
	String
	Bytes
	Chan
	Glif
	Regexp
	Enum
	Grup
	Func
	Struct
	Property
	Iterator
	Generator
	Any
	Noval

	Module
)

type Mtype[T int64 | float64 | string | byte | bool] struct {
	t T
}
