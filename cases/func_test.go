package cases

import (
	"testing"

	"github.com/lesnikyan/lisapet-go/objects"
)

/*
TODO:
ok 1. func def
ok 2. func call
ok 3. func args
ok 4. return, return with res
ok 4.1 builtin functions, funcs len, iter
ok 4.2 named args
ok 5. default arg val
ok 6. var type
ok 6.1 arg type
ok 7. multi assign
ok 8. multi result: return a, b, 10

9. return from: match-case

11. variative count of args, triple-dot operator
12. constructors of builtin types: int(), list(), tuple(), etc
13. func overaload: by arg count, by arg types
14.1 multitype for var
14.2 multitype for args

*/

func TestConstructMaybe(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = some(1)
		`, "r", Tmay(1)},
		{`
		r:maybe = some(2)
		`, "r", Tmay(2)},
		{`
		r = some('Hello')
		`, "r", Tmay("Hello")},
		{`
		r = some([1,2,3])
		`, "r", Tmay(Anis(1, 2, 3))},
		{`
		struct Abc a: int
		a1 = Abc{a:5}
		r = some(a1)
		`, "r", Tmay(Stf("Abc", dk{"a": 5}))},

		{`
		r:maybe = none
		`, "r", objects.None()},
		{`
		a1 = some('in-some')
		r = some(a1)
		`, "r", Tmay(Tmay("in-some"))},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestConstructByte(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r:byte = byte(1)
		`, "r", byte(1)},
		{`
		r = byte(100)
		`, "r", byte(100)},
		{`
		r = byte(255)
		`, "r", byte(255)},
		{`
		r = byte(0x105)
		`, "r", byte(5)},
		{`
		r = byte(0xf)
		`, "r", byte(0xf)},
		{`
		r = byte(true)
		`, "r", byte(1)},
		{`
		r = byte(false)
		`, "r", byte(0)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestConstructString(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = string('Hello')
		`, "r", "Hello"},
		{`
		r = string(65)
		`, "r", "65"},
		{`
		bb = 0x[fa]
		r = string(bb[0])
		`, "r", "250"},
		{`
		r = string([65, 66, 67, 97, 115, 116])
		`, "r", "ABCast"},
		{`
		bb = 0x[48 65 6c 6c 6f 20 62 79 74 65 73 21]
		r = string(bb)
		`, "r", "Hello bytes!"},
		{`
		r = string([c6 90 20 c6 8d 20 c6 80 20 c6 8b 20 c6 95 20 c6 a9 20 c6 b1 20 c6 b3 20 c6 9b])
		`, "r", "Ɛ ƍ ƀ Ƌ ƕ Ʃ Ʊ Ƴ ƛ"},
		{`
		r = string([g'G', g'L', g'i', g'P', g'h', g'S'])
		`, "r", "GLiPhS"},
		{`
		r = string(null)
		`, "r", "null"},
		{`
		r = string(true)
		`, "r", "true"},
		{`
		r = string(false)
		`, "r", "false"},
		{`
		r = string('')
		`, "r", ""},
		{`
		r = string()
		`, "r", ""},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFuncCompatibledArgs(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func foo(x: int)
			x * 2
		#
		r = foo(true)
		`, "r", int64(2)},
		{`
		func foo(x: int, n=2)
			x * n + 1
		#
		r = foo(false)
		`, "r", int64(1)},
		{`
		func foo(x: float, y: float)
			x * y
		#
		r = foo(3, 5)
		`, "r", float64(15)},
		{`
		func foo(x: int, y: int = 2)
			x + y
		#
		r = foo(null, 3)
		`, "r", int64(3)},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFuncDefTypedArgs(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func foo(x: int)
			x * 2
		#
		r = foo(23)
		`, "r", int64(46)},
		{`
		func foo(x: int, n=2)
			x * n -1
		#
		r = foo(23)
		`, "r", int64(45)},
		{`
		func foo(x: int, y: int)
			x * y
		#
		r = foo(23, 3)
		`, "r", int64(69)},
		{`
		func foo(x: int, y: int = 2)
			x * y
		#
		r = foo(23, 3)
		`, "r", int64(69)},
		{`
		# all cases
		func foo(a:int, b:int=2, c=100)
			a * b + c
		#
		r = []
		r <- foo(1)
		r <- foo(1,5)
		r <- foo(11, 6, 200)
		r <- foo(3, 9)
		r <- foo(3, c=30)
		r <- foo(4, c=120, b=2)
		`, "r", Anis(102, 105, 266, 127, 36, 128)},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestFuncDefNamedArgs(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func foo(x, n = 2)
			x * n
		#
		r = foo(23)
		`, "r", int64(46)},
		{`
		func foo(a, b=2, c=100)
			a * b + c
		#
		r = []
		r <- foo(1)
		r <- foo(1,5)
		r <- foo(11, 6, 200)
		r <- foo(3, 9)
		r <- foo(3, c=30)
		r <- foo(4, c=120, b=2)
		`, "r", Anis(102, 105, 266, 127, 36, 128)},
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
	}
}
