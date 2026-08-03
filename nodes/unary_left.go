package nodes

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
)

// ===========

type UnaryLeft struct {
	right base.Expression
	Oper  *Oper
	res   any
}

func (op *UnaryLeft) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *UnaryLeft) SetLeft(xp base.Expression) {

}

func (op *UnaryLeft) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *UnaryLeft) Do(cx base.Context) error {

	op.right.Do(cx)
	rvv := GetExprVal(op.right, cx)
	var res any
	var ok bool
	// fmt.Println("UnaryLeft.Do:", op.Oper, rvv)
	res, ok = LeftOper(op.Oper.Id, rvv)
	if !ok {
		return errors.New("Error in unary-Left oper") // TODO: add more informative error
	}
	op.res = res
	return nil
}

func LeftOper(opid Opid, arg any) (any, bool) {
	switch a := arg.(type) {
	case int64:
		switch opid {
		case OpPlus:
			return a, true
		case OpMinus:
			return -a, true
		case OpBitNot:
			return ^a, true
		}
	case float64:
		switch opid {
		case OpPlus:
			return a, true
		case OpMinus:
			return -a, true
		}
	case bool:
		switch opid {
		case OpNot:
			return !a, true
		}
	case string:
		switch opid {
		case OpBitNot:
			// TODO: string formatting: ~"template"
			return "", true
		}
	}
	return nil, false
}
