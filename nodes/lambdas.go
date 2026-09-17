package nodes

import (
	"github.com/lesnikyan/lisapet-go/base"
	ob "github.com/lesnikyan/lisapet-go/objects"
)

// =====
type AppendOper struct {
	Target any // *T of [T ob.ListVal | ob.DitcVal]
	Src    any // right arg
}

// Expression that make and return Lambda-obiect
type LambdaExpr struct {
	Args       []base.Expression
	BlockNodes []base.Expression

	res *ob.Function
}

// ser Args: var, colon, CommaSeq
func (md *LambdaExpr) SetLeft(base.Expression) {
	count := 1
	// if nor single vcar count = len of seq
	md.Args = make([]base.Expression, +count)
}

// set block expression: oper, value, semicolon seq, brackets
func (md *LambdaExpr) SetRight(base.Expression) {
	// need process passed expr
}

func (md *LambdaExpr) Do(base.Context) error {
	// TODO: make lambda (Function)
	return nil
}

func (md *LambdaExpr) Get() *base.Val {
	// TODO: return lambda
	return nil
}
