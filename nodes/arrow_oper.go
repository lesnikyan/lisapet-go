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

func (op *LeftArrow) GetIterAssign() *base.Val {
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

func (op *LeftArrow) Do(cx base.Context) error {
	err := op.right.Do(cx)
	if err != nil {
		return err
	}
	err = op.left.Do(cx)
	if err != nil {
		return err
	}
	// TODO:
	// 1) loop-assign: for n <- nn
	// actually its a generator-like expression

	// 2) append: nn <- v

	return nil
}

type ArrowAssign[T ob.ListVal | ob.TupleVal] struct {
	Src    *T          // means right arg
	Target []*base.Var // left arg, can be several

}
