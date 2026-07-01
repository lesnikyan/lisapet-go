package nodes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	ob "github.com/lesnikyan/lisapet-go/objects"
)

type MockExpr struct {
}

func (cs *MockExpr) Get() any { return nil }

func (cs *MockExpr) Do(ctx base.Context) error { return nil }

type OperAssign struct {
	Oper  string
	left  base.Expression
	right base.Expression
	res   any
}

func (op *OperAssign) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperAssign) SetRight(xp base.Expression) {
	op.right = xp
}

func SeqInfo(seq any) string {
	sep := " "
	subinf := []string{}
	switch sqq := seq.(type) {
	case *SequenceComma:
		sep = ", "
		for _, elm := range sqq.Subs {
			subinf = append(subinf, OperArgsInfo(elm))
		}
	}
	return strings.Join(subinf, sep)
}

func OperArgsInfo(expr any) string {
	var aL string
	var aR string
	var opStr string
	switch exx := expr.(type) {
	case *OperAssign:
		aL = OperArgsInfo(exx.left)
		aR = OperArgsInfo(exx.right)
		opStr = exx.Oper
	case *OperBin:
		aL = OperArgsInfo(exx.left)
		aR = OperArgsInfo(exx.right)
		opStr = exx.Oper.Sign
	case *OperColon:
		aL = OperArgsInfo(exx.left)
		aR = OperArgsInfo(exx.right)
		opStr = exx.Oper.Sign
	case *Brackets:
		subs := OperArgsInfo(exx.Sub)
		return fmt.Sprintf("(%v)", subs)
	case *TupleExpr:
		subs := OperArgsInfo(exx.Seq)
		return fmt.Sprintf("tuple(%v)", subs)
	case *ListExpr:
		subs := OperArgsInfo(exx.Seq)
		return fmt.Sprintf("list[%v]", subs)
	case *DictExpr:
		subs := OperArgsInfo(exx.Seq)
		return fmt.Sprintf("dict{%v}", subs)
	case *SequenceComma:
		return SeqInfo(expr)
	case *ValExpr:
		return fmt.Sprintf("%v", exx.Val)
	case *VarExpr:
		return fmt.Sprintf("(%s:)", exx.name)
	default:
		return "<non-oper object>"
	}
	return fmt.Sprintf("<%s>{L: %s, R: %s}", opStr, aL, aR)
}

func (op *OperAssign) Do(cx base.Context) error {
	err1 := op.left.Do(cx)
	if err1 != nil {
		return err1
	}
	err2 := op.right.Do(cx)
	if err2 != nil {
		return err2
	}
	lop := op.left.Get()
	rval := op.right.Get()
	switch target := lop.(type) {
	case *base.Var:
		target.Val = ob.GetVal(rval)
	}
	return nil
}

func (op *OperAssign) Get() any {
	return op.res // make sense for last expression in the Block
}

// func GetOperArgs(oper any) (base.Expression, base.Expression) {
// 	switch opp := oper.(type) {
// 	case *OperAssign, *OperBin:
// 		return opp.left, opp.right
// 	}

// }

type BrType int

const (
	UnknownBr BrType = 100
	RoundBr   BrType = 101
	SquareBr  BrType = 102
	CurlyBr   BrType = 103
)

var brTypeMap = map[string]BrType{
	"(": RoundBr,
	"[": SquareBr,
	"{": CurlyBr,
}

func GetBrType(s string) BrType {
	t, ok := brTypeMap[s]
	if ok {
		return t
	}
	return UnknownBr
}

type Brackets struct {
	Type BrType
	Sub  base.Expression
}

func (br *Brackets) Get() any {
	return br.Sub.Get()
}

func (br *Brackets) Do(ctx base.Context) error {
	return br.Sub.Do(ctx)
}

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

func (op *OperColon) Get() any {
	return op.res
}

func (op *OperColon) Do(ctx base.Context) error {
	return nil
}

// ===========

type OperBin struct {
	left  base.Expression
	right base.Expression
	Oper  *Oper
	res   any
	TODO  bool
}

func (op *OperBin) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperBin) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperBin) Get() any {
	return op.res
}

func (op *OperBin) Do(ctx base.Context) error {
	if op.TODO {
		return nil
	}
	op.left.Do(ctx)
	op.right.Do(ctx)
	lop := op.left.Get()
	lvv := ob.GetVal(lop)
	rop := op.right.Get()
	rvv := ob.GetVal(rop)
	var res any
	var ok bool
	switch val := lvv.(type) {
	case int64:
		res, ok = binOperInt(op.Oper.Id, val, rvv)
	case float64:
		res, ok = binOperFloat(op.Oper.Id, val, rvv)
	case string:
		res, ok = binOperString(op.Oper.Id, val, rvv)
	case *ob.ListVal:
		res, ok = binOperList(op.Oper.Id, val, rvv)
	}
	if !ok {
		return errors.New("Error in bin oper") // TODO: add more informative error
	}
	op.res = res
	return nil
}

type SequenceComma struct {
	Subs []base.Expression
	res  any
	// TODO: res type - of collection, func def args, func call argc, struct ded, struct constr, hmm...
	// []any ?
}

func (cs *SequenceComma) Get() any {
	// TODO
	return cs.res
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
	// TODO: res type - of collection, func def args, func call argc, struct ded, struct constr, hmm...
	// []any ?
}

func (cs *SequenceSemicolon) Get() any {
	// TODO
	return cs.res
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
	return nil
}

func (cs *SequenceSemicolon) Add(elem base.Expression) {
	cs.Subs = append(cs.Subs, elem)
}

func (cs *SequenceSemicolon) SetSubs(elems []base.Expression) {
	cs.Subs = elems
}

// func (op *OperBin) Plus(ctx *ob.Context) (any, bool) {
// 	op.left.Do(ctx)
// 	op.right.Do(ctx)
// 	lop := op.left.Get()
// 	lvv := ob.GetVal(lop)
// 	rop := op.right.Get()
// 	rvv := ob.GetVal(rop)

// 	switch val := lvv.(type) {
// 	case int64:
// 		return plusIntN(val, rvv)
// 	case float64:
// 		return plusFloatN(val, rvv)
// 	case string:
// 		return plusStringN(val, rvv)
// 	case ob.ListVal:
// 		// TODO: ValList
// 		return &ob.ListVal{}, false
// 	}

// 	return nil, false
// }

// func (op *OperBin) DoOper(ctx *ob.Context) (any, bool) {
// 	op.left.Do(ctx)
// 	op.right.Do(ctx)
// 	lop := op.left.Get()
// 	lvv := ob.GetVal(lop)
// 	rop := op.right.Get()
// 	rvv := ob.GetVal(rop)

// 	switch val := lvv.(type) {
// 	case int64:
// 		return plusIntN(val, rvv)
// 	case float64:
// 		return plusFloatN(val, rvv)
// 	case string:
// 		return plusStringN(val, rvv)
// 	case ob.ListVal:
// 		// TODO: ValList
// 		return &ob.ListVal{}, false
// 	}

// 	return nil, false
// }

// func (op *OperBin) Do1(ctx *ob.Context) error {
// 	// var err error
// 	// var res any = nil
// 	var res any
// 	var ok bool
// 	switch op.Oper.Id {
// 	case OpAssign:
// 		res, ok = op.Plus(ctx)
// 	}
// 	if !ok {
// 		return errors.New("Error in bin oper") // TODO: add more informative error
// 	}
// 	op.res = res
// 	return nil
// }
