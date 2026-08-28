package nodes

import (
	"errors"
	"fmt"

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

		// case *objects.TupleVal:
		// 	return int64(len(a.Elems)), nil

	}
	return false, errors.New("constr string: incorrect arg type")
}

func constr_byte(cx base.Context, args []any) (any, error) {
	if len(args) != 1 {
		return 0, errors.New("constr string: incorrect count of args")
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
		return 0, errors.New("constr string: incorrect count of args")
	}
	if args[0] != nil {
		return objects.Some(args[0]), nil
	}
	// switch a := args[0].(type) {
	// }
	return false, errors.New("constr byte: incorrect arg type")
}
