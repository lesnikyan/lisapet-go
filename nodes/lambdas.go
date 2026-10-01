package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
)

// Expression that make and return Lambda-obiect

type LambdaExp struct {
	Args []base.Expression
	Body base.Expression

	FDef *FuncDef
	res  any
}

func (mb *LambdaExp) SetLeft(args base.Expression) {
	// fmt.Printf("lambda setLeft: incorrect args type: %T \n", args)
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
	case *BackSlash:
		// \ x, y, -> expr
		exSub := sub.Right
		mb.SetLeft(exSub)
		return
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

// ====

type DollarOper struct {
	Left  base.Expression
	Right base.Expression

	fCall *FuncCall

	res any
}

func (op *DollarOper) SetLeft(xp base.Expression) {
	op.Left = xp
	op.fCall.Src = op.Left
}
func (op *DollarOper) SetRight(xp base.Expression) {
	op.Right = xp
	op.fCall.args = []base.Expression{op.Right}
}

func (op *DollarOper) Get() *base.Val {
	return op.fCall.Get()
}

func (op *DollarOper) Do(cx base.Context) error {
	err := op.fCall.Do(cx)
	if err != nil {
		return err
	}
	// op.res = op.fCall.Get()
	return nil
}

func EmptyDollarOper() *DollarOper {
	fCall := &FuncCall{}
	return &DollarOper{fCall: fCall}
}
