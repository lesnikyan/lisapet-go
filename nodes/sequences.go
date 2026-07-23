package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	ob "github.com/lesnikyan/lisapet-go/objects"
)

// ===========

type OperColon struct {
	left  base.Expression
	right base.Expression
	Oper  *Oper
	res   any
}

func (op *OperColon) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperColon) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperColon) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *OperColon) GetPair() *ColonPair {
	return &ColonPair{Left: op.left, Right: op.right}
}

func (op *OperColon) Do(cx base.Context) error {
	err1 := op.left.Do(cx)
	if err1 != nil {
		return err1
	}
	err2 := op.right.Do(cx)
	if err2 != nil {
		fmt.Println("OpAssign.R error", err2)
		return err2
	}

	return nil
}

type SequenceComma struct {
	Subs []base.Expression
	res  any
	// TODO: res type - of collection, func def args, func call argc, struct ded, struct constr, hmm...
	// []any ?
}

func (cs *SequenceComma) Get() *base.Val {
	// TODO
	return base.NewVal(cs.res)
}

func (cs *SequenceComma) Do(ctx base.Context) error {
	// TODO: loop processing over subs
	var err error
	for _, sub := range cs.Subs {
		err = sub.Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cs *SequenceComma) Add(elem base.Expression) {
	cs.Subs = append(cs.Subs, elem)
}

func (cs *SequenceComma) SetSubs(elems []base.Expression) {
	// subs := make([]base.Expression, len(elems))
	// copy(subs, elems)
	cs.Subs = elems
}

type SequenceSemicolon struct {
	Subs []base.Expression
	res  any
}

func (cs *SequenceSemicolon) Get() *base.Val {
	return base.NewVal(cs.res)
}

func (cs *SequenceSemicolon) Do(ctx base.Context) error {
	// TODO: loop processing over subs
	var err error
	for _, sub := range cs.Subs {
		err = sub.Do(ctx)
		if err != nil {
			return err
		}
	}
	// TODO: result of last sub expr
	return nil
}

func (cs *SequenceSemicolon) Add(elem base.Expression) {
	cs.Subs = append(cs.Subs, elem)
}

func (cs *SequenceSemicolon) SetSubs(elems []base.Expression) {
	cs.Subs = elems
}

// =====

// a .. b
type Dots2Expr struct {
	Left  base.Expression
	Right base.Expression

	res any
}

func (op *Dots2Expr) SetLeft(xp base.Expression) {
	op.Left = xp
}
func (op *Dots2Expr) SetRight(xp base.Expression) {
	op.Right = xp
}

func (op *Dots2Expr) Get() *base.Val {
	return nil
}

func (op *Dots2Expr) Do(cx base.Context) error {
	err := op.Left.Do(cx)
	if err != nil {
		return err
	}
	err = op.Right.Do(cx)
	if err != nil {
		return err
	}
	return nil
}

func (op *Dots2Expr) GetNumSeq() *NumSeqExpr {
	return &NumSeqExpr{Left: op.Left, Right: op.Right}
}

func NewDots2(lArg base.Expression, rArg base.Expression) *Dots2Expr {
	return &Dots2Expr{Left: lArg, Right: rArg}
}

type NumSeqExpr struct {
	Left  base.Expression
	Right base.Expression

	res *ob.NumSeqGen
}

// [a1<, a2> .. b]
func (sq *NumSeqExpr) Do(cx base.Context) error {
	err := sq.Left.Do(cx)
	if err != nil {
		return err
	}
	err = sq.Right.Do(cx)
	if err != nil {
		return err
	}
	// TODO: implement <n,m> case of left arg
	var lNum int64 = 0
	lv := sq.Left.Get()
	if lv != nil {
		lNum = lv.V.(int64)
	}
	var rNum int64 = 0
	rv := sq.Right.Get()
	if rv != nil {
		rNum = rv.V.(int64)
	}
	sq.res = ob.NewNumSeqGen(lNum, rNum, 1)
	return nil
}

func (sq *NumSeqExpr) Get() *base.Val {
	return base.NewVal(sq.res)
}

// a...
type Dots3Right struct {
}
