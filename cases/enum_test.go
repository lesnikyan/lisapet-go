package cases

import (
	"testing"
)

func TestMatchEnum(t *testing.T) {
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
		enum Colors
			red=0xff0000
			green=0x00ff00
			blue=0x0000ff
			white=0xffffff
			black=0
			orange=0xff8800
		r = 'colors2: %06x, %06x, %06x' << (Colors.orange, Colors.white, Colors.red)
		`, "r", "colors2: ff8800, ffffff, ff0000"},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
