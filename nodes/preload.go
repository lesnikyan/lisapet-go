package nodes

import "github.com/lesnikyan/lisapet-go/base"

func PreloadFuncs(cx base.Context) {
	BuiltFunc(cx, "len", Len_blin, nil)
	BuiltFunc(cx, "join", Join_blin, nil)
	BuiltFunc(cx, "split", Split_blin, nil)
	BuiltFunc(cx, "replace", Replace_blin, nil)
	BuiltFunc(cx, "print", Print_blin, nil)
	BuiltFunc(cx, "iter", Iter_blin, nil)
	BuiltFunc(cx, "some", constr_some, nil)
}

func PreloadConstr(cx base.Context) {
	BuiltConstr(cx, "string", constr_string, nil)
	BuiltConstr(cx, "byte", constr_byte, nil)
	BuiltConstr(cx, "int", constr_int, nil)
	BuiltConstr(cx, "float", constr_float, nil)
	BuiltConstr(cx, "bool", constr_bool, nil)
	BuiltConstr(cx, "glif", constr_glif, nil)
	BuiltConstr(cx, "bytes", constr_bytes, nil)
	BuiltConstr(cx, "list", constr_list, nil)
	BuiltConstr(cx, "tuple", constr_tuple, nil)
	BuiltConstr(cx, "dict", constr_dict, nil)
}

func AddType(cx base.Context, name string, id base.TypeId) {
	tp := base.BaseType(name, id)
	cx.AddType(tp)
}

func PreloadTypes(cx base.Context) {
	AddType(cx, "any", base.TypeAny)
	AddType(cx, "null", base.TypeNull)
	AddType(cx, "bool", base.TypeBool)
	AddType(cx, "int", base.TypeInt)
	AddType(cx, "byte", base.TypeByte)
	AddType(cx, "float", base.TypeFloat)
	AddType(cx, "glif", base.TypeGlif)
	AddType(cx, "string", base.TypeString)
	AddType(cx, "list", base.TypeList)
	AddType(cx, "tuple", base.TypeTuple)
	AddType(cx, "dict", base.TypeDict)
	AddType(cx, "bytes", base.TypeBytes)
	AddType(cx, "maybe", base.TypeMaybe)
	// AddType(cx, "enum", )
	// AddType(cx, "grup", )
}

func BuiltMethods(cx base.Context) {
	BuiltMethod(cx, "string", "split", stringSplit)
	BuiltMethod(cx, "string", "replace", stringReplace)
	BuiltMethod(cx, "string", "join", stringJoin)
	BuiltMethod(cx, "string", "bytes", stringBytes)
	BuiltMethod(cx, "string", "glifs", stringGlifs)
	BuiltMethod(cx, "string", "has", stringHas)
	BuiltMethod(cx, "string", "lines", stringLines)
	BuiltMethod(cx, "string", "trim", stringTrim)
}
