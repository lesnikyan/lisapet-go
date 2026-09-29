package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

// Expression that make and return Lambda-obiect
// type LambdaExpr struct {
// 	Args       []base.Expression
// 	BlockNodes []base.Expression

// 	res *ob.Function
// }

// // ser Args: var, colon, CommaSeq
// func (md *LambdaExpr) SetLeft(base.Expression) {
// 	count := 1
// 	// if nor single vcar count = len of seq
// 	md.Args = make([]base.Expression, +count)
// }

// // set block expression: oper, value, semicolon seq, brackets
// func (md *LambdaExpr) SetRight(base.Expression) {
// 	// need process passed expr
// }

// func (md *LambdaExpr) Do(base.Context) error {
// 	// TODO: make lambda (Function)
// 	return nil
// }

// func (md *LambdaExpr) Get() *base.Val {
// 	// TODO: return lambda
// 	return nil
// }

type LambdaExp struct {
	Args []base.Expression
	Body base.Expression

	FDef *FuncDef
	res  any
}

func (mb *LambdaExp) SetLeft(args base.Expression) {
	var arxp []base.Expression
	switch sub := args.(type) {
	case *VarExpr:
		arxp = []base.Expression{sub}
	case *SequenceComma:
		arxp = sub.Subs
	case *Brackets:
		brSub := sub.Sub
		mb.SetLeft(brSub)
		return
	case *TupleExpr:
		arxp = sub.Seq.Subs
	// case BackSlash:
	// \ x, y, -> expr
	default:
		panic(fmt.Sprintf("lambda setLeft: incorrect args type: %T", sub))
	}
	mb.SetArgs(arxp)
}

func (mb *LambdaExp) SetArgs(args []base.Expression) {
	mb.Args = args
}

func (mb *LambdaExp) SetRight(xp base.Expression) {
	mb.Body = xp
}

var _lambdaId = 0

func LambdaId() int {
	id := _lambdaId
	_lambdaId++
	return id
}

func (mb *LambdaExp) MakeDef() {
	if mb.FDef == nil {
		if mb.Args != nil {
			fdef := NewFuncDef(fmt.Sprintf("lambda#%000000d", LambdaId()), mb.Args)
			fdef.Add(mb.Body)
			mb.FDef = fdef
		}
	}
}

func (mb *LambdaExp) Get() *base.Val {
	return mb.FDef.Get()
}

func (mb *LambdaExp) Do(cx base.Context) error {
	// fmt.Printf("Lambla.Do\n")
	mb.MakeDef()
	err := mb.FDef.Do(cx)
	if err != nil {
		return err
	}
	return nil
}
