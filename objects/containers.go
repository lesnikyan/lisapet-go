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
		return nil, errors.New("list: get elem: index out of range")
		// return nil, ErrorBadKey
	}
	return &base.Val{V: v.Elems[key]}, nil
}

func (v *ListVal) Add(elem any) {
	v.Elems = append(v.Elems, elem)
}

func (v *ListVal) Delete(index int64) (*base.Val, error) {
	elem, err := v.GetElem(index)
	if err != nil {
		return nil, err
	}
	v.Elems = append(v.Elems[:index], v.Elems[index+1:]...)
	return elem, nil
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

func AnyList[T any](vals []T) []any {
	vv := make([]any, len(vals))
	for i, v := range vals {
		vv[i] = v
	}
	return vv
}

func NewAnyListVal[T any](vals []T) *ListVal {
	vv := AnyList(vals)
	return NewListVal(vv)
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
func (v *DictVal) Delete(key any) (*base.Val, error) {
	elem, err := v.GetElem(key)
	if err != nil {
		return nil, err
	}
	delete(v.Vmap, key)
	return elem, nil
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
	None bool
	Val  any
}

func (mb *Maybe) IsNone() bool {
	return mb.None
}

// type Some struct {
// 	Val any
// }

// type None struct {
// }

func Some(v any) *Maybe {
	return &Maybe{Val: v}
}
func None() *Maybe {
	return &Maybe{}
}
