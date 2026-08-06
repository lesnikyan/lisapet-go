package objects

import "github.com/lesnikyan/lisapet-go/base"

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
	case float64:
		return NewTypeInf("float", base.TypeFloat)
	case bool:
		return NewTypeInf("bool", base.TypeBool)
	case string:
		return NewTypeInf("string", base.TypeString)
	case *Null:
		return NewTypeInf("null", base.TypeNull)
	case *ListVal:
		return NewTypeInf("list", base.TypeList)
	case *TupleVal:
		return NewTypeInf("tuple", base.TypeTuple)
	case *DictVal:
		return NewTypeInf("dict", base.TypeDict)
	case base.FuncVal:
		return NewTypeInf("function", base.TypeFunction)
		// case *Maybe:
		// 	return NewTypeInf("int", base.TypeInt)
	}

	return NewTypeInf("undefined", base.TypeUndefined)
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

	case base.TypeInt:
		switch vv := val.(type) {
		case int64:
			return val
		case bool:
			return bool2int(vv)
		case null:
			return int64(0)
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
