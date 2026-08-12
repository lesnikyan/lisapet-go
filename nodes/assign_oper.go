package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
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
	// fmt.Printf("Op=Do %T, %v \n", op.right, op.right)
	err2 := op.right.Do(cx)
	if err2 != nil {
		fmt.Println("OpAssign.R error", err2)
		return err2
	}
	rval := GetExprVal(op.right, cx)
	// fmt.Printf("Op=Do#2, Rexp (%T, %v),  rval (%T, %v) \n", op.right, op.right, rval, rval)
	switch lexp := op.left.(type) {
	case *OperColon:
		lexp.Usage = ColonType
	}
	err1 := op.left.Do(cx)
	if err1 != nil {
		return err1
	}
	return AssignVal(cx, op.left, rval)
	// return nil
}

// actual for variable, argument, struct method in left of assign
// prepare and convert val
// result: isTypeOk, convertedVal
func PrepareVal(expType base.TypeId, val any) (bool, any) {
	tv := objects.TypeByVal(val)
	if tv.Id == expType {
		return true, val
	}
	if !base.TypeCompat(expType, tv.Id) {
		// fmt.Printf("PrepV=#2 Vla Not compatible Eq: %v == %v %v\n", tv.Id, expType, tv.Id == expType)
		return false, nil
	}
	cval := objects.ConvertByType(expType, val)
	return true, cval
}

func AssignVal(cx base.Context, lexpr base.Expression, rval any) error {
	// fmt.Printf(" = AssignVal=#0: (%T %v) = (%T, %v) \n", lexpr, lexpr, rval, rval)
	var leftObj any
	switch lexp := lexpr.(type) {
	case *VarExpr:
		// fmt.Printf("OpAsg=#00 VarExrp: %T %v \n", lexp, lexp)
		leftObj = GetVar(lexp, cx)
	case *OperColon:
		xres := lexp.Get()
		if xres == nil {
			return errors.New("oper assign: no target by colon-expression")
		}
		leftObj = xres.V
		// fmt.Printf("OpAsg=#1 OperColon: %T %v \n", leftObj, leftObj)
	case *ColElemExpr:
		leftObj = lexp.Get().V
		// fmt.Printf("OpAsg=#0 ColEl Ex: (%T %v) Elm (%T, %v) \n", lexp, lexp, leftObj, leftObj)
	case *SequenceComma:
		// a, b, c = expr
		// expr: comma-sequence, list, tuple
		targets := make([]*base.Var, len(lexp.Subs))
		for i, ex := range lexp.Subs {
			vex, ok := ex.(*VarExpr)
			if !ok {
				return errors.New("assign: multival, no var in left sequence")
			}
			vr := vex.GetOrNewVar(cx)
			targets[i] = vr
		}
		leftObj = targets
	}
	return SetValTo(leftObj, rval)
}

func SetValTo(leftObj any, rval any) error {
	// fmt.Printf("SetValTo L: %T, %v = R: %T, %v \n", leftObj, leftObj, rval, rval)
	// valTtype := objects.TypeByVal(rval)
	switch target := leftObj.(type) {
	// TODO: col[key] = val
	case *ColElem:
		// since collections is untyped,  we'll check type
		target.Set(rval)
		// TODO: obj.member = val
	case *base.Var:
		// fmt.Printf("SetValTo#2 L: (%T, %v) = R: %T \n", target, target, rval)
		// cval := rval
		if target.StrictType {
			typeOk, cval := PrepareVal(target.Type.Id, rval)
			if !typeOk {
				// bad val
				// vt := objects.TypeByVal(rval)
				// fmt.Printf("SetValTo#4 val conv error: var %v, val: %v, vtype: %s \n", target, rval, vt.Name)
				return errors.New("oper assign: incorrecttype of value in right operand")
			}
			rval = cval
		}
		target.Val = rval
	case []*base.Var:
		var valSet []any
		switch vals := rval.(type) {
		case []any:
			// just vals
			valSet = vals
		case *objects.ListVal:
			valSet = vals.Elems
		case *objects.TupleVal:
			valSet = vals.Elems
		default:
			return errors.New("oper assign: multi, bad type of right val set")
		}
		if len(target) != len(valSet) {
			return errors.New("oper assign: multi, incorrect count of vals")
		}
		for i, vr := range target {
			err := SetValTo(vr, valSet[i])
			if err != nil {
				return err
			}
		}
		// case *ObjectMember:
		// need check type
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

// if math-assign operator was overloaded
// means a += b not equals a = a + b
func isOverloadedMasgn(a any) bool {
	switch a.(type) {
	case *objects.ListVal:
		return true
	case *objects.DictVal:
		return true
	}
	return false
}

func (op *OperBinAssign) OverloadedDo(cx base.Context, lval any, rval any) error {
	switch lval.(type) {
	default:
		ApplyOper(lval, rval, op.Oper.Id)
	}
	return nil
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
	lvv := GetExprVal(op.left, cx)
	rvv := GetExprVal(op.right, cx)
	if isOverloadedMasgn(lvv) {
		return op.OverloadedDo(cx, lvv, rvv)
	}
	subOper, ok := opMAsg[op.Oper.Id]
	if !ok {
		return errors.New("sub oper of math-assign hasn't found")
	}
	res, ok := ApplyOper(lvv, rvv, subOper)
	AssignVal(cx, op.left, res)
	return nil
}
