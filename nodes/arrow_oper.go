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

func (op *LeftArrow) GetIterAssign() *IterAssign {
	return op.iter
}

// append to list, or update dict
func (op *LeftArrow) DoAppend(cx base.Context) error {
	// fmt.Printf("LAr.DoAppend#0  \n")
	lvv := GetExprVal(op.left, cx)
	rvv := GetExprVal(op.right, cx)
	// fmt.Printf("LAr.DoAppend#0 (%T, %v), (%T, %v) \n", lvv, lvv, rvv, rvv)
	switch targ := lvv.(type) {
	case *ob.ListVal:
		// append to List
		// fmt.Println("Append ListVal")
		targ.Add(rvv)
		op.res = targ
	case ob.Bytes:
		// little hack for append to []bytes
		switch vrx := op.left.(type) {
		case *VarExpr:
			res, err := targ.Add(rvv)
			if err != nil {
				return err
			}
			if vrx.vr == nil {
				return errors.New("err: trying append to unassigned var")
			}
			vrx.vr.Val = res
		}
		op.res = targ
	case *ob.DictVal:
		// update Dict
		switch rval := rvv.(type) {
		// rvv should be a tuple: (k, v) or dict
		case *ob.TupleVal:
			if rval.Len() != 2 {
				return errors.New("incorrect right tuple in dict:append")
			}
			// fmt.Printf("{} <- () [%T, %v],, [%T, %v] \n", rval.Elems[0], rval.Elems[0], rval.Elems[1], rval.Elems[1])
			k := rval.Elems[0]
			v := rval.Elems[1]
			targ.Set(k, v)
		case *ob.DictVal:
			// append dict
			for k, v := range rval.Vmap {
				targ.Set(k, v)
			}
		}
		op.res = targ
	default:
		// fmt.Printf("!!trying append to non-collection: (%T, %v) v(%T, %v) \n", targ, targ, rvv, rvv)
		return errors.New("trying append to non-collection")
	}
	return nil
}

// make iterator
func (op *LeftArrow) DoIter(cx base.Context) error {
	// actually its a generator-like expression
	targ := GetExprTarget(op.left, cx)
	var src any
	// fmt.Printf(" >>>  LArr.Do0.iter (%T, %v) Src: (%T, %v) \n", op.left, op.left, op.right, op.right)
	switch srcExp := op.right.(type) {
	case *NumSeqExpr:
		src = srcExp.Get().V
	default:
		src = GetExprVal(op.right, cx)
	}
	iter := MakeIter(src)
	// fmt.Printf(" >>>  LArr.Do3.iter (%T, %v) Src: (%T, %v) \n", iter, iter, src, src)
	op.iter = NewIterAssign(iter, targ)
	// fmt.Printf(" %T \n", op.iter)
	return nil
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
		return op.DoIter(cx)
	} else {
		// 2) append: nn <- v
		return op.DoAppend(cx)
	}
}

// type NumGenIter struct {
// 	Src *ob.NumSeqGen
// }

// func (it *NumGenIter) Init() {
// 	it.Src.Init()
// }

// func (it *NumGenIter) Finished() bool {
// 	return it.Src.Finished()
// }

// func (it *NumGenIter) Next() (base.Pair, error) {
// 	return it.Src.Next()
// }

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
	var vals []any
	vv, err := it.Src.Next()
	if err != nil {
		// fmt.Printf(" -- ItAs#0  (%T, %v) \n", err, err)
		return err
	}
	vals = vv
	// fmt.Printf(" -- ItAs#1 <%T> %d (%T, %v) \n", it.Src, len(vals), vals[0], vals[0])
	// check val ciunt for multi source: for a, b, c <- nn, cc, vv
	switch len(it.Target) {
	case 0:
		return errors.New("IterAssign: incorrect empty result of iter")
	case 1:
		// TODO: novar: for _ <- src
		v := vals[1]
		it.Target[0].Val = v
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
