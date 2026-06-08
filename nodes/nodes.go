package nodes

import ob "github.com/lesnikyan/lisapet-go/objects"

type Expression interface {
	Do(*ob.Context)
	Get() any
}

// super-expression, expression that can have sub-expression
type SupExpr interface {
	Add(sub Expression)
}

// type TNode struct {
// }

// type VarExpr struct {
// 	Name string
// }

// type Expr interface {
// 	Do(ctx ob.Context)
// 	Get() ob.Val
// }

// type Block struct {
// 	subs []Expr
// 	res  ob.TVal
// }

// func (bb *Block) Do(ctx ob.Context) {
// 	for _, expr := range bb.subs {
// 		expr.Do(ctx)
// 	}
// }

// func (bb *Block) Get() ob.TVal {
// 	return bb.res
// }
