package nodes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

// === List funcs ===
// TODO: len, foldl

func Len_blin(cx base.Context, args []any) (any, error) {
	if len(args) != 1 {
		return 0, errors.New("func len: incorrect count of args")
	}
	switch a := args[0].(type) {
	case string:
		return int64(len(a)), nil

	case *objects.ListVal:
		return int64(len(a.Elems)), nil

	case *objects.TupleVal:
		return int64(len(a.Elems)), nil

	case *objects.DictVal:
		return int64(len(a.Vmap)), nil

	}
	return false, errors.New("func len: incorrect arg type")
}

// === String funcs ===
// TODO: join(srcList, delim), split(src, sep), replace(src, old, new)

func Join_blin(cx base.Context, args []any) (any, error) {
	if len(args) != 2 {
		return nil, errors.New("func join: incorrect count of args")
	}
	src, ok := args[0].(*objects.ListVal)
	if !ok {
		return nil, errors.New("func join: 1-st arg must be a list")
	}
	sep, ok := args[1].(string)
	if !ok {
		return nil, errors.New("func join: 2-nd arg must be a string")
	}
	vv := make([]string, len(src.Elems))
	for i, n := range src.Elems {
		s, ok := n.(string)
		if !ok {
			return nil, errors.New("func join: list should contain strings")
		}
		vv[i] = s
	}
	res := strings.Join(vv, sep)
	return res, nil
}

// === Type conversion ===
// TODO: tolist(), toint(), tostr()

// === Other ===

func PreparePrint(args []any) []any {
	pp := make([]any, len(args))
	for i, n := range args {
		var t any = n
		switch v := n.(type) {
		case *objects.ListVal:
			t = v.Elems
		case *objects.TupleVal:
			t = v.Elems
		case *objects.DictVal:
			t = v.Vmap
		default:
			t = n
		}
		pp[i] = t
	}
	return pp
}

func Print_blin(cx base.Context, args []any) (any, error) {
	pp := PreparePrint(args)
	fmt.Println(pp...)
	return nil, nil
}
