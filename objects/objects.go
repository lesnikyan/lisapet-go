package objects

import (
	"fmt"
	"regexp"

	"github.com/lesnikyan/lisapet-go/base"
)

// func f1() {
// 	n := 1
// 	reflect.TypeOf(n)
// }

type Val interface {
	GetVal() any
}

// null value
type Null struct {
}

// internal value, can be produced from EmptyExpr
type EmptyVal struct {
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
	// fmt.Printf("GetVal#0: %T, %v\n", v, v)
	switch vv := v.(type) {
	case *base.Var:
		return vv.Val
	case *base.Val:
		// fmt.Printf("GetVal#Val: %T, %v\n", vv.V, vv.V)
		return vv.V
	default:
		return vv
	}
}

func SetVarVal(target *base.Var, rval any) error {
	if target.StrictType {
		typeOk, cval := PrepareVal(target.Type, rval)
		if !typeOk {
			// bad val
			// fmt.Printf("SetValTo#4 val conv error: var %v, val: %v, vtype: %s \n", target, rval, target.Type.Name)
			return fmt.Errorf("oper assign: incorrecttype of value in right operand ")
		}
		rval = cval
	}
	target.Val = rval
	return nil
}

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
