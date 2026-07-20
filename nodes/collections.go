package nodes

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

type ListExpr struct {
	Seq *SequenceComma
	res *objects.ListVal
}

func (cs *ListExpr) Get() *base.Val {
	return base.NewVal(cs.res)
}

func (cs *ListExpr) Do(ctx base.Context) error {
	res := make([]any, len(cs.Seq.Subs))
	for i, ex := range cs.Seq.Subs {
		err := ex.Do(ctx)
		if err != nil {
			return err
		}
		res[i] = objects.GetVal(ex.Get())
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

//==

type TupleExpr struct {
	Seq *SequenceComma
	res *objects.TupleVal
}

func (cs *TupleExpr) Get() *base.Val {
	return base.NewVal(cs.res)
}

func (cs *TupleExpr) Do(ctx base.Context) error {
	res := make([]any, len(cs.Seq.Subs))
	for i, ex := range cs.Seq.Subs {
		err := ex.Do(ctx)
		if err != nil {
			return err
		}
		res[i] = GetExprVal(ex, ctx)
	}
	cs.res = objects.NewTupleVal(res)
	return nil
}

// L : R
type Pair [2]any

type ColonPair struct {
	Left  base.Expression
	Right base.Expression
	res   Pair
}

func (cp *ColonPair) Get() *base.Val {
	return base.NewVal(cp.res)
}

func (cp *ColonPair) GetPair() Pair {
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
	cp.res = Pair{GetExprVal(cp.Left, ctx), GetExprVal(cp.Right, ctx)}
	return nil
}

// {dict}

type DictExpr struct {
	Seq *SequenceComma
	res *objects.DictVal
}

func (cs *DictExpr) Get() *base.Val {
	return base.NewVal(cs.res)
}

func (cs *DictExpr) Do(ctx base.Context) error {
	res := make(map[any]any, len(cs.Seq.Subs)) // TODO: think about constraints of key type
	for _, ex := range cs.Seq.Subs {
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
