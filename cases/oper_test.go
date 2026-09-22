package cases

import (
	"testing"

	"github.com/lesnikyan/lisapet-go/objects"
)

/*
ok 1. ?>
ok 2. !?>
ok 3. a ? b : c
ok 4. a ?: b
ok 5. @!
ok 5.1 @defined(varname)
ok 6. a :: type
ok 7. v : int|float
8. a :: int|float
*/

func TestCheckTypeStruct(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		struct A a:int
		struct B b: string
		struct C (A) c: int
		#
		#ss = [A{}]
		ss = [A{}, B{}, C{}, null]
		r = []
		for i, x <- ss
			r <- i
			r <- x
			r <- x :: int
			r <- x :: A
			r <- x :: B
			r <- x :: C
		#
		r <- null :: null
		`, "r",
			// Anis(Stf("A", dk{"a": 0}), false, Stf("B", dk{"b": ""}), false, Stf("C", dk{"a": 0, "c": 0}), false, Tnull(), false)},
			Anis(0, Stf("A", dk{"a": 0}), false, true, false, false,
				1, Stf("B", dk{"b": ""}), false, false, true, false,
				2, Stf("C", dk{"a": 0, "c": 0}), false, true, false, true,
				3, Tnull(), false, false, false, false, true)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestCheckType(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		x = 1
		r = [x]
		r <- x :: int
		r <- x :: float
		r <- x :: string
		`, "r", Anis(1, true, false, false)},
		{`
		x = 1.2
		r = [x]
		r <- x :: int
		r <- x :: float
		`, "r", Anis(1.2, false, true)},
		{`
		x = 1.2
		r = [x]
		r <- x :: int
		r <- x :: float
		`, "r", Anis(1.2, false, true)},
		{`
		x = 'abc'
		r = [x]
		r <- x :: string
		r <- x :: glif
		`, "r", Anis("abc", true, false)},
		{`
		x = g'a'
		r = [x]
		r <- x :: string
		r <- x :: glif
		`, "r", Anis('a', false, false)},
		{`
		x = 00xa
		r = [x]
		r <- x :: byte
		r <- x :: int
		r <- x :: glif
		r <- x :: bytes
		`, "r", Anis(byte(0xa), true, false, false, false)},
		{`
		x = 0x[a]
		r = [x]
		r <- x :: bytes
		r <- x :: byte
		r <- x :: string
		`, "r", Anis(objects.Bytes{0xa}, true, false, false)},
		{`
		x = [1,2]
		r = [x]
		r <- x :: list
		r <- x :: tuple
		r <- x :: dict
		r <- x :: maybe
		`, "r", Anis(Anis(1, 2), true, false, false, false)},
		{`
		x = (1,2)
		r = [x]
		r <- x :: list
		r <- x :: tuple
		r <- x :: dict
		r <- x :: maybe
		`, "r", Anis(Tanis(1, 2), false, true, false, false)},
		{`
		x = {1: 2}
		r = [x]
		r <- x :: list
		r <- x :: tuple
		r <- x :: dict
		r <- x :: maybe
		`, "r", Anis(adk(dk{1: 2}), false, false, true, false)},
		{`
		x = some(5)
		r = [x]
		r <- x :: list
		r <- x :: tuple
		r <- x :: dict
		r <- x :: maybe
		`, "r", Anis(Tmay(5), false, false, false, true)},
		{`
		x = none
		r = [x]
		r <- x :: list
		r <- x :: tuple
		r <- x :: dict
		r <- x :: maybe
		`, "r", Anis(Tmay(nil), false, false, false, true)},
		{`
		r = [11]
		r <- 11 :: int
		r <- 1 :: float
		r <- 1 :: maybe
		r <- 1 :: bool
		`, "r", Anis(11, true, false, false, false)},
		{`
		func f1()
			1
		r = ['func']
		r <- f1 :: function
		r <- f1 :: null
		r <- f1 :: maybe
		r <- f1 :: bool
		`, "r", Anis("func", true, false, false, false)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestAtDelete(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		a = true
		r = [] 
		r <- @defined(a)
		@! a
		r <- @defined(a)
		`, "r", Anis(true, false)},
		{`
		a = [1,2,3]
		@! a[1]
		r = a
		`, "r", Anis(1, 3)},
		{`
		a = [1, 2, 3, 4, 5, 6, 7]
		for i <- [5, 3, 2]
			@! a[i]
		r = a
		`, "r", Anis(1, 2, 5, 7)},
		{`
		a = {1:11, 2:22, 3:33, 4:44}
		@! a[2]
		r = a
		`, "r", adk(dk{1: 11, 3: 33, 4: 44})},
		{`
		a = {'aa':1, 'bb':2, 'cc':3}
		@! a['aa']
		r = a
		`, "r", adk(dk{"bb": 2, "cc": 3})},
		{`
		a = {'aa':1, 'bb':2, 'cc':3, 'dd': 4, 'ee': 5, 'ff': 6}
		#
		for k <- ['bb', 'cc', 'ee']
			@! a[k]
		r = a
		`, "r", adk(dk{"aa": 1, "dd": 4, "ff": 6})},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestOperElvis(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		a = true
		r = a ?: "no"
		`, "r", true},
		{`
		a = false
		r = a ?: "no2"
		`, "r", "no2"},
		{`
		a = ''
		r = a ?: "no3"
		`, "r", "no3"},
		{`
		a = 'yes4'
		r = a ?: "no3"
		`, "r", "yes4"},
		{`
		a = 10
		r = a ?: "no"
		`, "r", int64(10)},
		{`
		a = 0
		r = a ?: 111
		`, "r", int64(111)},
		{`
		a = []
		r = a ?: ['no5']
		`, "r", Anis("no5")},
		{`
		a = (,)
		r = a ?: ('no6',)
		`, "r", Tanis("no6")},
		{`
		a = (12,)
		r = a ?: ('no6',)
		`, "r", Tanis(12)},
		{`
		struct A a: int
		#
		a = null
		b = A{a:5}
		r = []
		r <- a ?: 'no7'
		r <- b ?: 113
		`, "r", Anis("no7", Stf("A", dk{"a": 5}))},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestOperTernary(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		a = true
		r = a ? "yes" : "no"
		`, "r", "yes"},
		{`
		a = false
		r = a ? "yes" : "no"
		`, "r", "no"},
		{`
		a = 1
		b = 2
		r = a < b ? "yes 1" : "no 2"
		`, "r", "yes 1"},
		{`
		a = 3
		b = 2
		r = a < b ? "yes 1" : "no 2"
		`, "r", "no 2"},
		{`
		a = 3
		b = 2
		r = [a < b ? "yes 1" : "no 2", a > b ? "yes 1" : "no 2", 55]
		`, "r", Anis("no 2", "yes 1", 55)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestOperIn(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// in list
		{`
		a = 2
		nn = [1,2,3]
		r = a ?> nn
		`, "r", true},
		{`
		a = 5
		nn = [1,2,3]
		r = a ?> nn
		`, "r", false},
		{`
		nn = [1,2,3]
		r = 2 !?> nn
		`, "r", false},
		{`
		nn = [1,2,3]
		r = 5 !?> nn
		`, "r", true},
		// in tuple
		{`
		a = 2
		nn = (1,2,3)
		r = a ?> nn
		`, "r", true},
		{`
		a = 5
		nn = (1,2,3)
		r = a ?> nn
		`, "r", false},
		{`
		nn = (1,2,3)
		r = 2 !?> nn
		`, "r", false},
		{`
		nn = (1,2,3)
		r = 5 !?> nn
		`, "r", true},
		// in dict
		{`
		a = 'a'
		dd = {'a':1, 'b':2}
		r = a ?> dd
		`, "r", true},
		{`
		a = 'X'
		r = a ?> {'a':1, 'b':2}
		`, "r", false},
		{`
		dd = {'a':1, 'b':2}
		r = 'Y' !?> dd
		`, "r", true},
		{`
		r = 5 !?> {'a':1, 'b':2}
		`, "r", true},
		// in maybe
		{`
		a = 2
		mm = some(2)
		r = a ?> mm
		`, "r", true},
		{`
		a = 5
		mm = some(2)
		r = a ?> mm
		`, "r", false},
		{`
		a = 5
		mm = none
		r = a ?> mm
		`, "r", false},
		{`
		mm = some(2)
		r = 2 !?> mm
		`, "r", false},
		{`
		mm = some(2)
		r = 5 !?> mm
		`, "r", true},
		{`
		mm = none
		r = 5 !?> mm
		`, "r", true},
		{`
		nn = [1,2,3]
		nums = [1, 5, 10, 3] 
		r = []
		for n <- nums
			if n ?> nn
				r <- n
			else if n !?> nn
				r <- n * 100
		`, "r", Anis(1, 500, 1000, 3)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestOperUnclosedBrackets(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = (1 +
		10 - 5)
		`, "r", int64(6)},
		{`
		r = (100 + 2 *
		10 - 2 **
		3)
		`, "r", int64(112)},
		{`
		r = [1, 2,
		33]
		`, "r", Anis(1, 2, 33)},
		{`
		r = (22, 33,
		44,
		55)
		`, "r", Tanis(22, 33, 44, 55)},
		{`
		r = {
			'a': 11,
			'b': 12
		}
		`, "r", adk(dk{"a": 11, "b": 12})},
		{`
		func foo(a, b, c)
			a + b + c
		r = foo(1000,
		200,
		34)
		`, "r", int64(1234)},
		{`
		func foo(a,
			b)
			a + b
		r = foo(1000, 200)
		`, "r", int64(1200)},
		{`
		func foo(a, 
			b,
			c=1)
			a + b + c
		r = foo(1000, 200)
		`, "r", int64(1201)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestOperUnpackCollection(t *testing.T) {
	// ok 1. a, b = list
	// ok 2. a, b = tuple
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

// func TestDev1(t *testing.T) {
// 	tdata := []struct {
// 		src   string
// 		vname string
// 		res   any
// 	}{
// 		{`
// 		a = 1
// 		b = 2
// 		r = a + b
// 		`, "r", int64(1234)},
// 	}
// 	for i, tt := range tdata {
// 		RunTCodeVarExp(t, i, tt)
// 	}
// }
