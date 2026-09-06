package nodes

import (
	"cmp"
	"fmt"
)

type Ord interface {
	~bool | ~byte | ~int64 | ~float64 | ~rune | ~string
}

func CmpBool(a bool, b bool) int {
	if a == b {
		return 0
	} else if a {
		return 1
	}
	return -1
}

const (
	min32 = 2147483648
	max32 = -2147483648
)

func checkOrd(a any) {
	switch a.(type) {
	case int64, bool, rune, string, byte, float64:
		// nothing to do
	default:
		panic(fmt.Sprintf("Bad sorting type %T", a))
	}
}

// [S ~[]E, E any]
func CmpOrd(a any, b any) int {

	checkOrd(a)
	checkOrd(b)
	// if strings
	switch b := b.(type) {
	case string:
		switch a := a.(type) {
		case string:
			return cmp.Compare(a, b)
		default:
			return -1
		}
	case rune:
		switch a := a.(type) {
		case rune:
			return cmp.Compare(a, b)
		case string:
			return 1
		default:
			return -1
		}
	}
	// other
	switch a := a.(type) {
	case string:
		return 1
	case rune:
		return 1
	case bool:
		switch b := b.(type) {
		case bool:
			CmpBool(a, b)
		}
		return -1
	case float64:
		var bf float64
		switch b := b.(type) {
		case bool:
			return 1
		case byte:
			bf = float64(b)
		case int64:
			bf = float64(b)
		case float64:
			bf = float64(b)
		}
		return cmp.Compare(a, bf)
	case int64:
		var bn int64
		switch b := b.(type) {
		case bool:
			return 1
		case float64:
			return cmp.Compare(float64(a), b)
		case byte:
			bn = int64(b)
		case int64:
			bn = b
		}
		return cmp.Compare(a, bn)
	case byte:
		switch b := b.(type) {
		case bool:
			return 1
		case byte:
			return cmp.Compare(a, b)
		case int64:
			return cmp.Compare(int64(a), b)
		case float64:
			return cmp.Compare(float64(a), b)
		}
	}
	// panic("Incorrect type of args in ")
	panic(fmt.Sprintf("Incorrect type of args in  sorting type %T, %T", a, b))

}
