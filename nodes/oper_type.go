package nodes

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"github.com/lesnikyan/lisapet-go/base"
	ob "github.com/lesnikyan/lisapet-go/objects"
	obb "github.com/lesnikyan/lisapet-go/objects"
)

func SameType(a any, b any) bool {
	if obb.TypeIdByVal(a) == obb.TypeIdByVal(b) {
		// TODO: think about Undefined
		return true
	}
	return false
}

func EqualInternTypes(a any, b any) bool {
	switch av := a.(type) {
	case *ob.ListVal:
		bv, ok := b.(*ob.ListVal)
		if !ok {
			return false
		}
		return slices.Equal(av.Elems, bv.Elems)
	case *ob.TupleVal:
		bv, ok := b.(*ob.TupleVal)
		if !ok {
			return false
		}
		return slices.Equal(av.Elems, bv.Elems)
	case *ob.DictVal:
		bv, ok := b.(*ob.DictVal)
		if !ok {
			return false
		}
		return maps.Equal(av.Vmap, bv.Vmap)
	}
	// other types
	return a == b
}

// TODO: implement ==, != for all types.
func EqCompare(a any, b any, oper Opid) bool {
	notSamet := !SameType(a, b)
	switch oper {
	case OpEqual:
		if notSamet {
			return false
		} else {
			// TODO: struct, enum, grup
			return EqualInternTypes(a, b)
		}
	case OpNotEqual:
		if notSamet {
			return true
		} else {
			// TODO: struct, enum, grup
			return !EqualInternTypes(a, b)
		}
	}
	return false
}

type MixedTypeExpr struct {
	Subs []base.Expression
	res  *base.MixedType
}

func (mx *MixedTypeExpr) Do(cx base.Context) error {
	mx.res = nil
	if len(mx.Subs) == 0 {
		return errors.New("mixed type expr: empty sub list")
	}
	res := make([]*base.Type, len(mx.Subs))
	for i, sx := range mx.Subs {
		switch vrx := sx.(type) {
		case *VarExpr:
			tname := vrx.GetName()
			tp := cx.GetType(tname)
			if tp == nil {
				return fmt.Errorf("mixed type: type `%s` not found", tname)
			}
			res[i] = tp
		}
	}
	mx.res = &base.MixedType{Types: res}
	return nil
}

func (mx *MixedTypeExpr) Get() *base.Val {
	if mx.res == nil {
		return nil
	}
	return base.NewVal(mx.res)
}

// ====

func ExprFromMixed(expr *OperBin) []base.Expression {
	res := []base.Expression{}
	for expr.left != nil {
		res = append(res, expr.right)
		switch bb := expr.left.(type) {
		case *OperBin:
			// fmt.Printf(" #1 oBin-- %T \n", bb)
			expr = bb
		case *VarExpr, *ValExpr:
			// fmt.Printf(" #2 varX-- %T \n", bb)
			res = append(res, bb)
			return res
		}
	}
	return nil
}

func MixedSubs(expr *OperBin, cx base.Context) (*base.MixedType, error) {
	xx := ExprFromMixed(expr)
	res := make([]*base.Type, len(xx))
	for i, ex := range xx {
		// fmt.Printf(" sub:-- %d) %T \n", i, ex)
		switch tupx := ex.(type) {
		case *VarExpr:
			err := tupx.Do(cx)
			if err != nil {
				// fmt.Println("sub type: error", err)
				return nil, err
			}
			ntype := cx.GetType(tupx.name)
			if ntype == nil {
				return nil, fmt.Errorf("sub type: can't find type: ")
			}
			res[i] = ntype

		case *ValExpr:
			// fmt.Printf(" -- # -- %T, %v\n", rvar, rvar.Val)
			v := tupx.Get()
			if v == nil {
				return nil, fmt.Errorf("sub type: bad type val: %T ", tupx)
			}
			switch v.V.(type) {
			case *obb.Null:
				ntype := cx.GetType("null")
				res[i] = ntype
			default:
				return nil, errors.New("sub type: incorrect value expr instead of type")
			}
		}
	}
	mt := &base.MixedType{Types: res}
	return mt, nil
}
