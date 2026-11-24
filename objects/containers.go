package objects

import (
	"errors"
	"fmt"

	// "fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

var ErrorBadKey = errors.New("Incorrect key of collection")

type ListVal struct {
	Elems []any
}

func (v *ListVal) posid(id int) int {
	if id < 0 {
		id = len(v.Elems) + id
	}
	return id
}

func (v *ListVal) Set(ikey int64, elem any) error {
	key := v.posid(int(ikey))
	// if key < 0 {
	// 	key = len(v.Elems) + key
	// }
	if key >= (len(v.Elems)) {
		// return ErrorBadKey
		panic("list: set elem: index out of range")
	}
	v.Elems[int(key)] = elem
	return nil
}
func (v *ListVal) GetElem(ikey int64) (*base.Val, error) {
	// key := int(ikey)
	// if key < 0 {
	// 	key = len(v.Elems) + key
	// }
	key := v.posid(int(ikey))
	if int(key) >= len(v.Elems) {
		return nil, fmt.Errorf("list: get elem: index %d out of range, len = %d ", key, len(v.Elems))
		// return nil, ErrorBadKey
	}
	return &base.Val{V: v.Elems[key]}, nil
}

func (v *ListVal) Add(elem any) {
	v.Elems = append(v.Elems, elem)
}

func (v *ListVal) Delete(index int64) (*base.Val, error) {
	key := v.posid(int(index))
	elem, err := v.GetElem(index)
	if err != nil {
		return nil, err
	}
	v.Elems = append(v.Elems[:key], v.Elems[key+1:]...)
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
	// fmt.Printf(" MB None: %v \n", mb.None)
	return mb.None
}

func Some(v any) *Maybe {
	return &Maybe{Val: v}
}

func None() *Maybe {
	return &Maybe{None: true}
}

// ****************************************

// == Enum

type Enum struct {
	Name string
	Vals map[string]int64
}

func (en *Enum) GetElem(name string) *base.Val {
	val, ok := en.Vals[name]
	if !ok {
		return nil
	}
	return base.NewVal(val)
}

func (en *Enum) HasElem(name string) bool {
	_, ok := en.Vals[name]
	return ok
}

func (en *Enum) GetName() string {
	return en.Name
}

func (en *Enum) GetMember(name string) *EnumMember {
	val := en.GetElem(name)
	if val == nil {
		return nil
	}
	return &EnumMember{Name: name, Val: val}
}

// --

type EnumMember struct {
	Enum *Enum
	Name string

	Val *base.Val
}

func (em *EnumMember) GetName() string {
	return em.Name
}

func (em *EnumMember) Get() *base.Val {
	if em.Val == nil {
		return nil
	}
	return em.Val
}

func (em *EnumMember) Set(val any) error {
	return fmt.Errorf("enum element cant be changed")
}

// GetName() string
// Get() *Val
// Set(any) error
