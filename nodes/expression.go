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
	return base.NewVal(vex.Val)
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
