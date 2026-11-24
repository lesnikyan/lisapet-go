package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

func GetExprTarget(v base.Expression, cx base.Context) any {
	switch vv := v.(type) {
	case *VarExpr:
		vr := GetVar(vv, cx)
		return vr
	case *SequenceComma:
		r := make([]*base.Var, len(vv.Subs))
		for i, vsub := range vv.Subs {
			vrex, ok := vsub.(*VarExpr)
			if !ok {
				panic("part of target comma separated expr is not a variable")
			}
			vr := GetVar(vrex, cx)
			r[i] = vr
		}
		return r
		// TODO: take vars, return []*Var

	case *ColElemExpr:
		colVal := vv.ColRes.Get()
		if colVal == nil {
			return nil
		}
		return colVal.V // *ColElem: obj[k]
	}
	return nil
}

func GetExprVal(v base.Expression, cx base.Context) any {
	// fmt.Printf("GetExprVal#0: %T, %v\n", v, v)
	var eVal *base.Val
	switch vv := v.(type) {
	case *VarExpr:
		// vr := GetVar(vv, cx)
		elem := vv.GetElem()
		// fmt.Printf("GetExprVal#VarEx: %T, %v >> %T, %v \n", vv, vv, elem, elem)
		if elem == nil {
			return nil
		}
		switch el := elem.V.(type) {
		case *base.Var:
			// fmt.Printf("GExVr#1 Var: %T, %v \n", el, el)
			val := el.Val
			return val
		case *base.Type:
			return el
		case *objects.Function, *NFunc:
			return el
		case *objects.Enum, *Grup:
			return el
		default:
			panic(fmt.Sprintf("Not a valid element found by word-lexem: %T", el))
		}
	case *ValExpr:
		eVal = vv.Get()
		// fmt.Printf("GetExprVal#Val: %T, %v\n", eVal.V, eVal.V)
	case *ColElemExpr:
		// fmt.Printf("GetExp.ColElem#Val: %T, %v\n", eVal, eVal)
		eVal = vv.ColRes.Get()
	case *FuncCall:
		eVal := vv.Get()
		if eVal == nil {
			// return with no value
			return NullV()
		}
		// fmt.Printf("GetExp.FuncCall#Val: %T, %v >> %T, %v \n", eVal, eVal, eVal.V, eVal.V)
		return eVal.V
	case *NumSeqExpr:
		res := vv.res.GetList()
		// fmt.Printf("GetExp.NumSeq# len: %v\n", len(res.Elems))
		return res
	case *SequenceComma:
		vals := vv.GetVals()
		if vals == nil {
			return nil
		}
		return vals.V
		// return vals
	case *StructConstr:
		return vv.Get().V
	case *Brackets:
		return vv.Get().V
	case *OperColon:
		return vv.GetPair()
	case *EmptyExpr:
		return &objects.EmptyVal{}
	case *OperDot:
		// fmt.Printf("GetExprVal# operDot: %T, %v\n", vv, vv)
		// fmt.Printf("GetExprVal# operDot.left: %T, %v\n", vv.left, vv.left)
		// fmt.Printf("GetExprVal# `a.b` : %T, %v\n", vv.Get(), vv.Get())
		// fmt.Printf("GetExprVal#  member: %T, %v\n", vv.GetMember(), vv.GetMember())
		xmember := vv.GetMember()
		// fmt.Printf("GetExprVal#  member: %T, %v\n", xmember, xmember)
		switch member := xmember.(type) {
		case *objects.StrMember:
			var mval *base.Val
			if member != nil {
				mval = member.Get()
			}
			if mval == nil {
				// fmt.Printf("GetExprVal# no member value %v\n", mval)
				return nil
			}
			// fmt.Printf("GetExprVal# memVal: %T, %v\n", mval.V, mval.V)
			switch val := mval.V.(type) {
			case *objects.Method:
				val.Inst = member.Obj
				return val

			default:
				return val
			}

		case *GrupMember:
			// fmt.Printf("GetExprVal# Dot Grup mem (%s): %T, %v\n", member.Name, member, member)
			if member == nil {
				panic("grup member: is nil")
			}
			mval := member.Get()
			// fmt.Printf("GetExprVal# Dot Grup mem val: %T, %v\n", mval.V, mval.V)
			if mval == nil {
				panic("grup member val: is nil")
			}
			switch mv := mval.V.(type) {
			case *base.Var:
				return mv.Val
			case *objects.Function:
				return mv
			case *objects.Enum:
				return mv
			case *base.Type:
				return mv
			case *Grup:
				return mv
			default:
				panic(fmt.Sprintf("grup member: incorrect type %T", mval))
			}
		case *objects.EnumMember:
			// fmt.Printf("GetExprVal# Dot Enum mem (%s): %T, %v\n", member.Name, member, member)
			mval := member.Get()
			if mval == nil {
				panic(fmt.Sprintf("enum member: no such element (%s) in enum %s ", member.Name, member.Enum.Name))
			}
			return mval.V
		default:
			// fmt.Printf("GetExprVal# Dot. def: %T:  %v ; %T:  %v\n", vv, vv, member, member)
			var mval *base.Val
			if member == nil {
				// fmt.Printf("GetExprVal# no member %v\n", member)
				mres := vv.Get()
				if mres == nil {
					return nil
				}
				mval = mres
			}
			if mval == nil {
				// fmt.Printf("GetExprVal# no member value %v\n", mval)
				return nil
			}
			// fmt.Printf("GetExprVal# memVal: %T, %v\n", mval.V, mval.V)
			switch val := mval.V.(type) {
			case *MFunc:
				// val.Inst
				return val
				// default:
				// 	return val
			}
		}

	default:
		// fmt.Printf("GetExprVal#100: (%T) %T, %v\n", vv, vv.Get(), vv.Get())
		exv := vv.Get()
		if exv == nil {
			return nil
		}
		// fmt.Printf("GetExprVal#101: %T:  %v\n", exv.V, exv.V)
		return exv.V
	}
	if eVal != nil {
		return eVal.V
	}
	return nil
}

func GrupMemVal(elem any) any {
	return nil
}
