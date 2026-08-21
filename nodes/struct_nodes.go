package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

// struct StrName(Parent) a, b, c:int
type StructDefExpr struct {
	Name   string
	Fields []base.Expression // var, var:type
	Subs   []base.Expression

	res *objects.StructDef
}

func (op *StructDefExpr) IsParent() bool {
	// don't used as Block by default
	return true
}

func (cs *StructDefExpr) Add(sub base.Expression) {
	if cs.Subs == nil {
		cs.Subs = []base.Expression{}
	}
	cs.Subs = append(cs.Subs, sub)
}

func (se *StructDefExpr) Get() *base.Val {
	if se.res == nil {
		return nil
	}
	return base.NewVal(se.res)
}

func (se *StructDefExpr) StrDefArgs(cx base.Context) ([]*objects.StructField, error) {
	anyT := cx.GetType("any")
	if anyT == nil {
		return nil, errors.New("struct def: any type not defined")
	}
	// fsource := make([]base.Expression, lfd+len(se.Subs))
	// copy(fsource[:lfd], se.Subs)
	lfd := len(se.Fields)
	fsource := make([]base.Expression, lfd)
	copy(fsource, se.Fields)
	fsource = append(fsource, se.Subs...)
	fields := make([]*objects.StructField, len(fsource))
	for i, ex := range fsource {
		var name string
		var ftype *base.Type
		switch fex := ex.(type) {
		case *VarExpr:
			// untyped field, type = any
			name = fex.name
			ftype = anyT
		case *OperColon:
			lexp, ok := fex.Left.(*VarExpr) // field
			if !ok {
				// strange case
				return nil, errors.New("struct def: bad syntax of field name in `field:type`")
			}
			name = lexp.name
			rexp, ok := fex.Right.(*VarExpr) // type
			if !ok {
				// strange case
				return nil, errors.New("struct def: bad syntax of type name in `field:type`")
			}
			et := cx.GetElem(rexp.name)
			if et == nil {
				// type not found
			}
			ft, ok := et.V.(*base.Type)
			if !ok {
				// not type
				return nil, errors.New("struct def: bad type name in `field:type`")
			}
			ftype = ft
		default:
			return nil, fmt.Errorf("struct def: incorrect expression (%T) instead of field", fex)
		}
		var defv any
		// simple check, if Id in range of user-defined types
		if ftype.Id > base.TypeStructBase {
			defv = &objects.Null{}
		} else {
			defv = base.DefaultVal(ftype.Id)
		}
		fld := &objects.StructField{Name: name, Type: ftype, DefVal: defv}
		fields[i] = fld
	}
	return fields, nil
}

func (se *StructDefExpr) Do(cx base.Context) error {
	se.res = nil
	// convert Fields to []*StructField
	fields, err := se.StrDefArgs(cx)
	if err != nil {
		return err
	}
	sdef := objects.NewSructDef(se.Name, fields)
	stype := base.DefineUserType(se.Name, sdef)
	sdef.Id = stype.Id
	cx.AddType(stype)
	return nil
}

// /
// StrName{a:2, b:5}
type StructConstr struct {
	Name   string
	Args   []*OperColon
	ArgMap map[string]any
	Subs   []*OperColon

	res *objects.StructInst
}

func (cs *StructConstr) Add(sub base.Expression) {
	if cs.Subs == nil {
		cs.Subs = []*OperColon{}
	}
	if arg, ok := sub.(*OperColon); ok {
		cs.Subs = append(cs.Subs, arg)
	}
}

func (se *StructConstr) Get() *base.Val {
	if se.res == nil {
		return nil
	}
	return base.NewVal(se.res)
}

func (se *StructConstr) Do(cx base.Context) error {
	se.res = nil
	tel := cx.GetType(se.Name)
	// fmt.Printf("StructConstr.Do %T, %v name=`%s` \n", tel, tel, se.Name)
	if tel == nil {
		return errors.New("struct: type not found")
	}
	if tel.Def == nil {
		return errors.New("struct: type dont have definition")
	}
	sdef, ok := tel.Def.(*objects.StructDef)
	if !ok {
		return errors.New("struct: type definition not a struct")
	}
	lfd := len(se.Args)
	fsource := make([]*OperColon, lfd)
	copy(fsource, se.Args)
	fsource = append(fsource, se.Subs...)
	// fields := make([]*objects.StructField, len(fsource))
	// for i, ex := range fsource {
	args := make(map[string]any)
	for _, arg := range fsource {
		nexp, ok := arg.Left.(*VarExpr)
		if !ok {
			return fmt.Errorf("struct constr: bad field name expr: %T", arg.Left)
		}
		arn := nexp.GetName()
		arg.Right.Do(cx)
		rv := GetExprVal(arg.Right, nil)
		if rv == nil {
			// no val
			return errors.New("struct: no val of expression in: `field: val`")
		}
		sf, ok := sdef.Fields[arn]
		if !ok {
			return fmt.Errorf("struct: incorrect field name `%s`", arn)
		}
		tOk, val := objects.PrepareVal(sf.Type.Id, rv)
		if !tOk {
			return fmt.Errorf("struct: incorrect type for field `%s`", arn)
		}
		args[arn] = val
	}
	inst := sdef.NewInstance(args)
	inst.Type = tel
	se.res = inst

	return nil
}

func NewStructConstr(name string, args []*OperColon) *StructConstr {
	return &StructConstr{Name: name, Args: args}
}
