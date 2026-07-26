package cases

import (
	"errors"
	"fmt"

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

func CaseIf(tree *LineTree) (base.Expression, error) {
	return nil, nil
}

func CaseFunc(elems []*lang.Elem, kwRoot *LineTree) (base.Expression, error) {
	// var subElems []*lang.Elem
	// selems := elems
	elems = SkipSpaces(elems)
	if len(elems) < 4 {
		// func name()
		return nil, errors.New("too small lexem set for func definition")
	}

	var prefTree *LineTree
	var err1 error
	// detect method: if has `:` before `()`
	// fmt.Printf("CaseFunc#1: %s \n", FPrintElems(elems))
	sigIndex := 1 // start of signature (from name)
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
	// Name
	name := "f###"
	if elems[sigIndex].Type == Lt.Word {
		name = elems[sigIndex].Text
	}
	// Args
	sigTree, err2 := Line2tree(elems[sigIndex:], kwRoot)
	if err2 != nil {
		return nil, errors.New("func def: incorrect signature")
	}
	fNode := sigTree.Tree.rightNode
	if !fNode.IsBrackets {
		// no brackets, so - error
		// fmt.Printf("CaseFunc#5: %T %v \n", fNode, fNode)
		return nil, errors.New("func def without brackets")
	}
	PrintONode(sigTree.Tree, 0)
	var args []*nodes.VarExpr
	// inBrNode := fNode.rightNode
	if fNode.rightNode != nil || len(fNode.rightElems) > 0 {
		argExp, sok := OperSub(fNode.rightNode, fNode.rightElems)
		if !sok {
			// ?
		}
		switch subs := argExp.(type) {
		case *nodes.SequenceComma:
			if len(subs.Subs) > 0 {
				args = make([]*nodes.VarExpr, len(subs.Subs))
				for i, sub := range subs.Subs {
					if aex, ok2 := sub.(*nodes.VarExpr); ok2 {
						args[i] = aex
					}
				}
			}
		case *nodes.VarExpr:
			args = []*nodes.VarExpr{subs}
		}
	} else {
		// empty args
	}
	// expected sigTree: Brackets -> CommaSepSequence

	funcDef := nodes.NewFuncDef(name, args)
	if prefTree != nil {
		// method def
	}
	return funcDef, nil
}

var mockErr = errors.New("case mock error")

func KWordExp(elems []*lang.Elem) (base.Expression, error) {
	var subElems []*lang.Elem
	// elems = SkipSpaces(elems)
	if len(elems) > 1 {
		subElems = elems[1:]
	}
	rootNode := &OperNode{oper: "KW", leftElems: []*lang.Elem{elems[0]}, prior: 111}
	kwRoot := &LineTree{Tree: rootNode, Parents: []*OperNode{rootNode}}
	switch elems[0].Text {
	case kIf:
		ifTree, err := Line2tree(subElems[1:], kwRoot)
		if err != nil {
			// do smth
		}
		if !ifTree.Finished {
			// unclosed expression, need continue on next line...
		}
		// exp, err2 := CaseIf(ifRes)
		subNode := ifTree.Tree
		// PrintONode(subNode, 0)
		subExp, ok := OperSub(subNode.rightNode, subNode.rightElems)
		if !ok {
			return nil, mockErr
		}
		exp := nodes.NewIf(subExp)
		return exp, nil
	case kElse:
		//
		exp := nodes.NewElseNode()
		if len(SkipSpaces(elems)) == 1 {
			return exp, nil
		}
		// if has inner `if`
		subExp, err := KWordExp(subElems[1:])
		if err != nil {
			// do smth
		}
		sub, ok := subExp.(*nodes.IfNode)
		if !ok {
			// do smth
		}
		exp.SetSlide(sub)
		return exp, nil

	case kFor:
		forTree, err := Line2tree(subElems[1:], kwRoot)
		if err != nil {
			// do smth
		}
		if !forTree.Finished {
			// unclosed expression, need continue on next line...
		}
		subNode := forTree.Tree
		// PrintONode(subNode, 0)
		subExp, ok := OperSub(subNode.rightNode, subNode.rightElems)
		fmt.Printf("Case#For1,1: (%T, %v): %v \n", subExp, subExp, ok)
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
		}
		return exp, err
	case kWhile:
		forTree, err := Line2tree(subElems[1:], kwRoot)
		if err != nil {
			// do smth
		}
		if !forTree.Finished {
			// unclosed expression, need continue on next line...
		}
		subNode := forTree.Tree
		// PrintONode(subNode, 0)
		subExp, ok := OperSub(subNode.rightNode, subNode.rightElems)
		fmt.Printf("Case#While#1: (%T, %v): %v \n", subExp, subExp, ok)
		if !ok {
			return nil, mockErr
		}
		exp := nodes.NewWhile(subExp)
		return exp, nil

	case kFunc:
		return CaseFunc(elems, kwRoot)
	case kBreak:
		exp := &nodes.BreakExp{}
		return exp, nil
	case kContinue:
		exp := &nodes.ContinueExp{}
		return exp, nil
	case kReturn:
	case kMatch:
	case kEnum:
	case kGrup:
	case kStruct:
	case kImport:
	case kConst:
	case kRun:
	}
	return nil, nil
}
