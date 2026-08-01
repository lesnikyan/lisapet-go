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

/*
TODO:
ok 1. func def
ok 2. func call
ok 3. func args
ok 4. return, return with res
ok 4.1 builtin functions, funcs len, iter
ok 4.2 named args
5. default arg val
6. var type
6.1 arg type
7. multi assign
8. multi result: return a, b, 10

--. variative count of args, triple-dot operator
-- constructors of builtin types: int(), list(), tuple(), etc
--. func overaload: by arg count, by arg types

*/

func t1() {

}

func PreloadContext(cx base.Context) {
	nodes.PreloadFuncs(cx)

}

type TTst = struct {
	src   string // code
	vname string // var
	res   any    // exp
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
			fmt.Printf("tt#int#1  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
			assert.Equal(t2, tt.res, vobj)
		case string:
			fmt.Printf("tt#string#1  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
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

func _TestFuncDefNamedArgs(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFuncCallNamedArgs(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func foo(a, b)
			a * 100 + b
		#
		r = foo(2, b = 5)
		`, "r", int64(205)},
		{`
		func foo(a, b, c)
			a*100 + b*10 + c
		#
		r = foo(c=1, b=2, a=3)
		`, "r", int64(321)},
		{`
		func foo(a, b, c)
			a*1000 + b*100 + c+10
		#
		r = foo(1, c=3, b=4)
		`, "r", int64(1413)},
		{`
		func foo(a,b,c,d,e)
			a*10**4 + b * 10**3 + c * 10**2 + d * 10 + e
		#
		r = []
		r <- foo(1,2,3,4,5)
		r <- foo(a=2, b=3, c=4, d=5, e=6)
		r <- foo(e=1, d=2, c=3, b=4, a=5)
		r <- foo(1,2,c=3, e=7, d=5)
		r <- foo(1 ,2 ,3, e=9, d=6)
		`, "r", Anis(12345, 23456, 54321, 12357, 12369)},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFuncBuiltins(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = len("hello! 12345")
		`, "r", int64(12)},
		{`
		s = "hello! 123"
		r = len(s)
		`, "r", int64(10)},
		{`
		s = []
		r = len(s)
		`, "r", int64(0)},
		{`
		s = [1,2,3,4,5]
		r = len(s)
		`, "r", int64(5)},
		{`
		r = len([1,2,3,4,5,6,7])
		`, "r", int64(7)},
		{`
		# print, just call
		print(1)
		print(1,2,3)
		print("Hey Print!")
		print([5,6,7])
		aa = [11,22]
		bb = [33,44]
		print(aa, bb)
		r = 16)
		`, "r", int64(16)},
		{`
		a1 = join([], "")
		a2 = join(['a'], "")
		a3 = join(['q','w','e'], '=')
		s4 = ['aa','bb','11','22','33','44','55','66','77','88']
		a4 = join(s4, " ")
		r = [a1, a2, a3, a4]
		`, "r", Anis("", "a", "q=w=e", "aa bb 11 22 33 44 55 66 77 88")},
		// {``, "r",  Anis(11, )},
		// {``, "r",  Anis(11, )},
		// # r = [foo(5, 2), foo(3, 5), foo(10, 9)]
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
		// t.Run(fmt.Sprintf("Test, %d) %s >>", i, tt.src), func(t2 *testing.T) {
		// 	clines := par.SplitCode(tt.src[1:])
		// 	block, err := TreeBlock(clines)
		// 	assert.Nil(t, err)
		// 	// t.Log("--- --- --- Do ...")
		// 	ctx := obb.NewContext(nil)
		// 	PreloadContext(ctx)
		// 	terr := block.Do(ctx)
		// 	if terr != nil {
		// 		assert.Fail(t2, terr.Error())
		// 		return
		// 	}
		// 	var val any
		// 	vr := ctx.GetVar(tt.vname)
		// 	// t.Log("tt#vr", vr)
		// 	if vr == nil {
		// 		vel := ctx.GetElem(tt.vname)
		// 		if vel == nil {
		// 			assert.Fail(t2, "No expected Var or elem")
		// 			return
		// 		}
		// 		val = vel.V
		// 		fmt.Printf(" TT#0: %T %v\n", val, val)
		// 	} else {
		// 		val = vr.Val
		// 	}
		// 	fmt.Printf("tt#Var#1  vr(%T, %v)  val(%T, %v) \n", vr, vr, val, val)
		// 	switch vobj := val.(type) {
		// 	case *nodes.Function:
		// 		// fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
		// 		assert.Equal(t2, tt.res, vobj.GetName())
		// 	case *obb.ListVal:
		// 		fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
		// 		assert.Equal(t2, tt.res, vobj.Elems)
		// 	case *obb.TupleVal:
		// 		fmt.Printf("tt#TupleVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
		// 		tup, ok := tt.res.(*Tup)
		// 		if !ok {
		// 			assert.Fail(t2, fmt.Sprintf("Tuple has gotten but test exp: : %T", tt.res))
		// 		}
		// 		assert.Equal(t2, tup.elems, vobj.Elems)
		// 	case *obb.DictVal:
		// 		fmt.Printf("tt#DictVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Vmap))
		// 		tm, tok := tt.res.(map[any]any)
		// 		assert.True(t2, tok)
		// 		res := pres(vobj)
		// 		assert.Equal(t2, tm, res)
		// 	case int64:
		// 		fmt.Printf("tt#int#1  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
		// 		assert.Equal(t2, tt.res, vobj)
		// 	case string:
		// 		fmt.Printf("tt#string#1  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
		// 		assert.Equal(t2, tt.res, vobj)
		// 	case *obb.Null:
		// 		fmt.Printf("tt#Null:  (%T, %v)  <Null> result \n", vr, vr)
		// 		assert.Equal(t2, tt.res, vobj)
		// 	default:
		// 		fmt.Printf("tt#default:  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
		// 		assert.Fail(t2, "unknown test result")
		// 	}
		// 	// fmt.Println("tt3>", vr, vr.Name, vr.Val)
		// })
	}
	// var aa []any = []any{1, 2, 3}
	// fmt.Println(aa)
}

func TestFuncReturn(t *testing.T) {
	// ok 1. return in function
	// ok 2. return with value
	// ok 3. return from `if`
	// ok 4. return from `for`, `while`
	// ok 5. return from deep inner block: for/for/if/if
	// 5.1 return from: match-case
	// 6. return set: return a, b, c

	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func foo()
			return 12
		#
		r = foo()
		`, "r", int64(12)},
		{`
		func foo(x)
			return x * 10
		#
		r = foo(5)
		`, "r", int64(50)},
		{`
		func foo(x)
			if x > 5
				return x * 10
			x * 100
		#
		r = foo(2)
		`, "r", int64(200)},
		{`
		func foo(x)
			if x > 5
				return x * 10
			x * 100
		#
		r = foo(7)
		`, "r", int64(70)},
		{`
		func foo(x)
			if x == 1
				return x + 1000
			else if x < 0
				return x * -10
			else
				return x * 100
		#
		r = [foo(-2), foo(1), foo(5)]
		`, "r", Anis(20, 1001, 500)},
		{`
		func foo(x, y)
			res = []
			for n <- [0 .. x]
				if n == 13
					return 3000
				if y == n
					return n + 100
			return -1 * x
		#
		a = foo(1,2)
		r = [foo(4,5), foo(3,2), foo(15, 14)]
		`, "r", Anis(-4, 102, 3000)},
		{`
		func foo(x, y)
			res = []
			for n <- [0 .. x]
				if n == 13
					return 3000
				if y == n
					return n + 100
			-1 * x
		#
		a = foo(1,2)
		r = [foo(4,5), foo(3,2), foo(15, 14)]
		`, "r", Anis(-4, 102, 3000)},
		{`
		a = []
		func foo()
			for n <- [2 .. 10]
				if n == 6
					return
				a <- n
		foo()
		`, "a", Anis(2, 3, 4, 5)},
		// {``, "r",  Anis(11, )},
		// {``, "r",  Anis(11, )},
		// # r = [foo(5, 2), foo(3, 5), foo(10, 9)]
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
		// t.Run(fmt.Sprintf("Test, %d) %s >>", i, tt.src), func(t2 *testing.T) {
		// 	clines := par.SplitCode(tt.src[1:])
		// 	block, err := TreeBlock(clines)
		// 	assert.Nil(t, err)
		// 	// t.Log("--- --- --- Do ...")
		// 	ctx := obb.NewContext(nil)
		// 	terr := block.Do(ctx)
		// 	if terr != nil {
		// 		assert.Fail(t2, terr.Error())
		// 	}
		// 	var val any
		// 	vr := ctx.GetVar(tt.vname)
		// 	// t.Log("tt#vr", vr)
		// 	if vr == nil {
		// 		vel := ctx.GetElem(tt.vname)
		// 		if vel == nil {
		// 			assert.Fail(t2, "No expected Var or elem")
		// 		}
		// 		val = vel.V
		// 		fmt.Printf(" TT#0: %T %v\n", val, val)
		// 	} else {
		// 		val = vr.Val
		// 	}
		// 	fmt.Printf("tt#Var#1  vr(%T, %v)  val(%T, %v) \n", vr, vr, val, val)
		// 	switch vobj := val.(type) {
		// 	case *nodes.Function:
		// 		// fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
		// 		assert.Equal(t2, tt.res, vobj.GetName())
		// 	case *obb.ListVal:
		// 		fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
		// 		assert.Equal(t2, tt.res, vobj.Elems)
		// 	case *obb.TupleVal:
		// 		fmt.Printf("tt#TupleVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
		// 		tup, ok := tt.res.(*Tup)
		// 		if !ok {
		// 			assert.Fail(t2, fmt.Sprintf("Tuple has gotten but test exp: : %T", tt.res))
		// 		}
		// 		assert.Equal(t2, tup.elems, vobj.Elems)
		// 	case *obb.DictVal:
		// 		fmt.Printf("tt#DictVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Vmap))
		// 		tm, tok := tt.res.(map[any]any)
		// 		assert.True(t2, tok)
		// 		res := pres(vobj)
		// 		assert.Equal(t2, tm, res)
		// 	case int64:
		// 		fmt.Printf("tt#int#1  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
		// 		assert.Equal(t2, tt.res, vobj)
		// 	case string:
		// 		fmt.Printf("tt#string#1  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
		// 		assert.Equal(t2, tt.res, vobj)
		// 	case *obb.Null:
		// 		fmt.Printf("tt#Null:  (%T, %v)  <Null> result \n", vr, vr)
		// 		assert.Equal(t2, tt.res, vobj)
		// 	default:
		// 		fmt.Printf("tt#default:  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
		// 		assert.Fail(t2, "unknown test result")
		// 	}
		// 	// fmt.Println("tt3>", vr, vr.Name, vr.Val)
		// })
	}
}

func TestFuncDef(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# func def
		func foo()
			1
		#
		r = foo
		`, "foo", "foo"},
		{`
		# func def
		func foo()
			12
		#
		r = foo()
		`, "r", int64(12)},
		{`
		func foo(a)
			a + 10
		#
		r = 1
		r = foo(5)
		`, "r", int64(15)},
		{`
		func foo(s)
			"fres " + s
		#
		r = ""
		r = foo('abc')
		`, "r", "fres abc"},
		{`
		func foo(a, b)
			a * b
		#
		r = 0
		r = foo(2, 7)
		`, "r", int64(14)},
		{`
		func foo(a, b, c)
			[a, b, c]
		#
		r = 0
		r = foo(2, 7, -4)
		`, "r", Anis(2, 7, -4)},
		{`
		func foo(a, b, c, d, e, f, g, h, i, j, k)
			(a, b, c, d, e, f, g, h, i, j, k)
		#
		r = 0
		r = foo(1,2,3,4,5,6,7,8,9,0,-1)
		`, "r", Tanis(1, 2, 3, 4, 5, 6, 7, 8, 9, 0, -1)},
		// {``, "r",  int64(1)},
		// {``, "r",  Anis(11, 12, 13, 14)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
		// t.Run(fmt.Sprintf("Test, %d) %s >>", i, tt.src), func(t2 *testing.T) {
		// 	clines := par.SplitCode(tt.src[1:])
		// 	block, err := TreeBlock(clines)
		// 	assert.Nil(t, err)
		// 	// t.Log("--- --- --- Do ...")
		// 	ctx := obb.NewContext(nil)
		// 	terr := block.Do(ctx)
		// 	if terr != nil {
		// 		assert.Fail(t2, terr.Error())
		// 	}
		// 	var val any
		// 	vr := ctx.GetVar(tt.vname)
		// 	// t.Log("tt#vr", vr)
		// 	if vr == nil {
		// 		vel := ctx.GetElem(tt.vname)
		// 		if vel == nil {
		// 			assert.Fail(t2, "No expected Var or elem")
		// 		}
		// 		val = vel.V
		// 		fmt.Printf(" TT#0: %T %v\n", val, val)
		// 	} else {
		// 		val = vr.Val
		// 	}
		// 	fmt.Printf("tt#Var#1  vr(%T, %v)  val(%T, %v) \n", vr, vr, val, val)
		// 	switch vobj := val.(type) {
		// 	case *nodes.Function:
		// 		// fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
		// 		assert.Equal(t2, tt.res, vobj.GetName())
		// 	case *obb.ListVal:
		// 		fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
		// 		assert.Equal(t2, tt.res, vobj.Elems)
		// 	case *obb.TupleVal:
		// 		fmt.Printf("tt#TupleVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Elems))
		// 		tup, ok := tt.res.(*Tup)
		// 		if !ok {
		// 			assert.Fail(t2, fmt.Sprintf("Tuple has gotten but test exp: : %T", tt.res))
		// 		}
		// 		assert.Equal(t2, tup.elems, vobj.Elems)
		// 	case *obb.DictVal:
		// 		fmt.Printf("tt#DictVal#1  (%T, %v)  (%T, %v) len: %d \n", tt.res, tt.res, vobj, vobj, len(vobj.Vmap))
		// 		tm, tok := tt.res.(map[any]any)
		// 		assert.True(t2, tok)
		// 		res := pres(vobj)
		// 		assert.Equal(t2, tm, res)
		// 	case int64:
		// 		fmt.Printf("tt#int#1  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
		// 		assert.Equal(t2, tt.res, vobj)
		// 	case string:
		// 		fmt.Printf("tt#string#1  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
		// 		assert.Equal(t2, tt.res, vobj)
		// 	case *obb.Null:
		// 		fmt.Printf("tt#Null:  (%T, %v)  <Null> result \n", vr, vr)
		// 		assert.Equal(t2, tt.res, vobj)
		// 	default:
		// 		fmt.Printf("tt#default:  (%T, %v)  (%T, %v) \n", vr, vr, vobj, vobj)
		// 		assert.Fail(t2, "unknown test result")
		// 	}
		// 	// fmt.Println("tt3>", vr, vr.Name, vr.Val)
		// })
	}
}
