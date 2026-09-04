package nodes

import (
	"errors"
	"fmt"
	"slices"
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
		// fmt.Printf(".join#001:  %T : %v \n", v, v)
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

func SeqJoin[T *objects.ListVal | *objects.TupleVal](inst T, args []any) (any, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("list.join: need 1 arg")
	}
	var sep string
	switch a0 := args[0].(type) {
	case string:
		sep = a0
	case objects.Glif:
		sep = string([]rune{a0})
	default:
		return nil, fmt.Errorf("Bad arg in list.join: %T", a0)
	}
	// fmt.Printf("arg in list.join: %T, %T:`%v` \n", inst, sep, sep)
	res, err := JoinElems(inst, sep)
	if err != nil {
		return nil, errors.Join(errors.New("list.join err"), err)
	}
	return res, nil
}

// ---- type List

func listJoin(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.ListVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance of list in list.join: %T", inst)
	}
	return SeqJoin(src, args)
}

func objFunc(farg any) (base.FuncVal, error) {
	switch fval := farg.(type) {
	case base.FuncVal:
		return fval, nil
	case *base.Type:
		// fmt.Println("Type as func!")
		if fval.Def != nil {
			panic("Custom type constructs not implemented yet!")
		}
		cons := fval.Construct
		if cons == nil {
			return nil, fmt.Errorf("list.map: agr `type` don't have callable construct")
		}
		return cons, nil
	default:
		return nil, fmt.Errorf("list.map: not a function arg %T", fval)
	}
}

func SeqMap(cx base.Context, src []any, farg any) ([]any, error) {
	fn, err := objFunc(farg)
	if err != nil {
		return nil, err
	}
	rr := make([]any, len(src))
	fargs := []any{nil}
	nargs := map[string]any{}
	for i, n := range src {
		// if n, ok := n.(base.Pair)
		fargs[0] = n
		fn.SetArgVals(fargs, nargs)
		err := fn.Do(cx)
		if err != nil {
			return nil, errors.Join(errors.New("error .map "), err)
		}
		fr := fn.Get()
		var nr any
		if fr == nil {
			nr = NullV()
		} else {
			nr = fr.V
		}
		rr[i] = nr
	}
	return rr, nil
}

func listMap(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.ListVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance of list in list.join: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("list.map: need 1 arg")
	}
	rr, err := SeqMap(cx, src.Elems, args[0])
	if err != nil {
		return nil, errors.Join(errors.New("error list.map res:"), err)
	}
	return objects.NewListVal(rr), nil
}

// list -> single val
func listFold(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.ListVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance of list in list.fold: %T", inst)
	}
	if len(args) < 2 {
		return nil, fmt.Errorf("list.fold: need 2 args: start val, func(a, b)")
	}
	start := args[0]
	var foldf base.FuncVal

	switch fval := args[1].(type) {
	case base.FuncVal:
		foldf = fval
	default:
		return nil, fmt.Errorf("list.fold: not a function arg %T", fval)
	}

	fargs := []any{nil, nil}
	nargs := map[string]any{}
	prev := start // prev and final
	for _, n := range src.Elems {
		fargs[0] = prev
		fargs[1] = n
		foldf.SetArgVals(fargs, nargs)
		err := foldf.Do(cx)
		if err != nil {
			return nil, errors.Join(errors.New("error in list.map "), err)
		}
		fr := foldf.Get()
		var nr any
		if fr == nil {
			nr = NullV()
		} else {
			nr = fr.V
		}
		prev = nr
	}
	return prev, nil
}

// [[1,2],[3,4]] >> [1,2,3,4]
func listFlat(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.ListVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance of list in list.flat: %T", inst)
	}
	rr := []any{}

	for _, n := range src.Elems {
		nn, ok := n.(*objects.ListVal)
		if !ok {
			// return nil, fmt.Errorf("Bad elem in list.flat: %T", n)
			rr = append(rr, n)
			continue
		}
		rr = append(rr, nn.Elems...)
	}
	return objects.NewListVal(rr), nil
}

func listSort(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.ListVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance of list in list.sort: %T", inst)
	}
	vals := src.Elems
	rr := make([]any, len(vals))
	for i, v := range vals {
		switch n := v.(type) {
		case bool, byte, int64, float64, rune, string:
			rr[i] = n
		default:
			return nil, fmt.Errorf("Not ordered type of elem in list.sort: %T", n)
		}
	}
	slices.SortFunc(rr, CmpOrd)
	return objects.NewListVal(rr), nil
}

func listReverse(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.ListVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance of list in list.reverse: %T", inst)
	}
	vals := src.Elems
	slen := len(vals)
	maxl := slen - 1
	rr := make([]any, slen)
	for i, v := range vals {
		rr[maxl-i] = v
	}
	return objects.NewListVal(rr), nil
}

// ---- type Tuple

func tupleJoin(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.TupleVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance of list in list.join: %T", inst)
	}
	return SeqJoin(src, args)
}

func tupleMap(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.TupleVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance of list in tuple.join: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("tuple.map: need 1 arg")
	}
	rr, err := SeqMap(cx, src.Elems, args[0])
	if err != nil {
		return nil, errors.Join(errors.New("error tuple.map res:"), err)
	}
	return objects.NewTupleVal(rr), nil
}

// ---- type Dict

// map keys and vals map(func(k, v) >> rk,rv)
func dictMap(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.DictVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance of list in dict.map: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("dict.map: need 1 arg")
	}
	fn, err := objFunc(args[0])
	if err != nil {
		return nil, err
	}
	// fmt.Printf("dmap#fn: %T, %v\n", fn, fn)
	rr := make(map[any]any, len(src.Vmap))
	fargs := make([]any, 2)
	nargs := map[string]any{}
	for k, v := range src.Vmap {
		// if n, ok := n.(base.Pair)
		fargs[0] = k
		fargs[1] = v
		fn.SetArgVals(fargs, nargs)
		err := fn.Do(cx)
		if err != nil {
			return nil, errors.Join(errors.New("error .map "), err)
		}
		fr := fn.Get()
		if fr != nil {
			// fmt.Printf("dmap#fres: %T, %v\n", fr.V, fr.V)
			vals, ok := fr.V.([]any)
			if !ok {
				return nil, errors.Join(errors.New("error dict.map 11"))
			}
			if len(vals) != 2 {
				return nil, errors.Join(errors.New("error dict.map 12"))
			}
			rk, rv := vals[0], vals[1]
			rr[rk] = rv
		}
	}
	return objects.NewDictVal(rr), nil
}

// map of keys
func dictKMap(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.DictVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance in dict.kmap: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("dict.kmap: need 1 arg")
	}

	fn, err := objFunc(args[0])
	if err != nil {
		return nil, err
	}
	// fmt.Printf("dkmap#fn: %T, %v\n", fn, fn)
	rr := make(map[any]any, len(src.Vmap))
	fargs := make([]any, 1)
	nargs := map[string]any{}
	for k, v := range src.Vmap {
		// if n, ok := n.(base.Pair)
		fargs[0] = k
		fn.SetArgVals(fargs, nargs)
		err := fn.Do(cx)
		if err != nil {
			return nil, errors.Join(errors.New("error .kmap "), err)
		}
		fr := fn.Get()
		if fr != nil {
			// fmt.Printf("dkmap#fres: %T, %v\n", fr.V, fr.V)
			rr[fr.V] = v
		}
	}
	return objects.NewDictVal(rr), nil
}

// map of vals
func dictVMap(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.DictVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance in dict.vmap: %T", inst)
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("dict.vmap: need 1 arg")
	}
	fn, err := objFunc(args[0])
	if err != nil {
		return nil, err
	}
	// fmt.Printf("vmap#fn: %T, %v\n", fn, fn)
	rr := make(map[any]any, len(src.Vmap))
	fargs := make([]any, 1)
	nargs := map[string]any{}
	for k, v := range src.Vmap {
		// if n, ok := n.(base.Pair)
		fargs[0] = v
		fn.SetArgVals(fargs, nargs)
		err := fn.Do(cx)
		if err != nil {
			return nil, errors.Join(errors.New("error .vmap "), err)
		}
		fr := fn.Get()
		if fr != nil {
			// fmt.Printf("vmap#fres: %T, %v\n", fr.V, fr.V)
			rr[k] = fr.V
		}
	}
	return objects.NewDictVal(rr), nil
}

func dictKeys(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.DictVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance in dict.keys: %T", inst)
	}
	nn := make([]any, len(src.Vmap))
	i := 0
	for k, _ := range src.Vmap {
		nn[i] = k
		i++
	}
	return objects.NewListVal(nn), nil

}

func dictVals(cx base.Context, inst any, args []any) (any, error) {
	src, ok := inst.(*objects.DictVal)
	if !ok {
		return nil, fmt.Errorf("Bad instance in dict.vals: %T", inst)
	}
	nn := make([]any, len(src.Vmap))
	i := 0
	for _, v := range src.Vmap {
		nn[i] = v
		i++
	}
	return objects.NewListVal(nn), nil
}
