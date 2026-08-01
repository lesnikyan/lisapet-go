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
	nmVals  map[string]any
	Args    []base.Expression

	Block  *BlockExpr
	defCtx base.Context

	res    any
	resVal *base.Val
}

// 1. positional args, 2. named args,
// 3. default arg vals, 4. variative count
// 5. ovreload by arg count, 6. overload by arg types

func (fn *Function) SetArgVals(vals []any, nvals map[string]any) {
	fmt.Printf(" Fu.SetArgVals  ovs=%d, nvs:%d  \n", len(vals), len(nvals))
	fn.argVals = vals
	fn.nmVals = nvals
}

// foo(<positional>, <variadic...>, <named=val>)
func (fn *Function) PrepareArgs(fcx base.Context) error {
	// minArgCount := len(fn.Args) // actual for odered args
	// if len(fn.argVals) <= minArgCount {
	// 	return errors.New("Func.prep args: wrong count of arg vals")
	// }
	// for i, argv := range fn.argVals {

	// }
	// args := make([]*base.Var, len(fn.Args))
	for i, ex := range fn.Args {
		fmt.Printf(" Fu.PArg#1 %d) (%T, %v)  \n", i, ex, ex)
		var vr *base.Var
		var val any
		switch vex := ex.(type) {
		case *VarExpr:
			// arg defined as ordered: foo(a, b, c)
			vex.NewVar(fcx)
			vr = vex.GetVar()
			if i < len(fn.argVals) {
				val = fn.argVals[i]
			} else {
				// val passed by name, after ordered ones: foo(ordered, named=val)
				nmval, ok := fn.nmVals[vr.Name]
				if !ok {
					return errors.New("Not enough arguments")
				}
				// named arg has been passed
				val = nmval
			}
		case *OperAssign:
			// arg defined with default value foo(named=123)

		case *OperColon:
			// typed arg - x : int

		case *TripleDots:
			// triple-dots arg - nn...
		}

		fmt.Printf(" Fu.PArg#N %d) (%T, %v) = (%T, %v)  \n", i, vr, vr, val, val)
		vr.Val = val
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

	// 2. result if return
	pup := fn.Block.GetPopUp()
	if pup != nil {
		switch pup.Parent {
		case NodeReturn:
			// return val
			fn.resVal = pup.Res
			if pup.Res != nil {
				fn.res = pup.Res.V
			}
			fmt.Printf(" - Fu.Do#4 br(%T : %v) fr(%T : %v) \n", fn.resVal, fn.resVal, fn.res, fn.res)
		}
		return nil
	}

	// 1. simple case: result of last expression
	r := fn.Block.Get()
	fn.resVal = r
	if r != nil {
		fn.res = r.V
		fmt.Printf(" - Fu.Do#3 br(%T : %v) fr(%T : %v) \n", r, r, fn.res, fn.res)
		return nil
	}

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
