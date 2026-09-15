package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

type GenLoop struct {
	IterExp *LeftArrow
	Loop    *ForSourceNode
	Subs    []base.Expression
	Guard   base.Expression
	SubLoop *GenLoop

	rex *GenResExpr
}

func (gn *GenLoop) Get() *base.Val {
	return nil
}

func (gn *GenLoop) DoSub(cx base.Context, post base.Expression) error {
	snd := gn.Loop
	err := snd.Init(cx)
	if err != nil {
		return err
	}
	for {
		if snd.assign.Finished() {
			break // correct finich
		}
		err = snd.assign.Next()
		if err != nil {
			return err
		}
		err = snd.Block.Do(cx)
		if err != nil {
			return err
		}
		if gn.Guard != nil {
			gn.Guard.Do(cx)
			gres := GetExprVal(gn.Guard, nil)
			gv, ok := gres.(bool)
			if !ok {
				return fmt.Errorf("generators guard should return bool but given : %T", gres)
			}
			if !gv {
				// skip result iteration by condition
				continue
			}
		}
		if gn.SubLoop == nil {
			err = post.Do(cx)
			if err != nil {
				return err
			}
		} else {
			err = gn.SubLoop.DoSub(cx, post)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gn *GenLoop) Do(cx base.Context) error {
	post := gn.rex
	err := gn.DoSub(cx, post)
	return err
}

func NewGenLoop(arrEx *LeftArrow, subs []base.Expression, guard base.Expression) *GenLoop {
	gn := &GenLoop{IterExp: arrEx, Guard: guard}
	gn.Loop = NewForSource(gn.IterExp)
	for _, sx := range subs {
		gn.Loop.Add(sx)
	}
	return gn
}

// ===
type GenResExpr struct {
	Exp base.Expression
	res []any
}

func (gg *GenResExpr) Get() *base.Val {
	return base.NewVal(gg.res)
}

func (gg *GenResExpr) Do(cx base.Context) error {
	if gg.res == nil {
		gg.res = []any{}
	}
	err := gg.Exp.Do(cx)
	if err != nil {
		return err
	}
	rval := GetExprVal(gg.Exp, nil)
	gg.res = append(gg.res, rval)
	return nil
}

//===

type ListComprh struct {
	ResEx    base.Expression
	LoopNode *GenLoop

	res *objects.ListVal
}

func (cm *ListComprh) Get() *base.Val {
	if cm.res == nil {
		return nil
	}
	return base.NewVal(cm.res)
}

func (cm *ListComprh) Do(cx base.Context) error {
	// fmt.Printf("ListCompr.Do\n")
	inCx := cx.SubContext()
	rr := []any{}
	post := &GenResExpr{Exp: cm.ResEx, res: rr}
	cm.LoopNode.rex = post
	cm.LoopNode.Do(inCx)
	// Run ResEx in the each iter of last sub loop
	rr = post.res
	cm.res = objects.NewListVal(rr)
	return nil
}

func NewListComprh(res base.Expression, loop *GenLoop) *ListComprh {
	return &ListComprh{ResEx: res, LoopNode: loop}
}
