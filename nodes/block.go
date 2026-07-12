package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

type BlockExpr struct {
	subs    []base.Expression
	res     any
	parMark bool
}

func (bk *BlockExpr) IsParent() bool {
	return bk.parMark
}

func (bk *BlockExpr) First() base.Expression {
	if len(bk.subs) > 0 {
		return bk.subs[0]
	}
	return nil
}

func (bk *BlockExpr) Do(cx base.Context) error {
	// var res any = nil
	var last base.Expression
	bk.res = nil
	if len(bk.subs) == 0 {
		return nil
	}
	for _, exp := range bk.subs {
		fmt.Printf("Bl.Do#0: %T, %v\n", exp, exp)
		err := exp.Do(cx)
		if err != nil {
			return err
		}
		last = exp
	}
	last = bk.subs[len(bk.subs)-1]
	if last != nil {
		bk.res = last.Get()
	}
	return nil
}

func (bk *BlockExpr) Add(sub base.Expression) {
	bk.subs = append(bk.subs, sub)
}

func (bk *BlockExpr) Get() *base.Val {
	return base.NewVal(bk.res)
}

func NewBlock() *BlockExpr {
	bk := &BlockExpr{}
	bk.subs = []base.Expression{}
	return bk
}
