package objects

import (
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

	Methods []*base.FuncVal
	methMap map[string]int
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
	st.Vals[name] = val
	return nil
}
