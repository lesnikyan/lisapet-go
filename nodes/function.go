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
	defVals map[string]*base.Val
	argVals []any
	nmVals  map[string]any
	dfnArgs []base.Expression
	Args    []*VarExpr

	Block  *BlockExpr
	defCtx base.Context

	res    any
	resVal *base.Val
}

// 1. positional args, 2. named args,
// 3. default arg vals, 4. variative count
// 5. ovreload by arg count, 6. overload by arg types

func (fn *Function) Init(cx base.Context) error {
	fn.defVals = make(map[string]*base.Val)
	nArgs := make([]*VarExpr, len(fn.dfnArgs))
	for i, ex := range fn.dfnArgs {
		switch arx := ex.(type) {
		case *VarExpr:
			nArgs[i] = arx

		case *OperAssign:
			// get var name
			vx, ok := arx.left.(*VarExpr)
			if !ok {
				return errors.New("func init: incorrect left-operand od default arg val")
			}
			name := vx.name
			// get default val
			err := arx.right.Do(cx)
			if err != nil {
				return err
			}
			v := arx.right.Get()
			if v == nil {
				return errors.New("func init: no value from default-val expression")
			}
			fn.defVals[name] = v
			nArgs[i] = vx

		case *OperColon:
			// typed arg - x : int

		case *TripleDots:
			// triple-dots arg - nn...
		}
	}
	fn.Args = nArgs
	return nil
}

func (fn *Function) SetArgVals(vals []any, nvals map[string]any) {
	// fmt.Printf(" Fu.SetArgVals  ovs=%d, nvs:%d  \n", len(vals), len(nvals))
	fn.argVals = vals
	fn.nmVals = nvals
}

// func (fn *Function) argVar(ex base.Expression, cx base.Context) *base.Var {
// 	fmt.Printf(" argVar#1 (%T, %v) \n", ex, ex)
// 	var vr *base.Var
// 	// var defval any
// 	switch vex := ex.(type) {
// 	case *VarExpr:
// 		// arg defined as ordered: foo(a, b, c)
// 		vex.NewVar(cx)
// 		vr = vex.GetVar()
// 		// case *OperAssign:
// 		// 	// arg defined with default value foo(named=123)
// 		// 	vr = fn.argVar(vex.left, cx)
// 		// case *OperColon:
// 		// 	// typed arg - x : int

// 		// case *TripleDots:
// 		// 	// triple-dots arg - nn...
// 	}
// 	return vr
// }

func (fn *Function) getVal(i int, name string) (any, error) {

	// var val any
	if i < len(fn.argVals) {
		// passed ordered arg
		// fmt.Printf(" Fu.try arg1 %d, %v  \n", i, name)
		return fn.argVals[i], nil
	}
	if fn.nmVals != nil {
		// passed named arg
		// fmt.Printf(" Fu.try arg2 %d, %v  \n", i, name)
		if nmval, ok := fn.nmVals[name]; ok {
			return nmval, nil
		}
	}
	// use default value
	defv, ok := fn.defVals[name]
	// fmt.Printf(" Fu.try arg3 %d, %v || %v %v \n", i, name, defv, ok)
	if ok {
		return defv.V, nil
	}

	return nil, errors.New("func prepare: no value for argument: " + name)
	// return nil, errors.New("fun: val not found")
}

// foo(<positional>, <variadic...>, <named=val>)
func (fn *Function) PrepareArgs(cx base.Context) error {
	// minArgCount := len(fn.Args) // actual for odered args
	// if len(fn.argVals) <= minArgCount {
	// 	return errors.New("Func.prep args: wrong count of arg vals")
	// }
	// for i, argv := range fn.argVals {

	// }
	// args := make([]*base.Var, len(fn.Args))
	for i, ex := range fn.Args {
		// fmt.Printf(" Fu.PArg#1 %d) (%T, %v)  \n", i, ex, ex)
		// var vr *base.Var

		// take var
		ex.NewVar(cx)
		vr := ex.GetVar()
		if vr == nil {
			return errors.New("func prepare: no arg var")
		}

		// get value
		val, err := fn.getVal(i, vr.Name)
		if err != nil {
			return err
		}

		// fmt.Printf(" Fu.PArg#N %d) (%T, %v) = (%T, %v)  \n", i, vr, vr, val, val)
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
			// fmt.Printf(" - Fu.Do#4 br(%T : %v) fr(%T : %v) \n", fn.resVal, fn.resVal, fn.res, fn.res)
		}
		return nil
	}

	// 1. simple case: result of last expression
	r := fn.Block.Get()
	fn.resVal = r
	if r != nil {
		fn.res = r.V
		// fmt.Printf(" - Fu.Do#3 br(%T : %v) fr(%T : %v) \n", r, r, fn.res, fn.res)
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
	return &Function{Name: name, dfnArgs: args, Block: block, defCtx: ctx}
}
