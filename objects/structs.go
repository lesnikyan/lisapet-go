package objects

import (
	"fmt"
	"slices"

	"github.com/lesnikyan/lisapet-go/base"
)

type StructField struct {
	Name string
	Type *base.Type

	DefVal any // 0, null, false, ""
	Val    any
}

// struct definition object
type StructDef struct {
	Name   string
	FNames []string
	Fields map[string]*StructField
	Id     base.TypeId

	Methods []*Method
	methMap map[string]int
}

func (sd *StructDef) HasMethod(name string) bool {
	// TODO: recursive search here and in parents
	if _, ok := sd.methMap[name]; ok {
		return true
	}
	return false
}

func (sd *StructDef) GetMethod(name string) *Method {
	if mind, ok := sd.methMap[name]; ok {
		return sd.Methods[mind]
	}
	return nil
}

func (sd *StructDef) NewInstance(args map[string]any) *StructInst {
	for name, f := range sd.Fields {
		_, ok := args[name]
		if !ok {
			args[name] = f.DefVal
		}
	}
	return &StructInst{Def: sd, Vals: args}
}

func NewSructDef(name string, fields []*StructField) *StructDef {
	// DefineUserType(name, def)
	fm := make(map[string]*StructField, len(fields))
	fnames := make([]string, len(fields))
	for i, fld := range fields {
		fm[fld.Name] = fld
		fnames[i] = fld.Name
	}
	return &StructDef{Name: name, Fields: fm, FNames: fnames}
}

///

type StructInst struct {
	Def *StructDef
	// Type *base.Type

	Vals map[string]any
}

func (st *StructInst) Get(name string) *base.Val {
	// fmt.Printf(" st.Get n=`%s` (%T, %v) \n", name, st.Def, st.Def)
	if !slices.Contains(st.Def.FNames, name) {
		panic("No such field in struct")
	}
	val, ok := st.Vals[name]
	if !ok {
		dval := st.Def.Fields[name].DefVal
		return base.NewVal(dval)
	}
	return base.NewVal(val)
}

func (st *StructInst) Set(name string, val any) error {
	// TODO: check field type and compatibility

	sf, ok := st.Def.Fields[name]
	if !ok {
		return fmt.Errorf("struct: incorrect field name `%s`", name)
	}
	tOk, val := PrepareVal(sf.Type.Id, val)
	if !tOk {
		return fmt.Errorf("struct: incorrect type for field `%s`", name)
	}
	st.Vals[name] = val
	return nil
}

func (st *StructInst) GetMember(name string) *StrMember {
	if slices.Contains(st.Def.FNames, name) {
		return &StrMember{Obj: st, Field: name}
	}
	if st.Def.HasMethod(name) {
		mt := st.Def.GetMethod(name)
		if mt != nil {
			return &StrMember{Obj: st, Method: mt}
		}
	}
	return nil
}

func (st *StructInst) GetMethod(name string) *Method {

	return nil
}

type Method struct {
	Func base.FuncVal
}

// type Member interface {
// 	*StructField
// 	*Method
// }

type StrMember struct {
	Obj    *StructInst
	Field  string
	Method *Method
}

func (mb *StrMember) Get() *base.Val {
	if mb.Field != "" {
		return mb.Obj.Get(mb.Field)
	}
	if mb.Method != nil {
		return base.NewVal(mb.Method)
	}
	return nil
}

func (mb *StrMember) Set(val any) error {
	if mb.Field != "" {
		return mb.Obj.Set(mb.Field, val)
	}
	return nil
}
