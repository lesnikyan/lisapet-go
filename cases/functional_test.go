package cases

import "testing"

/*
ok 1. x -> x + 10
ok 2. (x, y) -> x + y
ok 3. _ -> 1
ok 4. \x, y,  -> x + y
*/

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
