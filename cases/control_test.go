package cases

import (
	"fmt"
	"testing"

	obb "github.com/lesnikyan/lisapet-go/objects"
	par "github.com/lesnikyan/lisapet-go/parser"
	"github.com/stretchr/testify/assert"
)

// if,for,while expr /: expr; expr
// if expr; expr /: if expr /: expe; expr
func TestInlineControl(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = 1
		a = 1
		# if a == 1 /: if b < 0 /: for n <- nn /: for i=0; i<5; i += 1 /: a = 10; b = i + 6
		#if a == 1 /: for n <- nn /: a = 10; b = i + 6
		if a > 0 /: r = 5
		`, "r", int64(5)},
		{`
		r = 1
		a = 10
		if a > 0 /: if a < 20 /: r = 6
		`, "r", int64(6)},
		{`
		r = []
		a = 10
		if a > 0 /: for n <- [1 .. 5] /: r <- n
		`, "r", Anis(1, 2, 3, 4, 5)},
		{`
		r = []
		a = 10
		if a > 0 /: for n <- [1 .. 10] /: if n % 2 == 0 /: r <- n
		`, "r", Anis(2, 4, 6, 8, 10)},
		{`
		r = []
		a = 13
		for n <- [1 .. 20] /: if n < a /: if n % 2 == 0 /: r <- n
		`, "r", Anis(2, 4, 6, 8, 10, 12)},
		{`
		r = []
		a = 13
		for n <- [5 .. 20] /: if n < a /: if n % 2 == 0 /: b = n + 100; r <- b
		`, "r", Anis(106, 108, 110, 112)},
		{`
		r = []
		a = 13
		if nn = [3 .. 20] ; true /: for n <- nn /: if n < a /: if n % 2 != 0 /: b = n + 100; r <- b
		`, "r", Anis(103, 105, 107, 109, 111)},
		{`
		r = []
		a = 13
		if nn = [8 .. 15] ; true /: for n <- nn /: for m <- [1, 2, 5] /: if n % 2 != 0 /: b = n + m * 100; r <- b
		`, "r", Anis(109, 209, 509, 111, 211, 511, 113, 213, 513, 115, 215, 515)},
		{`
		r = []
		a = 13
		if nn = [8 .. 15] ; true /: for n <- nn /: for b <- [n + m * 200; m <- [1, 2, 5]; n % 2 != 0] /:  r <- b
		`, "r", Anis(209, 409, 1009, 211, 411, 1011, 213, 413, 1013, 215, 415, 1015)},
		{`
		r = []
		a = 13
		if nn = [8 .. 15]; i = 0 ; true /: while i < len(nn) /:  r <- nn[i]; i += 1
		`, "r", Anis(8, 9, 10, 11, 12, 13, 14, 15)},
		{`
		r = []
		a = 13
		if nn = [8 .. 15]; i = 0 ; true /: while i < len(nn) /: if c = i; i += 1; c % 3 == 0 /: r <- c;  r <- nn[c]
		`, "r", Anis(0, 8, 3, 11, 6, 14)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

// expr ; expr ; expr
func TestBlockInline(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = 10
		a = 1; b = 2;c = a + b 
		r += c
		`, "r", int64(13)},
		{`
		r = 10
		a = 1; b = 2; r += (a + b)
		`, "r", int64(13)},
		{`
		r = (a = 1; b = 2; a + b)
		`, "r", int64(3)},
		{`
		r = (
			a = 10;
			b = 6;
			a + b
		)
		`, "r", int64(16)},
		{`
		func foo(x)
			x + 10
		r = (a = 5; foo(a))
		`, "r", int64(15)},
		{`
		nn = [10]
		r = (a = 11; nn <- a)
		`, "r", Anis(10, 11)},
		{`
		nn = [10]
		r = (nn <- 12; nn - [0])
		`, "r", int64(10)},
		{`
		nn = [10]
		r = (nn<-1; nn<-2; nn<-3; nn<-4; nn <- 5)
		`, "r", Anis(10, 1, 2, 3, 4, 5)},
		{`
		r = (b = 2; a = 1; a *= b; a *= b; a *= b; a *= b; a *= b; a *= b; a *= b; a)
		`, "r", int64(128)},
		{`
		r = (a = 1; a *= 2; a *= 3; a *= 4; a *= 5; a *= 6; a *= 7; a)
		`, "r", int64(5040)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestControlIfNot(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = 1
		fff = false
		if ! fff
			r = 2
		else
			r = 3
		#
		`, "r", int64(2)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
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
		{`
		a = 1
		if a == 1
			a = 5
		`, int64(5)},
		{`
		n = 12
		a = 5
		if n == 12
			b = 11
			a = a + b
		`, int64(16)},
		{`
		a = 1
		if a == 1
			a = 2
		b = 10
		if a != 1
			a = a + b
		`, int64(12)},
		{`
		a = 1
		if b = 2; a == 1
			a = b + 10
		`, int64(12)},
		{`
		a = 1
		if a == 2
			a = 3
		else
			a = 4
		`, int64(4)}, // else
		{`
		a = 2
		if a == 2
			a = 3
		else
			a = 4
		`, int64(3)}, // if
		{`
		a = 5
		if a == 2
			a = 3
		else if a == 5
			a = 4
		`, int64(4)},
		{`
		a = 1
		if a == 1 || a == 2
			a = 3
		`, int64(3)},
		{`
		a = 1
		if a == 1
			if a == 1
				if a == 1
					if a == 1
						a = 7
		`, int64(7)},
		{`
		a = 1
		if a == 1
			a = 2
			if a == 2
				a = 3
				if a == 3
					a = 4
					if a == 4
						a = 5
		`, int64(5)},
		{`
		a = 1
		b = 10
		if a == 2
			a = 20
		else
			if a == 3
				a = 30
			else
				a = 112
		`, int64(112)},
		{`
		a = 1
		b = 10
		if a == 2
			a = 20
		else if a == 3
			a = 30
		else
			a = 111
		`, int64(111)},
		{`
		a = 3
		b = 10
		if a == 2
			a = 20
		else if a == 3
			a = 31
		else
			a = 111
		`, int64(31)},
		{`
		a = 7
		if a == 1
			a = 11
		else if a == 2
			a = 12
		else if a == 3
			a = 13
		else if a== 4
			a = 14
		else if a== 5
			a = 15
		else if a== 6
			a = 16
		else if a== 7
			a = 17
		else 
			a = 20
		`, int64(17)},
		{`
		a = 3
		b = 10
		if a == 1
			a = 11
		else if a == 3
			a = 13
			if b == 10
				a = 23
			else
				a = 24
		else 
			a = 12
		`, int64(23)},
		// {``, int64(1)},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("ExprDo, %s >>", tt.src), func(t2 *testing.T) {
			clines := par.SplitCode(tt.src[1:])
			block, err := TreeBlock(clines)
			assert.Nil(t, err)
			// t.Log("--- --- --- Do ...")
			ctx := obb.NewContext(nil)
			block.Do(ctx)
			vr := ctx.GetVar("a")
			// t.Log("tt#vr", vr)
			assert.Equal(t2, tt.res, vr.Val)
			// fmt.Println("tt3>", vr, vr.Name, vr.Val)
		})
	}
}
