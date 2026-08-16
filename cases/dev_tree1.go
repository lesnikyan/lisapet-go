package cases

import (
	"errors"
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
	"github.com/lesnikyan/lisapet-go/nodes"
)

var lineInterpretErr = errors.New("Incorrect interpretation of line")

type SplitState struct {
	Expr  base.Expression
	LTree *LineTree
	Done  bool
}

func Line2Expr(cline *lang.CLine, prevTree *LineTree) (*SplitState, error) {

	var err error
	var expr base.Expression
	var ok bool
	if IsLKWord(cline.Elems, prevTree) {
		// control or definition expression
		state, err := KWordExp(cline.Elems, prevTree)
		if err != nil {
			fmt.Println("Error LKWord!", err)
			return nil, err
		}
		if !state.Done {
			// unclosed
			return state, nil
		}
		expr = state.Expr
		// fmt.Println("L2E1>", cline.Src, expr, err)
	} else {
		if len(cline.Elems) == 1 {
			expr, ok := ProcSubElems(cline.Elems)
			if ok {
				res := &SplitState{Expr: expr, Done: true}
				return res, nil
			}
		}
		var res *LineTree
		res, err = Line2tree(cline.Elems, prevTree)
		if err != nil {
			fmt.Println("Error!", err)
			return nil, err
		}
		if !res.Finished {
			return &SplitState{Done: false, LTree: res}, nil
		}
		ltree := res.Tree
		// fmt.Println("L2E2>", cline.Src, res, err, "::", ltree, ":~", ltree.rightNode)
		// PrintONode(ltree, 0)
		operTree := ltree.rightNode
		if operTree == nil {
			return nil, nil
		}
		// fmt.Println("L2E3>", cline.Src, res, err, "r-oper:", operTree.oper)
		// PrintONode(operTree, 0)
		expr, ok = ProcExprTree(operTree)
		if !ok {
			return nil, lineInterpretErr
		}
	}
	res := &SplitState{Expr: expr, Done: true}
	return res, nil
}

type BlockLink struct {
	elem   base.Block
	indent int
}

type ExprLink struct {
	elem   base.Expression
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
	var prevState *SplitState
	var prevExpr *ExprLink
	var nExpr *ExprLink
	for _, cline := range clines {
		if len(cline.Elems) == 0 {
			continue
		}
		var ltree *LineTree
		if prevState != nil {
			// next part of expression
			ltree = prevState.LTree
		} else {
			// new expression
			cind = cline.Indent
		}
		// fmt.Println("\n>>>>", cline.Src, "bLen:", len(parents), fmt.Sprintf("nBlock: %T", nblock.elem), "indent:", cind)
		curState, err := Line2Expr(cline, ltree)

		if err != nil {
			// fmt.Println("Error of line expr!", err)
			return nil, err
		}
		if curState == nil {
			// comment, empty line
			continue
		}
		if !curState.Done {
			// unclosed brackets or another splitted expression
			prevState = curState
			// fmt.Printf(">>TreeBlock. State Not Done")
			continue
		}
		prevState = nil
		expr := curState.Expr
		if expr == nil {
			// possibly: commented line
			continue
		}
		prevExpr = nExpr
		nExpr = &ExprLink{elem: expr, indent: cind}
		// tp := fmt.Sprintf("%T", expr)
		// fmt.Println("tt2>", tp, nodes.OperArgsInfo(expr))
		elseInd := false // if expr is `else`

		// fmt.Println("Tree,Indent:", nblock.indent, cind, " back lvl:", cind <= nblock.indent, "pLen:", len(parents))
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
					nblock = parents[i]
					parents = parents[:i+1]
					break
				}
			}
		}
		// fmt.Println("Tree,Indent2:", "pLen:", len(parents))

		// fmt.Printf("tree.Block %T: %v .line: (%T: %v)  \n", nblock.elem, nblock.elem, expr, expr)
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
				// fmt.Printf("tree. expr base.Block: %T is parrent \n", texp)
				bl := &BlockLink{elem: texp, indent: cind}
				parents = append(parents, bl)
				nblock = bl
			}
		default:
			// fmt.Printf("tree.def nblock: %T: %v ex (%T: %v) ; %d >> %d  \n", nblock.elem, nblock.elem, texp, texp, nblock.indent, cind)
			if prevExpr != nil && cind > prevExpr.indent {
				// fmt.Printf("tree.def sub-expression indent %T > %T > %T \n", nblock.elem, prevExpr.elem, expr)
				switch upar := prevExpr.elem.(type) {
				case *nodes.OperAssign:
					// fmt.Printf("tree.def prev:OperAssign Right: %T  \n", upar.Right)
					if IsConstruct(upar.Right) {
						// if construct: [] (,) {} T{}
						bl := &BlockLink{elem: upar, indent: prevExpr.indent}
						parents = append(parents, bl)
						nblock = bl
					}
				}
			}

			// fmt.Printf("tree.def Add... : %T  \n", nblock.elem)
			switch parn := nblock.elem.(type) {
			case *nodes.OperAssign:
				parn.AddToRight(expr)
			default:
				parn.Add(expr)
			}
		}
	}
	return top, nil
}

func IsConstruct(ex base.Expression) bool {
	switch ex.(type) {
	case *nodes.ListExpr, *nodes.TupleExpr, *nodes.DictExpr:
		return true
	case *nodes.StructConstr:
		return true
	}
	return false
}
