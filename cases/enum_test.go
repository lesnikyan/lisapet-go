package cases

import (
	"testing"
)

// TODO: test of error if trying set new value to const

func TestGrupDefStruct(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = []
		grup Abc
			num = 1
			struct Q a: int
		#
		a1 = Abc.Q{a:1}
		r = [a1]
		`, "r", Anis(
			Stf("Q", dk{"a": 1}),
		)},
		{`
		r = []
		grup Abc
			num = 1
			struct Q a: int
		#
		a1 = Abc.Q{}
		r = [a1]
		`, "r", Anis(
			Stf("Q", dk{"a": 0}),
		)},
		{`
		grup Abc
			struct Q a: int, b: string
		#
		r = []
		a1 = Abc.Q{}
		r <- a1
		r <- Abc.Q{a:12, b:'b-val'}
		r <- Abc.Q(13, 'c-val')
		`, "r", Anis(
			Stf("Q", dk{"a": 0, "b": ""}),
			Stf("Q", dk{"a": 12, "b": "b-val"}),
			Stf("Q", dk{"a": 13, "b": "c-val"}),
		)},
		{`
		r = []
		grup Abc
			num = 1
			struct Q nn: list
		#
		a1 = Abc.Q{nn:[11,22,33]}
		r = [a1]
		`, "r", Anis(
			Stf("Q", dk{"nn": Anis(11, 22, 33)}),
		)},
		// method in grup
		{`
		r = []
		grup Abc
			num = 100
			struct Q a: int
			func q:Q foo(x)
				x * q.a
		#
		a1 = Abc.Q{a: 3}
		r = [a1.foo(2), a1.foo(5)]
		`, "r", Anis(6, 15)},
		// grup var in method
		{`
		r = []
		grup Abc
			num = 100
			struct Q a: int
			func q:Q foo(x)
				x * q.a + num
		#
		a1 = Abc.Q{a: 2}
		r = [a1.foo(2), a1.foo(5)]
		`, "r", Anis(104, 110)},
		// 2 methods
		{`
		r = []
		grup Abc
			num = 100
			struct Q a: int
			#
			func q:Q foo(x)
				x * q.a + num
			#
			func q:Q bar(s: string)
				[ ~"{s}-{n}";  n <- [1 .. q.a]]
		#
		a1 = Abc.Q{a: 3}
		r = [a1.bar('Ago'), a1.bar('Before')]
		`, "r", Anis(
			Anis("Ago-1", "Ago-2", "Ago-3"),
			Anis("Before-1", "Before-2", "Before-3"),
		)},
		// enum val in method
		{`
		r = []
		grup Abc
			enum N
				a = 100
			struct Q a: int
			func q:Q foo(x)
				x * q.a + N.a
		#
		a1 = Abc.Q{a: 3}
		r = [a1.foo(2), a1.foo(7)]
		`, "r", Anis(106, 121)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestGrupDef(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		grup Abc
			name = 'Vasya'
			num = 123
		#r = []
		#print('>>1 ', Abc.name)
		r = [Abc.name, Abc.num]
		`, "r", Anis("Vasya", 123)},
		// func in grup
		{`
		grup Abc
			func foo(x)
				x * 10
		
		r = [Abc.foo(3), Abc.foo(5)]
		`, "r", Anis(30, 50)},
		{`
		grup Abc
			func foo(x)
				x * 10
			Num = 7
		
		r = [Abc.foo(Abc.Num), Abc.foo(Abc.Num + 2)]
		`, "r", Anis(70, 90)},
		// func uses grup var
		{`
		grup Abc
			func foo(x)
				x * Num
			Num = 7
		
		r = [Abc.foo(2), Abc.foo(3)]
		`, "r", Anis(14, 21)},
		{`
		grup Abc
			func foo(x)
				x * Abc.Num
			Num = 7
		
		r = [Abc.foo(2), Abc.foo(4)]
		`, "r", Anis(14, 28)},
		{`
		grup Abc
			Num = 7
			func foo(x)
				x * Num
		
		r = [Abc.foo(2), Abc.foo(5)]
		`, "r", Anis(14, 35)},
		// enum in grup
		{`
		grup Abc
			enum N
				a = 1
				b = 2
				c = 3
		#
		r = [Abc.N.b, Abc.N.c * 100]
		`, "r", Anis(2, 300)},
		{`
		grup Abc
			Num = 5
			Num2 = 50
			enum N
				a = 1
				b = 2
				c = 3
			func foo(x)
				x * Num
		#
		r = [Abc.foo(Abc.N.b) + 1000, 2000 + Abc.foo(Abc.N.c + Abc.Num2)]
		`, "r", Anis(1010, 2265)},
		{`
		grup Abc
			Num = 100
			enum N
				a = 1
				c = 3
			func foo(x)
				x * N.c + Num
		#
		r = [Abc.foo(3), Abc.foo(7), ]
		`, "r", Anis(109, 121)},
		// const in grup
		{`
		grup Abc
			const A = 11
			const B = 12
		#
		r = [Abc.A, Abc.B + 100]
		`, "r", Anis(11, 112)},
		// prevent reassign var in upper context
		{`
		r = [1]
		grup Abc
			const A = 11
			const B = 12
			r = [123]
		#
		r  <- Abc.A
		`, "r", Anis(1, 11)},
		// grup const, same name var upper
		{`
		r = [1]
		A = 2
		grup Abc
			const A = 11
			const B = 12
			r = [123]
		#
		r <- A
		r  <- Abc.A
		`, "r", Anis(1, 2, 11)},
		// list
		{`
		grup Abc
			nn = [11,22,33]
		#
		r = [Abc.nn[0]]
		r <- Abc.nn[2]
		`, "r", Anis(11, 33)},
		// tuple
		{`
		grup Abc
			nn = (10, 20, 'Cucumber')
		#
		r = [Abc.nn[0]]
		r <- Abc.nn[2]
		`, "r", Anis(10, "Cucumber")},
		// dict
		{`
		grup Abc
			cc = {'red':0xff0000, 'green':0x00ff00, 'blue':0x0000ff}
			days = {4:'Thu', 5:'Fri', 6:'Sat', 7:'Sun'}
				1: 'Mon'
				2: 'Tue'
				3: 'Wed'
		#
		r = [Abc.cc['red']]
		r <- Abc.days[2]
		r <- Abc.days[7]
		`, "r", Anis(0xff0000, "Tue", "Sun")},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestConstVal(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		const name = 'Vasya'
		const num = 123
		r = []
		r <- num
		r <- name
		`, "r", Anis(123, "Vasya")},
		// tuple can be const
		{`
		const x = (1, 2, "3")
		r = list(x)
		`, "r", Anis(1, 2, "3")},
		{`
		const ONE = 1
		const TWO = 2
		const THREE = 3.0
		func foo()
			return THREE
		r = [ONE, TWO, foo()]
		`, "r", Anis(1, 2, 3.0)},
		{`
		const ONE = 1
		const TWO = 2
		const THREE = 3.0
		func foo(x )
			return THREE * x
		r = [ONE + 10, TWO * 3, foo(5.7)]
		`, "r", Anis(11, 6, 17.1)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestEnum(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		enum Num zero, one, two, three, four , five
		r = []
		r <- Num.zero
		r <- Num.one
		r <- Num.four
		`, "r", Anis(0, 1, 4)},
		{`
		enum Num zero, one, two, three, four , five
		r = 0
		x = 3
		y = 4
		if Num.three == x
			if y < Num.five
				r = 10 
		`, "r", int64(10)},
		{`
		enum Num zero, one, two, three, four , five
		r = 0
		x = [Num.four]
		r = x[Num.zero]
		`, "r", int64(4)},
		{`
		enum Num zero, one, two, three, four , five
		r = 0
		x = (Num.four, Num.two)
		r = x[Num.zero] * 100 + x[Num.one]
		`, "r", int64(402)},
		{`
		enum Num zero, one, two, three, four , five
		r = Num.three * Num.five + 100 * Num.two
		`, "r", int64(215)},
		{`
		enum Num zero, one, two, three, four , five
		#
		func foo(x)
			Num.two + x
		#
		r = foo(Num.three * 100)
		`, "r", int64(302)},
		{`
		enum Num zero, one, two, three, four , five
		#
		func foo(x)
			Num.two + x
		#
		func bar(y)
			Num.five * y * 100
		r = foo(bar(Num.three))
		`, "r", int64(1502)},
		{`
		enum Colors red=0xff0000, green=0x00ff00, blue=0x0000ff
		r = 'colors: %06x, %06x, %06x' << (Colors.red, Colors.green, Colors.blue)
		`, "r", "colors: ff0000, 00ff00, 0000ff"},
		{`
		enum Num zero, one, two, three, four , five
		enum N10 decim = 10, cent=100, mille=1000
		r = 0
		x = {Num.four : Num.five * N10.decim}
		r = x[Num.four]
		`, "r", int64(50)},
		{`
		# counter afdter explicit val
		enum A a, b, c=10, d, e
		r = 0
		r = [A.a, A.c, A.e]
		`, "r", Anis(0, 10, 12)},
		{`
		enum Colors
			red=0xff0000
			green=0x00ff00
			blue=0x0000ff
			white=0xffffff
			black=0
			orange=0xff8800
		r = 'colors2: %06x, %06x, %06x' << (Colors.orange, Colors.white, Colors.red)
		`, "r", "colors2: ff8800, ffffff, ff0000"},

		{`
		enum Num zero, one, two, three, four , five
		enum N10 decim = 10, cent=100, mille=1000
		r = [Num.two .. N10.decim]
		`, "r", Anis(2, 3, 4, 5, 6, 7, 8, 9, 10)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
