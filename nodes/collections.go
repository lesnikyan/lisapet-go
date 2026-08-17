package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

type ListExpr struct {
	Seq  *SequenceComma
	Subs []base.Expression
	res  *objects.ListVal
}

func (cs *ListExpr) Add(sub base.Expression) {
	if cs.Subs == nil {
		cs.Subs = []base.Expression{}
	}
	fmt.Printf("[] List Add: (%T, %v)  \n", sub, sub)
	cs.Subs = append(cs.Subs, sub)
}

func (cs *ListExpr) Get() *base.Val {
	return base.NewVal(cs.res)
}

// func (op *ListExpr) IsParent() bool {
// 	// don't used as Block by default
// 	return false
// }

func (cs *ListExpr) Do(ctx base.Context) error {
	src := make([]base.Expression, len(cs.Seq.Subs))
	copy(src, cs.Seq.Subs)
	src = append(src, cs.Subs...)
	res := make([]any, len(src))
	// fmt.Println("[] List Do:", len(src))
	for i, ex := range src {
		err := ex.Do(ctx)
		if err != nil {
			return err
		}
		res[i] = GetExprVal(ex, ctx)
	}
	cs.res = objects.NewListVal(res)
	return nil
}

// ===========

type ColElem struct {
	Src  any // ListVal, DictVal, etc
	KVal any // int64, string, etc
}

func (cc *ColElem) Get() *base.Val {
	switch src := cc.Src.(type) {
	case *objects.ListVal:
		index, ok := cc.KVal.(int64)
		if !ok {
			// return errors.New("incorrect type of index in collection-elem expr")
			panic("incorrect type of index in collection-elem expr")
		}
		res, err := src.GetElem(index)
		if err != nil {
			return nil
		}
		return res

	case *objects.TupleVal:
		index, ok := cc.KVal.(int64)
		if !ok {
			panic("incorrect type of index in collection-elem expr")
		}
		res, err := src.GetElem(index)
		if err != nil {
			return nil
		}
		return res

	case string:
		index, ok := cc.KVal.(int64)
		if !ok {
			panic("incorrect type of index in collection-elem expr")
		}
		res := src[int(index)]
		return base.NewVal(res)

	case *objects.DictVal:
		res, err := src.GetElem(cc.KVal)
		if err != nil {
			return nil
		}
		return res
	}
	return nil
}

func (cc *ColElem) Set(val any) error {
	switch src := cc.Src.(type) {
	case *objects.ListVal:
		index, ok := cc.KVal.(int64)
		if !ok {
			return errors.New("incorrect type of index in collection-elem expr")
		}
		err := src.Set(index, val)
		if err != nil {
			return nil
		}
	case *objects.DictVal:
		src.Set(cc.KVal, val)
	}
	return nil
}

type ColElemExpr struct {
	Col    base.Expression // list, dict, tuple, string object
	Key    base.Expression // index or key
	ColRes *ColElem
}

func (cc *ColElemExpr) Get() *base.Val {
	return base.NewVal(cc.ColRes)
}

func (cc *ColElemExpr) Do(ctx base.Context) error {
	// cc.KVal = nil
	// cc.Src = nil
	err := cc.Col.Do(ctx)
	if err != nil {
		return err
	}
	err = cc.Key.Do(ctx)
	if err != nil {
		return err
	}
	kval := GetExprVal(cc.Key, ctx)
	// cc.KVal = kval
	sval := GetExprVal(cc.Col, ctx)
	// cc.Src = sval
	cc.ColRes = &ColElem{Src: sval, KVal: kval}
	return nil
}

// ----  Slice:  collection[ start : end]
type ColSlice struct {
	Col  base.Expression // list, dict, tuple, string object
	Inds *ColonPair      // index or key

	res any // ListVal, TupleVal, string; TODO: xBytes
}

func (cs *ColSlice) Get() *base.Val {
	if cs.res == nil {
		return nil
	}
	return base.NewVal(cs.res)
}

func (cs *ColSlice) Do(cx base.Context) error {
	cs.res = nil
	err := cs.Col.Do(cx)
	if err != nil {
		return errors.New("slice: bad collection expression")
	}
	// colV := cs.Col.Get()
	colv := GetExprVal(cs.Col, nil)
	if colv == nil {
		return errors.New("slice: bad collection expression")
	}
	err = cs.Inds.Do(cx)
	if err != nil {
		fmt.Printf("ColSlice#inds err: %v\n", err)
		return err
	}
	inds := cs.Inds.GetPair()
	var start int64
	var end int64
	i0 := inds[0]
	i1 := inds[1]
	switch i0 := i0.(type) {
	case int64:
		start = i0
	case *objects.EmptyVal:
		start = 0
	default:
		return errors.New("slice: bad index start")
	}
	switch i1 := i1.(type) {
	case int64:
		end = i1
	case *objects.EmptyVal:
		switch col := colv.(type) {
		case *objects.ListVal:
			end = int64(len(col.Elems))
		case *objects.TupleVal:
			end = int64(len(col.Elems))
		case string:
			end = int64(len(col))
		}
	default:
		return errors.New("slice: bad index end")
	}

	// fmt.Printf("ColSlice#src: %T, %v\n", colv, colv)
	switch col := colv.(type) {
	case *objects.ListVal:
		vals := col.Elems[int(start):int(end)]
		cs.res = objects.NewListVal(vals)
	case *objects.TupleVal:
		vals := col.Elems[int(start):int(end)]
		cs.res = objects.NewTupleVal(vals)
	case string:
		val := col[int(start):int(end)]
		cs.res = val
	}
	return nil
}

func NewSlice(col base.Expression, sub *OperColon) *ColSlice {
	inds := sub.GetPair()
	return &ColSlice{Col: col, Inds: inds}
}

//===== Tuple

type TupleExpr struct {
	Seq  *SequenceComma
	Subs []base.Expression
	res  *objects.TupleVal
}

// func (op *TupleExpr) IsParent() bool {
// 	// don't used as Block by default
// 	return false
// }

func (cs *TupleExpr) Add(sub base.Expression) {
	if cs.Subs == nil {
		cs.Subs = []base.Expression{}
	}
	cs.Subs = append(cs.Subs, sub)
}

func (cs *TupleExpr) Get() *base.Val {
	return base.NewVal(cs.res)
}

func (cs *TupleExpr) Do(ctx base.Context) error {
	src := make([]base.Expression, len(cs.Seq.Subs))
	copy(src, cs.Seq.Subs)
	src = append(src, cs.Subs...)
	res := make([]any, len(src))
	// fmt.Println("(,) Tuple Do:", len(src))
	for i, ex := range src {
		err := ex.Do(ctx)
		if err != nil {
			return err
		}
		res[i] = GetExprVal(ex, ctx)
	}
	cs.res = objects.NewTupleVal(res)
	return nil
}

//==== ColonPair

type ColonPair struct {
	Left  base.Expression
	Right base.Expression
	res   base.Pair
}

func (cp *ColonPair) Get() *base.Val {
	return base.NewVal(cp.res)
}

func (cp *ColonPair) GetPair() base.Pair {
	return cp.res
}
func (cp *ColonPair) Do(ctx base.Context) error {
	err1 := cp.Left.Do(ctx)
	if err1 != nil {
		return err1
	}
	err2 := cp.Right.Do(ctx)
	if err1 != nil {
		return err2
	}
	cp.res = base.Pair{GetExprVal(cp.Left, ctx), GetExprVal(cp.Right, ctx)}
	return nil
}

//==== {dict}

type DictExpr struct {
	Seq  *SequenceComma
	Subs []base.Expression
	res  *objects.DictVal
}

func (op *DictExpr) IsParent() bool {
	// don't used as Block by default
	return false
}

func (cs *DictExpr) Add(sub base.Expression) {
	if cs.Subs == nil {
		cs.Subs = []base.Expression{}
	}
	cs.Subs = append(cs.Subs, sub)
}

func (cs *DictExpr) Get() *base.Val {
	return base.NewVal(cs.res)
}

func (cs *DictExpr) Do(ctx base.Context) error {
	// TODO: think about constraints of key type
	src := make([]base.Expression, len(cs.Seq.Subs))
	copy(src, cs.Seq.Subs)
	src = append(src, cs.Subs...)
	res := make(map[any]any, len(src))
	for _, ex := range src {
		opcol, ok := ex.(*OperColon)
		if !ok {
			return errors.New("Incorest subelement of dict expression")
		}
		cpair := opcol.GetPair()
		// fmt.Printf("DictDo1 (%T, %v) \n", cpair, cpair)
		err := cpair.Do(ctx)
		if err != nil {
			return err
		}
		pairVal := cpair.GetPair()
		res[pairVal[0]] = pairVal[1]
	}
	cs.res = objects.NewDictVal(res)
	return nil
}
