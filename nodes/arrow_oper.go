package nodes

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
	ob "github.com/lesnikyan/lisapet-go/objects"
)

// **********************************
type LeftArrow struct {
	left   base.Expression
	right  base.Expression
	IsIter bool
	iter   *IterAssign // result if assign in loop
	append *AppendOper // result if append
	res    any
}

func (op *LeftArrow) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *LeftArrow) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *LeftArrow) AsIter() {
	op.IsIter = true
}

func (op *LeftArrow) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *LeftArrow) DoAppend(cx base.Context) error {
	lvv := GetExprVal(op.left, cx)
	rvv := GetExprVal(op.right, cx)
	switch targ := lvv.(type) {
	case *ob.ListVal:
		targ.Add(rvv)
	case *ob.DictVal:
		switch rval := rvv.(type) {
		// rvv should be a tuple(2) or dict
		case *ob.TupleVal:
			if rval.Len() != 2 {
				return errors.New("incorrect right tuple in dict:append")
			}
			k := rval.Elems[0]
			v := rval.Elems[1]
			targ.Set(k, v)
		case *ob.DictVal:
			// append dict
			for k, v := range rval.Vmap {
				targ.Set(k, v)
			}
		}
	default:
		return errors.New("trying append to non-collection")
	}
	return nil
}

func (op *LeftArrow) GetIterAssign() *IterAssign {
	return op.iter
}

func (op *LeftArrow) Do(cx base.Context) error {
	err := op.right.Do(cx)
	if err != nil {
		return err
	}
	err = op.left.Do(cx)
	if err != nil {
		return err
	}
	if op.IsIter {
		// 1) loop-assign: for n <- nn
		// actually its a generator-like expression
		targ := GetExprTarget(op.left, cx)
		src := GetExprVal(op.right, cx)
		// fmt.Printf(" >>>  LArr.Do0.iter (%T, %v) Exp: (%T, %v) \n", targ, targ, src, src)
		var iter SourceIter
		switch ss := src.(type) {
		case *ob.ListVal:
			iter = NewListIter(ss.Elems)
		case *ob.DictVal:
			iter = NewDictIter(ss.Vmap)
		}
		// fmt.Printf(" >>>  LArr.Do1.iter (%T, %v) Exp: (%T, %v) \n", iter, iter, src, src)
		op.iter = NewIterAssign(iter, targ)

	} else {
		// TODO:
		// 2) append: nn <- v

	}

	return nil
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
	Next() (Pair, error) // index:val | key:val
	Finished() bool      // true if iter has ended
}

// ====

type PairIntAny struct {
	A int
	B any
}

func NewIntAny(a int, b any) *PairIntAny {
	return &PairIntAny{a, b}
}

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

func (it *ListIter) Next() (Pair, error) {
	var r Pair
	if it.Finished() {
		return r, errors.New("trying to Next of Finished iterator")
	}
	i := it.index
	val := it.Src[it.index]
	it.index += 1
	return Pair{i, val}, nil
}

func NewListIter(val []any) *ListIter {
	return &ListIter{Src: val, maxInd: len(val) - 1}
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

func (it *DictIter) Next() (Pair, error) {
	var r Pair
	if it.iter.Finished() {
		return r, errors.New("trying to Next of Finished iterator")
	}
	ival, err := it.iter.Next()
	if err != nil {
		return r, err
	}
	ind, ok := ival[1].(int)
	if !ok {
		return r, errors.New("trying to use non-int index")
	}
	k := it.keys[ind]
	val := it.MSrc[k]
	return Pair{k, val}, nil
}

func NewDictIter(val map[any]any) *DictIter {
	return &DictIter{MSrc: val}
}

// ===

type IterAssign struct {
	Src        SourceIter // means right arg: ListVal, TupleVal, DictVal; next: set of them
	listSource []any
	Target     []*base.Var // left arg, can be several

	// iteration index
	index    int
	maxInd   int
	dictKeys []any
}

func (it *IterAssign) Init() {
	it.Src.Init()
}

func (it *IterAssign) Finished() bool {
	return it.Src.Finished()
}

func (it *IterAssign) Next() error {
	vals, err := it.Src.Next()
	if err != nil {
		return err
	}
	// check val ciunt for multi source: for a, b, c <- nn, cc, vv
	switch len(it.Target) {
	case 0:
		return errors.New("IterAssign: incorrect empty result of iter")
	case 1:
		// TODO: novar: for _ <- src
		v := vals[1]
		it.Target[0].Val = v
	case 2:
		it.Target[0].Val = vals[0]
		it.Target[1].Val = vals[1]
	default:
		if len(it.Target) != len(vals) {
			return errors.New("IterAssign: number of left and right args is not equal")
		}
		for i, v := range vals {
			it.Target[i].Val = v
		}
	}
	return nil
}

func NewIterAssign(src SourceIter, targ any) *IterAssign {
	var left []*base.Var
	switch tt := targ.(type) {
	case *base.Var:
		left = []*base.Var{tt}
	case []*base.Var:
		left = tt
	}
	return &IterAssign{Src: src, Target: left}
}

// =====
type AppendOper struct {
	Target any // *T of [T ob.ListVal | ob.DitcVal]
	Src    any // right arg
}
