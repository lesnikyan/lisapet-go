package cases

import "testing"

func TestMatchTypeVar(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		x = 23
		r = 0
		match x
			a :: int
				r = 100 + a
		`, "r", int64(123)},
		// simple vals
		{`
		nn = [11, 12, 3.5, 4.002, g'Q', g'S', 'hello', null]
		
		#nn = [44, g'Q']
		r = []
		for x <- nn
			match x
				a :: int
					r <- 100 + a
				b :: float
					r <- 200.0 + b
				c :: glif
					r <- "'%c'" << c
				s :: string
					r <- '%% %s %%' << s
				_
					r <- 99
		`, "r", Anis(111, 112, 203.5, 204.002, "'Q'", "'S'", "% hello %", 99)},
		// containers
		{`
		nn = [{11:111}, (2,3), [4,5,6], some(7), none, null]
		r = []
		for x <- nn
			match x
				a :: list
					r += a
				b :: tuple
					for t <- b
						r <- t
				d :: dict
					for k, v <- d
						r <- k
						r <- v
				c :: maybe
					if c.isSome()
						r <- c.get()
					else
						r <- 98
				_
					r <- 99
		`, "r", Anis(11, 111, 2, 3, 4, 5, 6, 7, 98, 99)},
		// structs
		{`
		struct A a: int
		struct B b: float
		struct C(A) c: string
		struct D(A) d: bool
		#
		nn = [A(11), A(22), B(3.2), C(4, 'cats'), D(55, false), null, 19]
		r = []
		for x <- nn
			match x
				c :: C
					r <- 'C.c:%s' << c.c
				b :: B
					r <- 200.0 + b.b
				a :: A
					r <- 100 + a.a
				n :: null
					r <- n 
				_
					r <- 99
		`, "r", Anis(111, 122, 203.2, "C.c:cats", 155, Tnull(), 99)},
		// functions
		{`
		struct A a: int
		#
		func foo(x)
			x * 5
		func bar(x)
			x + 100
		func ns:A sum(x)
			x + ns.a
		#
		a1 = A(300)
		nn = [foo, bar, a1.sum, \x -> x * 3, 19]
		r = []
		for x <- nn
			match x
				f :: function
					r <- f(12)
				_
					r <- 999
		`, "r", Anis(60, 112, 312, 36, 999)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestMatchReturn(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// match / return
		{`
		func foo(x)
			match x
				:: int
					return 'num %d' << x
				"hello"
					return "greeting"
				:: string
					return 'str: %s' << x
				_
					return '_48'
		nn = [1, "quatro", 'hello', 2.5, 224]
		r = []
		for n <- nn
			r <- foo(n)
		`, "r", Anis("num 1", "str: quatro", "greeting", "_48", "num 224")},
		// match / if / return
		{`
		func foo(x)
			match x
				:: int
					if x % 2 == 0
						return 'x2 num %d' << x
					return 'num %d' << x
				"hello"
					return "greeting"
				:: string
					if len(x) > 8
						return 's: %s' << x[0:9]
					else
						return 'str: %s' << x
				_
					return '_49'
		nn = [1, "quatro", 'hello Barbarian', 2.5, 224]
		r = []
		for n <- nn
			r <- foo(n)
		`, "r", Anis("num 1", "str: quatro", "s: hello Bar", "_49", "x2 num 224")},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestMatchTypeNoVar(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		x = 111
		r = 0
		match x
			:: int
				r = 5
		`, "r", int64(5)},
		{`
		nn = [1, -5, 1.5, true, false, null, "", "hello", 00xf1, g'Q', mur(1)]
		r = []
		for x <- nn
			r <- x
			match x
				:: int
					r <- 5
				:: float
					r <- 6
				:: glif
					r <- 7
				:: bool
					r <- 9
				:: byte
					r <- 8
				:: string
					r <- 11
				:: null
					r <- 1100
				_
					r <- 99
		`, "r", Anis(1, 5, -5, 5, 1.5, 6, true, 9, false, 9, Tnull(), 1100, "", 11, "hello", 11, byte(0xf1), 8, 'Q', 7, mr(1), 99)},
		{`
		nn = [1, "", 0x[], [], (,), {}, 0x[ff 00 11], [3,4,5], (5,6,7), {'a':'ant'}]
		#nn = [1, "", [], (,), {}]
		r = []
		for x <- nn
			match x
				:: list
					r <- 11
				:: tuple
					r <- 12
				:: dict
					r <- 13
				:: string
					r <- 14
				:: bytes
					r <- 15
				_
					r <- 99
		`, "r", Anis(99, 14, 15, 11, 12, 13, 15, 11, 12, 13)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestMatchValue(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		x = 111
		r = 0
		match x
			111
				r = 5
		`, "r", int64(5)},

		{`
		# DEbug
		x = -111
		r = 0
		match x
			-111
				r = 5
		`, "r", int64(5)},

		{`
		x = "is"
		r = 0
		match x
			111
				r = 5
			"is"
				r = 8
		`, "r", int64(8)},
		{`
		x = 00x22
		r = 0
		match x
			111
				r = 5
			"is"
				r = 8
			00x22
				r = 16
		`, "r", int64(16)},
		{`
		nn = [1, -100, "#11", g'Q', g'W', 00xf7, false, true, [1,2,3], null, 1.5, 2.2]
		rr = []
		for x <- nn
			r = 0
			match x
				-100
					r = 1011
				"#11"
					r = 12
				00xf7
					r = 13
				true
					r = 15
				null
					r = 17
				1.5
					r = 18
				g'Q'
					r = 14
			rr <- r
		`, "rr", Anis(0, 1011, 12, 14, 0, 13, 0, 15, 0, 17, 18, 0)},
		{`
		nn = [1, 2, 3, 4]
		rr = []
		for x <- nn
			r = 0
			rr <- 1000 + x
			match x
				1
					rr <- 5
				3
					rr <- 6
				_
					rr <- 119
		`, "rr", Anis(1001, 5, 1002, 119, 1003, 6, 1004, 119)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
