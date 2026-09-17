package nodes

import (
	"errors"

	ob "github.com/lesnikyan/lisapet-go/objects"
)

// ====

func MakeIter(src any) SourceIter {
	var iter SourceIter
	switch ss := src.(type) {
	case *ob.ListVal:
		iter = NewListIter(ss.Elems)
	case *ob.TupleVal:
		iter = NewListIter(ss.Elems)
	case *ob.DictVal:
		iter = NewDictIter(ss.Vmap)
	case ob.Bytes:
		iter = NewBytesIter(ss)
	case *ob.NumSeqGen:
		iter = ss // &NumGenIter{Src: ss}
	case []any: // DEBUG
		// fmt.Printf("Multisource %T \n", ss)
		// for i, srcv := range ss {
		// 	fmt.Printf(" msrc#1 %d Src: (%T, %v) \n", i, srcv, srcv)
		// }
		rr := make([]SourceIter, len(ss))
		for i, sub := range ss {
			// fmt.Printf(" msrc#1 %d Src: (%T, %v) \n", i, sub, sub)
			rr[i] = MakeIter(sub)
		}
		iter = NewMultiIter(rr)
	default:
		panic("No iterable source in arrow-iter")
	}
	return iter
}

// ====

func MKeys(s map[any]any) []any {
	r := make([]any, len(s))
	i := 0
	for k := range s {
		r[i] = k
		i++
	}
	return r
}

// ====

type SourceIter interface {
	Init()
	Next() ([]any, error) // index:val | key:val
	Finished() bool       // true if iter has ended
}

// ====

type PairIntAny struct {
	A int
	B any
}

func NewIntAny(a int, b any) *PairIntAny {
	return &PairIntAny{a, b}
}

// ====

type MultiSourceIter struct {
	Sources []SourceIter
	// index   int
	// maxInd  int
}

func (ms *MultiSourceIter) Init() {
	for _, nsrc := range ms.Sources {
		nsrc.Init()
	}
}

func (ms *MultiSourceIter) Next() ([]any, error) {
	res := []any{}
	for _, nsrc := range ms.Sources {
		nn, err := nsrc.Next()
		if err != nil {
			return nil, err
		}
		switch nsrc.(type) {
		case *ListIter, *BytesIter, *ob.NumSeqGen:
			res = append(res, nn[1])
		case *DictIter:
			res = append(res, nn...)
		default:
			res = append(res, nn[1])
		}
	}
	return res, nil
}

func (ms *MultiSourceIter) Finished() bool {
	for _, nsrc := range ms.Sources {
		if nsrc.Finished() {
			return true
		}
	}
	return false
}

func NewMultiIter(iters []SourceIter) *MultiSourceIter {
	return &MultiSourceIter{Sources: iters}
}

// ====

type ListIter struct {
	Src    []any
	index  int
	maxInd int
}

func (it *ListIter) Init() {
	it.index = 0
	it.maxInd = len(it.Src) - 1
}

func (it *ListIter) Finished() bool {
	return it.index > it.maxInd
}

func (it *ListIter) Next() ([]any, error) {
	var r []any
	if it.Finished() {
		return r, errors.New("trying to Next of Finished iterator")
	}
	i := int64(it.index)
	val := it.Src[it.index]
	it.index += 1
	return []any{i, val}, nil
}

func NewListIter(val []any) *ListIter {
	return &ListIter{Src: val, maxInd: len(val) - 1}
}

// ===

type BytesIter struct {
	Src    ob.Bytes
	index  int
	maxInd int
}

func (it *BytesIter) Init() {
	it.index = 0
	it.maxInd = len(it.Src) - 1
}

func (it *BytesIter) Finished() bool {
	return it.index > it.maxInd
}

func (it *BytesIter) Next() ([]any, error) {
	var r []any
	if it.Finished() {
		return r, errors.New("trying to Next of Finished iterator")
	}
	i := int64(it.index)
	val := it.Src[it.index]
	it.index += 1
	return []any{i, val}, nil
}

func NewBytesIter(val ob.Bytes) *BytesIter {
	return &BytesIter{Src: val, maxInd: len(val) - 1}
}

// ===

type DictIter struct {
	MSrc map[any]any
	iter *ListIter
	keys []any
}

func (it *DictIter) Finished() bool {
	return it.iter.Finished()
}

func (it *DictIter) Init() {
	it.iter = NewListIter(MKeys(it.MSrc))
	it.iter.Init()
}

func (it *DictIter) Next() ([]any, error) {
	var r []any
	if it.iter.Finished() {
		return r, errors.New("trying to Next of Finished iterator")
	}
	ival, err := it.iter.Next()
	if err != nil {
		return r, err
	}
	k := ival[1]
	val := it.MSrc[k]
	return []any{k, val}, nil
}

func NewDictIter(val map[any]any) *DictIter {
	return &DictIter{MSrc: val}
}

// ===
