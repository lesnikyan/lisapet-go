package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

type MatchNode struct {
	ArgExp base.Expression
	Cases  []*MatchCaseNode
	PopUp  *base.PopUp

	res any // if match is a last expr in block and returns value
}

func (nd *MatchNode) IsParent() bool {
	return true
}

func (bk *MatchNode) GetPopUp() *base.PopUp {
	return bk.PopUp
}

func (mt *MatchNode) Get() *base.Val {
	if mt.res == nil {
		return nil
	}
	return base.NewVal(mt.res)
}

func (mt *MatchNode) AddCase(cs *MatchCaseNode) {
	mt.Cases = append(mt.Cases, cs)
}

func (mt *MatchNode) Add(exp base.Expression) {
	// fmt.Printf(" MtNode.Add: %T\n", exp)
	cs, ok := exp.(*MatchCaseNode)
	if ok {
		mt.AddCase(cs)
	}
}

func (mt *MatchNode) Do(cx base.Context) error {
	err := mt.ArgExp.Do(cx)
	if err != nil {
		return err
	}
	arg := GetExprVal(mt.ArgExp, nil)

	for i, cs := range mt.Cases {
		inCx := cx.SubContext()
		// cs.Pattern.PutArg(arg)
		// fmt.Printf(" MtNode.Do.for %d: %T >> %T: %v\n", i, cs, arg, arg)
		ptok, err := cs.Pattern.Match(inCx, arg)
		if err != nil {
			return errors.Join(fmt.Errorf("Error in match pattern %d: %T", i, cs.Pattern), err)
		}
		// cres := GetExprVal(cs.Pattern, nil)
		// cf, ok := cres.(bool)
		// if !ok {
		// 	return fmt.Errorf("Matching pattern must return bool, case %d returned %T", i, cres)
		// }
		if !ptok {
			// pattern not matched
			continue
		}
		// eval block
		err = cs.Block.Do(cx)
		if err != nil {
			return err
		}
		bres := GetExprVal(cs.Block, nil)
		if bres != nil {
			mt.res = bres
		}
		break
	}
	return nil
}

// ===

type MatchCaseNode struct {
	Pattern MatchPattern
	Block   *BlockExpr
	PopUp   *base.PopUp

	res any // if match returning result
}

func (mt *MatchCaseNode) Get() *base.Val {
	if mt.res == nil {
		return nil
	}
	return base.NewVal(mt.res)
}

func (mt *MatchCaseNode) Add(exp base.Expression) {
	// cs, ok := exp.(*MatchCaseNode)
	// fmt.Printf(" MatchCaseNode.Add: %T\n", exp)
	if mt.Block == nil {
		mt.Block = NewEmptyBlock()
	}
	mt.Block.Add(exp)
}

func (nd *MatchCaseNode) IsParent() bool {
	return true
}

func (bk *MatchCaseNode) GetPopUp() *base.PopUp {
	return bk.PopUp
}

func (mt *MatchCaseNode) Do(cx base.Context) error {
	return nil
}

func NewMatchCase(ptt MatchPattern) *MatchCaseNode {
	return &MatchCaseNode{Pattern: ptt, Block: NewEmptyBlock()}
}

// ====
