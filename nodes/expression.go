package nodes

import (
	"errors"

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
		v = base.CropStr(val, 1, 1)
	case *MString:
		// fmt.Println("ValExp:MT:", val.V
		v = base.CropStr(val.V, 3, 3)
	default:
		v = val
		// fmt.Println("ValExp: T:", val)
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
	Local   bool // should be found in current context only
	TypeExp base.Expression
	IsVar   bool
	// StrictType bool
	cxEl *base.ContextElem
}

func (ex *VarExpr) GetName() string {
	return ex.name
}

func (ex *VarExpr) SetName(name string) {
	ex.name = name
}

func (ex *VarExpr) AsType(cx base.Context) (*base.Type, error) {
	et := cx.GetElem(ex.name)
	if et == nil {
		// type not found
	}
	ft, ok := et.V.(*base.Type)
	if !ok {
		// not type
		return nil, errors.New("struct def: bad type name in `field:type`")
	}
	return ft, nil
}

func (ex *VarExpr) Do(cx base.Context) error {
	ex.IsVar = false
	ex.vr = nil
	ex.cxEl = nil
	var elem *base.ContextElem
	if ex.Local {
		elem = cx.GetLvlElem(ex.name)
	} else {
		elem = cx.GetElem(ex.name)
	}
	// fmt.Printf("Var.Do#0 VarExrp(%s): %T %v %v \n", ex.name, elem, elem, elem == nil)
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
