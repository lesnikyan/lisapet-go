package nodes

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

func plusIntN(a int64, b any) (any, bool) {
	switch b := b.(type) {
	case int64:
		return a + b, true
	case float64:
		return float64(a) + b, true
	}
	return nil, false
}

func plusFloatN(a float64, b any) (any, bool) {
	switch vb := b.(type) {
	case int64:
		return a + float64(vb), true
	case bool:
		// var bb float64
		// if vb {
		// 	bb = 1
		// }
		// floatv(vb)
		return a + floatv(vb), true
	case float64:
		return a + vb, true
	}
	return nil, false
}

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
