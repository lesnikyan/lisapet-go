package objects

import (
	"errors"

	// "fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

var ErrorBadKey = errors.New("Incorrect key of collection")

type ListVal struct {
	Elems []any
}

func (v *ListVal) Set(key int64, elem any) error {
	if key >= int64(len(v.Elems)) {
		// return ErrorBadKey
		panic("list: set elem: index out of range")
	}
	v.Elems[int(key)] = elem
	return nil
}
func (v *ListVal) GetElem(key int64) (*base.Val, error) {
	if key >= int64(len(v.Elems)) {
		panic("list: get elem: index out of range")
		// return nil, ErrorBadKey
	}
	return &base.Val{V: v.Elems[key]}, nil
}

func (v *ListVal) Add(elem any) {
	v.Elems = append(v.Elems, elem)
}

func (v *ListVal) Len() int64 {
	return int64(len(v.Elems))
}

func NewListVal(elems []any) *ListVal {
	if elems == nil {
		elems = []any{}
	}
	return &ListVal{Elems: elems}
}

// ****************************************

type TupleVal struct {
	Elems []any
}

func (v *TupleVal) GetElem(key int64) (*base.Val, error) {
	if key >= int64(len(v.Elems)) {
		return nil, ErrorBadKey
	}
	return &base.Val{V: v.Elems[key]}, nil
}

func (v *TupleVal) Add(elem any) {
	v.Elems = append(v.Elems, elem)
}

func (v *TupleVal) Len() int64 {
	return int64(len(v.Elems))
}

func NewTupleVal(elems []any) *TupleVal {
	if elems == nil {
		elems = []any{}
	}
	return &TupleVal{Elems: elems}
}

// ****************************************

type DictVal struct {
	Vmap map[any]any
}

func (v *DictVal) GetElem(key any) (*base.Val, error) {
	val, ok := v.Vmap[key]
	if !ok {
		return nil, ErrorBadKey
	}
	return &base.Val{V: val}, nil
}

func (v *DictVal) Set(key any, val any) {
	v.Vmap[key] = val
}

func (v *DictVal) Len() int64 {
	return int64(len(v.Vmap))
}

func NewDictVal(vm map[any]any) *DictVal {
	if vm == nil {
		vm = map[any]any{}
	}
	return &DictVal{Vmap: vm}
}

// ****************************************

type Maybe struct {
	IsNone bool
	Val    any
}

// type Some struct {
// 	Val any
// }

// type None struct {
// }
