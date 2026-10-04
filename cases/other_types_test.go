package cases

import "testing"

func TestMaybeMethods(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// isNone, isSome, get
		{`
		r = []
		ss = [none, some(1), some(true), some(false), some('Hello'), some(null)]
		for n <- ss
			r <- n.isNone() ? '#is' : '#not'
			if n.isSome()
				r <- n.get()
			else 
				r <- '#none'
		`, "r", Anis("#is", "#none", "#not", 1, "#not", true, "#not", false, "#not", "Hello", "#not", Tnull())},
		// map
		{`
		r = []
		func foo(x)
			'= %v' << x
		ss = [none, some(1), some(true), some(false), some('Hello'), some(null)]
		for n <- ss
			res = n.map(foo)
			if res.isSome()
				r <- res.get()
			else
				r <- '#none'
		`, "r", Anis("#none", "= 1", "= true", "= false", "= Hello", "= <null>")},
		// fold
		{`
		func foo(p, x)
			p + x * 10
		r = []
		ss = [none, some(1), some(5), some(-9)]
		for n <- ss
			r <- n.fold(0, foo)
		`, "r", Anis(0, 10, 50, -90)},
		// .maybe
		{`
		func foo(x)
			x + 10
		r = []
		ss = [none, some(1), some(5), some(-19)]
		for n <- ss
			r <- n.maybe(0, foo)
		`, "r", Anis(0, 11, 15, -9)},
		// filter
		{`
		func foo(x)
			x > 0
		r = []
		ss = [none, some(1), some(5), some(2), some(-19), some(20), some(21), some(105), ]
		for n <- ss
			t = n.filter(foo)
			if t.isSome()
				r <- t.get()
			else
				r <- '#n1'
			d = n.filter(\x -> x % 2 != 0)
			if d.isSome()
				r <- d.get()
			else
				r <- '#n2'
		`, "r", Anis("#n1", "#n2", 1, 1, 5, 5, 2, "#n2", "#n1", -19, 20, "#n2", 21, 21, 105, 105)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

/*

 */
