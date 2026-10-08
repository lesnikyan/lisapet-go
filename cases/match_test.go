package cases

import "testing"

func TestMatch(t *testing.T) {
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
