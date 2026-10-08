package cases

import "testing"

func _TestMatchTypeVar(t *testing.T) {
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
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func _TestMatchReturn(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
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
		nn = [1, "quatro"]
		r = []
		for n <- nn
			r <- foo(n)
		`, "r", int64(5)},
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
