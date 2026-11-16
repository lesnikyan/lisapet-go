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
	Elems  []*FPart
	FNodes []*FNode
}

// ' a {a + b : <+05.4f} b'

var _openBrs = "({["
var _closeBrs = "}])"
var _spaces = " \t"

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
			if strings.ContainsRune(_openBrs, c) {
				inBr++
			} else if strings.ContainsRune(_closeBrs, c) {
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
				expText := runes[expStart : expStart+subLen]
				form := ""
				if formLen > 0 {
					rawrr := runes[expStart+subLen+1 : expStart+subLen+formLen]
					rr := make([]rune, len(rawrr))
					i := 0
					for _, r := range rawrr {
						if strings.ContainsRune(_spaces, r) {
							continue
						}
						rr[i] = r
						i++
					}
					form = string(rr[0:i])
				}
				// fmt.Printf("  fpartN: `%s`, form: `%s` \n", string(expText), form)
				//...
				elems := parser.SplitLine(expText, &parser.SplitContext{})
				tres, err := Line2tree(elems, nil)
				if err != nil {
					return nil, err
				}
				ltree := tres.Tree
				// println("$== split format", ltree.oper)
				// PrintONode(ltree, 0)
				// +> as a special oper case
				rnodes = append(rnodes, &FNode{Form: form, Node: &OperNode{
					oper: "+>", rightNode: ltree.rightNode, rightElems: ltree.rightElems,
					subStr: form}})
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
