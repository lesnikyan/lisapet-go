package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
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

type AtDel struct {
	right base.Expression
	Oper  *Oper
	res   any
}

func (op *AtDel) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *AtDel) SetLeft(xp base.Expression) {

}

func (op *AtDel) Get() *base.Val {
	return objects.NullVal
}

func (op *AtDel) Do(cx base.Context) error {
	// fmt.Println("AtDel @!.Do")
	err := DelElem(cx, op.right)
	if err != nil {
		return err
	}
	return nil
}

func DelElem(cx base.Context, arg base.Expression) error {
	switch a := any(arg).(type) {
	case *VarExpr:
		name := a.GetName()
		return cx.DeleteElem(name)
	case *ColElemExpr:
		cc := a
		err := cc.Col.Do(cx)
		if err != nil {
			return err
		}
		err = cc.Key.Do(cx)
		if err != nil {
			return err
		}
		kval := GetExprVal(cc.Key, cx)
		// cc.KVal = kval
		sval := GetExprVal(cc.Col, cx)
		switch cont := sval.(type) {
		case *objects.ListVal:
			index, ok := kval.(int64)
			if !ok {
				return err
			}
			_, err := cont.Delete(index)
			if err != nil {
				return err
			}
			return nil
		case *objects.TupleVal:
			panic("operator @!: tuple is immutable type")
		case *objects.DictVal:
			_, err := cont.Delete(kval)
			return err
		}
	}
	return fmt.Errorf("operator @!: trying to delete incorrect type: %T", arg)
}
