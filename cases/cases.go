package cases

import (
	"fmt"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
)

// import (
// 	"github.com/lesnikyan/lisapet-go/lang"
// 	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
// 	"github.com/lesnikyan/lisapet-go/nodes"
// )

/*
Nearest plans:
-1. unary oper
1. tuple
2. list
3. dict
4. elem of collection: nn[key]
5. other base types: byte, glif, ..
6. context vars: add, get; parent ctx: find var
7. multiline brackets
8. multiline strings
9. const
9.1 typed var x: int
9.2 multiassign
9.3 unpack collecition: a,b,c = [1,2,3]
10. oprator math-assign: += -= .. ets
11. block
12. if
12.1 else
13. condition operators: ?> !?> ::
14. opertor ?:
15. bytes 0x[]
16. function: definition, call
17. builtin functions: print, len
18. func result
19. if result
20. for ; ;
21. while expr
22. continue, break
23. for n <- list, for k, v <- dict
24. comprehencion
25. generator
26. a, b, c <- gen|list|list+dict
27.
*/

type CaseRes struct {
	Expr base.Expression
	Subs [][]*lang.Elem
}

type LineTree struct {
	Tree     *OperNode   // parsed tree of operators
	Finished bool        // if all operators and brackates was finished
	Parents  []*OperNode // final chain of parent nodes - right branch of tree
}

// type LangCase interface {
// 	match(line lang.CLine) bool
// 	expr(line lang.CLine) *nodes.TNode
// }

// type CaseNum struct{}

// func (c *CaseNum) match(line lang.CLine) bool {
// 	return false
// }

// func (c *CaseNum) expr(line lang.CLine) *nodes.TNode {
// 	return nil
// }

// // another way:

// type CPart struct {
// 	Line  lang.CLine
// 	start int
// 	end   int
// }

// type TNode = nodes.TNode

// func CaseVar(part CPart) (bool, *TNode) {
// 	if len(part.Line.Elems) != 1 {
// 		return false, nil
// 	}
// 	if part.Line.Elems[0].Type != Lt.Word {
// 		return false, nil
// 	}
// 	// parse var
// 	return true, &TNode{}
// }

func PrintOpArg(side string, node *OperNode, elems []*lang.Elem, ind int) {
	if node != nil || elems == nil {
		fmt.Printf(" %s %s: ", strings.Repeat(" ▵", ind), side)
		if node != nil {
			PrintONode(node, ind)
		} else {
			fmt.Println("<->")
		}
	} else {
		fmt.Printf(" %s %s▷ %v\n", strings.Repeat(" .", ind), side, FPrintElems(elems))
	}
}

func PrintONode(node *OperNode, ind int) {
	//strings.Repeat(" ▵", ind)
	if node == nil {
		fmt.Println("Node", node, ind)
		return
	}
	fmt.Printf("⧐ %s %s \n", "", node.oper)
	PrintOpArg("L", node.leftNode, node.leftElems, ind+1)
	PrintOpArg("R", node.rightNode, node.rightElems, ind+1)
}
