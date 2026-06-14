package nodes

import "github.com/lesnikyan/lisapet-go/base"

// import ob "github.com/lesnikyan/lisapet-go/objects"

type IfExpr struct {
	res any
}

func (xp *IfExpr) Do(base.Context) error {
	return nil
}
func (xp *IfExpr) Get() any {
	return xp.res
}

type TODOExpr struct {
	res any
}

func (xp *TODOExpr) Do(base.Context) error {
	return nil
}
func (xp *TODOExpr) Get() any {
	return xp.res
}

// type Expression interface {
// 	Do(*ob.Context) error
// 	Get() any
// }
// type Block interface {
// 	Do(cx *ob.Context) error
// 	Get() any
// }

// // super-expression, expression that can have sub-expression
// type SupExpr interface {
// 	Add(sub Expression)
// }

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
