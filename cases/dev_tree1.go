package cases

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
	"github.com/lesnikyan/lisapet-go/nodes"
)

var lineInterpretErr = errors.New("Incorrect interpretation of line")

func Line2Expr(cline *lang.CLine) (base.Expression, error) {

	var err error
	var expr base.Expression
	var ok bool
	if IsLKWord(cline.Elems) {
		// control or definition expression
		expr, err = KWordExp(cline.Elems)
		if err != nil {
			fmt.Println("Error LKWord!", err)
			return nil, err
		}
		fmt.Println("L2E1>", cline.Src, expr, err)
	} else {
		if len(cline.Elems) == 1 {
			exp, ok := ProcSubElems(cline.Elems)
			if ok {
				return exp, nil
			}
		}
		var res *LineTree
		res, err = Line2tree(cline.Elems, nil)
		if err != nil {
			fmt.Println("Error!", err)
			return nil, err
		}
		ltree := res.Tree
		operTree := ltree.rightNode
		if operTree == nil {
			return nil, nil
		}
		fmt.Println("L2E2>", cline.Src, res, err, "r-oper:", operTree.oper)
		PrintONode(operTree, 0)
		expr, ok = ProcExprTree(operTree)
		if !ok {
			return nil, lineInterpretErr
		}
	}
	return expr, nil
}

type BlockLink struct {
	elem   base.Block
	indent int
}

// make executable block of expressions by parsed lines
func TreeBlock(clines []*lang.CLine) (*nodes.BlockExpr, error) {
	top := nodes.NewBlock()
	if len(clines) == 0 {
		return nil, errors.New("empty code to build tree")
	}
	firstIndent := clines[0].Indent
	var nblock *BlockLink = &BlockLink{elem: top, indent: firstIndent - 1} // current parent block
	parents := []*BlockLink{nblock}
	cind := firstIndent
	for _, cline := range clines {
		if len(cline.Elems) == 0 {
			continue
		}
		fmt.Println("\n>>>>", cline.Src, "bLen:", len(parents), fmt.Sprintf("nBlock: %T", nblock.elem))
		expr, err := Line2Expr(cline)
		if err != nil {
			fmt.Println("Error of line expr!", err)
			return nil, err
		}
		if expr == nil {
			// possibly: commented line
			continue
		}
		// tp := fmt.Sprintf("%T", expr)
		// fmt.Println("tt2>", tp, nodes.OperArgsInfo(expr))
		cind = cline.Indent
		elseInd := false // if expr is `else`

		fmt.Println("Tree,Indent:", nblock.indent, cind, " back lvl:", cind <= nblock.indent, "pLen:", len(parents))
		if cind <= nblock.indent {
			// end of prev block
			if _, ok := expr.(*nodes.ElseNode); ok {
				elseInd = true
			}
			pfound := false
			for i := len(parents) - 1; i >= 0; i-- {
				pfound = false
				if elseInd {
					if parents[i].indent == cind {
						pfound = true
					}
				} else if parents[i].indent < cind {
					pfound = true
				}
				// pfound = () || ()
				if pfound {
					// TODO: resolve cases: else, else if
					nblock = parents[i]
					parents = parents[:i+1]
					break
				}
			}
		}

		fmt.Printf("tree.Block %T: %v .line: (%T: %v)  \n", nblock.elem, nblock.elem, expr, expr)
		switch texp := expr.(type) { // cur expr
		case *nodes.ElseNode:
			// nblock is: if | else if
			switch prev := (nblock.elem).(type) { // parent block
			case *nodes.IfNode:
				prev.SetElse(texp)
			case *nodes.ElseNode:
				if !prev.Slided {
					// else after else without condition
					return nil, errors.New("Bad else structure N1")
				}
				prev.SlideIf.SetElse(texp)
			default:
				return nil, errors.New("Bad if-else structure")
			}
			bl := &BlockLink{elem: texp, indent: cind}
			// next `else` will be parent in chain
			parents = parents[:len(parents)-1]
			parents = append(parents, bl)
			nblock = bl

		case base.Block:
			nblock.elem.Add(expr)
			if texp.IsParent() {
				bl := &BlockLink{elem: texp, indent: cind}
				parents = append(parents, bl)
				nblock = bl
			}
		default:
			fmt.Printf("tree.def Add: %T: %v .Add (%T: %v)  \n", nblock.elem, nblock.elem, texp, texp)
			nblock.elem.Add(expr)
		}
	}
	return top, nil
}
