package nodes

import (
	"maps"
	"slices"

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
