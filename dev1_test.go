package main

import (
	"fmt"
	"testing"

	"github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/stretchr/testify/assert"
)

func Compare(a, any, b any) {

}

type wd struct{}

func valEls(src string, expT any) []*lang.Elem {
	var t Lt.Lt = Lt.None
	switch expT.(type) {
	case int:
		t = Lt.Num
	case float32:
		t = Lt.Num
	case string:
		t = Lt.Text
	case wd:
		t = Lt.Word
	}
	res := []*lang.Elem{&lang.Elem{Text: src, Type: t}}
	return res
}

func TestCaseVal(t *testing.T) {
	w := wd{}
	tdata := []struct {
		src string
		exp any
		tm  any
	}{
		{"10", int64(10), 1},
		{"1", int64(1), 1},
		{"0", int64(0), 1},
		{"5", int64(5), 1},
		{"1234567890", int64(1234567890), 1},
		{"12345.625", float64(12345.625), 1},
		{"0x123fab", int64(0x123fab), 1},
		{"0o10", int64(0o10), 1},
		{"0o12370", int64(0o12370), 1},
		{"0b1010", int64(0b1010), 1},
		{"0x001111", int64(0x001111), 1},
		{"0012345", int64(12345), 1},
		{"Hello!", "Hello!", ""},
		{"true", true, w},
		{"false", false, w},
		{"null", Null{}, w},
		// {"", 0},
		// {"", 0},
		// {"", 0},
		// {"", 0},
	}

	for _, tt := range tdata {
		t.Run(fmt.Sprintf("%s >> %v", tt.src, tt.exp), func(t2 *testing.T) {
			els := valEls(tt.src, tt.tm)
			res, ok := CaseVal(els)
			assert.True(t2, ok)
			exp := &ValExpr{Val: tt.exp}
			assert.Equal(t2, exp, res)
		})
	}
}
