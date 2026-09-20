package nodes

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
)

func servDefined(cx base.Context, args []any) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("func @defined: incorrect count of args")
	}
	vx, ok := args[0].(*VarExpr)
	if !ok {
		return nil, errors.New("func @defined: 1-st arg must be a string")
	}
	elem := cx.GetElem(vx.GetName())
	if elem == nil {
		// not found
		return false, nil
	}
	return true, nil
}
