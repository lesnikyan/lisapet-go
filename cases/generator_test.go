package cases

import "testing"

/*

ok 1. list comprehension
ok 2. dict comprehension
3. generator (: r; iter)
4. generator as an arg of constructor:
	4.1 list,
	4.2 tuple,
	4.3 dict(tuples),
	4.4 string(int|byte),
	4.5 bytes(int|byte|bytes)
*/

func TestGenDict(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		ss = "aa bb cc dd".split(' ')
		r = {k: n ; k, n <- ss, [3 .. 5]}
		`, "r", adk(dk{"aa": 3, "bb": 4, "cc": 5})},
		{`
		ss = "aa/11 bb/22 cc/33 dd/44".split(' ')
		r = {k: v ; n <- ss; k, v = n.split('/')}
		`, "r", adk(dk{"aa": "11", "bb": "22", "cc": "33", "dd": "44"})},
		{`
		ss = "aa bb cc dd".split(' ')
		r = {'%s%s>' << (k, c): n ; k, n <- ss, [3 .. 5]; c <- '~ ='.split(' ')}
		`, "r", adk(dk{"aa=>": 3, "aa~>": 3, "bb=>": 4, "bb~>": 4, "cc=>": 5, "cc~>": 5})},
		{`
		k1 = 'a b c'.split(' ')
		k2 = 'd e f'.split(' ')
		k3 = 'x y z'.split(' ')
		nn = [1,2,3]
		r = { k: n ; a, b, c, n <- k1, k2, k3, nn; k = '%s:%s:%s' << (a, b, c)}
		`, "r", adk(dk{"a:d:x": 1, "b:e:y": 2, "c:f:z": 3})},
		{`
		k1 = 'a b c'.split(' ')
		k2 = 'd e f'.split(' ')
		k3 = 'x y z'.split(' ')
		nn = [1,2,3]
		base = 0
		func foo()
			t = base
			base += 1
			return t
		r = { k: n + foo() ; a <- k1; b <- k2; c, n <- k3, nn; k = '%s:%s:%s' << (a, b, c)}
		`, "r", adk(dk{
			"a:d:x": 1, "a:d:y": 3, "a:d:z": 5, "a:e:x": 4, "a:e:y": 6, "a:e:z": 8, "a:f:x": 7, "a:f:y": 9, "a:f:z": 11,
			"b:d:x": 10, "b:d:y": 12, "b:d:z": 14, "b:e:x": 13, "b:e:y": 15, "b:e:z": 17, "b:f:x": 16, "b:f:y": 18, "b:f:z": 20,
			"c:d:x": 19, "c:d:y": 21, "c:d:z": 23, "c:e:x": 22, "c:e:y": 24, "c:e:z": 26, "c:f:x": 25, "c:f:y": 27, "c:f:z": 29})},

		// {``, "r", adk(dk{})},
		// {``, "r", adk(dk{})},
		// {``, "r", adk(dk{})},
		// {``, "r", adk(dk{})},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestComprInCompr(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = [k ;
			x <- [y;
				y <- [1..10]; y % 2  != 0
			];
			p <- [11, 12, 13]; k = x + p; k %2 == 1]
		`, "r", Anis(13, 15, 17, 19, 21)},
		{`
		r = [ x;
			x <- [1 .. 10]; x % 3 == 1]
		`, "r", Anis(1, 4, 7, 10)},
		{`
		t = [[y; y <- [x, x + 1, -x]]; x <- [5 .. 7]]
		r = t.flat()
		`, "r", Anis(5, 6, -5, 6, 7, -6, 7, 8, -7)},
		{`
		r = {k: v; _, 
			n <- [(a, b); 
				a, b <- ['aa','bb','cc'], [11, 12, 13]]; 
			k, v = n}
		`, "r", adk(dk{"aa": 11, "bb": 12, "cc": 13})},
		// {``, "r",  Anis()},
		// {``, "r",  Anis()},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

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
