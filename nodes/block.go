package nodes

import (
	"github.com/lesnikyan/lisapet-go/base"
)

type BlockExpr struct {
	subs []base.Expression
	res  any
}

func (bk *BlockExpr) Do(cx base.Context) error {
	// var res any = nil
	var last base.Expression
	bk.res = nil
	for _, exp := range bk.subs {
		err := exp.Do(cx)
		if err != nil {
			return err
		}
		last = exp
	}
	if last != nil {
		bk.res = last.Get()
	}
	return nil
}
func (bk *BlockExpr) Get() *base.Val {
	return base.NewVal(bk.res)
}
