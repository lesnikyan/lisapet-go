package nodes

import ob "github.com/lesnikyan/lisapet-go/objects"

type TNode struct {
}

type VarExpr struct {
	Name string
}

type Expr interface {
	Do(ctx ob.Context)
	Get() ob.TVal
}

type Block struct {
	subs []Expr
	res  ob.TVal
}

func (bb *Block) Do(ctx ob.Context) {
	for _, expr := range bb.subs {
		expr.Do(ctx)
	}
}

func (bb *Block) Get() ob.TVal {
	return bb.res
}
