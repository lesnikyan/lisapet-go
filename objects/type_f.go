package objects

import (
	"slices"

	"github.com/lesnikyan/lisapet-go/base"
)

type TypeInfo struct {
	Id   base.TypeId
	Name string
}

type null = *Null

func NewTypeInf(name string, id base.TypeId) *TypeInfo {
	return &TypeInfo{Name: name, Id: id}
}

func TypeByVal(val any) *TypeInfo {
	switch val.(type) {
	case int64:
		return NewTypeInf("int", base.TypeInt)
	case byte:
		return NewTypeInf("byte", base.TypeByte)
	case Bytes:
		return NewTypeInf("bytes", base.TypeBytes)
	case float64:
		return NewTypeInf("float", base.TypeFloat)
	case bool:
		return NewTypeInf("bool", base.TypeBool)
	case string:
		return NewTypeInf("string", base.TypeString)
	case *Null:
		return NewTypeInf("null", base.TypeNull)
	case *Maybe:
		return NewTypeInf("maybe", base.TypeMaybe)
	case *ListVal:
		return NewTypeInf("list", base.TypeList)
	case *TupleVal:
		return NewTypeInf("tuple", base.TypeTuple)
	case *DictVal:
		return NewTypeInf("dict", base.TypeDict)
	case Glif:
		return NewTypeInf("glif", base.TypeGlif)
	case base.Mur:
		return NewTypeInf("mur", base.TypeMur)
	case *Regexp:
		return NewTypeInf("regexp", base.TypeRegexp)
	case base.FuncVal:
		return NewTypeInf("function", base.TypeFunction)
		// case *Maybe:
		// 	return NewTypeInf("int", base.TypeMaybe)
	}

	return NewTypeInf("undefined", base.TypeUndefined)
}

func TypeIdByVal(val any) base.TypeId {
	// TODO: if StructInstance - return ID of struct definition
	switch val.(type) {
	case int64:
		return base.TypeInt
	case float64:
		return base.TypeFloat
	case bool:
		return base.TypeBool
	case byte:
		return base.TypeByte
	case string:
		return base.TypeString
	case Bytes:
		return base.TypeBytes
	case *Null:
		return base.TypeNull
	case *ListVal:
		return base.TypeList
	case *TupleVal:
		return base.TypeTuple
	case *DictVal:
		return base.TypeDict
	case base.FuncVal:
		return base.TypeFunction
	case *Maybe:
		return base.TypeMaybe
	case *Regexp:
		return base.TypeRegexp
	case base.Mur:
		return base.TypeMur
	}

	return base.TypeUndefined
}

func bool2int(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

// int = bool, null
// float = bool, null, int
// objects = null
func ConvertByType(exp base.TypeId, val any) any {
	switch exp {

	case base.TypeAny:
		return val

	case base.TypeBool:
		switch val.(type) {
		case null:
			return false
		}

	case base.TypeByte:
		switch vv := val.(type) {
		case null:
			return byte(0)
		case int64:
			return byte(vv)
		case bool:
			return byte(bool2int(vv))
		}

	case base.TypeInt:
		switch vv := val.(type) {
		case int64:
			return val
		case bool:
			return bool2int(vv)
		case null:
			return int64(0)
		case byte:
			return int64(vv)
		}

	case base.TypeFloat:
		switch vv := val.(type) {
		case int64:
			return float64(vv)
		case bool:
			return float64(bool2int(vv))
		case null:
			return float64(0)
		}

	// no hidden conversion for collections, just assign null if null
	case base.TypeList, base.TypeTuple, base.TypeDict, base.TypeFunction:
		switch vv := val.(type) {
		case null:
			return vv
		}

	}
	return nil
}

// if a is parent of b
func IsParent(a *StructDef, b *StructDef) bool {
	return slices.Contains(b.Parents, a)
}

// actual for variable, argument, struct method in left of assign
// prepare and convert val
// result: isTypeOk, convertedVal
func PrepareVal(expType *base.Type, val any) (bool, any) {
	if expType.Id > base.TypeStructBase {
		if null, ok := val.(*Null); ok {
			return true, null
		}
		if st, ok := val.(*StructInst); ok {
			// if struct value
			if st.Def.Type.Id == expType.Id {
				return true, val
			}
			// TODO: if parent
			if expDef, ok := expType.Def.(*StructDef); ok {
				okS := IsParent(expDef, st.Def)
				if okS {
					return true, val
				}
			}
		}

	}
	tv := TypeByVal(val)
	if tv.Id == expType.Id {
		return true, val
	}
	if !base.TypeCompat(expType.Id, tv.Id) {
		// fmt.Printf("PrepV=#2 Vla Not compatible Eq: %v == %v %v\n", tv.Id, expType, tv.Id == expType)
		return false, nil
	}
	cval := ConvertByType(expType.Id, val)
	return true, cval
}
