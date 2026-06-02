package objects

import dt "github.com/lesnikyan/lisapet-go/lang/datatype"

type Val struct {
	Type dt.DType
	// Val  any // not sure about universal type
}

type TVal interface {
	GetVal() any
	GetType() dt.DType
}

func (v *Val) GetType() dt.DType {
	return v.Type
}

// INT

type ValInt struct {
	Val
	value int64
}

func NewInt(v int64) *ValInt {
	return &ValInt{Val: Val{Type: dt.Int}, value: v}
}

func (vv *ValInt) GetVal() any {
	return vv.value
}

func (vv *ValInt) GetType() dt.DType {
	return vv.Val.Type
}

type ValFloat struct {
	Val
	value float64
}

type Null struct{}

type ValSingle struct {
	i  int64
	f  float64
	b  bool
	s  string
	bt byte
	g  rune
	// n  Null // no need value, TypeNull will be enough
	// Type dt.DType
}

type ValContainer[T ValSingle] struct {
	ss []ValSingle
}

type ListVal struct {
	ss []ValSingle
	cs []ListVal
	ns []any
}

type ValItem[T ValSingle] struct {
	val ValSingle
}

// type ValT[T int64 | float64 | bool | byte] struct {
// 	v T
// }

// type NN[T int64 | float64 | bool | byte | string | []byte] = T

type ValNT[T string | []byte | []int64] struct {
	v  T
	vv []*T
}

// VAR
type Var struct {
	Name string
	Type dt.DType
	val  TVal
}
