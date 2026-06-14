package nodes

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
	ob "github.com/lesnikyan/lisapet-go/objects"
)

type OperAssign struct {
	oper  string
	left  base.Expression
	right base.Expression
	res   any
}

func (op *OperAssign) Do(cx base.Context) error {
	err1 := op.left.Do(cx)
	if err1 != nil {
		return err1
	}
	err2 := op.right.Do(cx)
	if err2 != nil {
		return err2
	}
	lop := op.left.Get()
	rval := op.right.Get()
	switch target := lop.(type) {
	case *base.Var:
		target.Val = ob.GetVal(rval)
	}
	return nil
}

func (op *OperAssign) Get() any {
	return op.res // make sense for last expression in the Block
}

type OperBin struct {
	left  base.Expression
	right base.Expression
	Oper  *Oper
	res   any
}

func (op *OperBin) Get() any {
	return op.res
}

func (op *OperBin) Do(ctx base.Context) error {

	op.left.Do(ctx)
	op.right.Do(ctx)
	lop := op.left.Get()
	lvv := ob.GetVal(lop)
	rop := op.right.Get()
	rvv := ob.GetVal(rop)
	var res any
	var ok bool
	switch val := lvv.(type) {
	case int64:
		res, ok = binOperInt(op.Oper.Id, val, rvv)
	case float64:
		res, ok = binOperFloat(op.Oper.Id, val, rvv)
	case string:
		res, ok = binOperString(op.Oper.Id, val, rvv)
	case *ob.ListVal:
		res, ok = binOperList(op.Oper.Id, val, rvv)
	}
	if !ok {
		return errors.New("Error in bin oper") // TODO: add more informative error
	}
	op.res = res
	return nil
}

// func (op *OperBin) Plus(ctx *ob.Context) (any, bool) {
// 	op.left.Do(ctx)
// 	op.right.Do(ctx)
// 	lop := op.left.Get()
// 	lvv := ob.GetVal(lop)
// 	rop := op.right.Get()
// 	rvv := ob.GetVal(rop)

// 	switch val := lvv.(type) {
// 	case int64:
// 		return plusIntN(val, rvv)
// 	case float64:
// 		return plusFloatN(val, rvv)
// 	case string:
// 		return plusStringN(val, rvv)
// 	case ob.ListVal:
// 		// TODO: ValList
// 		return &ob.ListVal{}, false
// 	}

// 	return nil, false
// }

// func (op *OperBin) DoOper(ctx *ob.Context) (any, bool) {
// 	op.left.Do(ctx)
// 	op.right.Do(ctx)
// 	lop := op.left.Get()
// 	lvv := ob.GetVal(lop)
// 	rop := op.right.Get()
// 	rvv := ob.GetVal(rop)

// 	switch val := lvv.(type) {
// 	case int64:
// 		return plusIntN(val, rvv)
// 	case float64:
// 		return plusFloatN(val, rvv)
// 	case string:
// 		return plusStringN(val, rvv)
// 	case ob.ListVal:
// 		// TODO: ValList
// 		return &ob.ListVal{}, false
// 	}

// 	return nil, false
// }

// func (op *OperBin) Do1(ctx *ob.Context) error {
// 	// var err error
// 	// var res any = nil
// 	var res any
// 	var ok bool
// 	switch op.Oper.Id {
// 	case OpAssign:
// 		res, ok = op.Plus(ctx)
// 	}
// 	if !ok {
// 		return errors.New("Error in bin oper") // TODO: add more informative error
// 	}
// 	op.res = res
// 	return nil
// }
