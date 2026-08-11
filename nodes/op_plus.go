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

// ## + - * / ** ^/ == != < > <= >= %
// TODO: int: | & ^ << >>
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
		case OpRoot:
			return math.Pow(float64(b), 1/float64(a)), true
		case OpEqual:
			return a == b, true
		case OpNotEqual:
			return a != b, true
		case OpLess:
			return a < b, true
		case OpLessEqual:
			return a <= b, true
		case OpMore:
			return a > b, true
		case OpMoreEqual:
			return a >= b, true
		case OpPercent:
			return a % b, true
		case OpBitAnd:
			return a & b, true
		case OpBitOr:
			return a | b, true
		case OpXor:
			return a ^ b, true
		case OpBitLShift:
			return a << b, true
		case OpBitRShift:
			return a >> b, true
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
		case OpRoot:
			return math.Pow(float64(b), 1/float64(a)), true
		case OpEqual:
			return float64(a) == b, true
		case OpNotEqual:
			return float64(a) != b, true
		case OpLess:
			return float64(a) < b, true
		case OpLessEqual:
			return float64(a) <= b, true
		case OpMore:
			return float64(a) > b, true
		case OpMoreEqual:
			return float64(a) >= b, true
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
		case OpRoot:
			return math.Pow(float64(b), 1/float64(a)), true
		case OpEqual:
			return a == float64(b), true
		case OpNotEqual:
			return a != float64(b), true
		case OpLess:
			return a < float64(b), true
		case OpLessEqual:
			return a <= float64(b), true
		case OpMore:
			return a > float64(b), true
		case OpMoreEqual:
			return a >= float64(b), true
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
		case OpRoot:
			return math.Pow(b, 1/a), true
		case OpEqual:
			return a == b, true
		case OpNotEqual:
			return a != b, true
		case OpLess:
			return a < b, true
		case OpLessEqual:
			return a <= b, true
		case OpMore:
			return a > b, true
		case OpMoreEqual:
			return a >= b, true
		}
	}
	return nil, false
}

func strLShift(format string, vals []any) string {
	// TODO: prepare includes
	return fmt.Sprintf(format, vals...)
}

func binOperString(opid Opid, a string, b any) (any, bool) {
	switch b := b.(type) {
	// operators bitween two strings
	case string:
		switch opid {
		case OpPlus:
			return a + b, true
		case OpEqual:
			return strings.Compare(a, b) == 0, true
		case OpNotEqual:
			return a != b, true
		case OpBitLShift:
			return strLShift(a, []any{b}), true
		}
	default:
		switch opid {
		// formatting
		case OpBitLShift:
			var vals []any
			switch src := b.(type) {
			case *ob.ListVal:
				vals = src.Elems
			default:
				vals = []any{b}
			}
			return strLShift(a, vals), true
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
		case OpNotEqual:
			return a != b, true
		case OpAnd:
			return a && b, true
		case OpOr:
			return a || b, true
		}
	}
	return nil, false
}

func binOperByte(opid Opid, a byte, b any) (any, bool) { return nil, false }

func binOperGlyf(opid Opid, a rune, b any) (any, bool) { return nil, false }

func binOperList(opid Opid, a *ob.ListVal, b any) (any, bool) {

	switch opid {
	case OpPlus:
		// fmt.Printf("List/ <+> a:(%T), b:(%T)  \n", a, b)
		switch bval := b.(type) {
		case *ob.ListVal:
			alen := len(a.Elems)
			vals := make([]any, alen+len(bval.Elems))
			for i, v := range a.Elems {
				vals[i] = v
			}
			for i, v := range bval.Elems {
				vals[alen+i] = v
			}
			res := ob.NewListVal(vals)
			return res, true
		}
	case OpMinus:
		// fmt.Printf("List/ <-> a:(%T), b:(%T)  \n", a, b)
	case OpPlusAssign:
		// fmt.Printf("List/ <+=> a:(%T), b:(%T)  \n", a, b)
		switch bval := b.(type) {
		case *ob.ListVal:
			a.Elems = append(a.Elems, bval.Elems...)
			return a, true
		}
	}
	return nil, false

}
func binOperTuple(opid Opid, a *ob.TupleVal, b any) (any, bool) {
	switch opid {
	case OpPlus:
		// fmt.Printf("Tuple/ <+> a:(%T), b:(%T)  \n", a, b)
		switch bval := b.(type) {
		case *ob.TupleVal:
			alen := len(a.Elems)
			vals := make([]any, alen+len(bval.Elems))
			for i, v := range a.Elems {
				vals[i] = v
			}
			for i, v := range bval.Elems {
				vals[alen+i] = v
			}
			res := ob.NewTupleVal(vals)
			return res, true
		}
	}
	return nil, false
}

func binOperDict(opid Opid, a *ob.DictVal, b any) (any, bool) {
	return nil, false
}

func binOperType(opid Opid, a any, b any) (any, bool)            { return nil, false }
func binOperFunc(opid Opid, a Function, b any) (any, bool)       { return nil, false }
func binOperStruct(opid Opid, a ob.StructVal, b any) (any, bool) { return nil, false }
func binOperRegext(opid Opid, a ob.Regexp, b any) (any, bool)    { return nil, false }
func binOperAny(opid Opid, a any, b any) (any, bool)             { return nil, false }
