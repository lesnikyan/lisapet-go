package cases

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/lesnikyan/lisapet-go/nodes"
)

const (
	kIf       = "if"
	kFor      = "for"
	kWhile    = "while"
	kMatch    = "match"
	kFunc     = "func"
	kElse     = "else"
	kEnum     = "enum"
	kGrup     = "grup"
	kStruct   = "struct"
	kImport   = "import"
	kReturn   = "return"
	kBreak    = "break"
	kContinue = "continue"
	kConst    = "const"
	kRun      = "run"
)

// func CaseIf(tree *LineTree) (base.Expression, error) {
// 	return nil, nil
// }

func CaseFunc(elems []*lang.Elem, kwRoot *LineTree) (*SplitState, error) {
	elems = SkipSpaces(elems)
	var prefTree *LineTree
	var err1 error
	// detect method: if has `:` before `()`
	// fmt.Printf("CaseFunc#1: %s \n", FPrintElems(elems))
	sigIndex := 0 // start of signature (from name)
	if kwRoot.BracketsCount == 0 {
		sigIndex = 1 // start of signature (from name)
		if len(elems) > 6 {
			// possible mehod def: `func inst:Type Name()`
			if elems[2].Type == Lt.Oper && elems[2].Text == opColon {
				prefTree, err1 = Line2tree(elems[1:4], nil)
				if err1 != nil {
					return nil, errors.New("func def: incorrect prefix of method")
				}
				sigIndex = 4
				// TODO: add prefTree to kwRoot
			}
		}
	}
	// Name
	name := "f###"
	// Args
	sigTree, err2 := Line2tree(elems[sigIndex:], kwRoot)
	if err2 != nil {
		return nil, errors.New("func def: incorrect signature")
	}
	if !sigTree.Finished {
		return &SplitState{Done: false, LTree: sigTree}, nil
	}
	fNode := sigTree.Tree.rightNode
	if !fNode.IsBrackets {
		// no brackets, so - error
		// fmt.Printf("CaseFunc#5: %T %v \n", fNode, fNode)
		return nil, errors.New("func def without brackets")
	}
	// PrintONode(sigTree.Tree, 0)
	// PrintONode(fNode, 0)
	var args []base.Expression
	// inBrNode := fNode.rightNode
	if fNode.rightNode != nil || len(fNode.rightElems) > 0 {
		argExp, sok := OperSub(fNode.rightNode, fNode.rightElems)
		if !sok {
			// ?
		}
		switch subs := argExp.(type) {
		case *nodes.SequenceComma:
			if len(subs.Subs) > 0 {
				// fmt.Println("FuncCase#01", len(subs.Subs))
				args = make([]base.Expression, len(subs.Subs))
				for i, sub := range subs.Subs {
					switch aex := sub.(type) {
					case *nodes.VarExpr, *nodes.OperAssign, *nodes.OperColon, *nodes.TripleDots:
						args[i] = aex
					}
				}
			}
		case *nodes.VarExpr, *nodes.OperAssign, *nodes.OperColon, *nodes.TripleDots:
			args = []base.Expression{subs}
		}
	} else {
		// empty args
	}
	// expected sigTree: Brackets -> CommaSepSequence
	if len(fNode.leftElems) == 1 {
		// should be name
		name = fNode.leftElems[0].Text
	}
	// fmt.Println("FuncCase# name ", name)

	funcDef := nodes.NewFuncDef(name, args)
	if prefTree != nil {
		// method def
	}
	return &SplitState{Expr: funcDef, Done: true}, nil
}

// struct StrName(Parent) a, b, c:int
// struct StrName c:int
func CaseStructDef(elems []*lang.Elem, kwRoot *LineTree) (*SplitState, error) {
	// fmt.Printf("Case#Struct#0: '%v' : (%v) \n", elems[0].Text, Lt.TName(elems[0].Type))
	elems = SkipSpaces(elems)
	name := "s#"
	fIndex := 2 // field start index
	// normally struct definition placed in 1 line, or uses block-syntax
	if kwRoot.BracketsCount == 0 {
		if elems[1].Type != Lt.Word {
			return nil, errors.New("keyword struct: bad name in definition")
		}
		name = elems[1].Text
		// if brackets after name - split parent part and fields part
		if len(elems) > 2 {
			if elems[2].Type == Lt.Oper {
				if elems[2].Text != "(" {
					// bad syntax
					return nil, errors.New("Bad syntax of struct case")
				}
				// make Parent sub
				// fIndex = elem after brackets
			}
		} else {
			strDef := &nodes.StructDefExpr{Name: name, Fields: []base.Expression{}}
			return &SplitState{Expr: strDef, Done: true}, nil
		}
	}

	// expecting comma-separated sequence
	nTree, err := Line2tree(elems[fIndex:], kwRoot)
	if err != nil {
		// do smth
		return nil, err
	}

	if !nTree.Finished {
		return &SplitState{Done: false, LTree: nTree}, nil
	}

	// PrintONode(nTree.Tree, 0)
	var args []base.Expression
	subNode := nTree.Tree
	// PrintONode(subNode, 0)
	subExp, ok := OperSub(subNode.rightNode, subNode.rightElems)
	// fmt.Printf("Case#Struct#1: (%T, %v): %v \n", subExp, subExp, ok)
	if !ok {
		return nil, errors.New("struct: bad subs")
	}
	switch argExp := subExp.(type) {
	case *nodes.SequenceComma:
		// several fields
		args = make([]base.Expression, len(argExp.Subs))
		for i, arg := range argExp.Subs {
			args[i] = arg
		}
	case *nodes.VarExpr:
		// 1 field no type
		args = []base.Expression{argExp}
	case *nodes.OperColon:
		// 1 typed field
		args = []base.Expression{argExp}
	}
	strDef := &nodes.StructDefExpr{Name: name, Fields: args}
	return &SplitState{Expr: strDef, Done: true}, nil
}

var mockErr = errors.New("case mock error")

func KWordExp(elems []*lang.Elem, prevTree *LineTree) (*SplitState, error) {
	var subElems []*lang.Elem
	// elems = SkipSpaces(elems)
	if len(elems) > 1 {
		subElems = elems[1:]
	}
	var kwRoot *LineTree
	if prevTree != nil {
		kwRoot = prevTree
	} else {
		rootNode := &OperNode{oper: "KW", leftElems: []*lang.Elem{elems[0]}, prior: 111}
		kwRoot = &LineTree{Tree: rootNode, Parents: []*OperNode{rootNode}}
	}
	var kwText string
	if prevTree != nil && len(prevTree.Tree.leftElems) > 0 {
		kwText = prevTree.Tree.leftElems[0].Text
	} else {
		kwText = elems[0].Text
	}
	switch kwText {
	case kIf:
		ifTree, err := Line2tree(subElems, kwRoot)
		if err != nil {
			// do smth
		}
		if !ifTree.Finished {
			// unclosed expression, need continue on next line...
			return &SplitState{Done: false, LTree: ifTree}, nil
		}
		// exp, err2 := CaseIf(ifRes)
		subNode := ifTree.Tree
		// PrintONode(subNode, 0)
		subExp, ok := OperSub(subNode.rightNode, subNode.rightElems)
		if !ok {
			return nil, mockErr
		}
		exp := nodes.NewIf(subExp)
		return &SplitState{Expr: exp, Done: true}, nil
	case kElse:
		//
		exp := nodes.NewElseNode()
		if len(SkipSpaces(elems)) == 1 {
			return &SplitState{Expr: exp, Done: true}, nil
		}
		// if has inner `if`
		subExp, err := KWordExp(subElems[1:], prevTree)
		if err != nil {
			// do smth
		}
		sub, ok := subExp.Expr.(*nodes.IfNode)
		if !ok {
			// do smth
		}
		exp.SetSlide(sub)
		return &SplitState{Expr: exp, Done: true}, nil

	case kFor:
		forTree, err := Line2tree(subElems[1:], kwRoot)
		if err != nil {
			// do smth
		}
		if !forTree.Finished {
			// unclosed expression, need continue on next line...
			return &SplitState{Done: false, LTree: forTree}, nil
		}
		subNode := forTree.Tree
		// PrintONode(subNode, 0)
		subExp, ok := OperSub(subNode.rightNode, subNode.rightElems)
		// fmt.Printf("Case#For1,1: (%T, %v): %v \n", subExp, subExp, ok)
		if !ok {
			return nil, mockErr
		}
		var exp nodes.ForExpr
		switch subFor := subExp.(type) {
		case *nodes.SequenceSemicolon:
			exp = nodes.NewForCond(subFor)
		case *nodes.LeftArrow:
			exp = nodes.NewForSource(subFor)
		default:
			err = errors.New("bad sub-expr for `for` expression")
			return nil, err
		}
		return &SplitState{Expr: exp, Done: true}, err
	case kWhile:
		forTree, err := Line2tree(subElems[1:], kwRoot)
		if err != nil {
			// do smth
		}
		if !forTree.Finished {
			// unclosed expression, need continue on next line...
			return &SplitState{Done: false, LTree: forTree}, nil
		}
		subNode := forTree.Tree
		// PrintONode(subNode, 0)
		subExp, ok := OperSub(subNode.rightNode, subNode.rightElems)
		// fmt.Printf("Case#While#1: (%T, %v): %v \n", subExp, subExp, ok)
		if !ok {
			return nil, mockErr
		}
		exp := nodes.NewWhile(subExp)
		return &SplitState{Expr: exp, Done: true}, nil

	case kFunc:
		return CaseFunc(elems, kwRoot)
	case kBreak:
		exp := &nodes.BreakExp{}
		return &SplitState{Expr: exp, Done: true}, nil
	case kContinue:
		exp := &nodes.ContinueExp{}
		return &SplitState{Expr: exp, Done: true}, nil
	case kReturn:
		subElems = SkipSpaces(subElems)
		if len(subElems) == 0 {
			exp := nodes.NewReturn(nil)
			return &SplitState{Expr: exp, Done: true}, nil
		}
		subTree, err := Line2tree(subElems, kwRoot)
		if err != nil {
			// do smth
		}
		if !subTree.Finished {
			// unclosed expression, need continue on next line...
			return &SplitState{Done: false, LTree: subTree}, nil
		}
		subNode := subTree.Tree
		// PrintONode(subNode, 0)
		subExp, ok := OperSub(subNode.rightNode, subNode.rightElems)
		if !ok {
			return nil, errors.New("Return: bad sub")
		}
		exp := nodes.NewReturn(subExp)
		return &SplitState{Expr: exp, Done: true}, nil
	case kMatch:
	case kEnum:
	case kGrup:
	case kStruct:
		// struct StrName(Parent) a, b, c:int
		return CaseStructDef(elems, kwRoot)
	case kImport:
	case kConst:
	case kRun:
	}
	return nil, errors.New("keyword not found")
}
