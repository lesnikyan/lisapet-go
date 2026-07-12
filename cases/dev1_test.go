package cases

import (
	"fmt"
	"testing"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/lesnikyan/lisapet-go/nodes"
	obb "github.com/lesnikyan/lisapet-go/objects"
	par "github.com/lesnikyan/lisapet-go/parser"
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
		{"null", obb.Null{}, w},
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
			exp := &nodes.ValExpr{Val: tt.exp}
			assert.Equal(t2, exp, res)
		})
	}
}

func TestOperSplit2(t *testing.T) {
	tdata := []struct {
		src    string
		lowest int
		others []ints2
	}{
		// {"1 + 2", 1, []ints2{}},
		// {"n = 0xff", -1, []ints2{}},
		// {"11 + 33 - 44", 6, []ints2{}},
		// {"1 + 5 - 4 / 6", -1, []ints2{}},
		// {"1 + 5 * 6 - 8 / 2 ** 2 * 2 ^/ 49", -1, []ints2{}},
		// {"2 + (4 - 6)", -1, []ints2{}},
		// {"(4 - 6) + 5", -1, []ints2{}},
		// {"(3 + 5) * (4 - 6)", -1, []ints2{}},
		// {"4 * (4 - 6) + 5", -1, []ints2{}},
		// {"5 + 6 - 17 + 22", -1, []ints2{}},
		// {"6 * 5 / 10 * 11", -1, []ints2{}},
		// {"age = 2000 - 1955", -1, []ints2{}},
		// {"r = (1+2)*(3-4)/(5**6)", -1, []ints2{}},
		// {"a, b, c = 1, 2, 3", -1, []ints2{}},
		// {"r = (a * b + c; a,b <- aa, bb ;c = a + b)", -1, []ints2{}},
		// {"r = (: a * b + c; a,b <- aa, bb ;c = a + b)", -1, []ints2{}},
		// {"{'a': 123, 'b':4+5, 'c':60}", -1, []ints2{}},
		// {`[1, 2, nn..., -33, ~"{n}"]`, -1, []ints2{}},
		// {"aaa.bbb.ccc = 123", -1, []ints2{}},
		// {"obj.foo ~>", -1, []ints2{}},
		// {"f1 = obj.mem.foo ~>", -1, []ints2{}},
		// {"foo(1)", -1, []ints2{}},
		// {"obj.foo(1, 2)", -1, []ints2{}},
		// {"obj.foo ~> (123)", -1, []ints2{}},
		// {"foo ~> (123)", -1, []ints2{}},
		// {"ff ~> (1)(2)", -1, []ints2{}},
		// {"oob.foo ~> (1)(2)", -1, []ints2{}},
		// {"[1,2,3][4](5)", -1, []ints2{}},
		// {"[1,2,3]...", -1, []ints2{}},
		// {"re`[0-9]`Li", -1, []ints2{}},
		// {"0x[1e 2f]", -1, []ints2{}},
		// {"(x,y) -> x + y", -1, []ints2{}},

		// {"\\x,y -> x + y", -1, []ints2{}}, // TODO: resolve slash-leading lambda expression

		// {"", -1, []ints2{}},
		// {"", -1, []ints2{}},
		// {"", -1, []ints2{}},
		// {"", -1, []ints2{}},
		// {"", -1, []ints2{}},
		// {"", -1, []ints2{}},
		// {"", -1, []ints2{}},
		// {"", -1, []ints2{}},

		// Note: definitions and other Left-keywords is not a cases (fn = func(arg) - is possibly exception, not sure)
		// {"func Name(a:int, b:string)", -1, []ints2{}},
		// {"func inst:StrType Name(a:itn, b:list)", -1, []ints2{}},
	}
	for _, tt := range tdata {
		sctx := par.SplitContext{}
		line := par.SplitLine([]rune(tt.src), sctx)
		res, err := Line2tree(line, nil)
		t.Log("tt1>", tt.src, res, err)
		ltree := res.Tree
		PrintONode(ltree, 0)
	}
}

func TestCaseOfAssignSimpleVal(t *testing.T) {
	tdata := []struct {
		src     string
		resCase base.Expression
	}{
		// {"a = 123", &nodes.OperBin{}},
		// {"b = 'Hello1'", &nodes.OperBin{}},
		// {"c = 12.25", &nodes.OperBin{}},
		// {"d = x1 + 234", &nodes.OperBin{}},
		// {"dd: int = 1 + 4", &nodes.OperBin{}},
		// {"ee:int = 2 * 5 + (a - b) / c - 3 ** 2", &nodes.OperBin{}},
		// {"e2 = 1*2 + 2 ** 3", &nodes.OperBin{}},
		{"f = (1/(1/(1/(1/2))))", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("%s >>", tt.src), func(t2 *testing.T) {
			sctx := par.SplitContext{}
			line := par.SplitLine([]rune(tt.src), sctx)
			res, err := Line2tree(line, nil)
			ltree := res.Tree
			operTree := ltree.rightNode
			t.Log("tt1>", tt.src, res, err, "r-oper:", operTree.oper)
			PrintONode(operTree, 0)
			expr, ok := ProcExprTree(operTree)
			assert.True(t, ok)
			fmt.Println("tt2>", nodes.OperArgsInfo(expr))
		})
	}
}

func TestCaseSeqVal(t *testing.T) {
	tdata := []struct {
		src     string
		resCase base.Expression
	}{
		// {"(1,2,3,4,)", &nodes.OperBin{}},
		// {"[1,2,3,4]", &nodes.OperBin{}},
		// {"{'a':11, 'b':22, 'c':33}", &nodes.OperBin{}},
		// {"[[1,2,3]]", &nodes.OperBin{}},
		{"[[1,2,3], (4,5,6), {7:'Q7', 8:'Q8'}]", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("%s >>", tt.src), func(t2 *testing.T) {
			sctx := par.SplitContext{}
			line := par.SplitLine([]rune(tt.src), sctx)
			res, err := Line2tree(line, nil)
			ltree := res.Tree
			operTree := ltree.rightNode
			t.Log("tt1>", tt.src, res, err, "r-oper:", operTree.oper)
			PrintONode(operTree, 0)
			expr, ok := ProcExprTree(operTree)
			assert.True(t, ok)
			tp := fmt.Sprintf("%T", expr)
			fmt.Println("tt2>", tp, nodes.OperArgsInfo(expr))
		})
	}
}

func TestSimpleExpressionsDo(t *testing.T) {
	tdata := []struct {
		src     string
		resCase base.Expression
	}{
		// {"a = 101", &nodes.OperBin{}},
		// {"a = 2 + 3", &nodes.OperBin{}},
		// {"a = 5 * (3 + 4)", &nodes.OperBin{}},
		// {"a = `Hello, 1 2!`", &nodes.OperBin{}},
		// {"a = 1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9", &nodes.OperBin{}},
		// {"a = 1 * 2 * 3 * 4 * 5 * 6 * 7 * 8 * 9", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("ExprDo, %s >>", tt.src), func(t2 *testing.T) {

			sctx := par.SplitContext{}
			line := par.SplitLine([]rune(tt.src), sctx)
			res, err := Line2tree(line, nil)
			ltree := res.Tree
			operTree := ltree.rightNode
			t.Log("tt1>", tt.src, res, err, "r-oper:", operTree.oper)
			PrintONode(operTree, 0)

			expr, ok := ProcExprTree(operTree)
			assert.True(t, ok)
			tp := fmt.Sprintf("%T", expr)
			fmt.Println("tt2>", tp, nodes.OperArgsInfo(expr))

			ctx := obb.NewContext(nil)
			expr.Do(ctx)
			vr := ctx.GetVar("a")
			fmt.Println("tt3>", vr, vr.Name, vr.Val)
		})
	}
}

func TestVarInMathDo(t *testing.T) {
	tdata := []struct {
		src string
		res any
	}{
		{`
		b = 1
		a = b + 2`, int64(3)},
		{`
		b = 1
		c = 3
		d = 5
		a = b + c * d`, int64(16)},
		{`
		b = 7
		c = 4
		a = (2 + 3) * (b - c) * -2
		`, int64(-30)},
		{` a = -2 * (-3) * -(3 - - 1)`, int64(-24)},
		{`
		n = "Hello "
		m = 'example'
		a = n + m + '!'
		`, "Hello example!"},
		// {``, &nodes.OperBin{}},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("ExprDo, %s >>", tt.src), func(t2 *testing.T) {
			clines := par.SplitCode(tt.src[1:])
			ctx := obb.NewContext(nil)
			block, err := TreeBlock(clines)
			assert.Nil(t, err)
			block.Do(ctx)
			vr := ctx.GetVar("a")
			assert.Equal(t2, tt.res, vr.Val)
			fmt.Println("tt3>", vr, vr.Name, vr.Val)
		})
	}
}

func TestIfElseCase(t *testing.T) {
	tdata := []struct {
		src string
		res any
	}{
		// {`
		// a = 1
		// `, int64(1)},
		// {`
		// a = 1
		// if a == 1
		// 	a = 5
		// `, int64(5)},
		// {`
		// n = 12
		// a = 5
		// if n == 13
		// 	b = 11
		// 	a = a + b
		// `, int64(16)},
		// {`
		// a = 1
		// if a == 1
		// 	a = 2
		// b = 10
		// if a != 1
		// 	a = a + b
		// `, int64(12)},
		// {`
		// a = 1
		// if b = 2; a == 1
		// 	a = b + 10
		// `, int64(12)},
		// {`
		// a = 1
		// if a == 2
		// 	a = 3
		// else
		// 	a = 4
		// `, int64(4)}, // else
		// {`
		// a = 2
		// if a == 2
		// 	a = 3
		// else
		// 	a = 4
		// `, int64(3)}, // if
		// {`
		// a = 5
		// if a == 2
		// 	a = 3
		// else if a == 5
		// 	a = 4
		// `, int64(4)},
		// {`
		// a = 1
		// if a == 1 || a == 2
		// 	a = 3
		// `, int64(3)},
		// {`
		// a = 1
		// if a == 1
		// 	if a == 1
		// 		if a == 1
		// 			if a == 1
		// 				a = 7
		// `, int64(7)},
		// {`
		// a = 1
		// if a == 1
		// 	a = 2
		// 	if a == 2
		// 		a = 3
		// 		if a == 3
		// 			a = 4
		// 			if a == 4
		// 				a = 5
		// `, int64(5)},
		// {`
		// a = 1
		// b = 10
		// if a == 2
		// 	a = 20
		// else
		// 	if a == 3
		// 		a = 30
		// 	else
		// 		a = 112
		// `, int64(112)},
		{`
		a = 1
		b = 10
		if a == 2
			a = 20
		else if a == 3
			a = 30
		else
			a = 111
		`, int64(1)},
		// {``, int64(1)},
		// {``, int64(1)},
		// {``, int64(1)},
		// {``, int64(1)},
		// {``, int64(1)},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("ExprDo, %s >>", tt.src), func(t2 *testing.T) {
			clines := par.SplitCode(tt.src[1:])
			block, err := TreeBlock(clines)
			assert.Nil(t, err)
			t.Log("--- --- --- Do ...")
			ctx := obb.NewContext(nil)
			block.Do(ctx)
			vr := ctx.GetVar("a")
			t.Log("tt#vr", vr)
			assert.Equal(t2, tt.res, vr.Val)
			fmt.Println("tt3>", vr, vr.Name, vr.Val)
		})
	}
}

/*
Case in feature of multi-line expressions
*/
func TestCaseUnclosedBrackets(t *testing.T) {
	tdata := []struct {
		src     string
		resCase base.Expression
	}{
		{"a + b *( 1 - ", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
		// {"", &nodes.OperBin{}},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("Unclosed, %s >>", tt.src), func(t2 *testing.T) {

		})
	}
}
