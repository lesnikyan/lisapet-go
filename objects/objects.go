package objects

import "reflect"

func f1() {
	n := 1
	reflect.TypeOf(n)
}

type Block interface {
	Do(cx *Context)
	Get() any
}

type Val interface {
	GetVal() any
}

type Null struct {
}

type Var struct {
	Val  any
	Name string
	Type int
}

type ListVal struct {
	elems []any
}

type DictVal struct {
	vmap map[any]any
}

type Function struct {
	name  string
	block *Block
}

/**
* unpack variables, etc
 */
func GetVal(v any) any {
	switch vv := v.(type) {
	case Var:
		return vv.Val
	default:
		return vv
	}
}

// func Any2Val(obj any) any {
// 	switch src := obj.(type) {
// 	case *Var:
// 		return src.Val
// 	}
// 	return nil
// }
