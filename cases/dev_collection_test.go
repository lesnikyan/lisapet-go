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
1. slice: nn[a : b]
constructors:
2. list()
3. dict()
4. tuple()

*/

func TestCollDelElemOper(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		nn = [1,2,3,4,5]
		nn - [2]
		`, "nn", Anis(1, 2, 4, 5)},
		{`
		
		nn = [1,2,303,4,5]
		r = nn - [2]
		`, "r", int64(303)},
		// {``, "r",  Anis(11, )},
		// {``, "r",  Anis(11, )},
		// {``, "r",  Anis(11, )},
		// {``, "r",  int64(205)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestCollectionsAddOper(t *testing.T) {
	// tuple + tuple
	// list += list
	// dict += dict
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		nn = [1,2,3]
		x = 1
		x += 5
		nn += [15]
		nn += [16]
		`, "nn", Anis(1, 2, 3, 15, 16)},
		{`
		a = [1,2,3]
		b = [44, 55]
		r = a + b
		`, "r", Anis(1, 2, 3, 44, 55)},
		{`
		a = (1,2)
		b = (33,44)
		r = a + b
		`, "r", Tanis(1, 2, 33, 44)},
		{`
		r = (1,2)
		b = (33, 55)
		r += b
		`, "r", Tanis(1, 2, 33, 55)},
		{`
		a = {'a': 11}
		b = {'b':22, 'c':33}
		r = a + b
		`, "r", adk(dk{"a": 11, "b": 22, "c": 33})},
		{`
		r = {'a':11}
		b = {'b':222, 'c':333}
		r += b
		`, "r", adk(dk{"a": 11, "b": 222, "c": 333})},
		// {``, "r",  Anis(11, )},
		// {``, "r",  Anis(11, )},
		// {``, "r",  int64(205)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestListsCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# empty list
		nn = []
		`, "nn", Anynn([]any{})},
		{`
		# non-empty list
		nn = [1,2,3,4,5]
		`, "nn", Anynn([]int64{1, 2, 3, 4, 5})},
		{`
		# write to list
		nn = [1,2,3]
		nn[1] = 222
		`, "nn", Anynn([]int64{1, 222, 3})},
		{`
		# write to list in loop
		a = [0,0,0,0,0]
		for i=0; i < 5; i += 1
			a[i] = 10 + i
		`, "a", Anynn([]int64{10, 11, 12, 13, 14})},
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
		`, "a", Anynn([]int64{11, 12, 13, 14, 15})},
		{`
		# list[i] += v
		a = [0,0,0,0,100]
		for i=0; i < 5; i += 1
			a[i] += 10 + i
		`, "a", Anynn([]int64{10, 11, 12, 13, 114})},
		{`
		# list[i] += list[i]
		nn = [11,12,13,14,15]
		a = [0,0,0,0,100]
		for i=0; i < 5; i += 1
			a[i] += nn[i]
		`, "a", Anynn([]int64{11, 12, 13, 14, 115})},
		// {`` "nn", int64(1)},
		// {``, "nn", int64(1)},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("Test, %s >>", tt.src), func(t2 *testing.T) {
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

type Tup struct {
	elems []any
}

func tanynn[T any](vals []T) *Tup {
	vv := Anynn(vals)
	return &Tup{elems: vv}
}

func TestTupleCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# empty tuple
		nn = (,)
		`, "nn", tanynn([]any{})},
		{`
		# empty tuple
		nn = (1,)
		`, "nn", tanynn([]any{int64(1)})},
		{`
		# empty tuple
		nn = ('abc',)
		`, "nn", tanynn([]any{"abc"})},
		{`
		# empty tuple
		nn = (1,2)
		`, "nn", tanynn([]any{1, 2})},
		{`
		# non-empty tuple
		nn = (1,2,3,4,5)
		`, "nn", tanynn([]int64{1, 2, 3, 4, 5})},
		{`
		# tuple, different type elems
		nn = (1,2,true, false, "5", "hello")
		`, "nn", tanynn([]any{1, 2, true, false, "5", "hello"})},
		{`
		# read from tuple
		nn = (1,2,3,4,100])
		a = 0
		for i = 0; i < 5; i += 1
			a += nn[i]
		`, "a", int64(110)},
		{`
		# from tuple to list
		tt = (11, 12, 13, 14, 115)
		a = [0,0,0,0,0]
		for i = 0; i < 5; i += 1
			a[i] = tt[i]
		`, "a", Anynn([]any{11, 12, 13, 14, 115})},
		{`
		# list[i] += tuple[i]
		tt = (11, 12, 13, 14, 115)
		a = [10,20,30,40,50]
		for i = 0; i < 5; i += 1
			a[i] += tt[i]
		`, "a", Anynn([]any{21, 32, 43, 54, 165})},
		// {``, "nn", int64(1)},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("Test, %s >>", tt.src), func(t2 *testing.T) {
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
				fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
				assert.Equal(t2, tt.res, vobj.Elems)
			case *obb.TupleVal:
				fmt.Printf("tt#TupleVal#1  (%T, %v)  (%T, %v) len: %d \n", vr, vr, vobj, vobj, len(vobj.Elems))
				tup, tok := tt.res.(*Tup)
				assert.True(t2, tok)
				assert.Equal(t2, tup.elems, vobj.Elems)
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

func TestDictCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# empty dict
		dd = {}
		`, "dd", dk{}},
		{`
		# dict 1 elem
		dd = {'a': 11, }
		`, "dd", adk(dk{"a": 11})},
		{`
		# dict 1 elem
		dd = {'a': 12}
		`, "dd", adk(dk{"a": 12})},
		{`
		# dict
		dd = {'a': 11, 'b':22, 'c':33}
		`, "dd", adk(dk{"a": 11, "b": 22, "c": 33})},
		{`
		# dict str:str
		dd = {'a': 'abc'}
		`, "dd", dk{"a": "abc"}},
		{`
		# dict different type
		dd = {'a': 'abc', 2:22, 'c':[1,2,3]}
		`, "dd", adk(dk{"a": "abc", 2: 22, "c": Anynn([]int64{1, 2, 3})})},
		{`
		# dict from tuple
		tt = (1,2)
		dd = {tt[0] : tt[1]}
		`, "dd", adk(dk{1: 2})},
		{`
		# get elem
		dd = {'a':1, 'b': 20}
		a = dd['a']
		a += dd['b']
		`, "a", int64(21)},
		{`
		# set elem
		dd = {}
		dd['a'] = 11
		dd['b'] = 20
		dd['a'] = 111
		`, "dd", adk(dk{"a": 111, "b": 20})},
		{`
		# dd[k] += dd[k]
		dd = {'a':3, 'b': 2}
		dd['a'] += dd['b']
		`, "dd", adk(dk{"a": 5, "b": 2})},
		{`
		# read from dict
		nn = ['a','b','c']
		dd = {'a': 2, 'b': 3, 'c': 20}
		a = 0
		for i = 0; i < 3; i += 1
			k = nn[i]
			a += dd[k]
		`, "a", int64(25)},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("Test, %s >>", tt.src), func(t2 *testing.T) {
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
