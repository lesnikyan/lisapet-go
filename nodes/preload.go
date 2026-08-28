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
	AddType(cx, "bytes", base.TypeBytes)
	AddType(cx, "float", base.TypeFloat)
	AddType(cx, "string", base.TypeString)
	AddType(cx, "list", base.TypeList)
	AddType(cx, "tuple", base.TypeTuple)
	AddType(cx, "dict", base.TypeDict)
	AddType(cx, "maybe", base.TypeMaybe)
	AddType(cx, "glif", base.TypeGlif)
	// AddType(cx, "enum", )
	// AddType(cx, "grup", )

}
