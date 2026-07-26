package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

// import (
// 	"github.com/lesnikyan/lisapet-go/base"
// 	"github.com/lesnikyan/lisapet-go/nodes"
// )

type Function struct {
	Name string

	defArgs map[string]*base.Var
	argVals []any

	Block  *BlockExpr
	defCtx base.Context

	res    any
	resVal *base.Val
}

func (fn *Function) SetArgVals(vals []any) {
	fn.argVals = vals
}

func (fn *Function) Do(cx base.Context) error {
	fn.res = nil
	fn.resVal = nil

	fmt.Printf(" ---- Fu.Do#1 \n")

	// inner context
	inCx := fn.defCtx.SubContext()
	// set args

	// Do Block
	err := fn.Block.Do(inCx)
	if err != nil {
		return err
	}
	// take result
	// 1. simple case: result of last expression
	r := fn.Block.Get()
	fn.resVal = r
	if r != nil {
		fn.res = r.V
	}
	fmt.Printf(" - Fu.Do#3 br(%T : %v) fr(%T : %v) \n", r, r, fn.res, fn.res)

	// 2. result if return

	return nil
}

func (fn *Function) Get() *base.Val {
	// return base.NewVal(fn.res)
	return fn.resVal
}

func (fn *Function) GetName() string {
	return fn.Name
}

func NewFunction(name string, args []*base.Var, block *BlockExpr, ctx base.Context) *Function {
	margs := make(map[string]*base.Var)
	for _, n := range args {
		margs[n.Name] = n
	}
	return &Function{Name: name, defArgs: margs, Block: block, defCtx: ctx}
}
