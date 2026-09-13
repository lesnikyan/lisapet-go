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
	// dev funcs
	BuiltFunc(cx, "devNM", devNamed, nil)
	BuiltFunc(cx, "devORNM", devOrdNamed, nil)
	BuiltFunc(cx, "devDef33", devDefArgs33, nil)
	BuiltFunc(cx, "devDefOrd", devDefOrds, nil)
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
	AddType(cx, "regexp", base.TypeRegexp)
	// AddType(cx, "enum", )
	// AddType(cx, "grup", )
}

func BuiltMethods(cx base.Context) {
	// string
	BuiltMethod(cx, "string", "split", stringSplit)
	BuiltMethod(cx, "string", "replace", stringReplace)
	BuiltMethod(cx, "string", "join", stringJoin)
	BuiltMethod(cx, "string", "bytes", stringBytes)
	BuiltMethod(cx, "string", "glifs", stringGlifs)
	BuiltMethod(cx, "string", "has", stringHas)
	BuiltMethod(cx, "string", "lines", stringLines)
	BuiltMethod(cx, "string", "trim", stringTrim)
	// list
	BuiltMethod(cx, "list", "join", listJoin)
	BuiltMethod(cx, "list", "map", listMap)
	BuiltMethod(cx, "list", "fold", listFold)
	BuiltMethod(cx, "list", "flat", listFlat)
	BuiltMethod(cx, "list", "sort", listSort)
	BuiltMethod(cx, "list", "reverse", listReverse)
	// tuple
	BuiltMethod(cx, "tuple", "join", tupleJoin)
	BuiltMethod(cx, "tuple", "map", tupleMap)
	BuiltMethod(cx, "tuple", "join", tupleJoin)
	// dict
	BuiltMethod(cx, "dict", "map", dictMap)
	BuiltMethod(cx, "dict", "kmap", dictKMap)
	BuiltMethod(cx, "dict", "vmap", dictVMap)
	BuiltMethod(cx, "dict", "keys", dictKeys)
	BuiltMethod(cx, "dict", "vals", dictVals)
	// bytes
	BuiltMethod(cx, "bytes", "bits", bytesBits)
	BuiltMethod(cx, "bytes", "blocks", bytesBlocks)
	BuiltMethod(cx, "bytes", "map", bytesMap)
	BuiltMethod(cx, "bytes", "reverse", bytesReverse)
	BuiltMethod(cx, "bytes", "fold", bytesFold)
	BuiltMethod(cx, "bytes", "nums", bytesNums)
	// regexp
	BuiltMethod(cx, "regexp", "match", regexpMatch)
	BuiltMethod(cx, "regexp", "find", regexpFind)
	BuiltMethod(cx, "regexp", "findSubs", regexpFindSubs)
	BuiltMethod(cx, "regexp", "replace", regexpReplace)
	BuiltMethod(cx, "regexp", "split", regexpSplit)

}
