package nodes

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
)

type NFunc struct {
	Name  string
	args  []any                                  // passed args
	mvals map[string]any                         // passed named args
	fun   func(base.Context, []any) (any, error) // func-apapter called in Do()

	resV *base.Val
}

// Builtin function object
func (fn *NFunc) Do(cx base.Context) error {
	fn.resV = nil
	// TODO: named args, vary arg list, multi-result
	res, err := fn.fun(cx, fn.args)
	if err != nil {
		return err
	}
	fn.resV = base.NewVal(res)
	return nil
}

func (fn *NFunc) Get() *base.Val {
	if fn.resV == nil {
		return nil
	}
	return fn.resV
}

func (fn *NFunc) GetName() string {
	return fn.Name
}

func (fn *NFunc) SetArgVals(vals []any, mvals map[string]any) {
	fn.args = vals
	fn.mvals = mvals
}

/// ================= Usage Example ====================

// target func example
func mock_target(arg1 int, arg2 string) bool {
	res := len(arg2) == arg1
	return res
}

// mock example of function-adapter
func mock_adapter(cx base.Context, args []any) (any, error) {
	if len(args) != 2 {
		return false, errors.New("moch_target func: incorrect count of args")
	}
	a0, ok := args[0].(int64)
	if !ok {
		return false, errors.New("moch_target func: incorrect 1-st arg type")
	}
	a1, ok := args[1].(string)
	if !ok {
		return false, errors.New("moch_target func: incorrect 2-st arg type")
	}
	res := mock_target(int(a0), a1)
	return res, nil
}

/// ========================  Add new Builtin func ======================

func BuiltFunc(cx base.Context, name string, adapter func(base.Context, []any) (any, error), resType *base.Type) {
	nf := &NFunc{Name: name, fun: adapter}
	cx.AddFunc(nf)
}

func BuiltConstr(cx base.Context, name string, adapter func(base.Context, []any) (any, error), resType *base.Type) {
	nf := &NFunc{Name: name, fun: adapter}
	tp := cx.GetType(name)
	if tp == nil {
		panic("can't find type " + name)
	}
	tp.Construct = nf
	// ce := cx.GetElem(name)
	// ttp, ok := ce.V.(*base.Type)
	// if !ok {
	// 	panic("No such type # 1!!!")
	// }
	// fmt.Printf(" BCnstr `%s` %T con: %T \n", ttp.Name, ttp, ttp.Construct)
	// cx.AddFunc(nf)
}

// Builtin method object
type MFunc struct {
	Name string
	fun  func(base.Context, any, []any) (any, error) // func-apapter called in Do()

	inst  any
	args  []any          // passed args
	mvals map[string]any // passed named args

	resV *base.Val
}

func (fn *MFunc) Do(cx base.Context) error {
	fn.resV = nil
	// TODO: named args, vary arg list, multi-result
	res, err := fn.fun(cx, fn.inst, fn.args)
	if err != nil {
		return err
	}
	fn.resV = base.NewVal(res)
	return nil
}

func (fn *MFunc) Get() *base.Val {
	if fn.resV == nil {
		return nil
	}
	return fn.resV
}

func (fn *MFunc) GetName() string {
	return fn.Name
}

func (fn *MFunc) SetInst(inst any) {
	fn.inst = inst
}

func (fn *MFunc) SetArgVals(vals []any, mvals map[string]any) {
	fn.args = vals
	fn.mvals = mvals
}

func BuiltMethod(cx base.Context, typeName string, name string, adapter func(base.Context, any, []any) (any, error)) error {
	fn := &MFunc{Name: name, fun: adapter}
	tp := cx.GetType(typeName)
	if tp == nil {
		panic("can't find type " + name)
	}
	tp.AddMethod(fn)
	return nil
}
