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
	Name string
}

// --

type GrupDef struct {
	Name string
	Subs []base.Expression
}
