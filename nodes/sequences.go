package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
	ob "github.com/lesnikyan/lisapet-go/objects"
)

// ===========

type ColonUsage = int

const (
	_ ColonUsage = iota + 200

	ColonType // var:type
	ColonElem // key: val
)

type OperColon struct {
	Left  base.Expression
	Right base.Expression
	Oper  *Oper
	Usage ColonUsage

	res any
}

func (op *OperColon) AddToRight(xp base.Expression) {
	fmt.Println("OperAssign. AddTR:", xp)
	switch lex := op.Right.(type) {
	case base.SuperExpr:
		lex.Add(xp)
	}
}

func (op *OperColon) Add(sub base.Expression) {}

func (op *OperColon) SetLeft(xp base.Expression) {
	op.Left = xp
}
func (op *OperColon) SetRight(xp base.Expression) {
	op.Right = xp
}

func (op *OperColon) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *OperColon) GetPair() *ColonPair {
	lex := op.Left
	if lex == nil {
		lex = &EmptyExpr{}
	}
	rex := op.Right
	if rex == nil {
		rex = &EmptyExpr{}
	}
	return &ColonPair{Left: lex, Right: rex}
}

func (op *OperColon) DoVar(cx base.Context) error {
	// fmt.Printf(" ' : ' <DoVar \n")
	err1 := op.Left.Do(cx)
	if err1 != nil {
		return err1
	}
	varx, ok := op.Left.(*VarExpr)
	if !ok {
		return errors.New("colon: Not var in Left")
	}
	tupx, ok := op.Right.(*VarExpr)
	if !ok {
		return errors.New("colon: Not var in Left")
	}
	varx.NewVar(cx)

	err2 := tupx.Do(cx)
	if err2 != nil {
		fmt.Println("OperColon.R error", err2)
		return err2
	}
	// tname := tupx.name
	rtype := cx.GetType(tupx.name)
	// rval := tupx.Get()
	// fmt.Printf(" ' : ' DoVar2 (%T, %v) \n", tupx, tupx)
	// fmt.Printf(" ' : ' DoVar3 (%T, %v) \n", rtype, rtype)
	if rtype == nil {
		return errors.New("colon: can't find type")
	}
	// fmt.Printf(" ' : ' DoVar4 \n")

	if !ok {
		return errors.New("colon: Right part is not type")
	}
	varx.vr.Type = rtype // usage: check type for typed vars
	varx.vr.StrictType = true
	op.res = varx.vr

	// varx.vr.StrictType = true
	return nil
}

// in base case - declaration of typed var >>  x: int
func (op *OperColon) Do(cx base.Context) error {
	switch op.Usage {
	case ColonType:
		op.DoVar(cx)
	}
	return nil
}

// Commas

type SequenceComma struct {
	Subs []base.Expression
	res  any
	// TODO: res type - of collection, func def args, func call argc, struct ded, struct constr, hmm...
	// []any ?
}

func (cs *SequenceComma) Get() *base.Val {
	vals := cs.GetVals()
	return vals
}

func (cs *SequenceComma) GetVals() *base.Val {
	// return array of vals, used as internal value for other expressions
	vals := make([]any, len(cs.Subs))
	for i, ex := range cs.Subs {
		val := GetExprVal(ex, nil)
		// var val any
		if val == nil {
			val = &objects.Null{}
		}
		vals[i] = val
	}
	cs.res = vals
	return base.NewVal(cs.res)
}

func (cs *SequenceComma) Do(ctx base.Context) error {
	// TODO: loop processing over subs
	cs.res = nil
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

// Semicolons

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
	Left  base.Expression // any val result or pair `x, y`
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
	// fmt.Printf(" SqNum left: %T \n", sq.Left)
	var start int64 = 0
	var step int64 = 1
	switch lExp := sq.Left.(type) {
	case *SequenceComma:
		subs := lExp.Subs
		if len(subs) != 2 {
			return errors.New("NumGen incorrect num of left arg ")
		}
		v1 := subs[0].Get()
		if v1 == nil {
			return errors.New("NumGen incorrect 1-st agr of left ")
		}
		start = v1.V.(int64)
		v2 := subs[1].Get()
		if v2 == nil {
			return errors.New("NumGen incorrect 1-st agr of left ")
		}
		n2 := v2.V.(int64)
		step = n2 - start
	default:
		lv := GetExprVal(sq.Left, nil)
		if lv != nil {
			start = lv.(int64)
		}

	}
	var max int64 = 0
	rv := GetExprVal(sq.Right, nil)
	if rv != nil {
		max = rv.(int64)
	}
	sq.res = ob.NewNumSeqGen(start, max, step)
	return nil
}

func (sq *NumSeqExpr) Get() *base.Val {
	return base.NewVal(sq.res)
}

// a...
type Dots3Right struct {
}
