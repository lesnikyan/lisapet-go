package cases

import (
	"fmt"
	"testing"

	obb "github.com/lesnikyan/lisapet-go/objects"
	par "github.com/lesnikyan/lisapet-go/parser"
	"github.com/stretchr/testify/assert"
)

func RunTCodeErr(t *testing.T, i int, tt TTst) {
	t.Run(fmt.Sprintf("Test, %d) %s >>", i, TrunkString(tt.src)), func(t2 *testing.T) {
		clines := par.SplitCode(tt.src[1:])
		block, err := TreeBlock(clines)

		// parsing, interpretation errors
		assert.Nil(t, err) // check err
		// t.Log("--- --- --- Do ...")
		ctx := obb.NewContext(nil)
		PreloadContext(ctx)

		// runtime errors
		terr := block.Do(ctx)
		if terr != nil {
			assert.Fail(t2, terr.Error())
			return
		}
	})
}

func _TestSyntaxErrors(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{

		{`
		r = 1
		x = 2
		if x % 2 = 0
			r = 3
		`, "r", Anis()},

		{`
		r = [ x; x <- [5 .. 6]; x % 3 = 1]
		`, "r", Anis()},
		// {``, "r",  Anis()},
		// {``, "r",  Anis()},
		// {``, "r",  Anis()},
	}
	for i, tt := range tdata {
		// TODO: run code with error, check error type
		RunTCodeErr(t, i, tt)
	}
}
