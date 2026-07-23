package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

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
	name string
	vr   *base.Var
}

func (ex *VarExpr) Do(cx base.Context) error {
	vr := cx.GetVar(ex.name)
	fmt.Printf("Var.Do#0 VarExrp: %T %v %v \n", vr, vr, vr == nil)
	if vr == nil {
		ex.vr = nil
		return nil
	}
	ex.vr = vr
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

func (ex *VarExpr) NewVar(cx base.Context) {
	// TODO: take and set Type
	vr := &base.Var{Name: ex.name}
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
	fmt.Printf("GetExprVal#0: %T, %v\n", v, v)
	var eVal *base.Val
	switch vv := v.(type) {
	case *VarExpr:
		vr := GetVar(vv, cx)
		if vr == nil {
			return nil
		}
		return vr.Val
	case *ValExpr:
		eVal = vv.Get()
		// fmt.Printf("GetExprVal#Val: %T, %v\n", eVal.V, eVal.V)
		// return eVal.V
	case *ColElemExpr:
		eVal = vv.ColRes.Get()
		fmt.Printf("ColElem#Val: %T, %v\n", eVal, eVal)
		// return eVal.V
	case *NumSeqExpr:
		return vv.res.GetList()
	default:
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
