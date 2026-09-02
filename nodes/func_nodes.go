package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

// func definition node
type FuncDef struct {
	Name  string
	Args  []base.Expression
	Block *BlockExpr

	res *objects.Function
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

func (fd *FuncDef) MakeFunc(cx base.Context) (*objects.Function, error) {
	defCx := objects.NewContext(cx)

	fbk := &FuncBlock{Block: fd.Block, dfnArgs: fd.Args}
	err := fbk.Init(cx)
	// err := fn.Init(cx)
	if err != nil {
		return nil, errors.Join(errors.New("FuncDef.MakeFunc: init error"), err)
	}
	fn := objects.NewFunction(fd.Name, fbk, defCx)
	fd.res = fn
	return fn, nil
}

func (fd *FuncDef) Do(cx base.Context) error {
	fd.res = nil
	fn, err := fd.MakeFunc(cx)
	if err != nil {
		return err
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

// === MethodDef

type MethodDef struct {
	Func *FuncDef
	Inst *OperColon

	res *objects.Method
}

func (md *MethodDef) IsParent() bool {
	return true
}

func (md *MethodDef) Get() *base.Val {
	return base.NewVal(md.res)
}

func (md *MethodDef) Do(cx base.Context) error {
	texp, ok := md.Inst.Right.(*VarExpr)
	if !ok {
		return errors.New("MethodDef: instance expr has incorrect syntax")
	}
	tname := texp.GetName()
	tt := cx.GetType(tname)
	if tt == nil {
		return errors.New("MethodDef: type of instance not found")
	}
	iname := ""
	if nexp, ok := md.Inst.Left.(*VarExpr); ok {
		iname = nexp.GetName()
	}
	fd := md.Func
	fn, err := fd.MakeFunc(cx)
	if err != nil {
		return errors.Join(errors.New("FuncDef.MakeFunc: init error"), err)
	}
	met := objects.NewMethod(fn, tt, iname)
	// fmt.Printf(" MetodDef.Do: t(%T, %v) meth(%T, %v) \n", tt, tt, met, met)
	if tt.Def != nil {
		if sdef, ok := tt.Def.(*objects.StructDef); ok {
			sdef.AddMethod(met)
		}
	}
	md.res = met
	return nil
}

func (fd *MethodDef) Add(sub base.Expression) {
	fd.Func.Add(sub)
}

// for methods only
func (fd *MethodDef) SetObject(obj *VarExpr) {

}

func NewMethodDef(fun *FuncDef, inst *OperColon) *MethodDef {
	return &MethodDef{Func: fun, Inst: inst}
}

//==

var NoResult = base.NewVal(&objects.Null{})

// ==

// func call: smth([args])
type FuncCall struct {
	Src  base.Expression // should return function object
	args []base.Expression

	fun    base.FuncVal
	res    any
	resVal *base.Val
}

func (fc *FuncCall) Get() *base.Val {
	return fc.resVal
}

func (fc *FuncCall) getFunc(cx base.Context) error {
	fc.fun = nil // TODO: optimize for defined by name (not for obj in var)
	err := fc.Src.Do(cx)
	if err != nil {
		return err
	}
	fv := GetExprVal(fc.Src, cx)
	// fmt.Printf("FunCall#1: expr: %T elem: (%T, %v) \n", fc.Src, fv, fv)
	// fmt.Printf("FunCall#2: non func: (%T, %v) \n", fv, fv)
	if fv == nil {
		return errors.New("trying to call nil elem")
	}
	switch fn := fv.(type) {
	case *objects.Function:
		fc.fun = fn
	case *NFunc:
		fc.fun = fn
	case *MFunc:
		fc.fun = fn
	case *objects.Method:
		// fmt.Printf("getFunc, Method: %T, %v inst: %T, %v \n", fn, fn, fn.Inst, fn.Inst)
		fc.fun = fn
	case *base.Type:
		if fn.Construct == nil {
			return errors.New("trying to call undefined constructor of type " + fn.Name)
		}
		fc.fun = fn.Construct
	default:
		// fmt.Printf("Err FunCall: non func: (%T, %v) \n", fn, fn)
		return fmt.Errorf("FunCall: non func: (%T, %v) \n", fn, fn)
	}

	return nil
}

func (fc *FuncCall) DoArgs(cx base.Context) error {
	// do arg expr
	mvals := map[string]any{}
	vals := make([]any, len(fc.args))
	i := 0
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

	// fmt.Printf(" - FCall(%s).Do#4 br(%T : %v) \n", fc.fun.GetName(), r, r)
	fc.resVal = r
	if r == nil {
		r = NoResult
	}
	fc.resVal = r
	fc.res = r.V
	// fmt.Printf(" - FCall(%s).Do#5 br(%T : %v) fr(%T : %v) \n", fc.fun.GetName(), r, r, fc.res, fc.res)
	return nil
}

func NewFuncCall(fExp base.Expression, args []base.Expression) *FuncCall {
	return &FuncCall{Src: fExp, args: args}
}
