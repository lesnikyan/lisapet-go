package nodes

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

// func definition node
type FuncDef struct {
	Name  string
	Args  []base.Expression
	Block *BlockExpr

	res *Function
}

func (fd *FuncDef) IsParent() bool {
	return true
}

func (fd *FuncDef) Get() *base.Val {
	if fd.res == nil {
		return nil
	}
	return base.NewVal(fd.res)
}

func (fd *FuncDef) Do(cx base.Context) error {
	fd.res = nil
	defCx := objects.NewContext(cx)
	// args := make([]*base.Var, len(fd.Args))
	// for i, vex := range fd.Args {
	// 	vex.NewVar(defCx)
	// 	vr := vex.GetVar()
	// 	args[i] = vr
	// }
	fn := NewFunction(fd.Name, fd.Args, fd.Block, defCx)
	fd.res = fn
	err := fn.Init(cx)
	if err != nil {
		return errors.Join(err)
	}
	cx.AddFunc(fn)
	return nil
}

func (fd *FuncDef) Add(sub base.Expression) {
	fd.Block.Add(sub)
}

// for methods only
func (fd *FuncDef) SetObject(obj *VarExpr) {

}

func NewFuncDef(name string, args []base.Expression) *FuncDef {
	return &FuncDef{Name: name, Args: args, Block: NewBlock()}
}

//==

var NoResult = base.NewVal(&objects.Null{})

//==

// func call: smth([args])

type FuncCall struct {
	Src  base.Expression // should return function object
	args []base.Expression

	fun    base.FuncVal
	res    any
	resVal *base.Val
}

func (fc *FuncCall) Get() *base.Val {
	// if fc.res == nil {
	// 	return nil
	// }
	// return base.NewVal(fc.res)
	return fc.resVal
}

func (fc *FuncCall) getFunc(cx base.Context) error {
	fc.fun = nil // TODO: optimize for defined by name (not for obj in var)
	err := fc.Src.Do(cx)
	if err != nil {
		return err
	}
	fv := GetExprVal(fc.Src, cx)
	if fv == nil {
		return errors.New("trying to call nil elem")
	}
	switch fn := fv.(type) {
	case *Function:
		fc.fun = fn
	case *NFunc:
		fc.fun = fn
	default:
		// fmt.Printf("Err FunCall: non func: (%T, %v) \n", fn, fn)
		return errors.New("trying to call non-function ")
	}

	return nil
}

func (fc *FuncCall) DoArgs(cx base.Context) error {
	// do arg expr
	// namedN := 0
	mvals := map[string]any{}
	vals := make([]any, len(fc.args))

	i := 0
	// foo(ord1, ordN, variadic..., named1=val1, named2=val2)
	for _, vex := range fc.args {
		nmExp, ok := vex.(*OperAssign)
		// fmt.Printf("FunCall (Args1): exp:(%T, %v) isAssign: %v \n", vex, vex, ok)
		if !ok {
			// ordered arg
			err := vex.Do(cx)
			if err != nil {
				return err
			}
			val := GetExprVal(vex, cx)
			if val == nil {
				return errors.New("func call: no result of argument expression")
			}
			vals[i] = val
			i += 1
			continue
		}
		// variadic args

		// named arg
		lvar, ok := nmExp.Left.(*VarExpr)
		if !ok {
			return errors.New("func call (Args2): Named arg in func call without left part")
		}
		err := nmExp.Right.Do(cx)
		if err != nil {
			return err
		}
		lval := nmExp.Right.Get()
		// fmt.Printf("FunCall (Args3): r-exp:(%T, %v) lval: %v ?nil: %v \n", nmExp.right, nmExp.right, lval, lval == nil)
		if lval == nil {
			return errors.New("func call (Args): Named arg in func call without value")
		}
		argName := lvar.name
		mvals[argName] = lval.V
		// namedN += 1
	}
	// put args to func
	fc.fun.SetArgVals(vals[:i], mvals)
	return nil
}

func (fc *FuncCall) Do(cx base.Context) error {
	fc.res = nil
	err := fc.getFunc(cx)
	if err != nil {
		return err
	}
	// do args
	err = fc.DoArgs(cx)
	if err != nil {
		return err
	}
	// do func
	err = fc.fun.Do(cx)
	if err != nil {
		return err
	}
	// fmt.Printf(" - FCall.Do#3 br(%T : %v) fr(%T : %v) \n", r, r, fc.res, fc.res)
	r := fc.fun.Get()
	fc.resVal = r
	if r == nil {
		r = NoResult
	}
	fc.resVal = r
	fc.res = r.V
	// fmt.Printf(" - FCall.Do#3 br(%T : %v) fr(%T : %v) \n", r, r, fc.res, fc.res)
	return nil
}

func NewFuncCall(fExp base.Expression, args []base.Expression) *FuncCall {
	return &FuncCall{Src: fExp, args: args}
}
