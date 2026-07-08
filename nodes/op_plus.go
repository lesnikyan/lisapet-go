package nodes

import (
	"fmt"
	"math"
	"strings"

	ob "github.com/lesnikyan/lisapet-go/objects"
)

func floatv(a any) float64 {
	switch av := a.(type) {
	case int64:
		return float64(av)
	case bool:
		var r float64
		if av {
			r = 1.
		}
		return r
	}
	return 0
}

func binOperInt(opid Opid, a int64, b any) (any, bool) {
	fmt.Println("binInt:", opid, a, b)
	switch b := b.(type) {
	case int64:
		switch opid {
		case OpPlus:
			return a + b, true
		case OpMinus:
			return a - b, true
		case OpMult:
			return a * b, true
		case OpDiv:
			return float64(a) / float64(b), true
		case OpPow:
			return int64(math.Pow(float64(a), float64(b))), true
		case OpEqual:
			return a == b, true
		}
	case float64:
		switch opid {
		case OpPlus:
			return float64(a) + b, true
		case OpMinus:
			return float64(a) - b, true
		case OpMult:
			return float64(a) * b, true
		case OpDiv:
			return float64(a) / b, true
		case OpPow:
			return math.Pow(float64(a), b), true
		case OpEqual:
			return float64(a) == b, true
		case OpNot:
			return float64(a) != b, true
		}
	}
	return nil, false
}

func binOperFloat(opid Opid, a float64, b any) (any, bool) {
	switch b := b.(type) {
	case int64:
		switch opid {
		case OpPlus:
			return a + float64(b), true
		case OpMinus:
			return a - float64(b), true
		case OpMult:
			return a * float64(b), true
		case OpDiv:
			return a / float64(b), true
		case OpPow:
			return math.Pow(a, float64(b)), true
		case OpEqual:
			return a == float64(b), true
		case OpNot:
			return a != float64(b), true
		}
	case float64:
		switch opid {
		case OpPlus:
			return a + b, true
		case OpMinus:
			return a - b, true
		case OpMult:
			return a * b, true
		case OpDiv:
			return float64(a) / float64(b), true
		case OpPow:
			return math.Pow(float64(a), float64(b)), true
		case OpEqual:
			return a == b, true
		case OpNot:
			return a != b, true
		}
	}
	return nil, false
}

func binOperString(opid Opid, a string, b any) (any, bool) {
	switch b := b.(type) {
	case string:
		switch opid {
		case OpPlus:
			return a + b, true
		case OpEqual:
			return strings.Compare(a, b) == 0, true
		case OpNot:
			return a != b, true
		}
	}
	return nil, false
}

func binOperBool(opid Opid, a bool, b any) (any, bool) {
	switch b := b.(type) {
	case bool:
		switch opid {
		case OpEqual:
			return a == b, true
		case OpNot:
			return a != b, true
		}

	}
	return nil, false
}

func binOperByte(opid Opid, a byte, b any) (any, bool) { return nil, false }

func binOperGlyf(opid Opid, a rune, b any) (any, bool) { return nil, false }

func binOperList(opid Opid, a *ob.ListVal, b any) (any, bool)    { return nil, false }
func binOperTuple(opid Opid, a *ob.TupleVal, b any) (any, bool)  { return nil, false }
func binOperDict(opid Opid, a *ob.DictVal, b any) (any, bool)    { return nil, false }
func binOperType(opid Opid, a any, b any) (any, bool)            { return nil, false }
func binOperFunc(opid Opid, a ob.Function, b any) (any, bool)    { return nil, false }
func binOperStruct(opid Opid, a ob.StructVal, b any) (any, bool) { return nil, false }
func binOperRegext(opid Opid, a ob.Regexp, b any) (any, bool)    { return nil, false }
func binOperAny(opid Opid, a any, b any) (any, bool)             { return nil, false }

// func plusIntN(a int64, b any) (any, bool) {
// 	switch b := b.(type) {
// 	case int64:
// 		return a + b, true
// 	case float64:
// 		return float64(a) + b, true
// 	}
// 	return nil, false
// }

// func plusFloatN(a float64, b any) (any, bool) {
// 	switch vb := b.(type) {
// 	case int64:
// 		return a + float64(vb), true
// 	case bool:
// 		return a + floatv(vb), true
// 	case float64:
// 		return a + vb, true
// 	}
// 	return nil, false
// }

// func plusStringN(a string, b any) (any, bool) {
// 	switch vb := b.(type) {
// 	case string:
// 		return a + vb, true
// 	}
// 	return nil, false
// }

// func PlusInt(lvv any, rvv any) (any, bool) {

// 	li, lok := lvv.(int64)
// 	if !lok {
// 		return nil, false
// 	}
// 	ri, rok := rvv.(int64)
// 	if !rok {
// 		return nil, false
// 	}
// 	v := li + ri
// 	return v, true
// }
