package cases

import "testing"

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
		// {`
		// r = ["multitype"]
		// r <- "a" != 1
		// r <- 1 != 'a'
		// r <- 1 == 'a'
		// r <- 'a' == 1
		// `, "r", Anis("float", true, true, true, true, false, false, "a")},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
