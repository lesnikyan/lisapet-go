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

// == Grup

type Grup struct {
	Name string
}

// --

type GrupDef struct {
	Name string
	Subs []base.Expression
}
