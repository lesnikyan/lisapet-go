package cases

import (
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/lesnikyan/lisapet-go/parser"
)

/*
split formattin string
~' text {var} {var[keyIndex]} {obj.member} {funcCall()} {arg + arg} {{ just text }}'
~""
~''' '''
~""" """
~``
~``` ```


*/

// text or node
type FPart struct {
	Text string
	Node base.Expression // if nil - .Text is a value
}

type FNode struct {
	Form string
	Node *OperNode
}

type FString struct {
	Elems []*FPart
	// Nodes []*OperNode
	FNodes []*FNode
}

// ' a {a + b : <+05.4f} b'

func SplitFString(elem *lang.Elem) (*FString, error) {
	// fmt.Printf("SF-S. Text: `%v` \n", etext)
	pref := 0
	suf := 0
	switch elem.Type {
	case Lt.Text:
		pref = 1
		suf = 1
	case Lt.Mttext:
		pref = 3
		suf = 3
	}
	src := base.CropStr(elem.Text, pref, suf)
	runes := []rune(src)
	// fmt.Printf("SF-S. Text: `%v` \n", src)

	subStr := make([]rune, len(runes))
	subLen := 0
	formLen := 0
	exp := false
	cur := make([]rune, len(runes))
	curLen := 0
	res := []*FPart{}
	rnodes := []*FNode{}
	inBr := 0
	expStart := -1
	for i, c := range runes {
		// fmt.Printf(" %% -- %d %c\n", i, c)
		if !exp {
			if c == '{' {
				// open sub expr. TODO: resolve escaped {{}}
				exp = true
				next := string(cur[0:curLen])
				// fmt.Printf("  fpartS %d): `%s`\n", i, next)
				// res = append(res, &FPart{Text: next})
				nelems := []*lang.Elem{&lang.Elem{Type: Lt.Text, Text: next}}
				rnodes = append(rnodes, &FNode{Form: "", Node: &OperNode{oper: "T", rightElems: nelems}})
				curLen = 0
				formLen = 0
				subLen = 0
				expStart = i + 1
				continue
			}
			cur[curLen] = c
			curLen++
		} else {
			if strings.ContainsRune("({[", c) {
				inBr++
			} else if strings.ContainsRune("}])", c) {
				inBr--
			}
			if inBr > 0 {
				subLen++
				continue
			}
			if c == ':' && formLen == 0 {
				// separator between expression anf format
				formLen++
				continue
			}
			if c == '}' {
				exp = false
				// close sub expr
				// expText := subStr[0:subLen]
				expText := runes[expStart : expStart+subLen]
				form := ""
				if formLen > 0 {
					form = string(runes[expStart+subLen+1 : expStart+subLen+formLen])
				}
				// fmt.Printf("  fpartN: `%s`, form: `%s` \n", string(expText), form)
				//...
				elems := parser.SplitLine(expText, &parser.SplitContext{})
				// rootNode := &OperNode{oper: "FS", leftElems: []*lang.Elem{elems[0]}, prior: 111}
				// treeRoot := &LineTree{Tree: rootNode, Parents: []*OperNode{rootNode}}
				tres, err := Line2tree(elems, nil)
				if err != nil {
					return nil, err
				}
				ltree := tres.Tree
				// println("$== split format", ltree.oper)
				// PrintONode(ltree, 0)
				// operTree := ltree.rightNode
				// operElems := ltree.rightElems
				// +> as a special oper case
				rnodes = append(rnodes, &FNode{Form: form, Node: &OperNode{
					oper: "+>", rightNode: ltree.rightNode, rightElems: ltree.rightElems,
					subStr: form}})
				// subExp, ok := OperSub(ltree.rightNode, ltree.rightElems)
				// if !ok {
				// 	return nil, fmt.Errorf("format string: error in build oper tree from: `%s` ", string(expText))
				// }
				// res = append(res, &FPart{Node: subExp, Text: string(subStr[0:subLen])})
				subLen = 0
				continue
			}
			subStr[subLen] = c
			if formLen == 0 {
				subLen++
			} else {
				formLen++
			}
		}
	}
	if curLen > 0 {
		next := string(cur[0:curLen])
		nelems := []*lang.Elem{&lang.Elem{Type: Lt.Text, Text: next}}

		rnodes = append(rnodes, &FNode{Form: "", Node: &OperNode{oper: "T", rightElems: nelems}})
	}
	return &FString{Elems: res, FNodes: rnodes}, nil
}
