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

// type DeepField struct {
// 	Name string
// 	Field *StructField
// }

// struct definition object
type StructDef struct {
	Name    string
	FNames  []string
	Fields  map[string]*StructField
	Type    *base.Type
	Parents []*StructDef
	// fieldMap map[string]*StructDef

	// Methods []*Method
	Methods map[string]*Method
	// MtNames []string
}

// Note: if we redefine method in child  type, it closes all parents overloads

func (sd *StructDef) AddMethod(fn base.FuncVal) {
	if mt, ok := fn.(*Method); ok {
		sd.Methods[fn.GetName()] = mt
	}
}

func (sd *StructDef) HasMethod(name string) bool {
	// TODO: recursive search here and in parents
	// if _, ok := sd.methMap[name]; ok {
	// 	return true
	// }
	return false
}

// func (sd *StructDef) GetMethod(name string) *Method {

func (sd *StructDef) GetMethod(name string) (base.FuncVal, bool) {
	if mt, ok := sd.Methods[name]; ok {
		return mt, true
	}
	return nil, false
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

func (sd *StructDef) GetMethods() []base.FuncVal {
	pmt := make([]base.FuncVal, 0)
	for _, prt := range sd.Parents {
		mm := prt.GetMethods()
		pmt = append(pmt, mm...)
	}
	for _, fn := range sd.Methods {
		pmt = append(pmt, fn)
	}
	return pmt
}

func (sd *StructDef) GetFields() []*StructField {
	pmt := make([]*StructField, 0)

	// fmt.Printf("StDef GetF. pps: %d \n", len(sd.Parents))
	for _, prt := range sd.Parents {
		// fmt.Printf("StDef GetPf3. t: %s \n", prt.Name)
		mm := prt.GetFields()
		pmt = append(pmt, mm...)
	}
	for _, f := range sd.Fields {
		// fmt.Printf("StDef GetF4. f: %s : %s \n", k, f.Name)
		pmt = append(pmt, f)
	}
	return pmt
}

func (sd *StructDef) InitChild() {
	// fmt.Printf("StDef.InitCh1 t: %s \n", sd.Name)
	// fields
	ff := sd.GetFields()
	for _, fd := range ff {
		// fmt.Printf("StDef.InitCh2 f: %v \n", fd)
		sd.Fields[fd.Name] = fd
		sd.FNames = append(sd.FNames, fd.Name)
	}

	// methods
	mm := sd.GetMethods()
	for _, fn := range mm {
		if met, ok := fn.(*Method); ok {
			sd.Methods[met.GetName()] = met
		}
	}
}

func NewSructDef(name string, fields []*StructField, parents []*StructDef) *StructDef {
	// DefineUserType(name, def)
	fm := make(map[string]*StructField, len(fields))
	fnames := make([]string, len(fields))
	methMap := make(map[string]*Method)
	stype := &StructDef{Name: name, Fields: fm, FNames: fnames, Parents: parents, Methods: methMap}
	stype.InitChild()
	for i, fld := range fields {
		stype.Fields[fld.Name] = fld
		stype.FNames[i] = fld.Name
	}
	return stype
}

///

type StructInst struct {
	Def  *StructDef
	Type *base.Type

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
	sf, ok := st.Def.Fields[name]
	if !ok {
		return fmt.Errorf("struct: incorrect field name `%s`", name)
	}
	tOk, val := PrepareVal(sf.Type, val)
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
	fn, ok := st.Def.GetMethod(name)
	// fmt.Printf(" StrInst (%T, %v) GetMember n`%s` (%T, %v) \n", st.Type, st.Type, name, fn, ok)
	if ok {
		if met, mok := fn.(*Method); mok {
			return &StrMember{Obj: st, Method: met}
		}
	}
	return nil
}

// GetMethod(name string) FuncVal
// GetMethods() []FuncVal

// func (st *StructInst) GetMethod(name string) *Method {
func (st *StructInst) GetMethod(name string) *Method {
	fn, ok := st.Def.GetMethod(name)
	if ok {
		if mt, ok2 := fn.(*Method); ok2 {
			return mt
		}
	}
	return nil
}

// child: fieldMap <- parent.fields methodMap <- parent.methods

// ===============================

type Method struct {
	Func     base.FuncVal // *Function, *NFunc
	Type     *base.Type
	InstName string
	InstArg  *ArgExp
	Inst     any // base val or StructInst
}

func (mt *Method) SetInstVar(cx base.Context) error {
	vr := &base.Var{Name: mt.InstName, Type: mt.Type, StrictType: true}
	cx.AddVar(vr)
	return nil
}

func (mt *Method) Do(cx base.Context) error {
	return mt.Func.Do(cx)
}

func (mt *Method) InstVarN(name string) string {
	return fmt.Sprintf("inst##%s", name)
}

func (mt *Method) SetArgVals(vals []any, mvals map[string]any) {
	mvals[mt.InstVarN(mt.InstName)] = mt.Inst
	mt.Func.SetArgVals(vals, mvals)
}

func (mt *Method) GetName() string {
	return mt.Func.GetName()
}

func (mt *Method) Get() *base.Val {
	return mt.Func.Get()
}

func NewMethod(fn base.FuncVal, ntype *base.Type, iname string) *Method {
	// instArg :=
	mt := &Method{Func: fn, InstName: iname, Type: ntype}
	switch fnn := fn.(type) {
	case *Function:
		instArg := &ArgExp{Name: mt.InstVarN(iname), Type: ntype, StrictType: true}
		fnn.Block.AddArg(instArg)

		// case *NFunc:
	}

	return mt
	// return &Method{Func: fn, Type: ntype}
}

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
