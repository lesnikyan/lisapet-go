package cases

import "testing"

/*
ok 1. x -> x + 10
ok 2. (x, y) -> x + y
ok 3. _ -> 1
ok 4. \x, y,  -> x + y
ok 5. compose(foo, bar)
ok 6. foo * bar # compose oper
ok 7. foo $ arg # apply oper
8. carry()
9. foo~>
*/

// carry oper ~>
func TestFuncCarryOper(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func foo(x, y)
			x + 10 * y
		#
		f = foo ~>
		r = f(3)(2)
		`, "r", int64(23)},
		{`
		# 1 arg, just formal test
		func foo(x)
			x + 10
		#
		f = foo~>
		r = f(4)
		`, "r", int64(14)},
		{`
		# 3 args
		func foo(x, y, z)
			x * 100 + y * 10 + z
		#
		r = foo~>(3)(4)(5)
		`, "r", int64(345)},
		{`
		func foo(a, b, c, d)
			[a, b, c, d]
		#
		r = foo~>(1)(2)(3)(4)
		`, "r", Anis(1, 2, 3, 4)},
		{`
		func foo(a, b, c, d, e)
			[a, b, c, d, e*10]
		#
		r = foo~>(1)(2)(3)(4)(5)
		`, "r", Anis(1, 2, 3, 4, 50)},
		{`
		# 8 args
		func foo(a, b, c, d, e, f, g, h)
			[a, b, c, d, e+10, f+20, g+30, h+100]
		#
		r = foo~>(1)(2)(3)(4)(5)(6)(7)(8)
		`, "r", Anis(1, 2, 3, 4, 15, 26, 37, 108)},
		// combine with composition and $
		{`
		func foo(x, y)
			x + 10 * y
		#
		r = foo~>(3) $ 7
		`, "r", int64(73)},
		{`
		func foo(x, y)
			x + 10 * y
		#
		r = foo~> $ 3 $ 9
		`, "r", int64(93)},
		{`
		func foo(x, y)
			x + 10 * y
		func bar(x)
			10 + x
		#
		r = foo~>(3) * bar $ 7
		`, "r", int64(173)},
		{`
		func foo(x, y)
			x + 10 * y
		#
		r = foo~>(3) * (\x -> x + 20) $ 7
		`, "r", int64(273)},
		// carry method
		{`
		struct A a:int
		#
		func q:A foo(x, y)
			x + q.a * y
		#
		a1 = A{a:20}
		r = a1.foo~>(3)(5)
		`, "r", int64(103)},
		// lambda as arg
		{`
		func foo(f, x, y)
			f(y + x)
		#
		r = [] 
		f2 = foo~>(x -> x * 3)(20)
		for n <- [1 .. 5]
			r <- f2(n)
		`, "r", Anis(63, 66, 69, 72, 75)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

// carry(foo)
func TestFuncCarry(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func foo(x, y)
			x + 10 * y
		#
		f = carry(foo)
		r = f(3)(4)
		`, "r", int64(43)},
		{`
		# 1 arg, just formal test
		func foo(x)
			x + 10
		#
		f = carry(foo)
		r = f(4)
		`, "r", int64(14)},
		{`
		# 3 args
		func foo(x, y, z)
			x * 100 + y * 10 + z
		#
		f = carry(foo)
		r = f(3)(4)(5)
		`, "r", int64(345)},
		{`
		func foo(a, b, c, d)
			[a, b, c, d]
		#
		f = carry(foo)
		r = f(1)(2)(3)(4)
		`, "r", Anis(1, 2, 3, 4)},
		{`
		func foo(a, b, c, d, e)
			[a, b, c, d, e*10]
		#
		f = carry(foo)
		r = f(1)(2)(3)(4)(5)
		`, "r", Anis(1, 2, 3, 4, 50)},
		{`
		# 8 args
		func foo(a, b, c, d, e, f, g, h)
			[a, b, c, d, e+10, f+20, g+30, h+100]
		#
		f = carry(foo)
		r = f(1)(2)(3)(4)(5)(6)(7)(8)
		`, "r", Anis(1, 2, 3, 4, 15, 26, 37, 108)},
		// combine with composition and $
		{`
		func foo(x, y)
			x + 10 * y
		#
		f = carry(foo)
		r = f(3) $ 7
		`, "r", int64(73)},
		{`
		func foo(x, y)
			x + 10 * y
		#
		f = carry(foo)
		r = f $ 3 $ 9
		`, "r", int64(93)},
		{`
		func foo(x, y)
			x + 10 * y
		func bar(x)
			10 + x
		#
		f = carry(foo)
		r = f(3) * bar $ 7
		`, "r", int64(173)},
		// carry method
		{`
		struct A a:int
		#
		func q:A foo(x, y)
			x + q.a * y
		#
		a1 = A{a:20}
		f = carry(a1.foo)
		r = f(3)(5)
		`, "r", int64(103)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

// func $ arg
func TestFuncApplyOper(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func foo(x)
			x + 10
		#
		r = foo $ 5
		`, "r", int64(15)},
		// (lambda) $ arg
		{`
		foo = x -> x + 20
		r = foo $ 5
		`, "r", int64(25)},
		{`
		r = (x -> x + 20) $ 6
		`, "r", int64(26)},
		{`
		r = x -> x + 20 $ 7
		`, "r", int64(27)},
		// foo * bar $ arg
		{`
		r = (x -> x + 20) $ 6
		`, "r", int64(26)},
		{`
		func foo(x)
			x + 20
		func bar(x)
			x * 2
		r = foo * bar $ 3
		`, "r", int64(26)},
		{`
		func foo(x)
			x + 20
		func bar(x)
			x * 2
		r = bar * foo $ 3
		`, "r", int64(46)},
		// foo $ bar(x)
		{`
		func foo(x)
			x + 20
		func bar(x)
			x * 2
		r = foo $ bar(4)
		`, "r", int64(28)},
		{`
		func bar(x)
			x * 2
		r = \x -> x + 30 $ bar(4)
		`, "r", int64(38)},
		// foo $ x $ y # foo => func
		{`
		func foo(f)
			f(100)
		# 
		func bar(x)
			\y -> x + y
		#
		r = foo $ bar $ 7
		`, "r", int64(107)},
		{`
		func foo(f)
			f(100)
		# 
		func bar(x)
			\y -> x + y
		#
		r = foo $ (bar $ 8)
		`, "r", int64(108)},
		{`
		func foo(f)
			k = 100
			func (x)
				f(k + x)
		# 
		func bar(x)
			x * 2
		#
		r = (foo $ bar) $ 9
		`, "r", int64(218)},
		{`
		func foo(x)
			y = x * 10
			func Q(f)
				f(y)
		# 
		func bar(x)
			\y -> x + y
		#
		r = foo(3) $ bar $ 5
		`, "r", int64(35)},
		// onj.method $ arg
		{`
		struct A a: int
		func q: A foo(x)
			q.a * x
		#
		a1 = A{a:3}
		r = a1.foo $ 7
		`, "r", int64(21)},
		{`
		struct A a: int
		func q: A foo(x)
			q.a * x
		#
		func bar(x)
			x + 100
		#
		a1 = A{a:3}
		r = a1.foo * bar $ 7
		`, "r", int64(321)},
		{`
		struct A a: int
		func q: A foo(x)
			q.a * x
		#
		func bar(x)
			x + 100
		#
		a1 = A{a:3}
		r = bar * a1.foo $ 7
		`, "r", int64(121)},
		// builtin funcs
		{`
		r = 'q w e r'.split $ ' '
		`, "r", Anis("q", "w", "e", "r")},

		// map, fold
		{`
		nn = [1 .. 5]
		r = nn.map $ \x -> x * 10
		`, "r", Anis(10, 20, 30, 40, 50)},
		{`
		# fold can't apply because 2 args, carrying needed
		func xfold(nn)
			func R(x)
				func Q(f)
					nn.fold(x, f)
		nn = [1 .. 5]
		r = xfold(nn)(1000) $ \x,y -> x + y	
		`, "r", int64(1015)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

// composition, builtin oper
func TestFuncComposeOper(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func foo(x)
			x + 10
		func bar(x)
			x * 2
		f = foo * bar
		r = f(3)
		`, "r", int64(16)},
		{`
		func foo(x)
			x + 10
		func bar(x)
			x * 2
		func baz(x)
			x + 5
		f = foo * bar * baz
		r = f(3)
		`, "r", int64(26)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

// composition, builtin func
func TestFuncCompose(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		func foo(x)
			x + 10
		func bar(x)
			x * 2
		f = compose(foo, bar)
		r = []
		for n <- [1 .. 5]
			r <- f(n)
			`, "r", Anis(12, 14, 16, 18, 20)},
		{`
		func foo(x)
			x + 100
		func bar(x)
			x + 10
		func baz(x)
			x * 2
		f = compose(foo, bar, baz)
		r = []
		for n <- [1 .. 5]
			r <- f(n)
			`, "r", Anis(112, 114, 116, 118, 120)},
		{`
		func foo(x)
			x + 100
		func bar(x)
			x + 10
		func baz(x)
			x * 2
		f = compose(foo, bar, baz, \x -> x % 4)
		r = []
		for n <- [1 .. 5]
			r <- f(n)
			`, "r", Anis(112, 114, 116, 110, 112)},
		{`
		f = compose(\n -> n * 2, \x -> x + 10)
		r = []
		for n <- [1 .. 5]
			r <- f(n)
			`, "r", Anis(22, 24, 26, 28, 30)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

// #foo = \x, y -> x + y
// #r = foo(1)
func TestLambdaBackSlash(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = 1
		foo = \ x -> x + 20
		r = foo(5)
		`, "r", int64(25)},
		{`
		r = 1
		foo = \ x, y -> x + y * 10
		r = foo(6, 3)
		`, "r", int64(36)},
		{`
		func foo(f, x)
			f(x)
		r = []
		f1 = \ x -> x + 10
		r <- foo(f1, 3)
		r <- foo(\x -> x * 10, 4)
		`, "r", Anis(13, 40)},
		{`
		func foo(x, y, f)
			f(x, y)
		r = []
		f1 = \ x, y -> x * 10 + y
		r <- foo(1, 2, f1)
		r <- foo(5, 7, \x, y -> x * y)
		`, "r", Anis(12, 35)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestLambda(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		foo = x -> x + 10
		r = foo(5)
		`, "r", int64(15)},
		// lambda as an arg
		{`
		func foo(f, x)
			f(x)
		#
		f1 = (x) -> x + 30
		r = []
		r <- foo(f1, 2)
		r <- foo(x -> x + 100, 4)
		`, "r", Anis(32, 104)},
		// two args
		{`
		foo = (x, y) -> x + y
		r = foo(5, 20)
		`, "r", int64(25)},
		{`
		r = ((x, y) -> x + y)(3, 20)
		`, "r", int64(23)},
		// return lambda
		{`
		func foo(f, x)
			f(x)
		#
		func getf()
			x -> x + 40
		r = foo(getf(), 3)
		`, "r", int64(43)},
		{`
		func getf()
			x -> x + 50
		r = getf()(6)
		`, "r", int64(56)},
		{`
		func getf(y)
			x -> x + y
		r = getf(60)(7)
		`, "r", int64(67)},
		// static result
		{`
		foo = x -> 1
		r = []
		nn = [3, 4, 5]
		for n <- nn
			r <- n
			r <- foo(n)
		`, "r", Anis(3, 1, 4, 1, 5, 1)},
		// no args
		{`
		foo = _ -> 1
		r = []
		nn = [3, 4, 7]
		r <- foo()
		for n <- nn
			r <- n
			r <- foo(n)
		`, "r", Anis(1, 3, 1, 4, 1, 7, 1)},
		// inline block in lambda
		{`
		sum = (a, b) -> a + b
		foo = x -> (nn = [1 .. x]; nn.fold(0, sum))
		r = []
		for m <- [5, 6, 7]
			r <- foo(m)
		`, "r", Anis(15, 21, 28)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
