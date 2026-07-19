package cases

import (
	"fmt"
	"testing"

	obb "github.com/lesnikyan/lisapet-go/objects"
	par "github.com/lesnikyan/lisapet-go/parser"
	"github.com/stretchr/testify/assert"
)

/*
TODO:
ok 3. for i=0; i<5; i +=1
ok 1. create list: nn = [1,2,3]
ok 2.1 read list v = nn[0]
ok 2.2 set elem: nn[0] = 1
ok 5. create tuple: (1,2,3)
ok 5.1 tuple elem t[i]
ok 6. create dict: d1 = {'a':1, 'b': 2}
ok 6.1 dict elem: d1[k]
ok 4. LeftArrow:
ok 4.1 loop by list: for n <- nn
ok 4.2 append to list: nn <- v
6.2 loop by dict for k, v <- d1
6.3 append to dict: d1 <- (k, v)
*/

func TestForArrowAppendCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = []
		r <- 115
		`, "r", anynn([]any{115})},
		{`
		nn = [11,22,33,44,50,66,77]
		r = [8]
		for n <- nn
			r <- n
		`, "r", anynn([]int64{8, 11, 22, 33, 44, 50, 66, 77})},
		{`
		ii = [0,1,2,3,4,6]
		nn = [11,22,33,44,50,66,77]
		r = [8]
		for i <- ii
			r <- nn[i]
		`, "r", anynn([]int64{8, 11, 22, 33, 44, 50, 77})},
		// {``, "a", int64(1)},
		// {``, "a",  int64(1)},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("ArrowIter, %s >>", tt.src), func(t2 *testing.T) {
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
			val := vr.Val
			fmt.Printf("tt#Var#1  (%T, %v)  (%T, %v) \n", vr, vr, vr.Val, vr.Val)
			switch vobj := val.(type) {
			case *obb.ListVal:
				fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
				assert.Equal(t2, tt.res, vobj.Elems)
			case *obb.DictVal:
				fmt.Printf("tt#DictVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Vmap))
				tm, tok := tt.res.(map[any]any)
				assert.True(t2, tok)
				res := pres(vobj)
				assert.Equal(t2, tm, res)
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
func TestForArrowIterCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		a = 0
		nn = [1,2,3,4,5]
		for n <- nn
			a += n
		`, "a", int64(15)},
		{`
		ii = [0,1,2,3,4]
		nn = [11,22,33,44,50,66,77]
		r = [0,0,0,0,0]
		for i <- ii
			r[i] = nn[i]
		`, "r", anynn([]int64{11, 22, 33, 44, 50})},
		// {``, "a",  int64(1)},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("ArrowIter, %s >>", tt.src), func(t2 *testing.T) {
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
			val := vr.Val
			fmt.Printf("tt#Var#1  (%T, %v)  (%T, %v) \n", vr, vr, vr.Val, vr.Val)
			switch vobj := val.(type) {
			case *obb.ListVal:
				fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
				assert.Equal(t2, tt.res, vobj.Elems)
			case *obb.DictVal:
				fmt.Printf("tt#DictVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Vmap))
				tm, tok := tt.res.(map[any]any)
				assert.True(t2, tok)
				res := pres(vobj)
				assert.Equal(t2, tm, res)
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

func TestForCountCase(t *testing.T) {
	tdata := []struct {
		src string
		res any
	}{
		{`
		a = 0
		for i = 1; i < 6; i = i + 1
			a = a + i
		`, int64(15)},
		{`
		a = 0
		for i=0; i < 10; i += 1
			a += i
		`, int64(45)},
		{`
		a = 0
		for i=0; i < 10; i+= 1
			if i % 2 != 0
				a += i
		`, int64(25)},
		{`
		a = 0
		for i=1; i < 10; i += 1
			for j=1; j < 10; j += 1
				for k=1; k < 10; k += 1
					a += i + j + k
		`, int64(10935)},
		{`
		a = 0
		for i = 0; i < 10 ; i += 1
			for j=1; j < 10; j += 1
				if i % 2 > 0
					if j % 2 >0
						a += i + j
		`, int64(250)},
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
			vr := ctx.GetVar("a")
			// t.Log("tt#vr", vr)

			switch rr := vr.Val.(type) {
			case *obb.ListVal:
				// fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", vr, vr, col, col, len(col.Elems))
				assert.Equal(t2, tt.res, rr.Elems)
			default:
				assert.Equal(t2, tt.res, vr.Val)
				// fmt.Println("tt3>", vr, vr.Name, vr.Val)
			}
		})
	}

	// a := 0
	// for i := 1; i < 10; i++ {
	// 	for j := 1; j < 10; j++ {
	// 		// for k := 1; k < 10; k++ {
	// 		// a += i + j + k
	// 		if i%2 > 0 && j%2 > 0 {
	// 			a += i + j
	// 		}
	// 	}
	// }
	// fmt.Println("i+j+k", a)
}
