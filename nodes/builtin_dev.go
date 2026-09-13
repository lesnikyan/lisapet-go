package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	obb "github.com/lesnikyan/lisapet-go/objects"
)

// Functions used only for dev needs

// builtin func with named args
func devNamed(cx base.Context, args []any) (any, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("builtin named error:: No args")
	}
	// named, ok := args[0].(*NamedArgs)

	_, named := SplitNamed(args)
	if named == nil {
		return nil, fmt.Errorf("builtin named error: No named args args")
	}
	res := make(map[any]any, len(named))
	for key, val := range named {
		res[key] = fmt.Sprintf("%v", val)
	}
	return obb.NewDictVal(res), nil
}

// builtin func with ordered and named args
func devOrdNamed(cx base.Context, args []any) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("builtin ordered-named error: No enough args")
	}
	// var named *NamedArgs
	// ordd := make([]any, len(args))
	// k := 0
	// for _, arg := range args {
	// 	switch an := arg.(type) {
	// 	case *NamedArgs:
	// 		// only 1 expected
	// 		named = an
	// 	default:
	// 		ordd[k] = an
	// 		k++
	// 	}
	// }
	ordd, named := SplitNamed(args)
	if named == nil {
		return nil, fmt.Errorf("builtin ordered-named error: No named args")
	}
	res := make(map[any]any, len(named))
	for key, val := range named {
		res[key] = fmt.Sprintf("%v", val)
	}
	res["#ordered"] = obb.NewListVal(ordd)
	return obb.NewDictVal(res), nil
}

// builtin func with default args 3 ord, 3 default
func devDefArgs33(cx base.Context, args []any) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("builtin ordered-named error: No enough args")
	}
	ordd, named := SplitNamed(args)
	// named is optional here
	// if named == nil {
	// 	return nil, fmt.Errorf("builtin ordered-named error: No named args")
	// }
	if len(ordd) < 3 {
		return nil, fmt.Errorf("builtin ordered-named error: Not enough necessary args")
	}

	var a int64
	var b float64
	var c string
	var ok bool
	// switch len(ordd) {
	// case 1:
	// 	a = ordd[0].(int64)
	// case 2:
	// 	a = ordd[0].(int64)
	// 	b = ordd[1].(float64)
	// case 3:
	// }
	a, ok = ordd[0].(int64)
	if !ok {
		return nil, fmt.Errorf("builtin ordered-named error: arg `a` must be  int")
	}
	b, ok = ordd[1].(float64)
	if !ok {
		return nil, fmt.Errorf("builtin ordered-named error: arg `b` must be  float")
	}
	c, ok = ordd[2].(string)
	if !ok {
		return nil, fmt.Errorf("builtin ordered-named error: arg `c` must be  string")
	}

	defdd := int64(111)
	defee := float64(2.22)
	defff := "FooBar"
	// res := make(map[any]any, len(named))
	for key, val := range named {
		// res[key] = fmt.Sprintf("%v", val)
		switch key {
		case "dd":
			defdd, ok = val.(int64)
			if !ok {
				return nil, fmt.Errorf("builtin ordered-named error: arg `dd` must be  string")
			}
		case "ee":
			defee, ok = val.(float64)
			if !ok {
				return nil, fmt.Errorf("builtin ordered-named error: arg `ee` must be  string")
			}
		case "ff":
			defff, ok = val.(string)
			if !ok {
				return nil, fmt.Errorf("builtin ordered-named error: arg `ff` must be  string")
			}
		}
	}
	// res["#ordered"] = obb.NewListVal(ordd)
	res := fmt.Sprintf("a=%d, b=%0.2f, c=`%s`; dd=%x ee=%0.2f, ff=%s", a, b, c, defdd, defee, defff)
	return res, nil
}

// builtin func with default args 3 ord, 3 default
func devDefOrds(cx base.Context, args []any) (any, error) {
	// if len(args) < 1 {
	// 	return nil, fmt.Errorf("builtin ordered-named error: No enough args")
	// }
	ordd, named := SplitNamed(args)
	// named is optional here
	// if len(ordd) < 1 {
	// 	return nil, fmt.Errorf("builtin ordered-named error: Not enough necessary args")
	// }

	var a int64 = 10
	var b float64 = 1.0
	var c string = "---"
	var ok bool

	alen := len(ordd)
	if alen > 0 {
		a, ok = ordd[0].(int64)
		if !ok {
			return nil, fmt.Errorf("builtin ordered-named error: arg `a` must be  int")
		}
		if alen > 1 {

		}
		b, ok = ordd[1].(float64)
		if !ok {
			return nil, fmt.Errorf("builtin ordered-named error: arg `b` must be  float")
		}
		if alen > 2 {
			c, ok = ordd[2].(string)
			if !ok {
				return nil, fmt.Errorf("builtin ordered-named error: arg `c` must be  string")
			}
		}
	}
	sepv := ","
	if sep, ok := named["sep"]; ok {
		sepv, ok = sep.(string)
		if !ok {
			return nil, fmt.Errorf("builtin ordered-named error: optional arg `sep` must be a string")
		}
	}
	if alen < 3 {
		if sep, ok := named["c"]; ok {
			c, ok = sep.(string)
			if !ok {
				return nil, fmt.Errorf("builtin ordered-named error: optional arg `c` must be a string")
			}
		}
		if alen < 2 {
			if sep, ok := named["b"]; ok {
				b, ok = sep.(float64)
				if !ok {
					return nil, fmt.Errorf("builtin ordered-named error: optional arg `b` must be a float")
				}
			}
			if alen == 0 {
				if sep, ok := named["a"]; ok {
					a, ok = sep.(int64)
					if !ok {
						return nil, fmt.Errorf("builtin ordered-named error: optional arg `a` must be a int")
					}
				}
			}
		}
	}

	res := fmt.Sprintf("a=%d %[2]s b=%0.2f %[2]s c=`%[4]s`", a, sepv, b, c)
	return res, nil
}
