package nodes

import "github.com/lesnikyan/lisapet-go/base"

// import ob "github.com/lesnikyan/lisapet-go/objects"

type TODOExpr struct {
	res any
}

func (xp *TODOExpr) Do(base.Context) error {
	return nil
}
func (xp *TODOExpr) Get() *base.Val {
	return base.NewVal(xp.res)
}
