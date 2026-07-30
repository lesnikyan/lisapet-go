package nodes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

// === List funcs ===
// TODO: foldl(list|tuple|dict, function): any
// len(list|tuple|dict|string): int, iter([start,] max [[,step]] ): num-gen

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

func Iter_blin(cx base.Context, args []any) (any, error) {
	alen := len(args)
	if alen < 1 {
		return 0, errors.New("func iter: incorrect count of args")
	}
	a1, ok := args[0].(int64)
	if !ok {
		return 0, errors.New("func len: arg1 should be int")
	}
	step := int64(1)
	start := int64(0)
	max := a1
	if alen > 1 {
		start = a1
		a2, ok := args[1].(int64)
		if !ok {
			return 0, errors.New("func len: arg2 should be int")
		}
		max = a2
		if alen > 2 {
			a3, ok := args[2].(int64)
			if !ok {
				return 0, errors.New("func len: arg2 should be int")
			}
			step = a3
		}
	}
	return objects.NewNumSeqGen(start, max, step), nil
}

// === String funcs ===
// join(srcList, delim), split(src, sep), replace(src, old, new)

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

func Split_blin(cx base.Context, args []any) (any, error) {
	if len(args) != 2 {
		return nil, errors.New("func split: incorrect count of args")
	}
	src, ok := args[0].(string)
	if !ok {
		return nil, errors.New("func split: 1-st arg must be a string")
	}
	sep, ok := args[1].(string)
	if !ok {
		return nil, errors.New("func split: 2-nd arg must be a string")
	}
	ss := strings.Split(src, sep)
	res := objects.NewAnyListVal(ss)
	return res, nil
}

// replace(src, find, repl)
func Replace_blin(cx base.Context, args []any) (any, error) {
	if len(args) != 3 {
		return nil, errors.New("func replace: incorrect count of args")
	}
	aa := make([]string, 3)
	for i, n := range args {

		s, ok := n.(string)
		if !ok {
			return nil, errors.New("func replace: args must be string")
		}
		aa[i] = s
	}
	ss := strings.ReplaceAll(aa[0], aa[1], aa[2])
	return ss, nil
}

// === Type conversion ===
// TODO: tolist(), toint(), tostr()

// about glifs
// ? g'D' /D/ !D! {D} /\// /"/
// "abc".g(0) => g'a'

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
