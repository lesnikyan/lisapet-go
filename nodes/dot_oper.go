package nodes

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

type OperDot struct {
	left  base.Expression
	right base.Expression
	res   any // StructField, Method, imported thing
}

func (op *OperDot) SetLeft(xp base.Expression) {
	op.left = xp
}
func (op *OperDot) SetRight(xp base.Expression) {
	op.right = xp
}

func (op *OperDot) Get() *base.Val {
	switch mb := op.res.(type) {
	case *objects.StrMember:
		// fmt.Printf(" <.> Get mb:(%T, %v) \n", mb, mb)
		return mb.Get()
	}
	return nil
}

func (op *OperDot) GetMember() *objects.StrMember {
	switch mb := op.res.(type) {
	case *objects.StrMember:
		// fmt.Printf(" <.> Get mb:(%T, %v) \n", mb, mb)
		return mb
	}
	return nil
}

func (op *OperDot) Set(val *base.Val) error {
	switch mb := op.res.(type) {
	case *objects.StrMember:
		return mb.Set(val)
	}
	return nil
}

func (op *OperDot) Do(cx base.Context) error {
	op.res = nil
	op.left.Do(cx)

	// get object
	left := GetExprVal(op.left, cx)

	// get member name
	mname := ""
	switch nexp := op.right.(type) {
	case *VarExpr:
		// struct field
		mname = nexp.GetName()
	}

	// access to objects member
	switch obj := left.(type) {
	case *objects.StructInst:
		mb := obj.GetMember(mname)
		if mb != nil {
			op.res = mb // field found
			// fmt.Printf(" <.> Do n=`%s` (%T, %v) \n", mname, mb, mb)
			return nil
		} else {
			met, ok := obj.Type.GetMethod(mname)
			fmt.Printf(" <.> Do n=`%s` (%T, %v) \n", mname, met, ok)
			if ok {
				// if mth := obj.GetMethod(mname); mth != nil {
				op.res = met // method found
				return nil
			}

		}

		return fmt.Errorf("OperDot: member `%s` in struct not found", mname)

	default:
		return errors.New("OperDot: unknown type of object " + fmt.Sprintf("%T", obj))
	}
}
