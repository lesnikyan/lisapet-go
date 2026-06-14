package tree

import (
	// lang "github.com/lesnikyan/lisapet-go/"
	"errors"
	"slices"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	lang "github.com/lesnikyan/lisapet-go/lang"
	lt "github.com/lesnikyan/lisapet-go/lang/lt"

	// par "github.com/lesnikyan/lisapet-go/parser"
	ob "github.com/lesnikyan/lisapet-go/objects"
)

func list2Keys[T comparable](data []T) map[T]bool {
	r := make(map[T]bool, len(data))
	for _, k := range data {
		r[k] = true
	}
	return r
}

var noMatchErr = errors.New("No matches line")

type MatchingTree struct {
	subs   []*MatchingTree
	filter func([]*lang.Elem) bool
}

func FSingle(elems []*lang.Elem) bool {
	return len(elems) == 1
}

var keywords = strings.Split("func|if|for|while|match|enum|grup|struct|else|import|return|const|run", "|")
var kwMap = list2Keys(keywords)

// definition or control leading keyword like: func, if, struct, etc..
func FLKeyword(elems []*lang.Elem) bool {
	e1 := elems[0]
	if e1.Type != lt.Word {
		return false
	}
	_, ok := kwMap[e1.Text]
	return ok
}

var n = `
Filter tree

single
	single-command
	val
	var

solid
	funcCall()
	obj.member
	container[elem]
	brackets
		[list]
	str-prefix (re"")

non-solid
	L-keyword
	bin-oper
	un-oper

`

const (
	// Line cases
	LcEmpty  = 0
	LcSingle = 1
	LcOper   = 2
	LcSolid  = 3
	LcError  = 101
)

var solidOpers = strings.Split(". ~> ( ) [ ] { }", " ")

// simple detection from left without precedences
func CheckSingle(elems []*lang.Elem) int {
	var brn = 0 // depth in brackets
	opBr := "([{"
	clBr := "}])"
	for _, el := range elems {
		et := el.Text
		if strings.Contains(opBr, et) {
			brn++
			continue
		}
		if strings.Contains(clBr, et) {
			brn--
			if brn < 0 {
				return LcError
			}
			continue
		}
		if brn > 0 {
			// skip because in brackets
			continue
		}
		if el.Type != lt.Oper {
			continue
		}
		// out of brackets
		if el.Type == lt.Oper && slices.Contains(solidOpers, et) {
			return LcOper
		}
	}
	return LcSolid
}

/*
1. split elem list to solid parts: var/val, coll[elem], [coll, consrt], func(call), obj.field,  etc.
2. detect unary operators,
3. group around bin oberator by precedence
4. detect type and split parts
rrr = [1, 2, 3] + obj.prop1()[5] * (foo()) + data['elem'].bar().prop2[3].elem4

*/

func elem2node(elems []*lang.Elem) (base.Expression, error) {
	return nil, noMatchErr
}

func Build(code []*lang.CLine) *ob.Module {
	// lang.Elem{}
	curInd := 0
	root := ob.NewModule(nil)
	var curParent base.SupExpr = root
	waitingChild := false
	for _, cline := range code {
		elems := cline.Elems
		node, err := elem2node(elems)
		if err != nil {
			// doing with error
		}
		if cline.Indent > curInd {
			waitingChild = true
		}
		if waitingChild {
			curParent.Add(node)
		} else {

		}

	}
}
