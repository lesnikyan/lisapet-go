package nodes

/*
Unary operator in right side:
	arg...
	seq...
	func~>
*/
import "github.com/lesnikyan/lisapet-go/base"

type TripleDots struct {
	left base.Expression
	Oper *Oper
	res  any
}

func (op *TripleDots) SetLeft(xp base.Expression) {
	op.left = xp
}

func (op *TripleDots) SetRight(xp base.Expression) {}

func (op *TripleDots) Get() *base.Val {
	return base.NewVal(op.res)
}

func (op *TripleDots) Do(cx base.Context) error {
	return nil
}
