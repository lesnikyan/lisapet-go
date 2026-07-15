package cases

import (
	"errors"
	"fmt"
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

func FPrintElems(elems []*lang.Elem) string {
	ss := make([]string, len(elems))
	// println("lenEl", len(elems))
	for i := 0; i < len(elems); i++ {
		// println("?=", elems[i].Text)
		ss[i] = elems[i].Text
	}
	stt := strings.Join(ss, ",")
	return fmt.Sprintf("`%s`", stt)
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

const (
	opSemicol = ";"
	opDot     = "."
	opComm    = ","
	opColon   = ":"
	// op = ""
)

var binOpers = strings.Split("+ - * / ** ^/ % | & || && == !=", " ")
var binAssignOpers = strings.Split("+= -= *= /= %=", " ")

// func OperIndex(s)nodes.Opid{

// }
func OperByStr(s string) *nodes.Oper {
	return &nodes.Oper{Sign: s, Id: nodes.OperIndex(s)}
}

func ProcOperTree(rNode *OperNode) (base.OperExpr, bool) {
	var expr base.OperExpr
	oper := rNode.oper
	println("$$PROP#0:", oper)
	switch oper {
	case "=":
		expr = &nodes.OperAssign{Oper: oper}

	case "*", "/", "+", "-", "**", "^/", "<<", ">>", "%":
		expr = &nodes.OperBin{Oper: OperByStr(oper)}
	case "==", "!=", "<", "<=", ">", ">=", "&&", "||":
		expr = &nodes.OperBin{Oper: OperByStr(oper)}
	case "+=", "-=", "*=", "/=", "%=":
		expr = &nodes.OperBinAssign{Oper: OperByStr(oper)}
	case ":":
		expr = &nodes.OperColon{Oper: OperByStr(oper)}
	case "->":
		expr = &nodes.OperBin{} // lambda
	case "<-":
		expr = &nodes.OperBin{} // L-arrow
	case "$":
		expr = &nodes.OperBin{} // func-apply
	case "?:":
		expr = &nodes.OperBin{} // short-triple
	default:
		switch {
		case slices.Contains(binOpers, oper):
			expr = &nodes.OperBin{}
			// case slices.Contains(binAssignOpers, oper):
			// 	expr = &nodes.OperBin{} // BinAssign
		}

	}
	// println("$$PROP#1=:", expr)
	fmt.Printf("POT#2 %v (%T: %v) \n", oper, expr, expr)
	lArg, lok := OperSub(rNode.leftNode, rNode.leftElems)
	if lok {
		expr.SetLeft(lArg)
	} else {
		if slices.Contains(unary, oper) {
			expr = &nodes.UnaryLeft{Oper: OperByStr(oper)}
		}
	}
	rArg, rok := OperSub(rNode.rightNode, rNode.rightElems)
	if rok {
		expr.SetRight(rArg)
	}
	fmt.Println("POT#10:", nodes.OperArgsInfo(lArg), lok, nodes.OperArgsInfo(rArg), rok)

	return expr, expr != nil
}

func SubSeq(parent string, node base.Expression) (base.Expression, bool) {

	switch exx := node.(type) {
	case *nodes.SequenceComma:
		//collexrions
		switch parent {
		case "(":
			// tuple
			return &nodes.TupleExpr{Seq: exx}, true
		case "[":
			// list
			return &nodes.ListExpr{Seq: exx}, true
		case "{":
			// dict
			return &nodes.DictExpr{Seq: exx}, true
		}
	case *nodes.SequenceSemicolon:
		// generators, etc
		switch parent {
		case "[":
			// list comprehension
			return &nodes.MockExpr{}, false
		case "{":
			// dict compr
			return &nodes.MockExpr{}, false
		case "(:":
			// THINK: do we need tuple-comprehension? tuple([ ; ; ]) looks enough
			// generator
			return &nodes.MockExpr{}, false
		}
	default:
		// single elem in sequence?
		return nil, false
	}
	return nil, false
}

func BracketsExpr(rNode *OperNode) (base.Expression, bool) {
	oper := rNode.oper
	bt := nodes.GetBrType(rNode.oper)
	subs, rok := OperSub(rNode.rightNode, rNode.rightElems)
	lexp, lok := OperSub(rNode.leftNode, rNode.leftElems)
	fmt.Printf("BEx#1 < %s >  L(%T %v), R(%T %v)\n", oper, lexp, lok, subs, rok)
	seq, okc := subs.(*nodes.SequenceComma)
	switch oper {
	case "(":
		if lok {
			// if lok: func call
			return nil, false
		}
		if okc {
			// tuple here
		} else {
			return &nodes.Brackets{Sub: subs, Type: bt}, true
		}
		// sub-case of generator: `(: expr ; ..)`

	case "[":
		if okc {
			// list, not sure
		}
		// fmt.Printf("BEx#3 %s  %T %v\n", oper, lexp, lok)
		if lok {
			return &nodes.ColElem{Col: lexp, Key: subs}, true
		}
	}
	fmt.Printf("PET#2 %T\n", seq)
	if !okc {
		// other non-comma-separated cases
		// seq, ok := subs.(*nodes.SequenceComma)
		// possible empty or 1-elem in sequence

		subEx := []base.Expression{}
		if rok {
			subEx = append(subEx, subs)
		}
		sub1 := &nodes.SequenceComma{Subs: subEx}
		seq = sub1
	}
	// collection | comprehension case
	return SubSeq(oper, seq)
}

func ProcExprTree(rNode *OperNode) (base.Expression, bool) {
	// fmt.Println("DEB+++1", rNode.oper)
	oper := rNode.oper
	fmt.Printf("PET#0 %v\n", oper)
	switch rNode.oper {
	case "(", "[", "{":
		return BracketsExpr(rNode)
	case ",", ";":
		return ProcSequence(rNode)
	case "\\":
		// lambda \ arg -> ..
		return nil, false
	default:
		return ProcOperTree(rNode)
	}
}

func ProcSequence(rNode *OperNode) (base.Expression, bool) {
	var expr base.SequenceExpr
	switch rNode.oper {
	case ",":
		expr = &nodes.SequenceComma{}
	case ";":
		expr = &nodes.SequenceSemicolon{}
	}
	if expr == nil {
		return nil, false
	}

	subs, ok := SeqSubs(rNode)
	if ok {
		expr.SetSubs(subs)
	}
	return expr, true
}

// [[[1, 2], 3], 4],
func SeqSubs(rNode *OperNode) ([]base.Expression, bool) {
	node := rNode
	elems := []base.Expression{}
	// if node.leftNode == nil && len(node.leftElems) == 0 {
	// 	//
	// }
	// first, ok := OperSub(node.leftNode, node.leftElems)
	// if !ok {
	// 	// bad first elem
	// 	return nil, false
	// }
	last := rNode
	for node != nil && node.oper == rNode.oper {
		last = node
		sub, ok := OperSub(node.rightNode, node.rightElems)
		// if !ok {
		// 	// end of sequence
		// }
		if ok {
			elems = append(elems, sub)
		}
		// if node.leftNode == nil || node.leftNode.oper != node.oper {
		// 	// start of seq has found
		// 	first, ok := OperSub(node.leftNode, node.leftElems)
		// 	if ok {
		// 		elems = append(elems, first)
		// 	}
		// 	// TODO: do we need 1-st empty expression in sequences?
		// 	// ( ; a; b) ?
		// 	break
		// }
		node = node.leftNode
	}

	// if node.leftNode != nil ||  {
	// start of seq has found
	first, ok := OperSub(last.leftNode, last.leftElems)
	// fmt.Printf("SeqSub#3 %T %v %T", last, ok, first)
	if ok {
		elems = append(elems, first)
	}
	// }
	// elems = append(elems, first)
	slices.Reverse(elems)
	fmt.Println("Seq#1:", elems)
	// sub, ok := OperSub(node.leftNode, node.leftElems)
	return elems, true
}

func OperSub(node *OperNode, elems []*lang.Elem) (base.Expression, bool) {
	fmt.Println("OperSubs#1", node, len(elems), FPrintElems(elems))
	if node != nil {
		return ProcExprTree(node)
	} else if len(elems) > 0 {
		return ProcSubElems(elems)
	}
	return nil, false
}

func ProcSubElems(elems []*lang.Elem) (base.Expression, bool) {
	xpr, ok := CaseVal(elems)
	if ok {
		return xpr, true
	}
	return CaseVar(elems)
}

// func ProcOperSubNode(node *OperNode) (base.Expression, bool) {
// 	// possible cases: bin-oper, unary-oper,
// 	// brackets, collection,
// 	// func-call, collect[elem], struct-constr
// 	if node.oper == "" {
// 		// processing not opers
// 	}
// 	if strings.Contains(OBRS, node.oper) {
// 		// brackets
// 		switch node.oper {
// 		case "(":
// 			// grouping, generator
// 		case "[":
// 			// list, list-comprehension
// 		case "{":
// 			// dict
// 		case "\\":
// 			// lambda
// 		}
// 	}
// 	return nil, false
// }

func CaseBinOper(ee []*lang.Elem) (base.Expression, bool) {

	tres, err := Line2tree(ee, nil) // *OperNode, error
	if err != nil {
		// no oper fount out of brackets
		return nil, false
	}
	if !tres.Finished {
		return &UnclosedExpr{Prev: ee, Tree: tres}, false
	}
	rNode := tres.Tree
	// parts := [][]*lang.Elem{ee[:spres.Lowest], ee[spres.Lowest+1:]}
	expr, ok := ProcOperTree(rNode)
	if !ok {
		return nil, false
	}
	return expr, ok
	// foundT := ee[spres.Lowest].Text
	// return &CaseRes{Expr: expr, Subs: parts}, true
}

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
interpret sequence: find case, parse, make Expression node
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
			InterpretSups(resExp.(base.SupExpr), cres)
		}
		return &CaseRes{Expr: resExp}, nil
	}
	return nil, nil
}

func InterpretLine2(relems []*lang.Elem) {
	elems := SkipSpaces(relems)
	if IsLKWord(elems) {
		KWordExp(elems)
	} else {
		// operator operator expr
		CaseBinOper(elems)
	}
}

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

// if line have unclosed brackets || binary opers without right arg
type UnclosedExpr struct {
	Prev []*lang.Elem
	Tree *LineTree
	// Parents []*OperNode
}

func (q UnclosedExpr) Do(base.Context) error { return nil }
func (q UnclosedExpr) Get() *base.Val        { return nil }

func CaseIf(tree *LineTree) (base.Expression, error) {
	return nil, nil
}

func CaseFunc(sig *LineTree, pref *LineTree) (base.Expression, error) {
	return nil, nil
}

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

	case kFunc:
		if len(subElems) > 6 {
			// possible mehod def: `func inst:Type Name()`
			var prefTree *LineTree
			var err1 error
			sigIndex := 1
			if elems[2].Type == Lt.Oper && elems[2].Text == opColon {
				prefTree, err1 = Line2tree(elems[1:4], nil)
				if err1 != nil {
					//
				}
				sigIndex = 4
				// TODO: add prefTree to kwRoot
			}
			sigTree, err2 := Line2tree(elems[sigIndex:], kwRoot)
			if err2 != nil {
				//
			}
			return CaseFunc(sigTree, prefTree)
		}
	case kWhile:
	case kMatch:
	case kEnum:
	case kGrup:
	case kStruct:
	case kImport:
	case kReturn:
	case kBreak:
	case kContinue:
	case kConst:
	case kRun:
	}
	return nil, nil
}
