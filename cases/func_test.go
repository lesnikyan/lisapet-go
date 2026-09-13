package cases

import (
	"testing"
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
ok 12. constructors of builtin types: int(), list(), tuple(), etc

9. return from: match-case

11. variative count of args, triple-dot operator
13. func overaload: by arg count, by arg types
14.1 multitype for var
14.2 multitype for args
15. builtin methods: 'a b c'.split(' ')

*/

func TestFuncBuiltNamedArg(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// named only
		{`
		r = devNM(lorem=123, ipsum=1.25, ups=(1.5), dolor='Banana')
		`, "r", adk(dk{"dolor": "Banana", "ipsum": "1.25", "lorem": "123", "ups": "1.5"})},
		{`
		r=devNM(aaa=(1.25), bbb=2.5, cc=(3+4))
		`, "r", adk(dk{"aaa": "1.25", "bbb": "2.5", "cc": "7"})},
		// ordered and named
		{`
		r=devORNM(111, 'Bubble', 2.22, bbb=2.5, cc=(3+4))
		`, "r", adk(dk{"#ordered": Anis(111, "Bubble", 2.22), "bbb": "2.5", "cc": "7"})},
		{`
		r=devORNM(1,2,3,4,5, a=6, b=7, c=8, d=9, e=10, f=11, g=12, h=13, i=14, j=15, k=16)
		`, "r", adk(dk{"#ordered": Anis(1, 2, 3, 4, 5), "a": "6", "b": "7", "c": "8", "d": "9",
			"e": "10", "f": "11", "g": "12", "h": "13", "i": "14", "j": "15", "k": "16"})},
		// default args
		{`
		r=devDef33(10, 12.5, 'Com.port1')
		`, "r", "a=10, b=12.50, c=`Com.port1`; dd=6f ee=2.22, ff=FooBar"},
		{`
		r=devDef33(10, 12.5, 'Com.port2')
		`, "r", "a=10, b=12.50, c=`Com.port2`; dd=6f ee=2.22, ff=FooBar"},
		{`
		r=devDef33(10, 12.5, 'Com.port3', dd=0xf00ba11, ee=44.44, ff='FineBeer')
		`, "r", "a=10, b=12.50, c=`Com.port3`; dd=f00ba11 ee=44.44, ff=FineBeer"},
		// ordered default
		{`
		r=devDefOrd()
		`, "r", "a=10 , b=1.00 , c=`---`"},
		{`
		r=devDefOrd(11, 22.25, 'Homo')
		`, "r", "a=11 , b=22.25 , c=`Homo`"},
		{`
		r=devDefOrd(111, 22.25, 'Lama', sep='_/_')
		`, "r", "a=111 _/_ b=22.25 _/_ c=`Lama`"},
		// named as ordered
		{`
		r=devDefOrd(111, 22.25,  sep=';', c='Puma')
		`, "r", "a=111 ; b=22.25 ; c=`Puma`"},
		{`
		r=devDefOrd(sep=';', c='Pumba', b=33.55, a=10101)
		`, "r", "a=10101 ; b=33.55 ; c=`Pumba`"},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFuncNamedArgNestedBrackets(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func f1(a=1, b=2, c=3)
			[a, b, c]
		#
		r = f1(a=(1.25), b=2.5, c=(3,4))
		`, "r", Anis(1.25, 2.5, Tanis(3, 4))},
		{`
		func f1(a=1, b=2, c=3)
			[a, b, c]
		#
		func f2(a)
			a + 1000
		#
		r = f1(a=(2.25), b=4.5, c=f2(6))
		`, "r", Anis(2.25, 4.5, 1006)},
		// {``, "r",  Anis()},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFuncBuiltMethods(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		#
		s = 'a b c'
		r = s.split(' ')
		`, "r", Anis("a", "b", "c")},
		// {``, "r",  Anis(11, )},
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
		{`
		func foo(a=0, b:float=0)
			a + b
		#
		r = foo(a=9, b=1.25)
		`, "r", float64(10.25)},
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
		r = (16)
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
