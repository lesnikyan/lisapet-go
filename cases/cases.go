package cases

import (
	"github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/lesnikyan/lisapet-go/nodes"
)

type LangCase interface {
	match(line lang.CLine) bool
	expr(line lang.CLine) *nodes.TNode
}

type CaseNum struct{}

func (c *CaseNum) match(line lang.CLine) bool {
	return false
}

func (c *CaseNum) expr(line lang.CLine) *nodes.TNode {
	return nil
}

// another way:

type CPart struct {
	Line  lang.CLine
	start int
	end   int
}

type TNode = nodes.TNode

func CaseVar(part CPart) (bool, *TNode) {
	if len(part.Line.Elems) != 1 {
		return false, nil
	}
	if part.Line.Elems[0].Type != Lt.Word {
		return false, nil
	}
	// parse var
	return true, &TNode{}
}
