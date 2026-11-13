package cases

import (
	"fmt"
	"testing"

	"github.com/lesnikyan/lisapet-go/nodes"
	"github.com/stretchr/testify/assert"
)

/*
s +
c
d +
x +
b
o
f
e
g
%
0 +
.
, _
+ -
<
>
^

*/

func _TestFSFormatIntO(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		x = 0o7015
		r = ~"x1: {x:o}"
		`, "r", "x1: 10101"},
		// {`
		// x = -0b10101
		// r = ~"x1: {x:b}"
		// `, "r", "x1: -10101"},
		// {`
		// nn = [0o, 0b0, 0b1, -0b1, 2**31-3]
		// r = []
		// for x <- nn
		// 	r <- ~"x1: {x:o}"
		// 	r <- ~"x2: {x:012o}"
		// 	r <- ~"x3: {x:+011o}"
		// 	r <- ~"x4: {x:<+011o}"
		// 	r <- ~"x5: {x:>011o}"
		// 	r <- ~"x6: {x:^011o}"
		// 	r <- ~"x7: {x:>11o}"
		// `, "r", Anis(

		// )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFSFormatIntB(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		x = 0b10101
		r = ~"x1: {x:b}"
		`, "r", "x1: 10101"},
		{`
		x = -0b10101
		r = ~"x1: {x:b}"
		`, "r", "x1: -10101"},
		{`
		nn = [0b1001, -0b1001, 0b11001100, -0b1111111100000000, 0b0, 0b1, -0b1, 2**31-3]
		r = []
		for x <- nn
			r <- ~"x1: {x:b}"
			r <- ~"x2: {x:012b}"
			r <- ~"x3: {x:+011b}"
			r <- ~"x4: {x:<+011b}"
			r <- ~"x5: {x:>011b}"
			r <- ~"x6: {x:^011b}"
			r <- ~"x7: {x:>11b}"
		`, "r", Anis(
			"x1: 1001", "x2: 000000001001", "x3: +0000001001", "x4: +1001000000", "x5: 00000001001", "x6: 00010010000", "x7:        1001",
			"x1: -1001", "x2: -00000001001", "x3: -0000001001", "x4: -1001000000", "x5: 000000-1001", "x6: 000-1001000", "x7:       -1001",
			"x1: 11001100", "x2: 000011001100", "x3: +0011001100", "x4: +1100110000", "x5: 00011001100", "x6: 01100110000", "x7:    11001100",
			"x1: -1111111100000000", "x2: -1111111100000000", "x3: -1111111100000000", "x4: -1111111100000000",
			"x5: -1111111100000000", "x6: -1111111100000000", "x7: -1111111100000000",
			"x1: 0", "x2: 000000000000", "x3: +0000000000", "x4: +0000000000", "x5: 00000000000", "x6: 00000000000", "x7:           0",
			"x1: 1", "x2: 000000000001", "x3: +0000000001", "x4: +1000000000", "x5: 00000000001", "x6: 00000100000", "x7:           1",
			"x1: -1", "x2: -00000000001", "x3: -0000000001", "x4: -1000000000", "x5: 000000000-1", "x6: 0000-100000", "x7:          -1",
			"x1: 1111111111111111111111111111101", "x2: 1111111111111111111111111111101", "x3: +1111111111111111111111111111101",
			"x4: +1111111111111111111111111111101", "x5: 1111111111111111111111111111101", "x6: 1111111111111111111111111111101", "x7: 1111111111111111111111111111101")},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFSFormatIntX(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		x = 2357
		r = ~"x1: {x:x}"
		`, "r", "x1: 935"},
		{`
		x = 0xfe0d
		r = ~"x1:_#{x:6x};"
		`, "r", "x1:_#  fe0d;"},
		{`
		x = -0xfe0d
		r = ~"x1:_#{x:6x};"
		`, "r", "x1:_# -fe0d;"},
		{`
		x = 0xfe0d
		r = ~"x1:_#{x:+6x};"
		`, "r", "x1:_# +fe0d;"},
		{`
		x = 0xfa07
		r = ~"x1: {x:06x};"
		`, "r", "x1: 00fa07;"},
		{`
		x = 0xfe0d
		r = ~"x1:_#{x:<6x};"
		`, "r", "x1:_#fe0d  ;"},
		{`
		x = 0xfa07
		r = ~"x1: {x:<06x};"
		`, "r", "x1: fa0700;"},
		{`
		x = 0xfa07
		r = ~"x1: {x:^06x};"
		`, "r", "x1: 0fa070;"},
		{`
		x = -0xfa07
		r = ~"x1: {x:08x};"
		`, "r", "x1: -000fa07;"},
		{`
		x = -0xfa07
		r = ~"x1: {x:<08x};"
		`, "r", "x1: -fa07000;"},
		{`
		x = -0xfa07
		r = ~"x1: {x:>08x};"
		`, "r", "x1: 000-fa07;"},
		{`
		x = 0xfa07
		r = ~"x1: {x:<+08x};"
		`, "r", "x1: +fa07000;"},
		{`
		x = -0xfe0d
		r = ~"x1:_# {x:06x};"
		`, "r", "x1:_# -0fe0d;"},
		{`
		x = 0xfa07
		r = ~"x1: {x:>+08x};"
		`, "r", "x1: 000+fa07;"},
		{`
		x = 0xfa07
		r = ~"x1: {x:^+07x};"
		`, "r", "x1: 0+fa070;"},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFSFormatIntD(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		x = 2357
		r = ~"x1: {x:d}"
		`, "r", "x1: 2357"},
		{`
		x = 2351
		r = ~"x1:_{x:6d};"
		`, "r", "x1:_  2351;"},
		{`
		x = 2352
		r = ~"x1: {x:06d};"
		`, "r", "x1: 002352;"},
		{`
		x = 2352
		r = ~"x1:_{x:<6d};"
		`, "r", "x1:_2352  ;"},
		{`
		x = 2352
		r = ~"x1:_{x:<06d};"
		`, "r", "x1:_235200;"},
		{`
		x = 2352
		r = ~"x1:_{x:^08d};"
		`, "r", "x1:_00235200;"},
		{`
		x = -2353
		r = ~"x1: {x:d};"
		`, "r", "x1: -2353;"},
		{`
		x = -2353
		r = ~"x1:_{x:8d};"
		`, "r", "x1:_   -2353;"},
		{`
		x = -2353
		r = ~"x1:_{x:<8d};"
		`, "r", "x1:_-2353   ;"},
		{`
		x = -2353
		r = ~"x1: {x:>08d};"
		`, "r", "x1: 000-2353;"},
		{`
		x = -2353
		r = ~"x1: {x:<08d};"
		`, "r", "x1: -2353000;"},
		{`
		x = -2353
		r = ~"x1: {x:^09d};"
		`, "r", "x1: 00-235300;"},
		{`
		x = 2353
		r = ~"x12: {x:+d};"
		`, "r", "x12: +2353;"},
		{`
		x = 2353
		r = ~"x1: {x:-d};"
		`, "r", "x1: 2353;"},
		{`
		x = 2353
		r = ~"x1: {x:<+08d};"
		`, "r", "x1: +2353000;"},
		{`
		x = 2353
		r = ~"x1: {x:>+08d};"
		`, "r", "x1: 000+2353;"},
		{`
		x = 2353
		r = ~"x1: {x:^+09d};"
		`, "r", "x1: 00+235300;"},
		{`
		x = 2353
		r = ~"x1: {x:>+8d};"
		`, "r", "x1:    +2353;"},
		{`
		x = 2353
		r = ~"x1: {x:^+9d};"
		`, "r", "x1:   +2353  ;"},
		{`
		x = 00xf5
		r = ~"x1:{x:+8d};"
		`, "r", "x1:    +245;"},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFSFormatDefault(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		x = 2357
		r = ~"x1: {x}"
		`, "r", "x1: 2357"},
		{`
		x = "Abel"
		r = ~"x2: {x}"
		`, "r", "x2: Abel"},
		{`
		x = -12
		r = ~"x3: {x}"
		`, "r", "x3: -12"},
		{`
		x = g'Q'
		r = ~"x6: {x}"
		`, "r", "x6: Q"},
		{`
		x = true
		r = ~"x7: {x}"
		`, "r", "x7: true"},
		{`
		x = false
		r = ~"x8: {x}"
		`, "r", "x8: false"},
		{`
		x = none
		r = ~"x9: {x}"
		`, "r", "x9: none"},
		{`
		x = some(12)
		r = ~"x10: {x}"
		`, "r", "x10: some(12)"},
		{`
		x = 00xa7
		r = ~"x5: {x}"
		`, "r", "x5: a7"},
		// TODO: list, tuple, dict. Thinking
		// {`
		// x = 12.5
		// r = ~"x4: {x}"
		// `, "r", "x4:-?-"},
		// {`
		// x = 12
		// r = ~"x11: {x}"
		// `, "r", "-?-"},
		// {`
		// x = 12
		// r = ~"x12: {x}"
		// `, "r", "-?-"},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestFormSubFormat(t *testing.T) {
	tdata := []struct {
		val  any
		form string
		res  any
	}{
		{1, "<+5.4f", ""},
		{1, ">08.4f", ""},
		{1, "+06d", ""},
		{1, "+08x", ""},
		{1, ">06x", ""},
	}
	for i, tt := range tdata {
		t.Run(fmt.Sprintf("Test, %d) %s >>", i, tt.form), func(t2 *testing.T) {
			opt, err := nodes.FormOption(tt.form)
			assert.Nil(t2, err)
			assert.NotNil(t2, opt)
			// fmt.Printf("# opt.Type: %c \n", opt.Type)
		})
	}
}

func TestStringFormatSimple(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		x = 'abc'
		yyy = 3
		bb = 5
		r = ~"str: {x} : {yyy + bb}"
		`, "r", "str: abc : 8"},
		{`
		x = 'abc'
		name = 'Vasya'
		r = ~'''--- x ---
		?123 {name}'''
		`, "r", "--- x ---\n?123 Vasya"},
		{`
		x = 10.5
		yy = 5.0
		r = ~"nums: {x}, {yy}, {x + yy}, {x * yy}, {x / yy}"
		`, "r", "nums: 10.500000, 5.000000, 15.500000, 52.500000, 2.100000"},
		{`
		x = true
		yy = false
		r = ~"vals: {x}, {yy}, {x || yy}, {x && yy}, {!yy}, {x == yy}"
		`, "r", "vals: true, false, true, false, true, false"},
		{`
		cc = [g'A', g'Q', g'@', g'λ', g'Σ']
		c = cc[1]
		r = ~"vals: {cc[1]}"
		`, "r", "vals: Q"},
		{`
		cc = [g'A', g'Q', g'@', g'λ', g'Σ']
		r = ~"vals: {cc[0]}, {cc[1]}, {cc[2]}, {cc[3]}, {cc[4]}, {glif(0x03b8)}"
		`, "r", "vals: A, Q, @, λ, Σ, θ"},
		{`
		cc = [g'A', g'Q', g'@', g'λ', g'Σ']
		s = "QRYL"
		r = ~"vals: {cc[0]}, {cc[1]}, {cc[2]}, {cc[3]}, {cc[4]}, {glif(0x03b8)}, {s[3]}"
		`, "r", "vals: A, Q, @, λ, Σ, θ, L"},
		{`
		cc = [g'A', g'Q', g'@', g'λ', g'Σ']
		s = "QRYL"
		r = []
		for c <- cc
			r <- ~"!{c}!"
		`, "r", Anis("!A!", "!Q!", "!@!", "!λ!", "!Σ!")},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
