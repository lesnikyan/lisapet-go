package nodes

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
)

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
		v = val[1 : len(val)-1]
	default:
		v = val
	}
	return base.NewVal(v)
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

func (ex *VarExpr) Do(cx base.Context) error {
	ex.IsVar = false
	ex.vr = nil
	ex.cxEl = nil

	elem := cx.GetElem(ex.name)
	if elem == nil {
		return nil
	}
	// fmt.Printf("Var.Do#0 VarExrp: %T %v %v \n", vr, vr, vr == nil)
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
		if elem == nil {
			return nil
		}
		// var found:
		if vv.IsVar {
			vr := vv.GetVar()
			val := vr.Val
			return val
		}
		// other objects by name: func, type, enum, group
		switch el := elem.V.(type) {
		case *Function, *NFunc:
			return el
		}
		return "<? No val from VArExpr>"
	case *ValExpr:
		eVal = vv.Get()
		// fmt.Printf("GetExprVal#Val: %T, %v\n", eVal.V, eVal.V)
		// return eVal.V
	case *ColElemExpr:
		eVal = vv.ColRes.Get()
		// fmt.Printf("GetExp.ColElem#Val: %T, %v\n", eVal, eVal)
		// return eVal.V
	case *FuncCall:
		eVal := vv.Get()
		// fmt.Printf("GetExp.FuncCall#Val: %T, %v\n", eVal, eVal)
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
	default:
		// fmt.Printf("GetExprVal#10: %T, %v\n", vv.Get(), vv.Get())
		return vv.Get().V
	}
	if eVal != nil {
		return eVal.V
	}
	return nil
}

// // get point of Val
// func GetPVal(v any) any {
// 	fmt.Printf("GetVal#0: %T, %v\n", v, v)
// 	switch vv := v.(type) {
// 	case *base.Var:
// 		return vv.Val
// 	case *base.Val:
// 		fmt.Printf("GetVal#Val: %T, %v\n", vv.V, vv.V)
// 		return vv
// 	default:
// 		return vv
// 	}
// }
