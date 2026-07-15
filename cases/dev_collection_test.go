package cases

import (
	"fmt"
	"testing"

	obb "github.com/lesnikyan/lisapet-go/objects"
	par "github.com/lesnikyan/lisapet-go/parser"
	"github.com/stretchr/testify/assert"
)

// convert []T to []any
func anynn[T any](vals []T) []any {
	r := make([]any, len(vals))
	for i, n := range vals {
		r[i] = n
	}
	return r
}

// func LVals[T any](vals []*base.Val) []any {
// 	r := make([]any, len(vals))
// 	for i, n := range vals {
// 		r[i] = n.V
// 	}
// 	return r
// }

func TestListsCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# empty list
		nn = []
		`, "nn", anynn([]any{})},
		{`
		# non-empty list
		nn = [1,2,3,4,5]
		`, "nn", anynn([]int64{1, 2, 3, 4, 5})},
		{`
		# write to list
		nn = [1,2,3]
		nn[1] = 222
		`, "nn", anynn([]int64{1, 222, 3})},
		{`
		# write to list in loop
		a = [0,0,0,0,0]
		for i=0; i < 5; i += 1
			a[i] = 10 + i
		`, "a", anynn([]int64{10, 11, 12, 13, 14})},
		{`
		# read from list
		nn = [1,2,3,4,100]
		a = 0
		for i = 0; i < 5; i += 1
			a += nn[i]
		`, "a", int64(110)},
		{`
		# from list to list
		nn = [11,12,13,14,15]
		a = [0,0,0,0,0]
		for i = 0; i < 5; i += 1
			a[i] = nn[i]
		`, "a", anynn([]int64{11, 12, 13, 14, 15})},
		{`
		# list[i] += v
		a = [0,0,0,0,100]
		for i=0; i < 5; i += 1
			a[i] += 10 + i
		`, "a", anynn([]int64{10, 11, 12, 13, 114})},
		{`
		# list[i] += list[i]
		nn = [11,12,13,14,15]
		a = [0,0,0,0,100]
		for i=0; i < 5; i += 1
			a[i] += nn[i]
		`, "a", anynn([]int64{11, 12, 13, 14, 115})},
		// {`` "nn", int64(1)},
		// {``, "nn", int64(1)},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("ExprDo, %s >>", tt.src), func(t2 *testing.T) {
			clines := par.SplitCode(tt.src[1:])
			block, err := TreeBlock(clines)
			assert.Nil(t, err)
			// t.Log("--- --- --- Do ...")
			ctx := obb.NewContext(nil)
			block.Do(ctx)
			vr := ctx.GetVar(tt.vname)
			// t.Log("tt#vr", vr)
			if vr == nil {
				fmt.Printf(" TT#0: %T %v\n", vr, vr)
				return
			}
			fmt.Printf("tt#Var#1  (%T, %v)  (%T, %v) \n", vr, vr, vr.Val, vr.Val)
			switch vobj := vr.Val.(type) {
			case *obb.ListVal:
				fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", vr, vr, vobj, vobj, len(vobj.Elems))
				assert.Equal(t2, tt.res, vobj.Elems)
			case int64:
				fmt.Printf("tt#Var#1  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
				assert.Equal(t2, tt.res, vobj)
			default:
				fmt.Printf("tt#default:  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
			}
			// fmt.Println("tt3>", vr, vr.Name, vr.Val)
		})
	}
}
