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

type ArgExp struct {
	Name       string
	VExp       *VarExpr
	Type       *base.Type
	StrictType bool

	res any
}

func (ax *ArgExp) Do(cx base.Context) error {
	err := ax.VExp.Do(cx)
	if err != nil {
		return errors.Join(err, errors.New("ArgExp Do err1"))
	}
	return nil
}

func (ax *ArgExp) Get(cx base.Context) *base.Val {
	return ax.VExp.Get()
}

type Function struct {
	Name string

	// defArgs map[string]*base.Var
	defVals map[string]*base.Val
	argVals []any
	nmVals  map[string]any
	dfnArgs []base.Expression
	Args    []*ArgExp

	Block  *BlockExpr
	defCtx base.Context

	res    any
	resVal *base.Val
}

// 1. positional args, 2. named args,
// 3. default arg vals, 4. variative count
// 5. ovreload by arg count, 6. overload by arg types

func (fn *Function) InitArg(cx base.Context, ex base.Expression) (*ArgExp, error) {
	// fmt.Printf(" Fu.InitArg#0  argExp=(%T, %v) \n", ex, ex)
	switch arx := ex.(type) {
	case *VarExpr:
		rex := &ArgExp{Name: arx.name, VExp: arx}
		return rex, nil

	case *OperAssign:
		// var part
		var lvar *VarExpr
		var vtype *base.Type
		strict := false
		switch cvar := arx.Left.(type) {
		case *VarExpr:
			vtype = cx.GetType("any")
			lvar = cvar

		case *OperColon:
			// Typed Var with default val

			lExp, ok := cvar.Left.(*VarExpr)
			if !ok {
				return nil, errors.New("arg init: err2")
			}
			lvar = lExp

			err := cvar.Right.Do(cx)
			if err != nil {
				return nil, err
			}
			rvar, ok := cvar.Right.(*VarExpr)
			if !ok {
				// no type
				fmt.Printf(" Fu.InitArg.err111  n=%s tp=(%T, %v) \n", lvar.name, cvar.Right, cvar.Right)
				return nil, errors.New("arg init: err111, not a word in right of types arg ")
			}
			tt := cx.GetType(rvar.name)
			if !ok {
				return nil, errors.New("assign-colon: right part is not type")
			}
			vtype = tt
			strict = true
		}

		name := lvar.name
		rex := &ArgExp{Name: name, VExp: lvar, Type: vtype, StrictType: strict}

		// get default val
		err := arx.Right.Do(cx)
		if err != nil {
			return nil, err
		}
		v := arx.Right.Get()
		if v == nil {
			return nil, errors.New("arg init: no value from default-val expression")
		}
		fn.defVals[name] = v
		// return vx, nil
		return rex, nil

	case *OperColon:
		cvar := arx
		lvar, ok := cvar.Left.(*VarExpr)
		if !ok {
			return nil, errors.New("arg init: err22")
		}
		err := cvar.Right.Do(cx)
		if err != nil {
			return nil, err
		}
		rvar, ok := cvar.Right.(*VarExpr)
		if !ok {
			// no type
			// fmt.Printf(" Fu.InitArg.err221  n=%s tp=(%T, %v) \n", lvar.name, cvar.right, cvar.right)
			return nil, errors.New("arg init: err221, not a word in right of types arg ")
		}
		tt := cx.GetType(rvar.name)
		if tt == nil {
			// no type
			// fmt.Printf(" Fu.InitArg.err23  n=%s tp=(%T, %v) \n", lvar.name, cvar.right, cvar.right)
			return nil, errors.New("arg init: err23 ")
		}
		// vtype, ok := rval.V.(*base.Type)
		// tt, ok := rres.V.(*base.Type)

		// if !ok {
		// 	return nil, errors.New("arg init colon: right part is not type 2")
		// }
		vtype := tt
		name := lvar.name
		rex := &ArgExp{Name: name, VExp: lvar, Type: vtype, StrictType: true}
		// fmt.Printf(" Fu.InitArg#5  n=%s tp=(%T, %v) \n", name, tt, tt)
		return rex, nil

	case *TripleDots:
		// triple-dots arg - nn...
	}
	return nil, errors.New("arg init: incorrect end of init")
}

func (fn *Function) Init(cx base.Context) error {
	// fmt.Printf(" Fu.Init  deArgLen==%d  \n", len(fn.dfnArgs))
	fn.defVals = make(map[string]*base.Val)
	nArgs := make([]*ArgExp, len(fn.dfnArgs))
	for i, ex := range fn.dfnArgs {
		argx, err := fn.InitArg(cx, ex)
		if err != nil {
			return errors.Join(err, errors.New("func init"))
		}
		nArgs[i] = argx
	}
	fn.Args = nArgs
	return nil
}

func (fn *Function) SetArgVals(vals []any, nvals map[string]any) {
	// fmt.Printf(" Fu.SetArgVals  ovs=%d, nvs:%d  \n", len(vals), len(nvals))
	fn.argVals = vals
	fn.nmVals = nvals
}

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
	// fmt.Printf(" Fu.PArg#0  %d  \n", len(fn.Args))
	for i, arg := range fn.Args {
		// fmt.Printf(" Fu.PArg#1 %d) (%T, %v)  \n", i, arg, arg)
		// take var
		if arg.StrictType {
			arg.VExp.NewVarTyped(cx, arg.Type)

		} else {
			// cx.GetType("any")
			arg.VExp.NewVar(cx)
		}
		vr := arg.VExp.GetVar()
		// fmt.Printf(" Fu.PArg#2 %d) (%T, %v)  \n", i, vr, vr)
		if vr == nil {
			return errors.New("func prepare: no arg var")
		}
		// if arg.StrictType {
		// 	vr.StrictType = true
		// 	vr.Type = arg.Type
		// }

		// get value
		val, err := fn.getVal(i, arg.Name)
		if err != nil {
			return err
		}
		// fmt.Printf(" Fu.PArg#3 %d) (%T, %v)  \n", i, vr, vr)
		aerr := SetValTo(vr, val)
		if aerr != nil {
			return errors.Join(errors.New("func prep arg: assign arg error"), aerr)
		}

		// fmt.Printf(" Fu.PArg#N %d) (%T, %v) = (%T, %v)  \n", i, vr, vr, val, val)
		// vr.Val = val
	}
	return nil
}

func (fn *Function) Do(cx base.Context) error {
	fn.res = nil
	fn.resVal = nil

	// fmt.Printf(" ---- Fu.Do#1 \n")

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
		// fmt.Printf(" - Fu.Do#5 br(%T : %v) fr(%T : %v) \n", r, r, fn.res, fn.res)
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
	return &Function{Name: name, dfnArgs: args, Block: block, defCtx: ctx}
}
