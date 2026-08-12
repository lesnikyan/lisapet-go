package cases

import (
	"fmt"
	"testing"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/nodes"
	obb "github.com/lesnikyan/lisapet-go/objects"
	par "github.com/lesnikyan/lisapet-go/parser"
	"github.com/stretchr/testify/assert"
)

func PreloadContext(cx base.Context) {
	nodes.PreloadFuncs(cx)
	nodes.PreloadTypes(cx)

}

// convert []T to []any
func Anynn[T any](vals []T) []any {
	r := make([]any, len(vals))
	for i, n := range vals {
		var x any = n
		switch v := x.(type) {
		case int:
			x = int64(v)
		case []int64:
			x = Anynn(v)
		}
		r[i] = x
	}
	return r
}

func Anis(vals ...any) []any {
	return Anynn(vals)
}

// for tuples

func Tanis(vals ...any) *Tup {
	return &Tup{Anynn(vals)}
}

type TTst = struct {
	src   string // code
	vname string // var
	res   any    // exp
}

type dk = map[any]any

func tval(val any) any {
	var ak any = val
	switch tk := val.(type) {
	case int:
		ak = int64(tk)
	}
	return ak
}

func adk(src dk) dk {
	res := make(dk)
	for k, v := range src {
		var ak any = tval(k)
		var av any = tval(v)
		res[ak] = av
	}
	return res
}

// prepare test result
func pres(src any) any {
	switch val := src.(type) {
	case int64, string, bool, float64:
		return val
	case *obb.ListVal:
		r := make([]any, len(val.Elems))
		for i, vv := range val.Elems {
			r[i] = pres(vv)
		}
		return r
	case *obb.TupleVal:
		r := make([]any, len(val.Elems))
		for i, vv := range val.Elems {
			r[i] = pres(vv)
		}
		return r
	case *obb.DictVal:
		r := make(map[any]any)
		for k, v := range val.Vmap {
			r[pres(k)] = pres(v)
		}
		return r
	}
	return src
}

func Tnull() *obb.Null {
	return &obb.Null{}
}

func TrunkString(s string) string {
	sRns := []rune(s)
	n := 100
	if len(sRns) < 100 {
		n = len(sRns)
	}
	cr := sRns[:n]
	return string(cr)
}

func RunTCodeVarExp(t *testing.T, i int, tt TTst) {
	t.Run(fmt.Sprintf("Test, %d) %s >>", i, TrunkString(tt.src)), func(t2 *testing.T) {
		clines := par.SplitCode(tt.src[1:])
		block, err := TreeBlock(clines)
		assert.Nil(t, err)
		t.Log("--- --- --- Do ...")
		ctx := obb.NewContext(nil)
		PreloadContext(ctx)
		terr := block.Do(ctx)
		if terr != nil {
			assert.Fail(t2, terr.Error())
			return
		}
		var val any
		vr := ctx.GetVar(tt.vname)
		// t.Log("tt#vr", vr)
		if vr == nil {
			vel := ctx.GetElem(tt.vname)
			if vel == nil {
				assert.Fail(t2, "No expected Var or elem")
				return
			}
			val = vel.V
			fmt.Printf(" TT#0: %T %v\n", val, val)
		} else {
			val = vr.Val
		}
		fmt.Printf("tt#Var#1  vr(%T, %v)  val(%T, %v) \n", vr, vr, val, val)
		switch vobj := val.(type) {
		case *nodes.Function:
			// fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
			assert.Equal(t2, tt.res, vobj.GetName())
		case *obb.ListVal:
			fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
			assert.Equal(t2, tt.res, vobj.Elems)
		case *obb.TupleVal:
			fmt.Printf("tt#TupleVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
			tup, ok := tt.res.(*Tup)
			if !ok {
				assert.Fail(t2, fmt.Sprintf("Tuple has gotten but test exp: : %T", tt.res))
			}
			assert.Equal(t2, tup.elems, vobj.Elems)
		case *obb.DictVal:
			fmt.Printf("tt#DictVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Vmap))
			tm, tok := tt.res.(map[any]any)
			assert.True(t2, tok)
			res := pres(vobj)
			assert.Equal(t2, tm, res)
		case int64:
			fmt.Printf("tt#int  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
			assert.Equal(t2, tt.res, vobj)
		case float64:
			fmt.Printf("tt#float  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
			assert.Equal(t2, tt.res, vobj)
		case bool:
			fmt.Printf("tt#bool  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
			assert.Equal(t2, tt.res, vobj)
		case string:
			fmt.Printf("tt#string  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
			assert.Equal(t2, tt.res, vobj)
		case *obb.Null:
			fmt.Printf("tt#Null:  (%T, %v)  <Null> result \n", vr, vr)
			assert.Equal(t2, tt.res, vobj)
		default:
			fmt.Printf("tt#default:  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
			assert.Fail(t2, "unknown test result")
		}
		// fmt.Println("tt3>", vr, vr.Name, vr.Val)
	})
}
