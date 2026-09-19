package cases

import (
	"fmt"
	"testing"

	obb "github.com/lesnikyan/lisapet-go/objects"
	par "github.com/lesnikyan/lisapet-go/parser"
	"github.com/stretchr/testify/assert"
)

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
