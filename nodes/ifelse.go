package nodes

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
)

type IfNode struct {
	condition base.Expression
	preCond   *BlockExpr
	BlockIf   *BlockExpr
	BlockElse *ElseNode
	PopUp     *PopUp
	res       any
}

func (nd *IfNode) IsParent() bool {
	return true
}

func (nd *IfNode) Get() *base.Val {
	return nil
}

func (bk *IfNode) GetPopUp() *PopUp {
	return bk.PopUp
}

func (nd *IfNode) PreDo(cx base.Context) error {
	if nd.preCond == nil {
		return nil
	}
	return nd.preCond.Do(cx)
}

func (nd *IfNode) Do(cx base.Context) error {
	nd.res = nil
	nd.PopUp = nil
	// pre cond
	nd.PreDo(cx)
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
				evalBlock = nd.BlockElse.Block
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
	pup := evalBlock.GetPopUp()
	if pup != nil {
		// some aborting
		nd.PopUp = pup
	}
	res := evalBlock.Get()
	if res != nil {
		nd.res = res.V
	}
	return nil
}

func (nd *IfNode) Add(sub base.Expression) {
	nd.BlockIf.Add(sub)
	// return nil
}

// func (nd *IfNode) MakeElse() *BlockExpr {
// 	nd.BlockElse = NewBlock()
// 	return nd.BlockElse
// }

func (nd *IfNode) SetElse(block *ElseNode) {
	nd.BlockElse = block
}

func NewIf(cond base.Expression) *IfNode {
	var prev *BlockExpr
	// if cond is ;-separated expr
	// fmt.Println("# IF", cond)
	switch conExpr := cond.(type) {
	case *SequenceSemicolon:
		//
		// fmt.Println("# IF#2", len(conExpr.Subs))
		if len(conExpr.Subs) < 2 {
			//bad case: a=1;
			return nil
		}
		cond = conExpr.Subs[len(conExpr.Subs)-1]
		// fmt.Printf("#--IF (;) #1 (%T: %v), #2(%T: %v) \n", conExpr, conExpr, cond, cond)
		prev = NewBlock()
		for _, sub := range conExpr.Subs[:len(conExpr.Subs)-1] {
			prev.Add(sub)
		}
	default:

	}
	node := &IfNode{
		condition: cond,
		BlockIf:   NewBlock(),
		preCond:   prev,
	}
	return node
}

type ElseNode struct {
	Block   *BlockExpr
	SlideIf *IfNode
	Slided  bool
}

func (nd *ElseNode) IsParent() bool {
	return true
}

// for case: else if condition
func (nd *ElseNode) SetSlide(ndif *IfNode) {
	nd.SlideIf = ndif
	nd.Block.Add(ndif)
	nd.Slided = true
}

func (nd *ElseNode) Get() *base.Val {
	return nd.Block.Get()
}

func (bk *ElseNode) GetPopUp() *PopUp {
	return bk.Block.GetPopUp()
}

func (nd *ElseNode) Do(cx base.Context) error {
	return nd.Block.Do(cx)
}

func (nd *ElseNode) Add(sub base.Expression) {
	if nd.Slided {
		nd.SlideIf.Add(sub)
		return
	}
	nd.Block.Add(sub)
}

func NewElseNode() *ElseNode {
	return &ElseNode{Block: NewBlock()}
}
