package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

type OperTern struct {
	left  base.Expression
	right base.Expression
	Oper  *Oper
	res   any
	TODO  bool
}

func (op *OperTern) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperTern) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperTern) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *OperTern) Do(cx base.Context) error {
	// fmt.Printf("OperTern.Do#0: oper:%v (%T:%v) (%T:%v) \n	", op.Oper, op.left, op.left, op.right, op.right)
	op.left.Do(cx)

	lvv := GetExprVal(op.left, cx)
	cond, ok := lvv.(bool)
	if !ok {
		return fmt.Errorf("oper ternary: bad 	condition part, mast be bool value, but %T ", lvv)
	}

	colp, ok := op.right.(*OperColon)
	if !ok {
		return fmt.Errorf("oper ternary: bad values part, mast be colon, but %T ", op.right)
	}

	var valEx base.Expression
	if cond {
		valEx = colp.Left
	} else {
		valEx = colp.Right
	}
	err := valEx.Do(cx)
	if err != nil {
		return errors.Join(fmt.Errorf("oper ternary: val expr error"), err)
	}
	res := GetExprVal(valEx, nil)

	op.res = res
	return nil
}
