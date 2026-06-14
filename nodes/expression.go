package nodes

import "github.com/lesnikyan/lisapet-go/base"

// *** ValExpr
type ValExpr struct {
	Val any
}

func (vex *ValExpr) Do(cx base.Context) error {
	// do nothing
	return nil
}
func (vex *ValExpr) Get() any {
	return vex.Val
}

// *** VarExpr
type VarExpr struct {
	name string
	vr   *base.Var
}

func NewVarExpr(name string) *VarExpr {
	return &VarExpr{name: name}
}

func (ex *VarExpr) Do(cx base.Context) error {
	vr := cx.GetVar(ex.name)
	if vr == nil {
		vr := &base.Var{Name: ex.name}
		cx.AddVar(vr)
	}
	ex.vr = vr
	return nil
}
func (ex *VarExpr) Get() any {
	return ex.vr
}
