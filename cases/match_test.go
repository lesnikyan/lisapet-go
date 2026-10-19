package cases

import "testing"

func TestMatchTupleOptional(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		nn = [(1,), (2,), (2,3), (4,5), (5,6,7), (4,5,6,7), (4,5,6,7,88), (4,5,6,7,8,9), 5, null]
		r = []
		for x <- nn
			match x
				(2,?)
					r <- 222
				(a,?)
					r <- 100 + a
				(?,?,b)
					r <- 200 + b
				(a,?,?,?, b)
					r <- 1000 * a + b
				(?)
					r <- 333
				::tuple
					r <- 98
				_
					r <- 99
		`, "r", Anis(101, 222, 222, 104, 207, 4007, 4088, 98, 99, 99)},
		{`
		x = (4, 5)
		r = 0
		match x
			(*)
				r = 544
		`, "r", int64(544)},
		{`
		nn = [(1,), (1,2,3,4,5), (2,), (2,3), (4,5), (4,5,6), (4,4,4,4,4,4,4,4,5,66), (5,6,7), (4,5,6,7), (1,3,1,3,1,7), (4,5,6,7,8,9), 5, null]
		r = []
		for x <- nn
			match x
				(1,*)
					r <- 111
				(*,2)
					r <- 222
				(*,5,a)
					r <- 5000 + a
				(*,3,*,3,*,a)
					r <- 3000 + a
				(a,*,7,*,b)
					r <- 70000 + a*1000 + b
				(*,_,a)
					r <- 8000 + a
				(*)
					r <- 999
				_
					r <- 99
		`, "r", Anis(111, 111, 222, 8003, 8005, 5006, (5066), 8007, 8007, 111, 74009, 99, 99)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestMatchTuple(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// simple and basic case: vals, vars
		{`
		x = (4, 5)
		r = 0
		match x
			(4, 5)
				r = 555
		`, "r", int64(555)},
		{`
		x = (4, 5)
		r = 0
		match x
			(4, a)
				r = 500 + a
		`, "r", int64(505)},
		// one val
		{`
		x = (4,)
		r = 0
		match x
			(4,)
				r = 14
			_
				r = 99
		`, "r", int64(14)},
		// one var
		{`
		x = (17,)
		r = 0
		match x
			(4,)
				r = 14
			(a,)
				r = 100 + a
			_
				r = 99
		`, "r", int64(117)},
		// empty
		{`
		x = (,)
		r = 0
		match x
			(,)
				r = 100
			_
				r = 99
		`, "r", int64(100)},
		{`
		# empty tuple mp, no comma
		x = (,)
		r = 0
		match x
			()
				r = 1000
			_
				r = 99
		`, "r", int64(1000)},
		// vars
		{`
		nn = [(1,), (2,), tuple([3]), 5, null, (4,5), (5,6,7)]
		r = []
		for x <- nn
			match x
				(2,)
					r <- 222
				(a,)
					r <- 100 + a
				(a, b)
					r <- 1000 * a + b
				_
					r <- 99
		`, "r", Anis(101, 222, 103, 99, 99, 4005, 99)},
		// underscore sub pattern
		{`
		# empty tuple mp, no comma
		x = (4,)
		r = 0
		match x
			(1,)
				r = 100
			(_)
				r = 222
			_
				r = 99
		`, "r", int64(222)},
		{`
		nn = [(1), (2,3), (5, 1, 12), (7, 2, 21)]
		r = []
		for x <- nn
			match x
				(a, _)
					r <- 100 + a
				(7, a, b)
					r <- a * 100 + b
				(a, _, b)
					r <- a * 100 + b
				_
					r <- 99
		`, "r", Anis(99, 102, 512, 221)},
		{`
		nn = [(1), (2,3), (5, 1, 12), (7, 2, 21), (4,5,6,7)]
		r = []
		for x <- nn
			match x
				(a, _)
					r <- 100 + a
				(7, a, b)
					r <- a * 100 + b
				(a, _, b)
					r <- a * 100 + b
				(_, _, a, b)
					r <- a * 1000 + b
				_
					r <- 99
		`, "r", Anis(99, 102, 512, 221, 6007)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestMatchListStar(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// simple *
		{`
		x = [11]
		r = 0
		match x
			[*]
				r = 111
		`, "r", int64(111)},
		{`
		x = []
		r = 0
		match x
			[*]
				r = 110
		`, "r", int64(110)},
		{`
		nn = [[1], [2,3], [5, 1, 11], [1 .. 10], [1..100]]
		r = []
		for x <- nn
			match x
				[*]
					r <- len(x)
				_
					r <- 99
		`, "r", Anis(1, 2, 3, 10, 100)},
		{`
		nn = [[], [1], [2,3], [5, 1, 11], [7, 1, 11], [7, 1, 1,1,1,1], [8 .. 18], [9..100]]
		r = []
		for i, x <- nn
			match x
				[7, *]
					r <- -(7000 + len(x))
				[a, *]
					r <- a * 1000 + len(x)
				[*]
					r <- 100 + len(x)
				_
					r <- 99
		`, "r", Anis(100, 1001, 2002, 5003, -7003, -7006, 8011, 9092)},
		{`
		nn = [[], [11], [21,2], [31,2,3], [41,2,3,4], [51,2,3,4,5], [61,2,3,4,5,6]]
		r = []
		for i, x <- nn
			match x
				[*, a]
					r <- 1200+a
				:: list
					r <- 98
				_
					r <- 99
		`, "r", Anis(98, 1211, 1202, 1203, 1204, 1205, 1206)},
		{`
		nn = [[1], [2,4], [5, 1, 4], [7, 2, 21], [8,,1,1,1,1,1,1,1,22], [9,23]]
		r = []
		for x <- nn
			match x
				[a, *, 4]
					r <- 1000 + a
				[a, *, b]
					r <- a* 10000 + b
				_
					r <- 99
		`, "r", Anis(99, 1002, 1005, 70021, 80022, 90023)},
		{`
		nn = [[1], [2,4], [5, 1, 4], [7, 2, 21], [8,,1,1,1,1,1,1,1,22], [9,2,23], [9,2,3,24], [9,2,3,4,4,4,25], [], 10]
		r = []
		for x <- nn
			match x
				[a, *, 4]
					r <- 1000 + a
				[a, _, *, b]
					r <- a* 10000 + b
				[]
					r <- 97
				_
					r <- 99
		`, "r", Anis(99, 1002, 1005, 70021, 80022, 90023, 90024, 90025, 97, 99)},
		{`
		x = [3,2,1]
		r = 0
		match x
			[*, 1, *, 1, *, a]
				r = 111
			[a, *, 2, *,  b]
					r = 333
		`, "r", int64(333)},
		{`
		x = [1,5,1,6,1,7]
		r = 0
		match x
			[*, 1, *, 1, *, a]
				r = 111
		`, "r", int64(111)},
		{`
		nn = [[1..5], [1..10], [1,5,1,6,1,7], [2,1, 3, 1, 4], [3,4,5,11,12,17], [5,5,18], [1,1,3,4,5,6,11], , [1,2,3], [3,2,1], 10]
		r = []
		for x <- nn
			match x
				[*, 1, *, 1, *, a]
					r <- 1000 + a
				[a, *, 2, *,  b]
					r <- a* 10000 + b
				[*, 5, *, a]
					r <- 50000 + a
				_
					r <- 99
		`, "r", Anis(10005, 10010, 1007, 1004, 50017, 50018, 1011, 10003, 30001, 99)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestMatchListOptional(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// simple ?
		{`
		x = [11]
		r = 0
		match x
			[?]
				r = 111
		`, "r", int64(111)},
		{`
		x = []
		r = 0
		match x
			[?]
				r = 110
		`, "r", int64(110)},
		{`
		x = [1,2]
		r = 0
		match x
			[?]
				r = 110
			[?,?]
				r = 120
		`, "r", int64(120)},
		{`
		x = [4, 5, 6]
		r = 0
		match x
			[4, 5, ?]
				r = 114
		`, "r", int64(114)},
		{`
		x = [4, 5, 6]
		r = 0
		match x
			[4, ?, 6]
				r = 115
		`, "r", int64(115)},
		{`
		x = [4, 5, 6]
		r = 0
		match x
			[?, 5, ?]
				r = 116
		`, "r", int64(116)},
		{`
		x = [4, 5, 6, 7, 8, 9]
		r = 0
		match x
			[?, 5, ?, ?, ?, ?]
				r = 117
		`, "r", int64(117)},
		{`
		x = [4, 5]
		r = 0
		match x
			[?, 5, ?, ?, ?, ?]
				r = 118
		`, "r", int64(118)},
		{`
		#nn = [[1,2,3,4], [1,2,4], [1,3,4], [1,4,5,6,7], [1,4,8]]
		nn = [[4], [1,4], [1,2,3,4],  [1,2,4], [1,3,4],[1,4,5,6,7], [1,4,8], [1,1,1,4,4,4,4], [1,1,1,4,4,4,4,5]]
		r = []
		for i, x <- nn
			match x
				[1, ?, ?, 4, ?, ?, ?]
					r <- 1000*i + len(x)
				[?,4,?]
					r <- 98
				_
					r <- 99
		`, "r", Anis(98, 1002, 2004, 3003, 4003, 5005, 6003, 7007, 99)},
		{`
		nn = [[], [1], [1,2], [1,2,3], [1,2,3,4], [1,2,3,4,5], [1,2,3,4,5,6], [1 .. 10]]
		r = []
		for i, x <- nn
			match x
				[?]
					r <- 11
				[?,?]
					r <- 12
				[?,?,?]
					r <- 13
				[?,?,?,?]
					r <- 14
				[?,?,?,?,?,?,?,?]
					r <- 15
				:: list
					r <- 98
				_
					r <- 99
		`, "r", Anis(11, 11, 12, 13, 14, 15, 15, 98)},

		{`
		x = [5]
		r = 0
		match x
			[?,a]
				r = 1100 + a
		`, "r", int64(1105)},
		{`
		x = [2, 3]
		r = 0
		match x
			[?,a]
				r = 1200 + a
		`, "r", int64(1203)},
		{`
		x = [2, 3, 4]
		r = 0
		match x
			[?,a]
				r = 1200 + a
			[?,?,a]
				r = 1300 + a
		`, "r", int64(1304)},
		{`
		x = [3,4,5]
		r = 0
		match x
			[?,a]
				r = 1200 + a
			[?,a,b]
				r = 1000 + a*100 + b
		`, "r", int64(1405)},
		{`
		nn = [[], [1], [1,2], [1,2,3], [1,2,3,4], [1,2,3,4,5], [1,2,3,4,5,6]]
		r = []
		for i, x <- nn
			match x
				[?,a]
					r <- 12
				[?,?,a]
					r <- 100 + a
				[?,?,?,a]
					r <- 200 + a
				[?,?,?,?,?,?,?,a]
					r <- 300 + a
				:: list
					r <- 98
				_
					r <- 99
		`, "r", Anis(98, 12, 12, 103, 204, 305, 306)},
		{`
		nn = [[], [11], [21,2], [31,2,3], [41,2,3,4], [51,2,3,4,5], [61,2,3,4,5,6]]
		r = []
		for i, x <- nn
			match x
				[a, ?]
					r <- 1200+a
				[a,?,?]
					r <- 1000 + a
				[a,?,?,?]
					r <- 2000 + a
				[a,?,?,?,?,?,?,?]
					r <- 3000 + a
				:: list
					r <- 98
				_
					r <- 99
		`, "r", Anis(98, 1211, 1221, 1031, 2041, 3051, 3061)},
		{`
		nn = [[], [2,1,1,11], [1,1,1], [1,2,1,2], [2,2,2,1,1,3], [2,2,2,2,1,2,1,4], [1,2,2,1,2,2,5], 
			[2,2,2,1,3,3,1,4,4,6], [2,2,2,1,3,3,3,1,4,7]]
		r = []
		for i, x <- nn
			match x
				[?,?,?, 1,?,?, 1,?,?, a]
					r <- 1000 + a
				[?,?,?,?,?,?,1,?, 1,?, a]
					r <- 2000 + a
				[?,?,?,1,?,?,?,?,?, 1,?, a]
					r <- 3000 + a
				:: list
					r <- 98
		`, "r", Anis(98, 1011, 1001, 1002, 1003, 2004, 1005, 1006, 3007)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

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
		// underscore sub pattern
		{`
		nn = [[1], [2,3], [5, 1, 11], [7, 2, 21]]
		r = []
		for x <- nn
			match x
				[a, _]
					r <- 100 + a
				[7, a, b]
					r <- a * 100 + b
				[a, _, b]
					r <- a * 100 + b
				_
					r <- 99 
		`, "r", Anis(99, 102, 511, 221)},
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
