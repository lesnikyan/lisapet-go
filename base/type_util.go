package base

import "slices"

func TypeEqual(a *Type, b *Type) bool {
	return a.Id == b.Id
}

// type compatibilily: a = b
func TypeCompat(a TypeId, b TypeId) bool {
	if a == b {
		return true
	}
	switch a {
	case TypeAny:
		return true
	case TypeBool:
		return b == TypeNull
	case TypeInt:
		return slices.Contains([]TypeId{TypeBool, TypeNull}, b)
	case TypeFloat:
		return slices.Contains([]TypeId{TypeInt, TypeBool, TypeNull}, b)
	case TypeString:
		return b == TypeNull
	case TypeList, TypeTuple, TypeDict, TypeFunction:
		return b == TypeNull
	case TypeEnum, TypeGrup:
		return b == TypeNull
	case TypeGlif:
		return b == TypeInt
		// case TypeBytes:
		// 	return slices.Contains([]TypeId{TypeBool, TypeNull}, b)
		// case TypeMaybe:
		// 	return slices.Contains([]TypeId{TypeBool, TypeNull}, b)
	}
	return false
}
