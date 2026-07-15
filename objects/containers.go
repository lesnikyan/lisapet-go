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
		return ErrorBadKey
	}
	v.Elems[int(key)] = elem
	return nil
}
func (v *ListVal) GetElem(key int64) (*base.Val, error) {
	if key >= int64(len(v.Elems)) {
		return nil, ErrorBadKey
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

// ****************************************

type DictVal struct {
	Vmap map[any]any
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
