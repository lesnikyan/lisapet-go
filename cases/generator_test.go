package cases

import "testing"

/*

ok 1. list comprehension
2. dict comprehension
3. generator (: r; iter)
4. generator as an arg of constructor:
	4.1 list,
	4.2 tuple,
	4.3 dict(tuples),
	4.4 string(int|byte),
	4.5 bytes(int|byte|bytes)
*/

func TestGenListComprNested(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# 2 iters
		nn = [x + y ; x <- [0 .. 5]; y <- [10, 20, 30]]
		`, "nn", Anis(10, 20, 30, 11, 21, 31, 12, 22, 32, 13, 23, 33, 14, 24, 34, 15, 25, 35)},
		{`
		# 3 iters + cond
		nn = [x + y + z * 100 ; x <- [2 .. 5]; x != 5; y <- [10, 11, 12]; z <- (11, 13, 17); x + y != z && (x +y ) % 2 != 0]
		`, "nn", Anis(1113, 1713, 1113, 1713, 1115, 1315, 1715, 1115, 1315, 1715)},
		{`
		#  5 iters
		r = [k ; a <- [1,2,3] ; b <- [10, 11, 12]; c <- [20, 30]; d <- [100, 101]; e <- [1, -1]; k = a + b + c * d * e; a % 2 != 0 && k > -200]
		`, "r", Anis(2011, 2031, 3011, 3041, 2012, 2032, 3012, 3042, 2013, 2033, 3013, 3043, 2013, 2033, 3013, 3043, 2014, 2034, 3014, 3044, 2015, 2035, 3015, 3045)},
		{`
		# 5 iters
		r = [k ; a <- [1,2,3]; a % 2 != 0 ; b <- [10, 11, 12]; c <- [20, 30]; (a + b) % 10 == 1; d <- [100, 101]; e <- [1, -1]; k = a + b + c * d * e; k > 0]
		`, "r", Anis(2011, 2031, 3011, 3041)},
		{`
		#  5 iters
		r = [k ; a <- [1,2,3]; a != 1; b <- [10, 11, 12]; c <- [20, 30]; d <- [100, 101]; d + c > 120; e <- [1, -1]; k = a + b + c * d * e; a % 2 != 0 && k > 0]
		`, "r", Anis(2033, 3013, 3043, 2034, 3014, 3044, 2035, 3015, 3045)},
		{`
		#  5 iters, multiline
		r = [k ; 
			a <- [1,2,3]; a != 1; 
			b <- [10, 11, 14]; m = a + b;
			c <- [20, 30]; 
			d <- [100, 101]; d + c > 120; 
			e <- [1, -1]; k = m + c * d * e; 
			a % 2 != 0 && k > 0]
		`, "r", Anis(2033, 3013, 3043, 2034, 3014, 3044, 2037, 3017, 3047)},
		{`
		# compr to var, use in loop
		nn = [k; a <- [10, 20 .. 50]; b <- [1..9]; k = a + b; k % 2 != 0]
		r = []
		for n <- nn
			if n % 5 != 0
				r <- n
		`, "r", Anis(11, 13, 17, 19, 21, 23, 27, 29, 31, 33, 37, 39, 41, 43, 47, 49, 51, 53, 57, 59)},
		{`
		# conpr in loop
		r = []
		for n <- [k; a <- [10, 20 .. 40]; b <- [1..9]; k = a + b; k % 2 != 0]
			if n % 5 != 0
				r <- n
		`, "r", Anis(11, 13, 17, 19, 21, 23, 27, 29, 31, 33, 37, 39, 41, 43, 47, 49)},
		{`
		r = tuple([x ; x <- [1 .. 10]; x % 2 != 0])
		`, "r", Tanis(1, 3, 5, 7, 9)},
		{`
		ki = 0
		func next()
			n = ki
			ki += 1
			n
		#
		kk = "aa bb cc dd ee ff gg hh ii jj kk ll mm oo nn pp ss tt".split(' ')
		r = dict([(kk[next()], x) ; x <- [1 .. 10]; x % 2 != 0])
		`, "r", adk(dk{"aa": 1, "bb": 3, "cc": 5, "dd": 7, "ee": 9})},
		// {``, "r",  Anis()},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestGenListComprSimple(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# one simple loop by list
		nn = [x ; x <- [11, 12, 5, -2]]
		`, "nn", Anis(11, 12, 5, -2)},
		{`
		# one simple loop by num gen
		nn = [x ; x <- [0 .. 10]]
		`, "nn", Anis(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)},
		{`
		# one loop + guard
		nn = [x ; x <- [0 .. 5]; x > 2]
		`, "nn", Anis(3, 4, 5)},
		{`
		# with sub assignment
		nn = [a ; x <- [1 .. 5]; a = x+10;]
		`, "nn", Anis(11, 12, 13, 14, 15)},
		{`
		# with sub asg and condition
		nn = [a ; x <- [0 .. 5]; a = x+10; x > 2]
		`, "nn", Anis(13, 14, 15)},
		// {``, "r", Anis()},
		// {``, "r", Anis()},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
