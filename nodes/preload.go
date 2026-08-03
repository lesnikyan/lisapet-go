package nodes

import "github.com/lesnikyan/lisapet-go/base"

func PreloadFuncs(cx base.Context) {
	BuiltFunc(cx, "len", Len_blin, nil)
	BuiltFunc(cx, "join", Join_blin, nil)
	BuiltFunc(cx, "split", Split_blin, nil)
	BuiltFunc(cx, "replace", Replace_blin, nil)
	BuiltFunc(cx, "print", Print_blin, nil)
	BuiltFunc(cx, "iter", Iter_blin, nil)
}

func AddType(cx base.Context, name string) {
	tp := base.DefineType(name, false)
	cx.AddType(tp)
}

func PreloadTypes(cx base.Context) {
	AddType(cx, "any")
	AddType(cx, "null")
	AddType(cx, "int")
	AddType(cx, "float")
	AddType(cx, "bool")
	AddType(cx, "string")
	AddType(cx, "list")
	AddType(cx, "tuple")
	AddType(cx, "dict")
	// AddType(cx, "bytes", )
	// AddType(cx, "glif", )
	// AddType(cx, "enum", )
	// AddType(cx, "grup", )

}

func PreloadConstructors(cx base.Context) {

}
