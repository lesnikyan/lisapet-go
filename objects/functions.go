package objects

import (
	"github.com/lesnikyan/lisapet-go/base"
)

type ArgExp struct {
	Name string
	// VExp       *VarExpr
	Varp       *base.Var
	Type       *base.Type
	StrictType bool

	res any
}

func (ax *ArgExp) Get(cx base.Context) *base.Val {
	// return ax.VExp.Get()
	return nil
}

type Function struct {
	Name string

	// defArgs map[string]*base.Var
	defVals map[string]*base.Val
	argVals []any
	nmVals  map[string]any
	dfnArgs []base.Expression
	Args    []*ArgExp
	mtInst  *ArgExp // obj instance for method

	// Block  *BlockExpr
	Block  FSubBlock
	defCtx base.Context

	res    any
	resVal *base.Val
}

func (fn *Function) SetArgVals(vals []any, nvals map[string]any) {
	// fmt.Printf(" Fu.SetArgVals  ovs=%d, nvs:%d  \n", len(vals), len(nvals))
	fn.Block.SetArgVals(vals, nvals)
}

func (fn *Function) Do(cx base.Context) error {
	fn.res = nil
	fn.resVal = nil

	// fmt.Printf(" ---- ob.Fu.Do#1 \n")

	// inner context
	inCx := fn.defCtx.SubContext()
	err := fn.Block.Do(inCx)
	if err != nil {
		return err
	}
	fn.resVal = fn.Block.Get()
	return nil
}

func (fn *Function) Get() *base.Val {
	// return base.NewVal(fn.res)
	return fn.resVal
}

func (fn *Function) GetName() string {
	return fn.Name
}

func NewFunction(name string, block FSubBlock, ctx base.Context) *Function {
	return &Function{Name: name, Block: block, defCtx: ctx}
}

// / =======================
type FSubBlock interface {
	SetArgVals(vals []any, nvals map[string]any)
	PrepareArgs(cx base.Context) error
	Do(cx base.Context) error
	Get() *base.Val
	GetPopUp() *base.PopUp
}
