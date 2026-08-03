package datatype

type DType int

const (
	Undefined DType = iota
	Type
	Noval
	Null
	Any
	Num // all numeric, have operators + - * / ** ^/
	Bool
	Int
	Float
	Ratio   // (int, int) means" int/int
	Complex // (float + float*j) or (ratio + ratio * j)
	Byte    // int8
	String
	Bytes     // []byte
	Glif      // rune
	Container // list dict tuple maybe
	List
	Tuple
	Dict
	Maybe
	Regexp
	Func
	Struct
	Property
	Iterator
	Generator
	Chan
	Enum
	Grup

	Module
)

type Mtype[T int64 | float64 | string | byte | bool] struct {
	t T
}
