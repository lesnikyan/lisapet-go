package nodes

import (
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

type OperAssign struct {
	oper  string
	left  Expression
	right Expression
	res   any
}

func (op *OperAssign) Do(cx *ob.Context) {
	op.left.Do(cx)
	op.right.Do(cx)
	lop := op.left.Get()
	rval := op.right.Get()
	switch target := lop.(type) {
	case *ob.Var:
		target.Val = ob.GetVal(rval)
	}
}

func (op *OperAssign) Get() any {
	return op.res // make sense for last expression in the Block
}

type OperBin struct {
	left  Expression
	right Expression
	oper  *Oper
}

func (op *OperBin) Plus(ctx *ob.Context) (any, bool) {
	op.left.Do(ctx)
	op.right.Do(ctx)
	lop := op.left.Get()
	lvv := ob.GetVal(lop)
	rop := op.right.Get()
	rvv := ob.GetVal(rop)

	switch val := lvv.(type) {
	case int64:
		return plusIntN(val, rvv)
	case float64:
		return plusFloatN(val, rvv)
	case string:
		// TODO: string
		return "", false
	case ob.ListVal:
		// TODO: ValList
		return &ob.ListVal{}, false
	}

	return nil, false
}

func (op *OperBin) Do(ctx ob.Context) {

}
