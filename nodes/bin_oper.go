package nodes

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	ob "github.com/lesnikyan/lisapet-go/objects"
)

type MockExpr struct {
}

func (cs *MockExpr) Get() *base.Val { return nil }

func (cs *MockExpr) Do(ctx base.Context) error { return nil }

func SeqInfo(seq any) string {
	sep := " "
	subinf := []string{}
	switch sqq := seq.(type) {
	case *SequenceComma:
		sep = ", "
		for _, elm := range sqq.Subs {
			subinf = append(subinf, OperArgsInfo(elm))
		}
	}
	return strings.Join(subinf, sep)
}

func OperArgsInfo(expr any) string {
	var aL string
	var aR string
	var opStr string
	switch exx := expr.(type) {
	case *OperAssign:
		aL = OperArgsInfo(exx.Left)
		aR = OperArgsInfo(exx.Right)
		opStr = exx.Oper
	case *OperBin:
	case *OperDot:
		aL = OperArgsInfo(exx.left)
		aR = OperArgsInfo(exx.right)
		opStr = "."
	case *UnaryLeft:
		aL = "unar"
		aR = OperArgsInfo(exx.right)
		opStr = exx.Oper.Sign
	case *OperColon:
		aL = OperArgsInfo(exx.Left)
		aR = OperArgsInfo(exx.Right)
		opStr = exx.Oper.Sign
	case *LeftArrow:
		aL = OperArgsInfo(exx.left)
		aR = OperArgsInfo(exx.right)
		opStr = "<-"
	case *Brackets:
		subs := OperArgsInfo(exx.Sub)
		return fmt.Sprintf("(%v)", subs)
	case *TupleExpr:
		subs := OperArgsInfo(exx.Seq)
		return fmt.Sprintf("tuple(%v)", subs)
	case *ListExpr:
		subs := OperArgsInfo(exx.Seq)
		return fmt.Sprintf("list[%v]", subs)
	case *DictExpr:
		subs := OperArgsInfo(exx.Seq)
		return fmt.Sprintf("dict{%v}", subs)
	case *SequenceComma:
		return SeqInfo(expr)
	case *ValExpr:
		return fmt.Sprintf("%v", exx.Val)
	case *VarExpr:
		return fmt.Sprintf("(%s:)", exx.name)
	default:
		return "<non-oper object>"
	}
	return fmt.Sprintf("<%s>{L: %s, R: %s}", opStr, aL, aR)
}

// ===========

type OperBin struct {
	left  base.Expression
	right base.Expression
	Oper  *Oper
	res   any
	TODO  bool
}

func (op *OperBin) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperBin) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperBin) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *OperBin) Do(cx base.Context) error {
	if op.TODO {
		return nil
	}
	// fmt.Printf("OperBin.Do#0: oper:%v (%T:%v) (%T:%v) \n	", op.Oper, op.left, op.left, op.right, op.right)
	op.left.Do(cx)
	op.right.Do(cx)
	// lop := op.left.Get()
	lvv := GetExprVal(op.left, cx)
	// rop := op.right.Get()
	rvv := GetExprVal(op.right, cx)

	res, ok := ApplyOper(lvv, rvv, op.Oper.Id)
	if !ok {
		// fmt.Printf("Error in bin oper: L(%v) <%s> R(%v) \n", lvv, op.Oper.Sign, rvv)
		errm := fmt.Errorf("Error in bin oper: L(%T: %v) <%s> R(%T: %v) ", op.left, op.left, op.Oper.Sign, op.right, op.right)
		return errm // TODO: add more informative error
	}
	op.res = res
	return nil
}

func ApplyOper(left any, right any, oper Opid) (any, bool) {
	if oper == OpEqual || oper == OpNotEqual {
		return EqCompare(left, right, oper), true
	}
	var res any
	var ok bool
	// fmt.Printf("ApplyOper#0: <%v> (%T, %v) (%T, %v) \n", oper, left, left, right, right)
	// Do operators by type of left operand
	switch val := left.(type) {
	case int64:
		res, ok = binOperInt(oper, val, right)
	case float64:
		res, ok = binOperFloat(oper, val, right)
	case byte:
		res, ok = binOperByte(oper, val, right)
	case bool:
		res, ok = binOperBool(oper, val, right)
	case string:
		res, ok = binOperString(oper, val, right)
	case *ob.ListVal:
		res, ok = binOperList(oper, val, right)
	case *ob.TupleVal:
		res, ok = binOperTuple(oper, val, right)
	case *ob.DictVal:
		res, ok = binOperDict(oper, val, right)
	case ob.Bytes:
		res, ok = binOperBytes(oper, val, right)
	case *ob.Regexp:
		res, ok = binOperRegexp(oper, val, right)

	}
	return res, ok
}

// ===============

type OperIn struct {
	left  base.Expression
	right base.Expression
	Oper  *Oper
	res   any
	TODO  bool
}

func (op *OperIn) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperIn) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperIn) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *OperIn) Do(cx base.Context) error {
	if op.TODO {
		return nil
	}
	// fmt.Printf("OperIn.Do#0: oper:%v (%T:%v) (%T:%v) \n	", op.Oper, op.left, op.left, op.right, op.right)
	op.left.Do(cx)
	op.right.Do(cx)
	// lop := op.left.Get()
	lvv := GetExprVal(op.left, cx)
	// rop := op.right.Get()
	rvv := GetExprVal(op.right, cx)

	res, err := CheckIn(lvv, rvv, op.Oper)
	if err != nil {
		// fmt.Printf("Error in OperIn: L(%v) <%s> R(%v) \n", lvv, op.Oper.Sign, rvv)
		errm := fmt.Errorf("Error in OperIn: L(%T: %v) <%s> R(%T: %v) ", op.left, op.left, op.Oper.Sign, op.right, op.right)
		return errors.Join(errm, err)
	}
	op.res = res
	return nil
}

func CheckIn(left any, right any, oper *Oper) (bool, error) {
	var res bool
	switch container := right.(type) {
	case *ob.ListVal:
		res = slices.Contains(container.Elems, left)
	case *ob.TupleVal:
		res = slices.Contains(container.Elems, left)
	case *ob.DictVal:
		_, ok := container.Vmap[left]
		res = ok
	case *ob.Maybe:
		res = !container.IsNone() && container.Val == left
	default:
		return false, fmt.Errorf("Error in CheckIn: bad right arg: `%T` ", right)
	}

	switch oper.Id {
	case OpIn:
		return res, nil
	case OpNotIn:
		return !res, nil
	}
	return false, fmt.Errorf("Error in CheckIn: bad oper: `%s` ", oper.Sign)
}
