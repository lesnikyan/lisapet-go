package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

type BlockExpr struct {
	subs     []base.Expression
	res      any
	parMark  bool
	aborted  bool      // after: return, break
	abortRes *base.Val // after return
	PopUp    *PopUp
}

func (bk *BlockExpr) IsParent() bool {
	return bk.parMark
}

func (bk *BlockExpr) IsAborted() bool {
	return bk.aborted
}

func (bk *BlockExpr) GetPopUp() *PopUp {
	return bk.PopUp
}

func (bk *BlockExpr) Add(sub base.Expression) {
	bk.subs = append(bk.subs, sub)
}

func (bk *BlockExpr) Get() *base.Val {
	return base.NewVal(bk.res)
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
	bk.PopUp = nil
	if len(bk.subs) == 0 {
		return nil
	}
	hasRes := true // most expressions has result
	var popUp NodeType = 0
	for _, exp := range bk.subs {
		fmt.Printf("Bl.Do#0: %T, %v\n", exp, exp)
		err := exp.Do(cx)
		if err != nil {
			return err
		}
		last = exp
		stop := false
		switch cur := exp.(type) {
		case *BreakExp:
			popUp = NodeBreak
			hasRes = false
			stop = true
		case *ContinueExp:
			popUp = NodeContinue
			hasRes = false
			stop = true
		case *ReturnExp:
			popUp = NodeReturn
			hasRes = false
			if cur.Sub != nil {
				hasRes = true
			}
			stop = true
		case *IfNode:
			inPup := cur.GetPopUp()
			if inPup != nil {
				bk.PopUp = inPup
				bk.res = cur.Get()
				return nil
			}
			// other resulting expressions: func def, func call, operators, value, if-else, match, etc
		case LoopNode:
			hasRes = false
		default:
			hasRes = true
			// case *ValExpr:
			// 	hasRes = true
			// case *OperAssign, *OperBin:
			// 	hasRes = true
			// case *ListExpr, *TupleExpr, *DictExpr:
			// 	hasRes = true
			// case *SequenceComma, *SequenceSemicolon:
			// 	hasRes = true
		}
		if stop {
			break
		}
	}
	// if popup break: break, continue, return
	if popUp > 0 {
		var popRes *base.Val
		if hasRes {
			popRes = last.Get()
		}
		bk.PopUp = &PopUp{
			Parent: popUp,
			Res:    popRes,
		}
		return nil
	}
	// if block has result
	if !hasRes {
		return nil
	}
	last = bk.subs[len(bk.subs)-1]
	if last != nil {
		bk.res = last.Get()
	}
	return nil
}

func NewBlock() *BlockExpr {
	bk := &BlockExpr{}
	bk.subs = []base.Expression{}
	return bk
}
