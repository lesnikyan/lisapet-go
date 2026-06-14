package objects

import (
	"reflect"
	"regexp"

	"github.com/lesnikyan/lisapet-go/base"
)

func f1() {
	n := 1
	reflect.TypeOf(n)
}

// type Expression interface {
// 	Do(*Context) error
// 	Get() any
// }
// type Block interface {
// 	Add(sub Expression)
// 	Do(cx *Context) error
// 	Get() any
// }

type Val interface {
	GetVal() any
}

type Null struct {
}

// type Var struct {
// 	Val  any
// 	Name string
// 	Type int
// }

type ListVal struct {
	elems []any
}

type TupleVal struct {
	elems []any
}

type DictVal struct {
	vmap map[any]any
}

type Function struct {
	name  string
	block *base.Block
}

type StructVal struct {
	fields []*base.Var
}

type Regexp struct {
	pattern *regexp.Regexp
}

/**
* unpack variables, etc
 */
func GetVal(v any) any {
	switch vv := v.(type) {
	case base.Var:
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

type Module struct {
	block base.Block
	ctx   *Context
}

func NewModule(ctx *Context) *Module {
	return &Module{ctx: ctx}
}

func (md *Module) Do(cx *Context) error {
	return nil
}
func (md *Module) Get() any {
	return nil
}

func (md *Module) Add(sub base.Expression) {
	md.block.Add(sub)
}
