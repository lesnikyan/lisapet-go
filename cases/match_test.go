package cases

import "testing"

func TestMatchList(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// simple and basic case: vals, vars
		{`
		x = [4, 5, 6, 7]
		r = 0
		match x
			[4, 5, a, b]
				r = 1000 + a * 10 +b 
		`, "r", int64(1067)},
		// one val
		{`
		x = [4]
		r = 0
		match x
			[4]
				r = 14 
			_
				r = 99 
		`, "r", int64(14)},
		// one var
		{`
		x = [17]
		r = 0
		match x
			[4]
				r = 14 
			[a]
				r = 100 + a 
			_
				r = 99 
		`, "r", int64(117)},
		// empty list
		{`
		x = []
		r = 0
		match x
			[]
				r = 100
			_
				r = 99
		`, "r", int64(100)},
		// one var
		{`
		nn = [[1], [2], [3], 5, null, [4,5], [5,6,7]]
		r = []
		for x <- nn
			match x
				[2]
					r <- 222 
				[a]
					r <- 100 + a 
				[a, b]
					r <- 1000 * a + b 
				_
					r <- 99 
		`, "r", Anis(101, 222, 103, 99, 99, 4005, 99)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestMatchMultitype(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// no brackets
		{`
		x = 23
		r = 0
		match x
			a :: int | float
				r = a + 100
		`, "r", int64(123)},
		// in brackets
		{`
		x = 22.4
		r = 0
		match x
			a :: (int | float)
				r = a + 100
		`, "r", float64(122.4)},
		{`
		nn = [-88, 61.5, 0, 22, g'Q', 00x41, null]
		r = []
		for x <- nn
			match x
				a :: (int | float)
					r <- a + 100
				b :: (list | tuple)
					r += list(b)
				c :: (byte | int | glif)
					r <- glif(c)
				s :: (string | bytes)
					r <- string(s)
				_
					r <- 99
		`, "r", Anis(12, 161.5, 100, 122, 'Q', 'A', 99)},
		{`
		struct A a: int
		struct B a: int
		struct C c: int
		struct D c: int, d : int
		struct E(D) e: int
		nn = [A(11), B(12), C(13), D(14, 24), E(15, 25, 16), null, none, 22]
		r = []
		for x <- nn
			match x
				a :: (A | B)
					r <- a.a 
				c :: C | D
					r <- c.c
					match c
						e :: E
							r <- e.e
						d :: D
							r <- d.d
				:: null | maybe
					r <- 98
				_
					r <- 99
		`, "r", Anis(11, 12, 13, 14, 24, 15, 16, 98, 98, 99)},
		{`
		nn = [101, 00x22, g'W', 1.4, "Hello", true, null, re'[0-9]']
		r = []
		ss = "aa2bb4cc"
		for x <- nn
			match x
				a :: int | float | byte | glif | bool | string | regexp
					match a
						b :: regexp
							r += b.split(ss)
							continue
					r <- a
				_
					r <- 99000099
		`, "r", Anis(101, byte(0x22), 'W', 1.4, "Hello", true, 99000099, "aa", "bb", "cc")},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

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
