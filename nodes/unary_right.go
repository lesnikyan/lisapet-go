package nodes

import (
	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

/*
Unary operator in right side:
	arg...
	seq...
	func~>
*/

// left...
type TripleDots struct {
	Left base.Expression
	res  any
}

func (op *TripleDots) SetLeft(xp base.Expression) {
	op.Left = xp
}

func (op *TripleDots) SetRight(xp base.Expression) {}

func (op *TripleDots) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *TripleDots) Len() int {
	switch coll := op.res.(type) {
	case *objects.ListVal:
		return len(coll.Elems)
	case *objects.TupleVal:
		return len(coll.Elems)
	case *objects.Maybe:
		if coll.IsNone() {
			return 0
		} else {
			return 1
		}
	}
	return 0
}

func (op *TripleDots) Do(cx base.Context) error {
	// fmt.Printf("Oper... Do#0 \n")
	err := op.Left.Do(cx)
	if err != nil {
		return err
	}
	val := GetExprVal(op.Left, nil)
	op.res = val
	return nil
}

// ====

// left ~>
type TildArrow struct {
	left base.Expression
	res  any
}

func (op *TildArrow) SetLeft(xp base.Expression) {
	op.left = xp
}

func (op *TildArrow) SetRight(xp base.Expression) {}

func (op *TildArrow) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *TildArrow) Do(cx base.Context) error {
	return nil
}
