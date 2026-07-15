package cases

import (
	"fmt"
	"testing"

	"github.com/lesnikyan/lisapet-go/base"
	obb "github.com/lesnikyan/lisapet-go/objects"
	par "github.com/lesnikyan/lisapet-go/parser"
	"github.com/stretchr/testify/assert"
)

func anys[T any](vals []T) []any {
	r := make([]any, len(vals))
	for i, n := range vals {
		r[i] = n
	}
	return r
}

func LVals[T any](vals []*base.Val) []any {
	r := make([]any, len(vals))
	for i, n := range vals {
		r[i] = n.V
	}
	return r
}

func TestListsCase(t *testing.T) {
	tdata := []struct {
		src string
		res any
	}{
		{`
		nn = [1,2,3,4,5]
		`, anys([]int64{1, 2, 3, 4, 5})},
		// {``, int64(1)},
		// {``, int64(1)},
		// {``, int64(1)},
		// {``, int64(1)},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("ExprDo, %s >>", tt.src), func(t2 *testing.T) {
			clines := par.SplitCode(tt.src[1:])
			block, err := TreeBlock(clines)
			assert.Nil(t, err)
			// t.Log("--- --- --- Do ...")
			ctx := obb.NewContext(nil)
			block.Do(ctx)
			vr := ctx.GetVar("nn")
			// t.Log("tt#vr", vr)
			if vr == nil {
				fmt.Printf(" TT#0: %T %v\n", vr, vr)
				return
			}
			switch col := vr.Val.(type) {
			case *obb.ListVal:
				fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", vr, vr, col, col, len(col.Elems))
				assert.Equal(t2, tt.res, col.Elems)
			default:
				fmt.Printf("tt#default:  (%T, %v)  (%T, %v) \n", vr, vr, col, col)
			}
			// fmt.Println("tt3>", vr, vr.Name, vr.Val)
		})
	}
}
