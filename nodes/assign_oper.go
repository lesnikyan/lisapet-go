package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

type OperAssign struct {
	Oper  string
	left  base.Expression
	right base.Expression
	res   any
}

func (op *OperAssign) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperAssign) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperAssign) Get() *base.Val {
	return base.NewVal(op.res) // make sense for last expression in the Block
}

func (op *OperAssign) Do(cx base.Context) error {
	fmt.Printf("Op=Do %T, %v \n", op.right, op.right)
	err2 := op.right.Do(cx)
	if err2 != nil {
		fmt.Println("OpAssign.R error", err2)
		return err2
	}
	rval := GetExprVal(op.right, cx)
	fmt.Printf("Op=Do#2, Rexp (%T, %v),  rval (%T, %v) \n", op.right, op.right, rval, rval)
	// if rr, ok := rval.(*objects.ListVal); ok {
	// 	fmt.Printf("Op=Do#2, R-list len = %v \n", len(rr.Elems))
	// }
	err1 := op.left.Do(cx)
	if err1 != nil {
		return err1
	}
	AssignVal(cx, op.left, rval)
	return nil
}

func AssignVal(cx base.Context, lexpr base.Expression, rval any) error {
	// rval := GetExprVal(rexpr, cx)
	var leftObj any
	switch lexp := lexpr.(type) {
	case *VarExpr:
		// fmt.Printf("OpAsg=#0 VarExrp: %T %v \n", lexp, lexp)
		leftObj = GetVar(lexp, cx)
	case *ColElemExpr:
		leftObj = lexp.Get().V
		fmt.Printf("OpAsg=#0 ColEl Ex: (%T %v) Elm (%T, %v) \n", lexp, lexp, leftObj, leftObj)
	}
	fmt.Printf("OP= L: %v = R: %T \n", leftObj, rval)
	switch target := leftObj.(type) {
	// TODO: col[key] = val
	case *ColElem:
		target.Set(rval)
		// TODO: obj.member = val
	case *base.Var:
		fmt.Printf("OP=#2 L: %T = R: %T \n", target, rval)
		target.Val = rval
	}
	return nil
}

// ===========

type OperBinAssign struct {
	Oper  *Oper
	left  base.Expression
	right base.Expression
	res   any
}

func (op *OperBinAssign) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperBinAssign) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperBinAssign) Get() *base.Val {
	return base.NewVal(op.res)
}

var opMAsg = map[Opid]Opid{
	OpPlusAssign:    OpPlus,
	OpMinusAssign:   OpMinus,
	OpMultAssign:    OpMult,
	OpDivAssign:     OpDiv,
	OpPercentAssign: OpPercent,
}

func (op *OperBinAssign) Do(cx base.Context) error {
	err1 := op.left.Do(cx)
	if err1 != nil {
		return err1
	}
	err2 := op.right.Do(cx)
	if err2 != nil {
		fmt.Println("OpAssign.R error", err2)
		return err2
	}
	// rval := GetExprVal(op.right, cx)

	subOper, ok := opMAsg[op.Oper.Id]
	if !ok {
		return errors.New("sub oper of math-assign hasn't found")
	}
	lvv := GetExprVal(op.left, cx)
	rvv := GetExprVal(op.right, cx)
	res, ok := ApplyOper(lvv, rvv, subOper)
	AssignVal(cx, op.left, res)
	return nil
}
