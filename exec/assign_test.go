package exec

import (
	"testing"

	ob "github.com/lesnikyan/lisapet-go/objects"
	par "github.com/lesnikyan/lisapet-go/parser"
	tr "github.com/lesnikyan/lisapet-go/tree"
)

func TestAssignValue(t *testing.T) {
	code := `
	x = 5
	y = 1.25
	s = "hello3"
	yes = true
	no = false
	no2 = null
	`
	code = par.Crop(code)
	// parse
	clines := par.SplitCode(code)

	// build tree
	var mod ob.Module
	mod = *tr.Build(clines)

	// make context
	ctx := ob.NewContext(nil)
	print(ctx, mod)
	// execute

	// ver vars from context

	// check values

}

func TestMultiAssign(t *testing.T) {

}

func TestAssignExprRes(t *testing.T) {

}
