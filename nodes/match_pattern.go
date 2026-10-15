package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

type MatchPattern interface {
	// should return bool
	// set local vars
	Do(base.Context) error
	Get() *base.Val
	// PutArg(arg any)
	Match(cx base.Context, val any) (bool, error)
}

// ====

type MCaseVal struct {
	Val any
}

func (mc *MCaseVal) Do(base.Context) error {
	return nil
}
func (mc *MCaseVal) Get() *base.Val {
	return nil
}

func (mc *MCaseVal) Match(cx base.Context, arg any) (bool, error) {
	// fmt.Printf("MCaseVal (%v). Match: %T: %v \n", mc.Val, arg, arg)
	if !SameType(mc.Val, arg) {
		return false, nil
	}
	return EqualInternTypes(mc.Val, arg), nil
}

// ====

type MCaseUnder struct {
}

func (mc *MCaseUnder) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseUnder) Get() *base.Val {
	return nil
}

func (mc *MCaseUnder) Match(cx base.Context, arg any) (bool, error) {
	// fmt.Printf("MCaseUnder  Match: %T: %v \n", arg, arg)
	return true, nil
}

// ====

type MCaseType struct {
	TypeExp base.Expression
	Var     *MCaseVar
	TypeVal *base.Type
}

func (mc *MCaseType) Do(cx base.Context) error {
	vt, err := ExprGetType(mc.TypeExp, cx)
	if err != nil {
		return err
	}
	mc.TypeVal = vt
	return nil
}
func (mc *MCaseType) Get() *base.Val {
	return nil
}

func (mc *MCaseType) Match(cx base.Context, arg any) (bool, error) {
	// fmt.Printf("MCaseType  Match: %T: %v :: type: %v \n", arg, arg, mc.TypeExp)
	err := mc.Do(cx)
	if err != nil {
		return false, err
	}
	// fmt.Printf("MCaseType  Match: tval= %T, typeId: %v \n", mc.TypeVal, mc.TypeVal.Id)
	ok, err := CheckTypeEqual(arg, mc.TypeVal)
	// fmt.Printf("MCaseType  Match: ok= %v :: err= %v \n", ok, err)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	// assign var if presents
	if mc.Var != nil {
		ok, err = mc.Var.Match(cx, arg)
		if err != nil {
			return false, err
		}

	}
	return ok, nil
}

// ====

type MCaseVar struct {
	Expr *VarExpr
}

func (mc *MCaseVar) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseVar) Get() *base.Val {
	return nil
}

func (mc *MCaseVar) Match(cx base.Context, arg any) (bool, error) {
	// fmt.Printf("MCaseVar  Match: %T: %v \n", arg, arg)
	// mc.Do(cx)
	mc.Expr.NewVar(cx)
	vr := mc.Expr.GetVar()
	if vr == nil {
		return false, fmt.Errorf("var patter couldn't add var %s", mc.Expr.GetName())
	}
	vr.Val = arg
	// fmt.Printf("MCaseVar  Match: arg= %v :: vr = %v \n", arg, vr)
	return true, nil
}

// ====

// *
type MCaseStar struct {
}

func (mc *MCaseStar) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseStar) Get() *base.Val {
	return nil
}

func (mc *MCaseStar) Match(cx base.Context, arg any) (bool, error) {
	return true, nil
}

// ====

// ?
type MCaseQMark struct {
}

func (mc *MCaseQMark) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseQMark) Get() *base.Val {
	return nil
}

func (mc *MCaseQMark) Match(cx base.Context, arg any) (bool, error) {
	return true, nil
}

// ====

type MCaseList struct {
	Subs []MatchPattern

	minLen int
	maxLen int
}

func (mc *MCaseList) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseList) Get() *base.Val {
	return nil
}

func (mc *MCaseList) Match(cx base.Context, arg any) (bool, error) {
	// min := 0
	// max := 0
	// noMin := false
	// noMax := false
	arglist, ok := arg.(*objects.ListVal)
	if !ok {
		return false, nil
	}
	vals := arglist.Elems
	// vid := 0 // vals index
	subs := mc.Subs
	slen := len(subs)
	// vlen := len(vals)
	// fmt.Printf("MCaseList  Match: %T: %v len: %d \n", arg, arg, slen)
	sid := 0 // subs index
	for _, v := range vals {
		// fmt.Printf(" --- LC: i=%d, v=%v, sid=%d\n", i, v, sid)
		if sid >= slen {
			return false, nil
		}
		sp := subs[sid]
		switch mp := sp.(type) {
		case *MCaseVal:
			ok, err := mp.Match(cx, v)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
			sid++
		case *MCaseVar:
			ok, err := mp.Match(cx, v)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
			sid++
		case *MCaseUnder:
			sid++
		case *MCaseStar:
			// noMax = true
		case *MCaseQMark:

		default:
			// noMin = true
		}
	}
	if sid < slen {
		return false, nil
	}
	// fmt.Printf("MCaseList  Match: arg= %v :: vr = %v \n", arg, vr)
	return true, nil
}

func NewMCaseList(subs []MatchPattern) *MCaseList {
	// min := 0
	// max := 0
	// noMin := false
	// noMax := false
	// for i, sp := range subs {
	// 	switch mp := sp.(type) {
	// 	case *MCaseVal:
	// 	case *MCaseVar:
	// 	case *MCaseUnder:
	// 	case *MCaseStar:
	// 		noMax = true
	// 	case *MCaseQMark:

	// 	default:
	// 		noMin = true
	// 	}
	// }

	return &MCaseList{Subs: subs}
}
