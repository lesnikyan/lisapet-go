package nodes

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	ob "github.com/lesnikyan/lisapet-go/objects"
	obb "github.com/lesnikyan/lisapet-go/objects"
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

	// Elvis operator
	if op.Oper.Id == OpElvis {
		res := ElvisRes(lvv, rvv)
		op.res = res
		return nil
	}

	// Other bin operators
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

func ElvisRes(left any, right any) any {
	var con bool
	switch a := left.(type) {
	case bool:
		con = a
	case int64:
		con = a != 0
	case float64:
		con = a != 0.0
	case *ob.Null:
		con = false
	case *ob.StructInst:
		con = true
	case *ob.ListVal:
		con = a.Len() > 0
	case *ob.TupleVal:
		con = a.Len() > 0
	case *ob.DictVal:
		con = a.Len() > 0
	case string:
		con = a != ""
	case byte:
		con = a != byte(0x0)
	case *ob.Maybe:
		con = !a.IsNone()

	}
	if con {
		return left
	}
	return right
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

// ====

type OperType struct {
	left  base.Expression
	right base.Expression
	res   any
	TODO  bool
}

func (op *OperType) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperType) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperType) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *OperType) Do(cx base.Context) error {
	// fmt.Printf("OperType.Do#0: oper :: (%T:%v) (%T:%v) \n", op.left, op.left, op.right, op.right)
	op.left.Do(cx)
	op.right.Do(cx)
	// lop := op.left.Get()
	lvv := GetExprVal(op.left, cx)
	var rtype *base.Type
	switch rvar := op.right.(type) {
	case *VarExpr:
		tname := rvar.GetName()
		rtype = cx.GetType(tname)
		if rtype == nil {
			return fmt.Errorf("Error in OperType: type `%s` not found ", tname)
		}
	case *ValExpr:
		// fmt.Printf(" -- # -- %T, %v\n", rvar, rvar.Val)
		v := rvar.Get()
		if v == nil {
			return fmt.Errorf("Error in OperType: bad type val: %T ", rvar)
		}
		switch v.V.(type) {
		case *obb.Null:
			rtype = cx.GetType("null")
		}
	case *MixedTypeExpr:
		err := rvar.Do(cx)
		if err != nil {
			return err
		}
		res, err := CheckTypeEqual(lvv, rtype)
		if err != nil {
			return err
		}
		op.res = res
		return nil
	default:
		return fmt.Errorf("Error in OperType: bad type expression: %T ", op.right)
	}

	// fmt.Printf("OperType.Do#5: L(%T, %v) <::> R(%T, %v) \n", lvv, lvv, rtype, rtype)
	res, err := CheckTypeEqual(lvv, rtype)
	// fmt.Printf("-- OperType.Do#6: val(%T, %v) :: exp: %v >> %v \n", lvv, obb.TypeIdByVal(lvv), rtype.Id, res)
	if err != nil {
		// fmt.Printf("Error in OperType: L(%v) <%s> R(%v) \n", lvv, op.Oper.Sign, rvv)
		errm := fmt.Errorf("Error in OperType: L(%T: %v) <::> R(%T: %v) ", op.left, op.left, op.right, op.right)
		return errors.Join(errm, err)
	}
	op.res = res
	return nil
}

func CheckTypeEqual[MT *base.Type | *base.MixedType](val any, expType MT) (bool, error) {
	switch mt := any(expType).(type) {
	case *base.MixedType:
		for _, sub := range mt.Types {
			res, err := CheckTypeEqual(val, sub)
			if err != nil {
				return false, err
			}
			if res {
				return true, nil
			}
		}
		return false, nil
	case *base.Type:
		switch tval := val.(type) {
		case *ob.StructInst:
			stype := tval.Type
			if stype == nil {
				return false, fmt.Errorf("Error in type check: struct type: not defined ")
			}
			if stype.Id == mt.Id {
				// simple case - the same type
				return true, nil
			}
			// check parent
			return tval.Def.HasParent(mt.Id), nil
		default:
			vtype := obb.TypeIdByVal(val)
			return vtype == mt.Id, nil
		}
	}
	return false, fmt.Errorf("Error in type check: sttrange case ")
}
