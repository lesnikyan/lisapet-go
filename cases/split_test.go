package cases

import (
	"fmt"
	"testing"

	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	par "github.com/lesnikyan/lisapet-go/parser"
	"github.com/stretchr/testify/assert"
)

// test blank
func _TestSplitMatchCase(t *testing.T) {
	tdata := []struct {
		src   string
		opstr string
		ttype Lt.Lt
	}{
		{`var @ list`,
			"@", Lt.Oper},
	}
	for i, tt := range tdata {
		t.Run(fmt.Sprintf("Test, %d) %s >>", i, TrunkString(tt.src)), func(t2 *testing.T) {
			// println("")
			clines := par.SplitCode(tt.src)
			cline := clines[0]
			var prevtree *LineTree
			// fmt.Println(">>>>", cline.Src, "Elems:", FPrintElems(cline.Elems))

			res, err := Line2tree(cline.Elems, prevtree)
			if err != nil {
				t.Logf("Err i:%d Line2tree : %v ", i, err)
			}
			assert.Nil(t, err)
			ltree := res.Tree
			operTree := ltree.rightNode
			// PrintONode(operTree, 0)
			// fmt.Printf("TT-Split#0 i: %d node: %v\n", i, operTree.oper)
			assert.Equal(t, operTree.oper, tt.opstr)
		})
	}
}

func TestSplitAt(t *testing.T) {
	tdata := []struct {
		src   string
		opstr string
		ttype Lt.Lt
	}{
		{`@! var`,
			"@!", Lt.Oper},
		{`@! elem[index]`,
			"@!", Lt.Oper},
		{`r = @defined(a)`,
			"=", Lt.Oper},
		{`@defined(a)`,
			"@", Lt.Oper},
		{`!@defined(a)`,
			"!", Lt.Oper},
		{`var @ list`,
			"@", Lt.Oper},
	}
	for i, tt := range tdata {
		t.Run(fmt.Sprintf("Test, %d) %s >>", i, TrunkString(tt.src)), func(t2 *testing.T) {
			// println("")
			clines := par.SplitCode(tt.src)
			cline := clines[0]
			var prevtree *LineTree
			// fmt.Println(">>>>", cline.Src, "Elems:", FPrintElems(cline.Elems))

			res, err := Line2tree(cline.Elems, prevtree)
			if err != nil {
				t.Logf("Err i:%d Line2tree : %v ", i, err)
			}
			assert.Nil(t, err)
			ltree := res.Tree
			operTree := ltree.rightNode
			// PrintONode(operTree, 0)
			// fmt.Printf("TT-Split#0 i: %d node: %v\n", i, operTree.oper)
			assert.Equal(t, operTree.oper, tt.opstr)
		})
	}
}
