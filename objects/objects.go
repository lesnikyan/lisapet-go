package objects

import (
	"fmt"
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

type Function struct {
	Name  string
	Block *base.Block
}

type StructVal struct {
	Fields []*base.Var
}

type Regexp struct {
	Pattern *regexp.Regexp
}

/**
* unpack variables, etc
 */
func GetVal(v any) any {
	fmt.Printf("GetVal#0: %T, %v\n", v, v)
	switch vv := v.(type) {
	case *base.Var:
		return vv.Val
	case *base.Val:
		fmt.Printf("GetVal#Val: %T, %v\n", vv.V, vv.V)
		return vv.V
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
func (md *Module) Get() *base.Val {
	return nil
}

func (md *Module) Add(sub base.Expression) {
	md.block.Add(sub)
}
