package cases

import (
	"fmt"
	"testing"

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
4. return, return with res
4.1 builtin functions, funcs len, iter
4.2 named args
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

func TestFuncReturn(t *testing.T) {
	// 1. return in function
	// 2. return with value
	// 3. return from `if`
	// 4. return from `for`, `while`
	// 5. return from deep inner block: for/for/if/if
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
		t.Run(fmt.Sprintf("Test, %d) %s >>", i, tt.src), func(t2 *testing.T) {
			clines := par.SplitCode(tt.src[1:])
			block, err := TreeBlock(clines)
			assert.Nil(t, err)
			// t.Log("--- --- --- Do ...")
			ctx := obb.NewContext(nil)
			terr := block.Do(ctx)
			if terr != nil {
				assert.Fail(t2, terr.Error())
			}
			var val any
			vr := ctx.GetVar(tt.vname)
			// t.Log("tt#vr", vr)
			if vr == nil {
				vel := ctx.GetElem(tt.vname)
				if vel == nil {
					assert.Fail(t2, "No expected Var or elem")
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
		t.Run(fmt.Sprintf("Test, %d) %s >>", i, tt.src), func(t2 *testing.T) {
			clines := par.SplitCode(tt.src[1:])
			block, err := TreeBlock(clines)
			assert.Nil(t, err)
			// t.Log("--- --- --- Do ...")
			ctx := obb.NewContext(nil)
			terr := block.Do(ctx)
			if terr != nil {
				assert.Fail(t2, terr.Error())
			}
			var val any
			vr := ctx.GetVar(tt.vname)
			// t.Log("tt#vr", vr)
			if vr == nil {
				vel := ctx.GetElem(tt.vname)
				if vel == nil {
					assert.Fail(t2, "No expected Var or elem")
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
}
