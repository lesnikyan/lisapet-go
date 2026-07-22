package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

type ForExpr interface {
	IsParent() bool
	Get() *base.Val
	Do(cx base.Context) error
	Add(sub base.Expression)

	Init(cx base.Context) error
	// Post(cx base.Context) error
	Loop(cx base.Context) error
}

type LoopNode interface {
	IsParent() bool
	Do(cx base.Context) error
	Add(sub base.Expression)
	Loop(cx base.Context) error
}

type ForCondNode struct {
	InitEx base.Expression
	CondEx base.Expression
	PostEx base.Expression
	Block  *BlockExpr

	SourceCollection any // list, dict, tuple, comprehension, generator, slice, etc
	SourceGen        any

	IsCond   bool // for ex; ex; ex
	IsSource bool // for a <- b
	stopped  bool

	res any
}

func (nd *ForCondNode) IsParent() bool {
	return true
}

func (nd *ForCondNode) Get() *base.Val {
	return nil
}

func (bk *ForCondNode) GetPopUp() *PopUp {
	return bk.Block.GetPopUp()
}

func (nd *ForCondNode) Init(cx base.Context) error {
	err := nd.InitEx.Do(cx)
	return err
}

func (nd *ForCondNode) Check(cx base.Context) (bool, error) {
	err := nd.CondEx.Do(cx)
	if err != nil {
		return false, err
	}
	val := nd.CondEx.Get()
	if val == nil {
		return false, errors.New("No val of for-loop condition")
	}
	return (val.V).(bool), nil
}

func ForStop(bk *BlockExpr) bool {
	pup := bk.GetPopUp()
	if pup == nil {
		return false
	}
	// TODO: return
	switch pup.Parent {
	case NodeBreak:
		fmt.Println("#For.Break")
		return true
	case NodeContinue:
		fmt.Println("#For.Continue")
		return false
	case NodeReturn:
		fmt.Println("#For.Return")
		return true
	}
	return false
}

func (nd *ForCondNode) Loop(cx base.Context) error {
	// TODO: implement in loop Block: break, continue, return
	// stop := false
	for {
		// check condition
		ok, err := nd.Check(cx)
		if err != nil {
			return err
		}
		if !ok {
			break // correct finish
		}
		// Do block
		err = nd.Block.Do(cx)
		if err != nil {
			return err
		}

		// pup := nd.Block.GetPopUp()
		// if pup != nil {
		// 	// TODO: return, break
		// 	switch pup.Parent {
		// 	case NodeBreak:
		// 		fmt.Println("#For.Break")
		// 		stop = true
		// 	case NodeContinue:
		// 		fmt.Println("#For.Continue")
		// 	case NodeReturn:
		// 		fmt.Println("#For.Return")
		// 		stop = true
		// 	}
		// }
		stop := ForStop(nd.Block)
		// block.aborted(); // break | return
		// block.Brocken(); // break case
		// block.HasReturn; block.GetReturned()
		if stop {
			break
		}
		// Do post
		err = nd.PostEx.Do(cx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (nd *ForCondNode) Do(cx base.Context) error {
	fmt.Println("# FOR")

	inCx := objects.NewContext(cx)
	// Init
	nd.Init(inCx)

	// Loop
	err := nd.Loop((inCx))
	return err
}

func (nd *ForCondNode) Add(sub base.Expression) {
	nd.Block.Add(sub)
}

func NewForCond(subEx *SequenceSemicolon) *ForCondNode {
	if len(subEx.Subs) != 3 {
		// bad case. check and resolve: (for ; true ; )
	}
	init := subEx.Subs[0]
	cond := subEx.Subs[1]
	post := subEx.Subs[2]
	return &ForCondNode{InitEx: init, CondEx: cond, PostEx: post, Block: NewBlock()}
}

// ***************************************************

type ForSourceNode struct {
	// InitEx base.Expression
	// CondEx base.Expression
	// PostEx base.Expression
	IterExpr *LeftArrow
	Block    *BlockExpr

	// SourceCollection any // list, dict, tuple, comprehension, generator, slice, etc
	// iter             SourceIter
	assign *IterAssign

	res any
}

func (nd *ForSourceNode) IsParent() bool {
	return true
}

func (nd *ForSourceNode) Get() *base.Val {
	return nil
}

func (bk *ForSourceNode) GetPopUp() *PopUp {
	return bk.Block.GetPopUp()
}

func (nd *ForSourceNode) Init(cx base.Context) error {
	err := nd.IterExpr.Do(cx)
	if err != nil {
		return err
	}
	nd.assign = nd.IterExpr.GetIterAssign()
	nd.assign.Init()
	return nil
}

func (nd *ForSourceNode) Loop(cx base.Context) error {
	for {
		if nd.assign.Finished() {
			break // correct finich
		}
		err := nd.assign.Next()
		if err != nil {
			return err
		}
		err = nd.Block.Do(cx)
		if err != nil {
			return err
		}
		stop := ForStop(nd.Block)
		if stop {
			break
		}
	}
	return nil
}

func (nd *ForSourceNode) Do(cx base.Context) error {
	fmt.Println("# FOR")

	inCx := objects.NewContext(cx)
	// Init
	nd.Init(inCx)

	// Loop
	nd.Loop((inCx))
	return nil
}

func (nd *ForSourceNode) Add(sub base.Expression) {
	nd.Block.Add(sub)
}

func NewForSource(subEx *LeftArrow) *ForSourceNode {
	subEx.IsIter = true
	return &ForSourceNode{Block: NewBlock(), IterExpr: subEx}
}

type WhileNode struct {
	CondEx base.Expression
	Block  *BlockExpr
}
