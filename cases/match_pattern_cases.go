package cases

import (
	"fmt"
	"slices"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
	"github.com/lesnikyan/lisapet-go/nodes"
)

func SimplePattern(exp base.Expression) (nodes.MatchPattern, error) {
	switch tex := exp.(type) {
	case *nodes.ValExpr:
		return &nodes.MCaseVal{Val: nodes.GetExprVal(tex, nil)}, nil
	case *nodes.VarExpr:
		// sub pattern: var
		if tex.GetName() == "_" {
			// fmt.Println("mt var=", tex.GetName())
			// common sub pattern
			return &nodes.MCaseUnder{}, nil
		}
	}
	return nil, fmt.Errorf("No expression case for simple matching pattern by %T type ", exp)
}

func ValOperPattern(expr base.Expression) (nodes.MatchPattern, error) {
	// fmt.Printf("valOper, expr: %T \n", expr)
	switch pex := expr.(type) {
	case *nodes.ValExpr:
		return SimplePattern(expr)
	case *nodes.UnaryLeft:
		pex.Do(nil)
		switch pex.Oper.Id {
		case nodes.OpMinus:
			val := nodes.GetExprVal(pex, nil)
			// fmt.Println("raw val=", val)
			return &nodes.MCaseVal{Val: val}, nil
		}
	case *nodes.OperBin:
		switch pex.Oper.Id {
		case nodes.OpDuoColon:
			// type pattern
		case nodes.OpBitOr:
			// one | two
			//

			//  case nodes.OpDot:
			// 	println("MtC oper dot")
			// val := nodes.GetExprVal(pex, nil)
		}
	}
	return nil, nil
}

func ProcSubPattern(rNode *OperNode, elems []*lang.Elem) (nodes.MatchPattern, error) {
	return nil, nil
}

var _mtPatOps = strings.Split(". -", " ")

func ProcMatchPattern(rNode *OperNode, elems []*lang.Elem) (nodes.MatchPattern, error) {
	// fmt.Println("ProcMPT", FPrintElems(elems), "node:", rNode)
	if rNode == nil {
		if len(elems) == 0 {
			// empty tree
			return nil, fmt.Errorf("No elems or nodes for matching pattern case")
		}
		expr, ok := OperSub(nil, elems)
		if !ok {
			return nil, fmt.Errorf("Bad elem set for matching pattern: `%s`", FPrintElems(elems))
		}
		return SimplePattern(expr)
	}
	if slices.Contains(_mtPatOps, rNode.oper) {
		// println("val oper: ", rNode.oper)
		// PrintONode(rNode, 0)
		expr, ok := OperSub(rNode, nil)
		if !ok {
			return nil, fmt.Errorf("Bad oper node for matching pattern: `%s`", rNode.oper)
		}
		return ValOperPattern(expr)
	}
	switch rNode.oper {
	case "[":
		// list
	case "(":
		// tuple | some() | expr in brackets
	case "{":
		// dict | Struct{}
	case "|":
		// multicase
	case "::":
		// typed val
		// println("oper ::")
		// PrintONode(rNode, 0)
		expr, ok := OperSub(rNode.rightNode, rNode.rightElems)
		if !ok {
			return nil, fmt.Errorf("Bad sub expr in :: oper")
		}
		return &nodes.MCaseType{TypeExp: expr}, nil
	case "@":
		// val assign
	case ":?":
		// extra guard

		// case "_":
		// most common pattern

	}
	return nil, fmt.Errorf("matching pattern not found %s", rNode.oper)
}
