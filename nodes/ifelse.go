package nodes

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
)

type IfNode struct {
	condition base.Expression
	preCond   *BlockExpr
	BlockIf   *BlockExpr
	BlockElse *BlockExpr
	res       any
}

func (nd *IfNode) Get() *base.Val {
	return nil
}

func (nd *IfNode) Do(cx base.Context) error {
	nd.res = nil
	// pre cond

	evalBlock := nd.BlockIf
	// cond
	if err := nd.condition.Do(cx); err != nil {
		return err
	}
	condRes := nd.condition.Get()
	if condRes == nil {
		return errors.New("nil res from if-condition")
	}
	switch condV := condRes.V.(type) {
	case bool:
		if !condV {
			if nd.BlockElse != nil {
				evalBlock = nd.BlockElse
			} else {
				// correct end of if when condition returns false
				return nil
			}
		}
	default:
		return errors.New("if-cond returns not bool")
	}

	// eval
	evalBlock.Do(cx)
	res := evalBlock.Get()
	if res != nil {
		nd.res = res.V
	}
	return nil
}

func (nd *IfNode) Add(sub base.Expression) error {
	nd.BlockIf.Add(sub)
	return nil
}

func (nd *IfNode) MakeElse() *BlockExpr {
	nd.BlockElse = NewBlock()
	return nd.BlockElse
}

func NewIf(cond base.Expression) *IfNode {
	var prev *BlockExpr
	// if cond is ;-separated expr
	switch conExpr := cond.(type) {
	case *SequenceSemicolon:
		//
		if len(conExpr.Subs) < 2 {
			//bad case: a=1;
		}
		cond = conExpr.Subs[len(conExpr.Subs)-1]
		prev = NewBlock()
		for _, sub := range conExpr.Subs[:len(conExpr.Subs)-1] {
			prev.Add(sub)
		}

	}
	node := &IfNode{
		condition: cond,
		BlockIf:   NewBlock(),
		preCond:   prev,
	}
	return node
}
