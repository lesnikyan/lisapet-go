package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

type anym map[any]any

type NamedArgs struct {
	Nvals map[string]any
}

func NewNamedArgs(v map[string]any) *NamedArgs {
	return &NamedArgs{Nvals: v}
}

// ---

type NFunc struct {
	Name   string
	args   []any                                  // passed args
	nmargs map[string]any                         // passed named args
	fun    func(base.Context, []any) (any, error) // func-apapter called in Do()
	serv   bool

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
	// fmt.Printf(" - NFunc.Do# f: %s  br(%T : %v) fr(%T : %v) || err: %v \n", fn.Name, fn.resV, fn.resV, res, res, err)
	return nil
}

func (fn *NFunc) IsServ() bool {
	return fn.serv
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

func (fn *NFunc) SetArgVals(vals []any, nmvals map[string]any) {
	if len(nmvals) > 0 {
		fn.nmargs = nmvals
		nm := NewNamedArgs(nmvals)
		vals = append(vals, nm)
	}
	// println("NFunc.SetArgVals", len(vals))
	fn.args = vals
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

func BuiltServeFunc(cx base.Context, name string, adapter func(base.Context, []any) (any, error), resType *base.Type) {
	nf := &NFunc{Name: name, fun: adapter, serv: true}
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

// ==== Builtin method object

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

func (fn *MFunc) SetArgVals(vals []any, nmvals map[string]any) {
	if len(nmvals) > 0 {
		fn.mvals = nmvals
		nm := NewNamedArgs(nmvals)
		vals = append(vals, nm)
	}
	fn.args = vals
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

// ==== Composed function

type Composed struct {
	Name string
	args []any // passed args
	// fun    func(base.Context, []any) (any, error) // func-apapter called in Do()
	Funcs []base.FuncVal

	resV *base.Val
}

// Builtin function object
func (fn *Composed) Do(cx base.Context) error {
	fn.resV = nil
	args := fn.args
	var res any
	for i, f := range fn.Funcs {
		// fmt.Printf(" - Composed.Do# fn<%s> comz f:(%T : %v) \n", fn.Name, f, f)
		nmd := map[string]any{}
		f.SetArgVals(args, nmd)
		err := f.Do(cx)
		if err != nil {
			return err
		}
		vv := f.Get()
		if vv == nil {
			return fmt.Errorf("func %s in composed returned nil in %d iter", f.GetName(), i)
		}
		res = vv.V
		args[0] = vv.V
	}
	fn.resV = base.NewVal(res)
	// fmt.Printf(" - Composed.Do# f: %s  br(%T : %v) fr(%T : %v) || err: %v \n", fn.Name, fn.resV, fn.resV, res, res, err)
	return nil
}

func (fn *Composed) IsServ() bool {
	return false
}

func (fn *Composed) Get() *base.Val {
	if fn.resV == nil {
		return nil
	}
	return fn.resV
}

func (fn *Composed) GetName() string {
	return fn.Name
}

func (fn *Composed) SetArgVals(vals []any, nmvals map[string]any) {
	// composed can't have its own named args
	fn.args = vals
}

func NewComposed(name string, funcs []base.FuncVal) *Composed {
	return &Composed{Name: name, Funcs: funcs}
}
