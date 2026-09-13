package nodes

func vals2anis[T any](vals []T) []any {
	res := make([]any, len(vals))
	for i, v := range vals {
		res[i] = v
	}
	return res
}

func SplitNamed(args []any) ([]any, map[string]any) {
	var named map[string]any
	ordd := make([]any, len(args))
	k := 0
	for _, arg := range args {
		switch an := arg.(type) {
		case *NamedArgs:
			if named == nil {
				named = an.Nvals
			} else {
				for k, v := range an.Nvals {
					named[k] = v
				}
			}
		// case TripleDots...?
		default:
			ordd[k] = an
			k++
		}
	}
	return ordd[0:k], named
}
