package nodes

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"unicode/utf8"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

func constr_string(cx base.Context, args []any) (any, error) {
	// fmt.Printf("constr string: len=%d %T, %v \n", len(args), args, args)
	if len(args) != 1 {
		if len(args) == 0 {
			return "", nil
		}
		return 0, errors.New("constr string: incorrect count of args")
	}
	switch a := args[0].(type) {
	case string:
		return a, nil
	case int64:
		return fmt.Sprintf("%d", a), nil
	case byte:
		// rr := []rune{rune(a)}
		// return string(rr), nil
		return fmt.Sprintf("%d", a), nil
	case objects.Glif:
		rr := []rune{a}
		return string(rr), nil

	case objects.Bytes:
		return string(a), nil

	case *objects.Null:
		return "null", nil

	case bool:
		var s string
		if a {
			s = "true"
		} else {
			s = "false"
		}
		return s, nil

	case *objects.ListVal:
		rr := make([]rune, len(a.Elems))
		for i, s := range a.Elems {
			var r rune
			switch x := s.(type) {
			case int64:
				r = rune(x)
			case byte:
				r = rune(x)
			case objects.Glif:
				r = x
			}
			rr[i] = r
		}
		return string(rr), nil
	}
	return false, errors.New("constr string: incorrect arg type")
}

func constr_byte(cx base.Context, args []any) (any, error) {
	if len(args) != 1 {
		return 0, errors.New("constr byte: incorrect count of args")
	}
	switch a := args[0].(type) {
	case int64:
		return byte(a % 0x100), nil
	case bool:
		s := byte(0)
		if a {
			s = byte(1)
		}
		return s, nil
	}
	return false, errors.New("constr byte: incorrect arg type")
}

func constr_some(cx base.Context, args []any) (any, error) {
	if len(args) != 1 {
		return 0, errors.New("constr some: incorrect count of args")
	}
	if args[0] != nil {
		return objects.Some(args[0]), nil
	}
	// switch a := args[0].(type) {
	// }
	return false, errors.New("constr byte: incorrect arg type")
}

var validIntBases = []int{2, 8, 10, 16}

func constr_int(cx base.Context, args []any) (any, error) {
	if len(args) < 1 {
		return 0, errors.New("constr int: incorrect count of args")
	}
	ordd, named := SplitNamed(args)
	switch a := ordd[0].(type) {
	case int64:
		return a, nil
	case bool:
		s := int64(0)
		if a {
			s = int64(1)
		}
		return s, nil
	case *objects.Null:
		return int64(0), nil
	case string:
		base := 10
		if bb, ok := named["base"]; ok {
			if bint, ok := bb.(int64); ok {
				bval := int(bint)
				if slices.Contains(validIntBases, bval) {
					base = bval
				}
			}
		}
		return strconv.ParseInt(a, base, 64)
	case float64:
		return int64(a), nil
	case objects.Glif:
		return int64(a), nil
	case byte:
		return int64(a), nil
	case objects.Bytes:
		v := 0
		// TODO: convert bytes to int
		return int64(v), nil
	}
	return 0, nil
}

func constr_float(cx base.Context, args []any) (any, error) {
	if len(args) != 1 {
		return 0, errors.New("constr float: incorrect count of args")
	}
	switch a := args[0].(type) {
	case int64:
		return float64(a), nil
	case bool:
		s := float64(0)
		if a {
			s = float64(1)
		}
		return s, nil
	case *objects.Null:
		return float64(0), nil
	case string:
		return strconv.ParseFloat(a, 64)
	case float64:
		return a, nil
	case objects.Glif:
		return float64(a), nil
	case byte:
		return float64(a), nil
	case objects.Bytes:
		v := 0
		// TODO: convert bytes to int, (?)
		return float64(v), nil
	}
	return 0, nil
}

func constr_bool(cx base.Context, args []any) (any, error) {
	if len(args) != 1 {
		return 0, errors.New("constr bool: incorrect count of args")
	}
	switch a := args[0].(type) {
	case int64:
		return a != 0, nil
	case bool:
		return a, nil
	case *objects.Null:
		return false, nil
	case string:
		return strconv.ParseBool(a)
	case float64:
		return a != 0, nil
	case objects.Glif:
		return a != 0, nil
	case byte:
		return a != 0, nil
	case objects.Bytes:
		v := len(a) != 0
		return v, nil
	}
	return 0, nil
}

func constr_glif(cx base.Context, args []any) (any, error) {
	if len(args) != 1 {
		return 0, errors.New("constr glif: incorrect count of args")
	}
	switch a := args[0].(type) {
	case int64:
		return objects.Glif(a), nil
	case string:
		rr := []objects.Glif(a)
		return rr[0], nil
	case byte:
		return objects.Glif(a), nil
	case objects.Bytes:
		r, _ := utf8.DecodeRune(a)
		return objects.Glif(r), nil
	}

	return 0, nil
}

func nums2bytes(vv []any) (any, error) {
	bb := make(objects.Bytes, len(vv))
	for i, n := range vv {
		switch n := n.(type) {
		case byte:
			bb[i] = n
		case int64:
			bb[i] = byte(n % 256)
		default:
			return nil, errors.New("bytes constr: incorrect val in list arg")
		}
	}
	return bb, nil
}

func constr_bytes(cx base.Context, args []any) (any, error) {
	if len(args) == 0 {
		return objects.Bytes{}, nil
	}
	switch a := args[0].(type) {
	case int64:
		return make(objects.Bytes, a), nil
	case byte:
		return make(objects.Bytes, a), nil
	case string:
		return objects.Bytes(a), nil
	case objects.Glif:
		// TODO: maybe should split rune (int32) byte-by-byte. thinking
		rr := []rune{a}
		return objects.Bytes(string(rr)), nil
	case objects.Bytes:
		bb := make(objects.Bytes, len(a))
		copy(bb, a)
		return bb, nil
	case *objects.ListVal:
		// mast contains bytes / ints
		return nums2bytes(a.Elems)
	case *objects.TupleVal:
		return nums2bytes(a.Elems)
	}
	return nil, errors.New("bytes constr: incorrect args")
}

func makeSequence(vals []any) ([]any, error) {
	switch len(vals) {
	case 0:
		return []any{}, nil
	case 1:
		switch a := vals[0].(type) {
		case int64:
			mx := int(a)
			vals := make([]any, a)
			for i := range mx {
				vals[i] = 0
			}
			return vals, nil
		case byte:
			mx := int(a)
			vals := make([]any, a)
			for i := 0; i < mx; i++ {
				vals[i] = 0
			}
			return vals, nil
		case string:
			gg := []objects.Glif(a)
			vals = make([]any, len(gg))
			for i, v := range gg {
				vals[i] = v
			}
			return vals, nil
		case *objects.ListVal:
			bb := make([]any, len(a.Elems))
			copy(bb, a.Elems)
			return bb, nil
		case *objects.TupleVal:
			bb := make([]any, len(a.Elems))
			copy(bb, a.Elems)
			return bb, nil
		case objects.Bytes:
			bb := make([]any, len(a))
			for i, n := range a {
				bb[i] = n
			}
			return bb, nil
			// Iterator
			// Generator
		}
	}
	return nil, errors.New("sequence constr: incorrect incoming args")
}

func constr_list(cx base.Context, args []any) (any, error) {
	vv, err := makeSequence(args)
	if err != nil {
		return nil, err
	}
	return objects.NewListVal(vv), nil
}

func constr_tuple(cx base.Context, args []any) (any, error) {
	vv, err := makeSequence(args)
	if err != nil {
		return nil, err
	}
	return objects.NewTupleVal(vv), nil
}

func constr_dict(cx base.Context, args []any) (any, error) {
	dd := make(map[any]any)
	switch len(args) {
	case 0:
		// empty dict
		return objects.NewDictVal(dd), nil
	case 1:
		switch a := args[0].(type) {
		case *objects.DictVal:
			maps.Copy(dd, a.Vmap)
			return objects.NewDictVal(dd), nil
		case *objects.ListVal:
			// list with pairs of key-val: [(k, v),..]
			for _, el := range a.Elems {
				// fmt.Printf(" Dcon#1 %T %v \n", el, el)
				n, ok := el.(*objects.TupleVal)
				if !ok {
					return nil, errors.New("dict constr: incorrect type of elem in list arg")
				}
				if len(n.Elems) != 2 {
					return nil, errors.New("dict constr: incorrect size of (k,v) tuple in list arg")
				}
				k, v := n.Elems[0], n.Elems[1]
				dd[k] = v
			}
		case *NamedArgs:
			for key, val := range a.Nvals {
				dd[key] = val
			}
		default:
			return nil, errors.New("dict constr: incorrect arg")
		}
	case 2:
		// two sets (list or tuple) with keys ans vals
		// dict(keys, vals): dict([], []); dict((,), (,))
		var ks []any
		var vs []any
		switch e0 := args[0].(type) {
		case *objects.ListVal:
			ks = e0.Elems
		case *objects.TupleVal:
			ks = e0.Elems
		}
		switch e1 := args[1].(type) {
		case *objects.ListVal:
			vs = e1.Elems
		case *objects.TupleVal:
			vs = e1.Elems
		}
		// if number of values is more just ignore others
		if len(ks) > len(vs) {
			return nil, errors.New("dict constr: too little vals")
		}
		for i, k := range ks {
			dd[k] = vs[i]
		}
	default:
		return nil, errors.New("dict constr: incorrect args")
	}

	return objects.NewDictVal(dd), nil
}
