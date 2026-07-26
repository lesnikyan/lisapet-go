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
	Args  []*VarExpr
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
	args := make([]*base.Var, len(fd.Args))
	for i, vex := range fd.Args {
		vex.NewVar(defCx)
		vr := vex.GetVar()
		args[i] = vr
	}
	fn := NewFunction(fd.Name, args, fd.Block, defCx)
	fd.res = fn
	cx.AddFunc(fn)
	return nil
}

func (fd *FuncDef) Add(sub base.Expression) {
	fd.Block.Add(sub)
}

// for methods only
func (fd *FuncDef) SetObject(obj *VarExpr) {

}

func NewFuncDef(name string, args []*VarExpr) *FuncDef {
	return &FuncDef{Name: name, Args: args, Block: NewBlock()}
}

//==

var NoResult = base.NewVal(objects.Null{})

//==

// func call: smth([args])

type FuncCall struct {
	Src  base.Expression // should return function object
	args []base.Expression

	fun    *Function
	res    any
	resVal *base.Val
}

func (fc *FuncCall) Get() *base.Val {
	if fc.res == nil {
		return nil
	}
	return base.NewVal(fc.res)
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
	fn, ok := fv.(*Function)
	if !ok {
		return errors.New("trying to call non-function")
	}
	fc.fun = fn
	return nil
}

func (fc *FuncCall) Do(cx base.Context) error {
	err := fc.getFunc(cx)
	if err != nil {
		return err
	}
	fc.fun.Do(cx)
	r := fc.fun.Get()
	fc.resVal = r
	if r == nil {
		r = NoResult
	}
	fc.res = r.V
	fmt.Printf(" - FCall.Do#3 br(%T : %v) fr(%T : %v) \n", r, r, fc.res, fc.res)
	return nil
}

func NewFuncCall(fExp base.Expression, args []base.Expression) *FuncCall {
	return &FuncCall{Src: fExp, args: args}
}
