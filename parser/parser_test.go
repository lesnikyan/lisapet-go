package parser

import (
	"fmt"
	"strings"
	"testing"

	lang "github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/stretchr/testify/assert"
)

func crop(code string) string {
	code = strings.ReplaceAll(code, "\r\n", "\n")
	if code[0] == '\n' {
		code, _ = strings.CutPrefix(code, "\n")
	}
	return code
}

func TestLines(t *testing.T) {
	code := `
	x = 1
	y = 2
	res = x + y
	`
	code = crop(code)
	lines := Lines(code)
	xss := []string{"x = 1", "y = 2", "res = x + y", ""}
	for i, line := range lines {
		// t.Log(">>", line)
		// exp := xss[i]
		assert.Equal(t, xss[i], line)
	}

}

func TestSplitCode(t *testing.T) {
	code := `
	x = 1
	y = 2
	res = x + y
	str = "Hello Vasya!"
	r = obj.foo(bar([{1:2}]))
	`
	code = crop(code)
	exp := [][]string{
		{"x?^word", " ?^space", "=?^oper", " ?^space", "1?^num"},
		{"y?^word", " ?^space", "=?^oper", " ?^space", "2?^num"},
		{"res?^word", " ?^space", "=?^oper", " ?^space", "x?^word", " ?^space", "+?^oper", " ?^space", "y?^word"},
		{"str?^word", " ?^space", "=?^oper", " ?^space", "\"Hello Vasya!\"?^text"},
		{"r?^word", " ?^space", "=?^oper", " ?^space", "obj?^word", ".?^oper", "foo?^word", "(?^oper", "bar?^word",
			"(?^oper", "[?^oper", "{?^oper", "1?^num", ":?^oper", "2?^num", "}?^oper", "]?^oper", ")?^oper", ")?^oper"},
		{},
	}
	clines := SplitCode(code)
	rts := make([][]string, len(clines))
	for i, cl := range clines {
		rw := make([]string, len(cl.Elems))
		for j, el := range cl.Elems {
			rw[j] = fmt.Sprintf("%s?^%s", el.Text, Lt.TName(el.Type))
		}
		strings.Join(rw, `","`)
		// t.Log(strings.Join(rw, `","`))
		rts[i] = rw
	}
	assert.Equal(t, exp, rts)

}

var num = Lt.Num
var wrd = Lt.Word
var opr = Lt.Oper
var non = Lt.None
var txt = Lt.Text
var qut = Lt.Quot
var esc = Lt.Esc
var spc = Lt.Space
var cmm = Lt.Comm

func TestElemType(t *testing.T) {

	tdata := []struct {
		c    rune
		prev Lt.Lt
		exp  Lt.Lt
	}{
		{'a', non, wrd},
		{'a', num, num},
		{'a', wrd, wrd},
		{'0', non, num},
		{'0', num, num},
		{'0', wrd, wrd},
		{'+', non, opr},
		{'+', num, opr},
		{'+', wrd, opr},
		{'.', non, opr},
		{'.', wrd, opr},
		{'.', txt, txt},
		{'0', txt, txt},
		{'`', non, qut},
		{'\'', non, qut},
		{'"', txt, qut},
		{'\\', non, opr},
		{'\\', txt, esc},
		{'\\', esc, txt},
		{'n', esc, txt},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("e-type '%s', <%s> ", string(tt.c), Lt.TName(tt.exp)), func(t2 *testing.T) {
			res := elemType(tt.c, tt.prev)
			assert.Equal(t2, tt.exp, res, "")
		})
	}

}

type elem = *lang.Elem

type tte struct {
	s string
	t Lt.Lt
}

func telems(ees []tte) []elem {
	var r = make([]elem, len(ees))
	for i, ee := range ees {
		r[i] = &lang.Elem{ee.s, ee.t}
	}
	return r
}

func TestSplitLine(t *testing.T) {
	tdata := []struct {
		cline   string
		pretype Lt.Lt
		strval  string
		exp     []tte
	}{
		{"123", non, "", ([]tte{{"123", num}})},
		{"qwe rty", non, "", ([]tte{{"qwe", wrd}, {" ", spc}, {"rty", wrd}})},
		{"123.05 + 0x9f", non, "", []tte{{"123.05", num}, {" ", spc}, {"+", opr}, {" ", spc}, {"0x9f", num}}},
		{`name = "Vasya Pupkin"`, non, "", []tte{{"name", wrd}, {" ", spc}, {"=", opr}, {" ", spc}, {`"Vasya Pupkin"`, txt}}},
		{"nums = [12, 34, 55.5]", non, "", []tte{{"nums", wrd}, {" ", spc}, {"=", opr}, {" ", spc}, {"[", opr}, {"12", num}, {",", opr},
			{" ", spc}, {"34", num}, {",", opr}, {" ", spc}, {"55.5", num}, {"]", opr}}},
		{"obj.foo()", non, "", []tte{{"obj", wrd}, {".", opr}, {"foo", wrd}, {"(", opr}, {")", opr}}},
		{`ss=["a","9",]`, non, "", []tte{{"ss", wrd}, {"=", opr}, {"[", opr}, {`"a"`, txt}, {",", opr}, {`"9"`, txt}, {",", opr}, {"]", opr}}},
		{"n # qwerty1", non, "", []tte{{"n", wrd}, {" ", spc}, {"# qwerty1", cmm}}},
		{"number1=10 # qwerty %^&*(!)", non, "", []tte{{"number1", wrd}, {"=", opr}, {"10", num}, {" ", spc}, {"# qwerty %^&*(!)", cmm}}},
		{"n=({[f~>()()]}...)#!@#$%^&*()_+=", non, "", []tte{{"n", wrd}, {"=", opr}, {"(", opr}, {"{", opr}, {"[", opr}, {"f", wrd}, {"~>", opr},
			{"(", opr}, {")", opr}, {"(", opr}, {")", opr}, {"]", opr}, {"}", opr}, {"...", opr}, {")", opr}, {"#!@#$%^&*()_+=", cmm}}},
		// {"", non, "", []tte{{"",0}}},
		// {"", non, "", []tte{{"",0}}},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("%s", tt.cline), func(t2 *testing.T) {
			ctx := SplitContext{Ltype: tt.pretype, strval: tt.strval}
			rline := Runes(tt.cline)
			res := SplitLine(rline, ctx)
			exp := telems(tt.exp)
			assert.Equal(t2, exp, res)
		})
	}
}
