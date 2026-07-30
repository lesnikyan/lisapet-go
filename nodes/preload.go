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

func PreloadTypes(cx base.Context) {

}

func PreloadConstructors(cx base.Context) {

}
