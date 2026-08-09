package cases

import "testing"

/*
1. math expr: (...\n...)
2. tuple, list, dict constructor
3. control expr sum-expression: if, for, while, func def, func call
4. generator and comprehension
*/
func TestOperUnclosedBrackets(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}

}

func TestOperUnpackCollection(t *testing.T) {
	// 1. a, b = list
	// 2. a,b = tuple
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		nn = [11, 12]
		a, b = nn
		r = [a, b]
		`, "r", Anis(11, 12)},
		{`
		nn = [11, 12, 13]
		a, b, c = nn
		r = [a, b, c]
		`, "r", Anis(11, 12, 13)},
		{`
		a, b, c = [11, 12, 13]
		r = [a, b, c]
		`, "r", Anis(11, 12, 13)},
		{`
		nn = (21, 22)
		a, b = nn
		r = [a, b]
		`, "r", Anis(21, 22)},
		{`
		nn = (21, 22, 23)
		a, b, c = nn
		r = [a, b, c]
		`, "r", Anis(21, 22, 23)},
		{`
		# long multiassign
		ss = split('a1,b1,c1,d1,e1,f1,g1,h1,i1,j1,k1,l1,m1,n1,o1,p1', ',')
		a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p = ss
		r = [a,b,c,d,e,f,g,h, '>>', i,j,k,l,m,n,o,p]
		`, "r", Anis("a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
			">>", "i1", "j1", "k1", "l1", "m1", "n1", "o1", "p1")},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestOperMultiAssign(t *testing.T) {
	// ok 1. a, b = 1, 2 # vals
	// ok 2. a, b = v1, v2 # vars
	// ok 3. a, b = f1(), f2() # func calls
	// ok 5. a, b = foo() # multi result from function
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		a, b = 1, 2
		r = [a,b]
		`, "r", Anis(1, 2)},
		{`
		a, b, c = 3, 4, 5
		r = [a, b, c]
		`, "r", Anis(3, 4, 5)},
		{`
		q = 6
		p = 7
		t = 8
		a, b, c = q, p, t
		r = [a, b, c]
		`, "r", Anis(6, 7, 8)},
		{`
		func f1(x)
			x + 10
		#
		func f2(x)
			x * 2
		#
		a, b, c = f1(2), f2(7), f2(23)
		r = [a, b, c]
		`, "r", Anis(12, 14, 46)},
		{`
		func foo()
			11, 22, 55
		#
		r = 11112
		a, b, c = foo()
		r = [a, b, c]
		`, "r", Anis(11, 22, 55)},
		{`
		# func multi result
		func foo(x)
			x , x + 10, x * 2
		#
		a, b, c = foo(15)
		r = [a, b, c]
		`, "r", Anis(15, 25, 30)},
		{`
		# func multi return
		func foo(x)
			if x < 0
				return x, x * -2
			x, x * 5
		#
		r = []
		a,b  = foo(3)
		r <- a
		r <- b
		c,d = foo(-7)
		r <- c
		r <- d
		`, "r", Anis(3, 15, -7, 14)},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestOperMath(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{

		{`
		r = []
		r <- 10 + 1
		r <- 10 * 2
		r <- 10 - 3
		r <- 10 / 4
		r <- 2 ** 5
		r <- 2 ^/ 9
		r <- 5 % 2
		r <- 17 % 6
		`, "r", Anis(11, 20, 7, 2.5, 32, 3.0, 1, 5)},

		{`
		r = []
		r <- 1 << 1
		r <- 1 << 2
		r <- 1 << 3
		r <- 0b10 | 0b01
		r <- 0b1111 & 0b0111
		r <- 0b100 >> 2
		r <- 0b11111111 >> 4
		r <- 0b0000 ^ 0b1001
		r <- 0b1010 ^ 0b0101
		`, "r", Anis(0b10, 0b100, 0b1000, 0b11, 0b111, 1, 0b1111, 0b1001, 0b1111)},
		{`
		r = ["int"]
		r <- 1 == 1
		r <- 1 == 2
		r <- 1 != 2
		r <- 1 < 2
		r <- 1 > 2
		r <- "a"
		r <- 1 <= 2
		r <- 1 <= 1
		r <- 1 >= 2
		r <- 1 >= 1
		`, "r", Anis("int", true, false, true, true, false, "a", true, true, false, true)},
		{`
		r = ["float"]
		r <- 1. == 1.
		r <- 1. == 2.
		r <- 1. != 2.
		r <- 1. < 2.
		r <- 1. > 2.
		r <- "a"
		r <- 1. <= 2.
		r <- 1. <= 1.
		r <- 1. >= 2.
		r <- 1. >= 1.
		`, "r", Anis("float", true, false, true, true, false, "a", true, true, false, true)},

		{`
		r = ["other"]
		r <- true == true
		r <- false == false
		r <- true != false
		r <- false != true
		r <- true != true
		r <- false != false
		r <- "a"
		r <- "" == ""
		r <- "a" != ""
		r <- "a" == "a"
		r <- "a" != "b"
		r <- "" != ""
		r <- "a" != "a"
		r <- "a" == "b"
		`, "r", Anis("other", true, true, true, true, false, false,
			"a", true, true, true, true, false, false, false)},
		{`
		r = ["multitype"]
		r <- "a" != 1
		r <- 1 != 'a'
		r <- 1 == 'a'
		r <- 'a' == 1
		`, "r", Anis("multitype", true, true, false, false)},
		{`
		r = ["collection"]
		r <- [] == []
		r <- [1,2,3] ==[1,2,3]
		r <- (,) == (,)
		r <- (1,2) == (1,2)
		r <- {} == {}
		r <- {1:11, 2:22} == {1:11, 2:22}
		r <- "a"
		r <- [] != [1]
		r <- [] == [1]
		r <- (,) != (2,)
		r <- (,) == (3,)
		r <- {} != {3:33}
		r <- {} == {4:44}
		r <- "b"
		r <- 1 != [1]
		r <- 2 != {5:55}
		r <- 3 == (3,'c')
		r <- (3) == (3,)
		`, "r", Anis("collection", true, true, true, true, true, true, "a",
			true, false, true, false, true, false,
			"b", true, true, false, false)},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
