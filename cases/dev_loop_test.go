package cases

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/lesnikyan/lisapet-go/base"
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
ok 6.2 loop by dict for k, v <- d1
ok 6.3 append to dict: d1 <- (k, v)
ok 7.1 break
ok 7.2 continue
ok 8. num list: [1..5], [1, 3 .. 7]
*/

func TestInterfaceNil(t *testing.T) {
	var ee base.Expression
	log.Printf("IfcNil#1 (%T, %v) ==%v !=%v \n", ee, ee, ee == nil, ee != nil)
	ee = nil
	log.Printf("IfcNil#1 (%T, %v) ==%v !=%v \n", ee, ee, ee == nil, ee != nil)
}

func TestWhileCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# while
		r = []
		a = 1
		while a < 5
			n = a + 10
			a += 1
			r <- n
		`, "r", anis(11, 12, 13, 14)},
		{`
		# while, more complex
		r = []
		x = 9
		ss = [1..20]
		i = 0
		L = 0
		while L < 6
			if i == 20
				break
			v = ss[i]
			if v % 2 > 0
				if v % 3 > 0
					if v % 5 > 0
						r <- v
			i += 1
		`, "r", anis(1, 7, 11, 13, 17, 19)},
		// {``, "a",  int64(1)},
	}
	for i, tt := range tdata {
		t.Run(fmt.Sprintf("ArrowIter, %d) %s >>", i, tt.src), func(t2 *testing.T) {
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
				assert.Fail(t2, "No expected Var")
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

func TestNumberSequenceCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# iter by var = [ .. ]
		r = [-1]
		ss = [1 .. 5]
		for n <- ss
			r <- n
		`, "r", anynn([]int64{-1, 1, 2, 3, 4, 5})},
		{`
		# iter by [ .. ]
		r = [-1]
		for n <- [2 .. 7]
			r <- n
		`, "r", anynn([]int64{-1, 2, 3, 4, 5, 6, 7})},
		{`
		# num-gen, custom step to var
		a = [1, 3 .. 7]
		`, "a", anis(1, 3, 5, 7)},
		{`
		# iter by [ .. ]
		r = [-2]
		for n <- [1, 3 .. 7]
			r <- n
		`, "r", anis(-2, 1, 3, 5, 7)},
		{`
		# num gen [ .. ] custom step, max is over
		r = [-3]
		for n <- [1, 4 .. 11]
			r <- n
		`, "r", anis(-3, 1, 4, 7, 10)},
		{`
		# num gen from negative - up
		r = [-50]
		for n <- [5, 3 .. -3]
			r <- n
		`, "r", anis(-50, 5, 3, 1, -1, -3)},
		{`
		# num gen from negative - down
		r = [-51]
		for n <- [-1, 4, .. 15]
			r <- n
		`, "r", anis(-51, -1, 4, 9, 14)},
		{`
		# num gen negative step
		r = [-52]
		for n <- [-1, -4, .. -15]
			r <- n
		`, "r", anis(-52, -1, -4, -7, -10, -13)},
		// {``, "a",  int64(1)},
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
				assert.Fail(t2, "No expected Var")
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
		{`
		# dict append with tuple
		dd = {}
		dd <- ('a', 11)
		`, "dd", adk(dk{"a": 11})},
		{`
		nn = ['a','b','c',]
		dd = {}
		for i, v <- nn
			dd <- (v, i)
		`, "dd", adk(dk{"a": 0, "b": 1, "c": 2})},
		{`
		nn = ['a','b','c', 'd', 'e', 'f']
		dd = {}
		for i, v <- nn
			if i % 2 != 0
				dd <- (v, i)
		`, "dd", adk(dk{"b": 1, "d": 3, "f": 5})},
		{`
		ss = {'aa':11, 'bb':12}
		dd = {}
		dd <- ss
		`, "dd", adk(dk{"aa": 11, "bb": 12})},
		{`
		dd = {1:11, }
		ss = {3:33, 4:44}
		dd <- {2:22}
		dd <- ss
		`, "dd", adk(dk{1: 11, 2: 22, 3: 33, 4: 44})},
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
				assert.Fail(t2, "No expected Var")
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
		# v <- list
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
		{`
		# i, v <- list
		ss = [10, 11, 12, 13, 14]
		a = 0
		for i, v <- ss
			a += i + v
		`, "a", int64(70)},
		{`
		# k,v <- dict
		dd = {1:100, 2:200, 3:300}
		a = []
		a <- 0
		for k,v <- dd
			a <- k
			a <- v
		`, "a", anynn([]int64{0, 1, 100, 2, 200, 3, 300})},
		{`
		# k,v <- dict
		dd = {1:100, 2:200, 3:300}
		a = {}
		a <- 0
		for k,v <- dd
			a <- (v, k)
		`, "a", adk(dk{100: 1, 200: 2, 300: 3})},
		{`
		# for <- continue, break
		nn = [1,2,30,40,50,6,7,8,9]
		rr = []
		for i, v <- nn
			if v > 10
				continue
			if i == 7
				rr <- 88
				break
			rr <- v
		`, "rr", anis(1, 2, 6, 7, 88)},
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
				// assert.Equal(t2, tm, res)
				assert.True(t2, reflect.DeepEqual(tm, res))
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
		{`
		# break
		a = 0
		for i=3; i < 10; i += 1
			a += i
			break
		`, int64(3)},
		{`
		# for-if
		a = 100
		for i=3; i < 6; i += 1
			if i == 4
				a += 1000
			a += i
		`, int64(1112)},
		{`
		# for-if: break
		a = 100
		for i=3; i <= 5; i += 1
			if i == 4
				break
			a += i
		`, int64(103)},
		{`
		# for--if--else-if: break
		a = 100
		for i=1; i <= 5; i += 1
			if i == 6
				break
			else if i == 3
				a += 20
				break
			a += i
		`, int64(123)},
		{`
		# for--if-else: break
		a = 100
		for i=1; i <= 5; i += 1
			if i < 4
				a += 10
			else
				a += 100
				break
			a += i
		`, int64(236)},
		{`
		# for-continue
		a = 100
		for i=1; i <= 5; i += 1
			a += i
			continue
		`, int64(115)},
		{`
		# for-continue 1
		a = 100
		for i=1; i <= 5; i += 1
			if i < 3
				continue
			# 3 4 5
			a += i
		`, int64(112)},
		{`
		# for-continue 2
		a = 100
		for i=1; i <= 5; i += 1
			if i % 2 == 0
				continue
			# 1 3 5
			a += i
		`, int64(109)},
		{`
		# continue, break
		a = 100
		for i=1; i <= 9 ; i += 1
			if i < 5
				continue
			if i == 8
				break
			# 5 6 7
			a += i
		`, int64(118)},
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
