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
	case *MFunc:
		return base.NewVal(op.res)
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

	// fmt.Printf(" <.> Do n=`%s` obj(%T, %v) \n", mname, left, left)
	if mname == "" {
		return errors.New("OperDot: method with empty name ")
	}
	// access to objects member
	switch obj := left.(type) {
	case *objects.StructInst:
		mb := obj.GetMember(mname)
		if mb != nil {
			op.res = mb // field found
			// fmt.Printf(" <. oper> StInst.field n=`%s` (%T, %v) \n", mname, mb, mb)
			return nil
		} else {
			met := obj.GetMethod(mname)
			// fmt.Printf(" <. oper> StInst.method n=`%s` (%T, %v) \n", mname, met)
			// if mth := obj.GetMethod(mname); mth != nil {
			if met != nil {
				met.Inst = obj
				op.res = met // method found
				return nil
			}
		}
		// fmt.Printf("OperDot: member `%s` in struct not found \n", mname)
		return fmt.Errorf("OperDot: member `%s` in struct not found", mname)

	default:
		tInf := objects.TypeByVal(obj)
		if tInf.Id == base.TypeUndefined {
			return fmt.Errorf("Undefined type of instance in method call: %T", obj)
		}
		tp := cx.GetType(tInf.Name)

		// fmt.Printf("OperDot3: instType `%s` : %T \n", mname, tp)
		if tp == nil {
			return fmt.Errorf("Undefined type `%T` in method call", obj)
		}
		fn, ok := tp.GetMethod(mname)
		if !ok {
			return fmt.Errorf("method call: method `%s` of type `%s` not found", mname, tp.Name)
		}
		// fmt.Printf("OperDot5: mt: `%s` : %T \n", mname, fn)
		switch mt := fn.(type) {
		case *MFunc:
			mt.SetInst(obj)
			// return &objects.StrMember{Obj: st, Method: mt}
			op.res = mt
			return nil
		default:
			return fmt.Errorf("OperDot: unknown type of method `%T`", obj)
		}
		// fmt.Printf("OperDot: unknown type of object %T \n", obj)
		// return errors.New("OperDot: unknown type of object " + fmt.Sprintf("%T", obj))
	}
}
