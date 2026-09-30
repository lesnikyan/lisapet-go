package nodes

import (
	"fmt"
	"slices"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
)

func compose(cx base.Context, args []any) (*Composed, error) {
	ffs := make([]base.FuncVal, len(args))
	mId := len(args) - 1
	cname := make([]string, len(args))

	for i, arg := range slices.Backward(args) {
		fun, ok := arg.(base.FuncVal)
		if !ok {
			return nil, fmt.Errorf("func compose: not a function arg: %T", arg)
		}
		// fmt.Printf("compose: %d, %s\n", i, fun.GetName())
		ffs[mId-i] = fun
		cname[i] = fun.GetName()
	}
	res := NewComposed(strings.Join(cname, "*"), ffs)
	return res, nil
}

func funcsCompose(cx base.Context, args []any) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("func compose: waits at least 2 args, had: %d ", len(args))
	}
	res, err := compose(cx, args)
	return res, err
}
