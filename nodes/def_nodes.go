package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

// --

type EnumDef struct {
	Name string
	Subs []base.Expression // VarExpr, Assign
	res  *objects.Enum
}

func (en *EnumDef) Do(cx base.Context) error {
	// define enum
	enum := &objects.Enum{Name: en.Name, Vals: map[string]int64{}}
	counter := int64(0)
	for _, exp := range en.Subs {
		switch sx := exp.(type) {
		case *VarExpr:
			sname := sx.GetName()
			enum.Vals[sname] = counter
			counter++
		case *OperAssign:
			xname, ok := sx.Left.(*VarExpr)
			if !ok {
				return fmt.Errorf("enum def: bad name expr of sub elem %T", sx.Left)
			}
			vex, ok := sx.Right.(*ValExpr)
			if !ok {
				return fmt.Errorf("enum def: bad valaue expr of sub elem %T", sx.Right)
			}
			vnum, ok := vex.Val.(int64)
			if !ok {
				return fmt.Errorf("enum def: bad number expr of sub elem %T", vex.Val)
			}
			enum.Vals[xname.GetName()] = vnum
			counter = vnum + 1
		}
	}

	//put enum into context
	cx.AddConst(en.Name, enum)

	return nil
}

func (en *EnumDef) Get() *base.Val {
	return nil
}

func (en *EnumDef) IsParent() bool {
	return true
}

func (en *EnumDef) Add(exp base.Expression) {
	en.Subs = append(en.Subs, exp)
}

func NewEnumDef(name string, subs []base.Expression) *EnumDef {
	if subs == nil {
		subs = []base.Expression{}
	}
	return &EnumDef{Name: name, Subs: subs}
}

// ==

type ConstNode struct {
	Assign *OperAssign
}

func (cn *ConstNode) Get() *base.Val {
	return nil
}

// func (cn *ConstNode) 	HasElem(name string) bool{}
// func (cn *ConstNode) 	GetElem(name string) *Val{}
// func (cn *ConstNode) 	GetName() string{}

func (cn *ConstNode) Do(cx base.Context) error {
	err := cn.Assign.Right.Do(cx)
	if err != nil {
		return err
	}
	val := GetExprVal(cn.Assign.Right, nil)
	typeId := objects.TypeIdByVal(val)
	vtype := cx.GetTypeById(typeId)
	if vtype == nil {
		return fmt.Errorf("const val: undefined type")
	}
	vname := ""
	switch lvar := cn.Assign.Left.(type) {
	case *VarExpr:
		vname = lvar.name
	// case *OperColon:
	//	typed name doesn't make sence for const
	default:
		return fmt.Errorf("const: left expression of assign should be a word, %T given", cn.Assign.Left)
	}
	// lvar, ok := cn.Assign.Left.(*VarExpr)
	// if !ok {
	// 	return fmt.Errorf("const: left expression of assign should be a word, %T given", cn.Assign.Left)
	// }
	constElem := &base.Var{Name: vname, Val: val, Type: vtype, Const: true}
	cx.AddVar(constElem)
	return nil
}

// == Grup

type Grup struct {
	Name    string
	Names   map[string]any
	Context base.Context // internal context
}

func (gg *Grup) HasElem(name string) bool {
	_, ok := gg.Names[name]
	return ok
}

func (gg *Grup) GetElem(name string) *base.Val {
	r := gg.Context.GetLvlElem(name)
	if r == nil {
		return nil
	}
	// fmt.Printf("Grup.GetElem %T, %v\n", r.V, r.V)
	return base.NewVal(r.V)
}

func (gg *Grup) GetMember(name string) *GrupMember {
	elem := gg.GetElem(name)
	if elem == nil {
		return nil
	}
	tmem := elem.V
	mem := &GrupMember{Grup: gg, Name: name, Target: tmem}
	return mem
}

func (gg *Grup) GetName() string {
	return gg.Name
}

func NewGrup(name string, subCx base.Context) *Grup {
	return &Grup{Name: name, Names: make(map[string]any), Context: subCx}
}

// --

type GrupMember struct {
	Grup   *Grup
	Name   string
	Target any // var, func, struct def, enum,
}

func (gm *GrupMember) GetName() string {
	return gm.Name
}

func (gm *GrupMember) Get() *base.Val {
	if gm.Target == nil {
		return nil
	}
	return base.NewVal(gm.Target)
}

func (gm *GrupMember) Set(val any) error {
	// ?
	switch tt := gm.Target.(type) {
	case *base.Var:
		// type ?
		tt.Val = val
	default:
		return fmt.Errorf("grup.member: can't set new value to %T", gm.Target)
	}
	return nil
}

// --

type GrupDef struct {
	Name string
	Subs []base.Expression
	Res  *Grup

	// do we need returning result of grup definition?
	// look like no,  because it make sense for func def,
	// but not for struct, vars, consts or enum - they should be accessible by defined names.
}

func (gg *GrupDef) Get() *base.Val {
	return nil
}

func (gg *GrupDef) IsParent() bool {
	return true
}

func (gg *GrupDef) Add(exp base.Expression) {
	if gg.Subs == nil {
		gg.Subs = []base.Expression{}
	}
	gg.Subs = append(gg.Subs, exp)
}

func (gg *GrupDef) Do(cx base.Context) error {
	// define Grup in paren context
	name := gg.Name
	gcx := cx.SubContext()
	grup := NewGrup(name, gcx)
	cx.AddConst(name, grup)

	// defina internal things in grups context
	for _, sub := range gg.Subs {
		// do sub definition, declaration or assignment
		if aop, ok := sub.(*OperAssign); ok {
			// assigning var or const,
			// prevent reassign var in upper context
			if gvar, ok := aop.Left.(*VarExpr); ok {
				gvar.Local = true
			}
		}
		err := sub.Do(gcx)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewGrupDef(name string) *GrupDef {
	return &GrupDef{Name: name}
}
