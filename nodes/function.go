package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

// import (
// 	"github.com/lesnikyan/lisapet-go/base"
// 	"github.com/lesnikyan/lisapet-go/nodes"
// )

type Function struct {
	Name string

	// defArgs map[string]*base.Var
	defArgs map[string]*base.Var
	argVals []any
	Args    []base.Expression

	Block  *BlockExpr
	defCtx base.Context

	res    any
	resVal *base.Val
}

// 1. positional args, 2. named args,
// 3. default arg vals, 4. variative count
// 5. ovreload by arg count, 6. overload by arg types

func (fn *Function) SetArgVals(vals []any, mvals map[string]any) {
	fn.argVals = vals
}

// foo(<positional>, <named>)
func (fn *Function) PrepareArgs(fcx base.Context) error {
	if len(fn.Args) != len(fn.argVals) {
		return errors.New("Func.prep args: wrong count of arg vals")
	}
	// for i, argv := range fn.argVals {

	// }
	// args := make([]*base.Var, len(fn.Args))
	for i, ex := range fn.Args {
		// fmt.Printf(" Fu.PArg#1 %d) (%T, %v) = %v \n", i, ex, ex, fn.argVals[i])
		switch vex := ex.(type) {
		case *VarExpr:
			vex.NewVar(fcx)
			vr := vex.GetVar()
			vr.Val = fn.argVals[i]
		// args[i] = vr
		case *OperAssign:
			// named|default arg - arg = 5

			// typed arg - x : int

			// triple-dots arg - nn...
		}
	}
	return nil
}

func (fn *Function) Do(cx base.Context) error {
	fn.res = nil
	fn.resVal = nil

	fmt.Printf(" ---- Fu.Do#1 \n")

	// inner context
	inCx := fn.defCtx.SubContext()
	// set args
	err := fn.PrepareArgs(inCx)
	if err != nil {
		return err
	}

	// Do Block
	err = fn.Block.Do(inCx)
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
	// fmt.Printf(" - Fu.Do#3 br(%T : %v) fr(%T : %v) \n", r, r, fn.res, fn.res)

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

func NewFunction(name string, args []base.Expression, block *BlockExpr, ctx base.Context) *Function {
	// margs := make(map[string]*base.Var)
	// for _, n := range args {
	// 	margs[n.Name] = n
	// }
	return &Function{Name: name, Args: args, Block: block, defCtx: ctx}
}
