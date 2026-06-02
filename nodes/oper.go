package nodes

import (
	dt "github.com/lesnikyan/lisapet-go/lang/datatype"
	ob "github.com/lesnikyan/lisapet-go/objects"
)

type Opid int

const (
	nn Opid = iota
	OpAssign
	OpPlus
	OpMinus
	OpMult
	OpDiv
	OpPow
	OpRoot
	OpEqual
	OpNotEqual
	OpMoreEqual
	OpLessEqual
	OpPlusAssign
	OpMinusAssign
	OpMultAssign
	OpDivAssign
	OpAnd
	OpOr
	OpNot
	OpBiAnd
	OpBinOr
	OpBinNot
	OpXor
	OpDot
	OpComma
	OpColon
	OpSemicolon
)

type Oper struct {
	sign string
	id   Opid
}

func NewOper(sg string, id Opid) *Oper {
	return &Oper{sign: sg, id: id}
}

type OperBin struct {
	left  Expr
	right Expr
	oper  *Oper
}

func (op *OperBin) Plus(ctx ob.Context) (ob.TVal, bool) {
	op.left.Do(ctx)
	op.right.Do(ctx)
	lop := ob.TVal(op.left.Get())
	lvv := lop.GetVal()
	rop := op.right.Get()
	rvv := rop.GetVal()

	switch lop.GetType() {
	case dt.Int:
		return PlusInt(lvv, rvv)
	}

	return nil, false
}

func PlusInt(lvv any, rvv any) (ob.TVal, bool) {

	li, lok := lvv.(int64)
	if !lok {
		return nil, false
	}
	ri, rok := rvv.(int64)
	if !rok {
		return nil, false
	}
	v := li + ri
	return ob.NewInt(v), true
}

func (op *OperBin) Do(ctx ob.Context) {

}
