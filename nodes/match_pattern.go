package nodes

import (
	"github.com/lesnikyan/lisapet-go/base"
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

func (mc *MCaseUnder) Do(base.Context) error {
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
	ok, err := CheckTypeEqual(arg, mc.TypeVal)
	// fmt.Printf("MCaseType  Match: ok= %v :: err= %v \n", ok, err)
	if err != nil {
		return false, err
	}
	return ok, nil
}
