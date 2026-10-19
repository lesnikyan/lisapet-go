package cases

import (
	"errors"
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
		vname := tex.GetName()
		if vname == "_" {
			// fmt.Println("mt var=", tex.GetName())
			// common sub pattern
			return &nodes.MCaseUnder{}, nil
		}
		return &nodes.MCaseVar{Expr: tex}, nil
	}
	return nil, fmt.Errorf("No expression case for simple matching pattern by %T type ", exp)
}

func VarPattern(node *OperNode, elems []*lang.Elem) (*nodes.MCaseVar, error) {
	if len(elems) != 1 {
		return nil, fmt.Errorf("Not a var expression in case matching pattern by elems: %s ", FPrintElems(elems))
	}
	expr, ok := OperSub(node, elems)
	if !ok {
		return nil, fmt.Errorf("No expression case for var matching pattern by elems: %s ", FPrintElems(elems))
	}
	varx, ok := expr.(*nodes.VarExpr)
	if !ok {
		return nil, fmt.Errorf("mathcing pattern: not a var exprettion: %T ", expr)
	}
	return &nodes.MCaseVar{Expr: varx}, nil
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

var _mtPatOps = strings.Split(". -", " ") // matching pattern operators of single value: negative, float

// func subMPElemSimple(node *OperNode, elems []*lang.Elem) (nodes.MatchPattern, error) {
// 	if len(elems) != 1 {
// 		return nil, fmt.Errorf("m-pattern: sub elem more 1 lexem: not found")
// 	}
// 	switch elems[0].Type {
// 	case Lt.Oper:
// 		// *, ?
// 	case Lt.Word:
// 		// var
// 		switch
// 	default:
// 		// values

// 	}
// }

// func subMPElem(node *OperNode, elems []*lang.Elem) (nodes.MatchPattern, error) {

// }

func MPSeqSubs(rNode *OperNode) ([]nodes.MatchPattern, error) {
	node := rNode
	elems := []nodes.MatchPattern{}
	last := rNode
	for node != nil && node.oper == rNode.oper {
		// fmt.Println("MPSeqSubs", FPrintElems(node.rightElems), "node:")
		// PrintONode(node.rightNode, 0)
		last = node
		if node.rightNode == nil && len(node.rightElems) == 0 {
			break
		}
		sub, err := ProcMatchPattern(node.rightNode, node.rightElems)
		// fmt.Printf("MPSeqSubs#1 %v %T err: %v\n", node.oper, sub, err)
		if err != nil {
			return nil, err
		}
		elems = append(elems, sub)
		node = node.leftNode
	}
	// nodes.MatchPattern
	if last.leftNode != nil || len(last.leftElems) != 0 {
		first, err := ProcMatchPattern(last.leftNode, last.leftElems)
		if err != nil {
			return nil, err
		}
		// fmt.Printf("SeqSub#3 %T ;; %T", last, first)
		elems = append(elems, first)
	}
	slices.Reverse(elems)
	// fmt.Println("Seq#1:", elems)
	return elems, nil
}

// func __MPList(node *OperNode) (nodes.MatchPattern, error) {
// 	subNode := node.rightNode
// 	elems := node.rightElems
// 	if subNode == nil {
// 		if len(elems) == 0 {
// 			// no subs
// 			return nodes.NewMCaseList(nil), nil
// 		}
// 		// 1 sub: val, var, _, ?, *
// 		expr, ok := OperSub(nil, elems)
// 		if !ok {
// 			return nil, fmt.Errorf("Bad sequnce elem in matching pattern: `%s`", FPrintElems(elems))
// 		}
// 		subOne, err := SimplePattern(expr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return nodes.NewMCaseList([]nodes.MatchPattern{subOne}), nil
// 	}
// 	// next if sub oper not empty
// 	switch subNode.oper {
// 	case "?":
// 		//  optional
// 		// fmt.Printf("QM elem ? \n")
// 		return nodes.NewMCaseList([]nodes.MatchPattern{&nodes.MCaseQMark{}}), nil
// 	case "*":
// 		// any num of optional
// 		return nodes.NewMCaseList([]nodes.MatchPattern{&nodes.MCaseStar{}}), nil
// 	case ",":
// 		// regular case: several elements
// 		subs, err := MPSeqSubs(subNode)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return nodes.NewMCaseList(subs), nil
// 	}
// 	panic("list match pattern: case not found")
// }

func SeqPatternSubs(subNode *OperNode, elems []*lang.Elem) ([]nodes.MatchPattern, error) {
	if subNode == nil {
		if len(elems) == 0 {
			// no subs
			return nil, nil
		}
		// 1 sub: val, var, _, ?, *
		expr, ok := OperSub(nil, elems)
		if !ok {
			return nil, fmt.Errorf("Bad sequnce elem in matching pattern: `%s`", FPrintElems(elems))
		}
		subOne, err := SimplePattern(expr)
		if err != nil {
			return nil, err
		}
		return []nodes.MatchPattern{subOne}, nil
	}
	// next if sub oper not empty
	switch subNode.oper {
	case "?":
		//  optional
		// fmt.Printf("QM elem ? \n")
		return []nodes.MatchPattern{&nodes.MCaseQMark{}}, nil
	case "*":
		// any num of optional
		return []nodes.MatchPattern{&nodes.MCaseStar{}}, nil
	case ",":
		// regular case: several elements
		subs, err := MPSeqSubs(subNode)
		if err != nil {
			return nil, err
		}
		return subs, nil
	}
	return nil, fmt.Errorf("sequence match pattern: case not found")
}

func MPList(node *OperNode) (nodes.MatchPattern, error) {
	subNode := node.rightNode
	elems := node.rightElems
	subs, err := SeqPatternSubs(subNode, elems)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("bad subs of list match-pattern"), err)
	}
	return nodes.NewMCaseList(subs), nil
}

func MPTuple(node *OperNode) (nodes.MatchPattern, error) {
	subNode := node.rightNode
	elems := node.rightElems
	if subNode == nil && len(elems) != 0 {
		// exceptions of single sub in parentheses
		mp, err := ProcMatchPattern(subNode, elems)
		if err != nil {
			return nil, err
		}
		switch pp := mp.(type) {
		case *nodes.MCaseUnder:
			return nodes.NewMCaseTuple([]nodes.MatchPattern{pp}), nil
		default:
			return mp, nil
		}
		// return nodes.NewMCaseTuple([]nodes.MatchPattern{mp}), nil
	}
	subs, err := SeqPatternSubs(subNode, elems)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("bad subs of tuple match-pattern"), err)
	}
	return nodes.NewMCaseTuple(subs), nil
}

func ProcMatchPattern(rNode *OperNode, elems []*lang.Elem) (nodes.MatchPattern, error) {
	// fmt.Println("ProcMPT", FPrintElems(elems), "node:")
	// PrintONode(rNode, 0)
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
		mpt, err := MPList(rNode)
		if err != nil {
			return nil, err
		}
		return mpt, nil
	case "(":
		// println("oper ( ..")
		// PrintONode(rNode, 0)
		// tuple | some() | expr in brackets
		// () - tuple
		// (*), (?), (_) - tuple
		// (5), (null), ([1,2]), ('hello') - not a tuple
		// (int|list), (a::int) - not a tuple
		if len(rNode.leftElems) == 0 && rNode.leftNode == nil {
			// tuple pattern
			mpt, err := MPTuple(rNode)
			if err != nil {
				return nil, err
			}
			return mpt, nil
		}
		// other cases with parentheses
	case "{":
		// dict | Struct{} | _{}
	case "|":
		// multicase
	case "::":
		// typed val
		// println("oper ::")
		// PrintONode(rNode, 0)
		expr, ok := OperSub(rNode.rightNode, rNode.rightElems)
		// fmt.Printf(" :: %T \n", expr)
		if !ok {
			return nil, fmt.Errorf("Bad right sub expr in :: oper")
		}
		var varPtr *nodes.MCaseVar
		if rNode.leftNode != nil || len(rNode.leftElems) > 0 {
			vpt, err := VarPattern(rNode.leftNode, rNode.leftElems)
			if err != nil {
				return nil, err
			}
			varPtr = vpt
		}
		return &nodes.MCaseType{TypeExp: expr, Var: varPtr}, nil
	case "?":
		//  optional one
		// fmt.Printf("QM ? \n")
		return &nodes.MCaseQMark{}, nil
	case "*":
		//  optional any number
		// fmt.Printf("Star * \n")
		return &nodes.MCaseStar{}, nil
	case "@":
		// val assign
	case ":?":
		// extra guard

		// case "_":
		// most common pattern

	}
	return nil, fmt.Errorf("matching pattern not found %s", rNode.oper)
}
