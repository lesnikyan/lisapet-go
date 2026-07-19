package objects

import (
	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/nodes"
)

type Function struct {
	Name string

	defArgs map[string]*base.Var
	args    []*base.Var

	Block nodes.BlockExpr

	res any
}

func (fn *Function) Do(cx base.Context) error {
	return nil
}

func (fn *Function) Get() *base.Val {

	return base.NewVal(fn.res)
}

func NewFunction()
