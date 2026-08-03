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

type TTst = struct {
	src   string // code
	vname string // var
	res   any    // exp
}

func Tnull() *obb.Null {
	return &obb.Null{}
}

func RunTCodeVarExp(t *testing.T, i int, tt TTst) {
	t.Run(fmt.Sprintf("Test, %d) %s >>", i, tt.src), func(t2 *testing.T) {
		clines := par.SplitCode(tt.src[1:])
		block, err := TreeBlock(clines)
		assert.Nil(t, err)
		// t.Log("--- --- --- Do ...")
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
