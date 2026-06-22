package cases

import (
	"errors"
	"slices"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/lesnikyan/lisapet-go/nodes"
)

var SubPartErr = errors.New("error in subpart of keywd case")

// just simple var = val
// (base.Expression, [][]*lang.Elem, bool)
// func CaseAssign(ee []*lang.Elem) (*CaseRes, bool) {
// 	// matching part
// 	spres, err := OperSplit(ee)
// 	if err != nil || spres.Lowest == -1 || ee[spres.Lowest].Text != "=" {
// 		return nil, false
// 	}
// 	// ind := spres.Lowest
// 	// if ee[ind].Text != "=" {
// 	// 	return nil, false
// 	// }
// 	// parsing part
// 	parts := [][]*lang.Elem{ee[:spres.Lowest], ee[spres.Lowest+1:]}
// 	return &CaseRes{Expr: &nodes.OperAssign{}, Subs: parts}, true
// }

func SkipSpaces(elems []*lang.Elem) []*lang.Elem {
	res := []*lang.Elem{}
	for _, ee := range elems {
		if ee.Type != Lt.Space {
			res = append(res, ee)
		}
	}
	return res
}

func _AssignSubs(elems []*lang.Elem, spres *SplittedRes) {
	if len(spres.Others) > 0 {
		if spres.Others[0][0] < spres.Lowest {
			// posibly multiassign
		} else {
			// leftEls := SkipSpaces(elems[:spres.Lowest])
			// Interpret(leftEls)
		}
	}
}

var binOpers = strings.Split("+ - * / ** ^/ % | & || && == !=", " ")
var binAssignOpers = strings.Split("+= -= *= /= %=", " ")

func ProcOperTree(rNode *OperNode) (base.BinOperExpr, bool) {
	var expr base.BinOperExpr
	oper := rNode.oper
	switch oper {
	case "=":
		expr = &nodes.OperAssign{}
	case "->":
		expr = &nodes.OperBin{} // lambda
	case "<-":
		expr = &nodes.OperBin{} // Larrow
	case "$":
		expr = &nodes.OperBin{} // func-apply
	case "?:":
		expr = &nodes.OperBin{} // short-triple
	// case "*","/","+","-","*","^/","|":
	default:
		switch {
		case slices.Contains(binOpers, oper):
			expr = &nodes.OperBin{}
		case slices.Contains(binAssignOpers, oper):
			expr = &nodes.OperBin{} // BinAssign
		}

	}

	return expr, expr != nil
}

func CaseBinOper(ee []*lang.Elem) (base.Expression, bool) {
	rNode, err := Line2tree(ee) // *OperNode, error
	if err != nil {
		// no oper fount out of brackets
		return nil, false
	}
	// parts := [][]*lang.Elem{ee[:spres.Lowest], ee[spres.Lowest+1:]}
	expr, ok := ProcOperTree(rNode)
	if !ok {
		return nil, false
	}
	return expr, ok
	// foundT := ee[spres.Lowest].Text
	// return &CaseRes{Expr: expr, Subs: parts}, true
}

// func CaseBinOper(ee []*lang.Elem) (*CaseRes, bool) {
// 	spres, err := Line2tree(ee)
// 	if err != nil || spres.Lowest == -1 {
// 		// no oper fount out of brackets
// 		return nil, false
// 	}
// 	parts := [][]*lang.Elem{ee[:spres.Lowest], ee[spres.Lowest+1:]}
// 	var expr base.Expression
// 	foundT := ee[spres.Lowest].Text
// 	switch foundT {
// 	case "=":
// 		expr = &nodes.OperAssign{}
// 	case "->":
// 		expr = &nodes.TODOExpr{}
// 	case "<-":
// 		expr = &nodes.TODOExpr{}
// 	case "$":
// 		expr = &nodes.TODOExpr{}
// 	case "?:":
// 		expr = &nodes.TODOExpr{}
// 	// case "*","/","+","-","*","^/","|":
// 	default:
// 		switch {
// 		case slices.Contains(binOpers, foundT):
// 			expr = &nodes.OperBin{}
// 		case slices.Contains(binAssignOpers, foundT):
// 			expr = &nodes.TODOExpr{} // BinAssign
// 		}

// 	}
// 	return &CaseRes{Expr: expr, Subs: parts}, true
// }

func InterpretSpRes(elems []*lang.Elem, spres *SplittedRes) {

}

/*
-- leading keywords
if condition
if expr; condition
else
else if condition
for
while
match
func name(args)
func obj:Type name(args)
struct name fields
break
continue
return
enum
grup

*/

// func List2Keys[T comparable](data []T) map[T]bool {
// 	r := make(map[T]bool, len(data))
// 	for _, k := range data {
// 		r[k] = true
// 	}
// 	return r
// }

// var keywords = strings.Split("func|if|for|while|match|enum|grup|struct|else|import|return|const|run", "|")
// var kwMap = List2Keys(keywords)

func IsLKWord(elems []*lang.Elem) bool {
	e1 := elems[0]
	if e1.Type != Lt.Word {
		return false
	}
	_, ok := kwMap[e1.Text]
	return ok
}

func CaseLKeyword(elems []*lang.Elem) (*CaseRes, bool) {
	if elems[0].Type != Lt.Word || !IsLKWord((elems)) {

	}
	var expr base.Expression
	parts := [][]*lang.Elem{elems[1:]}

	return &CaseRes{Expr: expr, Subs: parts}, true
}

func ParseKWSub(elems []*lang.Elem) ([]base.Expression, error) {
	return nil, nil
}

func SubIf(elems []*lang.Elem) ([]base.Expression, error) {
	// var err error = nil //  SubPartErr
	// detect subs
	// spres, sperr := Line2tree(elems)
	// if sperr != nil {
	// 	return nil, sperr
	// }
	// if spres.Lowest == -1 {
	// 	// no oper, just solid expr
	// 	return nil, nil // fix result
	// }
	// if elems[spres.Lowest].Text == ";" {
	// 	// has extra expression before condition
	// 	// 1. split elems to sub-expressions
	// 	// 2. Inrerpret each sub
	// 	return nil, nil // TODO: fix resilt
	// }
	expr := &nodes.IfExpr{}
	exprs := []base.Expression{expr}
	return exprs, nil
}

/*
keyword: kword [others]
bin-operator: left <oper> right
separated-sequence: a, b, c // a ; b ; c //
@spec-expr

obj.field
funcCall()
in-brackets: [], (), {}
collElem[key]
var, val

*/

var caseList = []func(ee []*lang.Elem) (*CaseRes, bool){
	CaseLKeyword, // CaseBinOper,
}

func InterpretSups(supr base.SupExpr, cres *CaseRes) error {
	// var supr base.SupExpr
	// supr = resExp.(base.SupExpr)
	for _, subp := range cres.Subs {
		subExp, err := InterpretLine(subp)
		if err != nil {
			return err
		}
		supr.Add(subExp.Expr)
	}
	return nil
}

/*
interpret sequence: fing case, parse, make Expression node
*/
func InterpretLine(elems []*lang.Elem) (*CaseRes, error) {
	var resExp base.Expression
	for _, cfun := range caseList {
		cres, found := cfun(elems)
		resExp = cres.Expr
		if !found {
			continue
		}
		if len(cres.Subs) > 0 {
			// var supr base.SupExpr
			// supr = resExp.(base.SupExpr)
			// for _, subp := range cres.Subs {
			// 	subExp, err := InterpretLine(subp)
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// 	supr.Add(subExp.Expr)
			// }
			InterpretSups(resExp.(base.SupExpr), cres)
		}
		return &CaseRes{Expr: resExp}, nil
	}
	return nil, nil
}
