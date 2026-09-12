package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

func NullV() *objects.Null {
	return &objects.Null{}
}

// *** EMPTY

type EmptyExpr struct {
}

func (op *EmptyExpr) Get() *base.Val {
	return nil
}

func (op *EmptyExpr) Do(cx base.Context) error {
	return nil
}

// *** ValExpr
type ValExpr struct {
	Val any
}

func (vex *ValExpr) Do(cx base.Context) error {
	// do nothing
	return nil
}
func (vex *ValExpr) Get() *base.Val {
	v := vex.Val
	switch val := v.(type) {
	case string:
		// fmt.Println("ValExp: T:", val)
		// v = val[1 : len(val)-1]
		v = base.CropStr(val, 1, 1)
	case *MString:
		// fmt.Println("ValExp:MT:", val.V
		// v = val.V[3 : len(val.V)-3]
		v = base.CropStr(val.V, 3, 3)
	default:
		v = val
	}
	return base.NewVal(v)
}

// ====

type MString struct {
	V string
}

// ====

type NumField struct {
	V []string
}

func (nf *NumField) Do(cx base.Context) error {
	// do nothing
	// nf.V = nil
	return nil
}

// func (nf *NumField) Add(n int64) error {
// 	if nf.V == nil {
// 		nf.V = []int64{n}
// 	}
// 	return nil
// }

func (nf *NumField) Get() *base.Val {
	if nf.V == nil {
		return nil
	}
	return base.NewVal(nf.V)
}

// *** VarExpr
type VarExpr struct {
	name    string
	vr      *base.Var
	TypeExp base.Expression
	IsVar   bool
	// StrictType bool
	cxEl *base.ContextElem
}

func (ex *VarExpr) GetName() string {
	return ex.name
}
func (ex *VarExpr) Do(cx base.Context) error {
	ex.IsVar = false
	ex.vr = nil
	ex.cxEl = nil

	elem := cx.GetElem(ex.name)
	// fmt.Printf("Var.Do#0 VarExrp: %T %v %v \n", elem, elem, elem == nil)
	if elem == nil {
		return nil
	}
	switch vr := elem.V.(type) {
	case *base.Var:
		ex.vr = vr
		ex.IsVar = true
		if ex.TypeExp != nil {
			tex := ex.TypeExp
			err := tex.Do(cx)
			if err != nil {
				return errors.New("VarExp Do err 5")
			}
			tres := tex.Get()
			if tres == nil {
				return errors.New("VarExp Do err 6")
			}
			vtype, ok := tres.V.(*base.Type)
			if !ok {
				return errors.New("VarExp Do err 7")
			}
			vr.Type = vtype
			vr.StrictType = true
		}
	}
	ex.cxEl = elem
	// ex.vr = vr
	return nil
}

func (ex *VarExpr) Get() *base.Val {
	if ex.vr == nil {
		return nil
	}
	return base.NewVal(ex.vr)
}

func (ex *VarExpr) GetVar() *base.Var {
	if ex.vr == nil {
		return nil
	}
	return ex.vr
}

func (ex *VarExpr) GetOrNewVar(cx base.Context) *base.Var {
	if ex.vr == nil {
		ex.NewVar(cx)
	}
	return ex.vr
}

func (ex *VarExpr) GetElem() *base.ContextElem {
	return ex.cxEl
}

func (ex *VarExpr) NewVar(cx base.Context) {
	vr := &base.Var{Name: ex.name}
	cx.AddVar(vr)
	ex.vr = vr
}

func (ex *VarExpr) NewVarTyped(cx base.Context, vtype *base.Type) {
	// TODO: take and set Type
	vr := &base.Var{Name: ex.name, Type: vtype, StrictType: true}
	// fmt.Printf("NewVarTyped#0 : %v // %v // %v \n", ex.name, vtype, vr)
	cx.AddVar(vr)
	ex.vr = vr
}

func NewVarExpr(name string) *VarExpr {
	return &VarExpr{name: name}
}

func GetVar(expr *VarExpr, cx base.Context) *base.Var {
	var vr *base.Var // Var, comma-sequence
	vr = expr.GetVar()
	// fmt.Printf("GetVar#0 VarExrp: %v %v \n", vr, vr == nil)
	if vr == nil {
		expr.NewVar(cx)
		vr = expr.GetVar()
	}
	return vr
}

func GetExprTarget(v base.Expression, cx base.Context) any {
	switch vv := v.(type) {
	case *VarExpr:
		vr := GetVar(vv, cx)
		return vr
	case *SequenceComma:
		// subs := vv.Subs
		r := make([]*base.Var, len(vv.Subs))
		for i, vsub := range vv.Subs {
			vrex, ok := vsub.(*VarExpr)
			if !ok {
				panic("part of target comma separated expr is not a variable")
			}
			vr := GetVar(vrex, cx)
			r[i] = vr
		}
		return r
		// TODO: take vars, return []*Var

	case *ColElemExpr:
		colVal := vv.ColRes.Get()
		if colVal == nil {
			return nil
		}
		return colVal.V // *ColElem: obj[k]
	}
	return nil
}

func GetExprVal(v base.Expression, cx base.Context) any {
	// fmt.Printf("GetExprVal#0: %T, %v\n", v, v)
	var eVal *base.Val
	switch vv := v.(type) {
	case *VarExpr:
		// vr := GetVar(vv, cx)
		elem := vv.GetElem()
		// fmt.Printf("GetExprVal#VarEx: %T, %v >> %T, %v \n", vv, vv, elem, elem)
		if elem == nil {
			return nil
		}
		switch el := elem.V.(type) {
		case *base.Var:
			// fmt.Printf("GExVr#1 Var: %T, %v \n", el, el)
			val := el.Val
			return val
		case *base.Type:
			return el
		case *objects.Function, *NFunc:
			return el
		}
		panic(fmt.Sprintf("Not a valid element found by word-lexem: %T", v))
	case *ValExpr:
		eVal = vv.Get()
		// fmt.Printf("GetExprVal#Val: %T, %v\n", eVal.V, eVal.V)
		// return eVal.V
	case *ColElemExpr:
		// fmt.Printf("GetExp.ColElem#Val: %T, %v\n", eVal, eVal)
		eVal = vv.ColRes.Get()
		// return eVal.V
	case *FuncCall:
		eVal := vv.Get()
		// fmt.Printf("GetExp.FuncCall#Val: %T, %v >> %T, %v \n", eVal, eVal, eVal.V, eVal.V)
		return eVal.V
	case *NumSeqExpr:
		res := vv.res.GetList()
		// fmt.Printf("GetExp.NumSeq# len: %v\n", len(res.Elems))
		return res
	case *SequenceComma:
		vals := vv.GetVals()
		if vals == nil {
			return nil
		}
		return vals.V
		// return vals
	case *StructConstr:
		return vv.Get().V
	case *Brackets:
		return GetExprVal(vv.Sub, cx)
	case *EmptyExpr:
		return &objects.EmptyVal{}
	case *OperDot:
		// fmt.Printf("GetExprVal# operDot: %T, %v\n", vv, vv)
		// fmt.Printf("GetExprVal# `a.b` : %T, %v\n", vv.Get(), vv.Get())
		// fmt.Printf("GetExprVal#  member: %T, %v\n", vv.GetMember(), vv.GetMember())
		member := vv.GetMember()
		var mval *base.Val
		if member == nil {
			// fmt.Printf("GetExprVal# no member %v\n", member)
			mres := vv.Get()
			if mres == nil {
				return nil
			}
			mval = mres
		} else {
			mval = member.Get()
		}

		if mval == nil {
			// fmt.Printf("GetExprVal# no member value %v\n", mval)
			return nil
		}
		// fmt.Printf("GetExprVal# memVal: %T, %v\n", mval.V, mval.V)
		switch val := mval.V.(type) {
		case *objects.Method:
			val.Inst = member.Obj
			return val
		case *MFunc:
			// val.Inst
			return val
		default:
			return val
		}

	default:
		// fmt.Printf("GetExprVal#100: (%T) %T, %v\n", vv, vv.Get(), vv.Get())
		return vv.Get().V
	}
	if eVal != nil {
		return eVal.V
	}
	return nil
}
