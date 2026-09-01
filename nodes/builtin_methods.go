package nodes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

func vals2anys[T any](vals []T) []any {
	res := make([]any, len(vals))
	for i, v := range vals {
		res[i] = v
	}
	return res
}

func stringSplit(cx base.Context, inst any, args []any) (any, error) {
	s, ok := inst.(string)
	if !ok {
		return nil, fmt.Errorf("Bad instance of string in string.split: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("No args of in string.split")
	}
	var sep string
	switch sepv := args[0].(type) {
	case string:
		sep = sepv
	case objects.Glif:
		sep = string([]rune{sepv})
	default:
		return nil, fmt.Errorf("Bad separator in string.split: %T", args[0])
	}
	ee := strings.Split(s, sep)
	res := objects.NewListVal(vals2anys(ee))
	return res, nil
}

type GS interface {
	string | objects.Glif
}

// func strReplaceMulti[K GS, T GS](s string, a K, b T) string {
// 	var rep string
// 	switch bb := any(b).(type) {
// 	case string:
// 		rep = bb
// 	case objects.Glif:
// 		rep = string([]rune{bb})
// 	}
// 	var old string
// 	switch aa := any(a).(type) {
// 	case string:
// 		old = aa
// 	case objects.Glif:
// 		old = string([]rune{aa})
// 	}
// 	return strings.ReplaceAll(s, old, rep)
// }

func stringReplace(cx base.Context, inst any, args []any) (any, error) {
	// args[0,1]: string|glif|Regexp, string|glif
	// args[0]: dict
	s, ok := inst.(string)
	if !ok {
		return nil, fmt.Errorf("Bad instance of string in string.replace: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("No args of in string.replace")
	}

	// replacement
	var rep string
	if len(args) > 1 {
		switch a1 := args[1].(type) {
		case string:
			rep = a1
		case objects.Glif:
			rep = string([]rune{a1})
		default:
			return nil, fmt.Errorf("string.replace: replacement mast be string")
		}
	}

	// searching elem
	switch old := args[0].(type) {
	case string:
		res := strings.ReplaceAll(s, old, rep)
		return res, nil
	case objects.Glif:
		res := strings.ReplaceAll(s, string([]rune{old}), rep)
		return res, nil
	case *objects.DictVal:
		return nil, fmt.Errorf("string.replace: dict arg not implemented")
	case *objects.Regexp:
		return nil, fmt.Errorf("string.replace: regexp arg not implemented")
	default:
		return nil, fmt.Errorf("string.replace: pattern mast be string or glif")
	}
}

// vals: ListVal | TupleVal
func JoinElems(src any, sep string) (any, error) {
	var vals []any
	switch a0 := src.(type) {
	case *objects.ListVal:
		vals = a0.Elems
	case *objects.TupleVal:
		vals = a0.Elems
	}
	ss := make([]string, len(vals))
	for i, v := range vals {
		switch n := v.(type) {
		case string:
			ss[i] = n
		case objects.Glif:
			ss[i] = string([]rune{n})
		default:
			return nil, fmt.Errorf(".join: element should be a string or glif, %T given", v)
		}
	}
	return strings.Join(ss, sep), nil
}

func stringJoin(cx base.Context, inst any, args []any) (any, error) {
	// args: list, tuple
	s, ok := inst.(string)
	if !ok {
		return nil, fmt.Errorf("Bad instance of string in string.join: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("sequence.join: need 2 args")
	}
	res, err := JoinElems(args[0], s)
	if err != nil {
		return nil, errors.Join(errors.New("string.join err"), err)
	}
	return res, nil
}

func stringBytes(cx base.Context, inst any, args []any) (any, error) {
	s, ok := inst.(string)
	if !ok {
		return nil, fmt.Errorf("Bad instance of string in string.bytes: %T", inst)
	}
	bb := objects.Bytes(s)
	return bb, nil
}

func stringGlifs(cx base.Context, inst any, args []any) (any, error) {
	s, ok := inst.(string)
	if !ok {
		return nil, fmt.Errorf("Bad instance of string in string.glifs: %T", inst)
	}
	rr := []objects.Glif(s)
	gg := make([]any, len(rr))
	for i, r := range rr {
		gg[i] = objects.Glif(r)
	}
	return objects.NewListVal(gg), nil
}

func stringHas(cx base.Context, inst any, args []any) (any, error) {
	// args[0]: string | glif
	s, ok := inst.(string)
	if !ok {
		return nil, fmt.Errorf("Bad instance of string in string.has: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("sequence.join: need 1 arg")
	}
	switch sub := args[0].(type) {
	case string:
		res := strings.Contains(s, sub)
		return res, nil
	case objects.Glif:
		res := strings.ContainsRune(s, sub)
		return res, nil
	default:
		return nil, fmt.Errorf("Bad instance of string in string.has: %T", inst)
	}
}

func stringLines(cx base.Context, inst any, args []any) (any, error) {
	s, ok := inst.(string)
	if !ok {
		return nil, fmt.Errorf("Bad instance of string in string.lines: %T", inst)
	}
	vv := []any{}
	for n := range strings.Lines(s) {
		vv = append(vv, n)
	}
	return objects.NewListVal(vv), nil
}

func stringTrim(cx base.Context, inst any, args []any) (any, error) {
	// args: string
	s, ok := inst.(string)
	if !ok {
		return nil, fmt.Errorf("Bad instance of string in string.trim: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("sequence.trim: need 1 arg")
	}

	var cuts string
	switch cc := args[0].(type) {
	case string:
		cuts = cc
	case objects.Glif:
		cuts = string([]rune{cc})
	default:
		return nil, fmt.Errorf("Bad arg of string.trim: %T", inst)
	}
	return strings.Trim(s, cuts), nil
}

// ---- type List

// func listJoin(cx base.Context, inst any, args []any) (any, error) {

// func listMap(cx base.Context, inst any, args []any) (any, error) {

// func listFold(cx base.Context, inst any, args []any) (any, error) {

// ---- type Tuple

// func tupleJoin(cx base.Context, inst any, args []any) (any, error) {

// func tupleMap(cx base.Context, inst any, args []any) (any, error) {
