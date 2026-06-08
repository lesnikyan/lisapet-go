package nodes

import ob "github.com/lesnikyan/lisapet-go/objects"

type Block interface {
	Do(cx *ob.Context)
	Get() any
}

type BlockExpr struct {
	subs []Expression
	res  any
}

func (bk *BlockExpr) Do(cx *ob.Context) {
	var res any = nil
	for _, exp := range bk.subs {
		exp.Do(cx)
	}
	bk.res = res
}
func (bk *BlockExpr) Get() any {
	return bk.res
}

// ValExpr
type ValExpr struct {
	Val any
}

func (vex *ValExpr) Do(cx *ob.Context) {
	// do nothing
}
func (vex *ValExpr) Get() any {
	return vex.Val
}

// VarExpr
type VarExpr struct {
	name string
	vr   *ob.Var
}

func (ex *VarExpr) Do(cx *ob.Context) {
	vr := cx.GetVar(ex.name)
	if vr == nil {
		vr := &ob.Var{Name: ex.name}
		cx.AddVar(vr)
	}
	ex.vr = vr
}
func (ex *VarExpr) Get() any {
	return ex.vr
}
