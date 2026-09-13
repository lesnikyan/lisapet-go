package cases

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/lesnikyan/lisapet-go/nodes"
)

var SubPartErr = errors.New("error in subpart of keywd case")

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

var NoNumElem = errors.New("Not number elem to parse")

func ParseInt(el *lang.Elem) (int64, error) {
	// le := node.leftElems[0]
	if el.Type != Lt.Num {
		return 0, NoNumElem
	}
	return strconv.ParseInt(el.Text, 10, 64)
}

func DotCaseOper(node *OperNode) *nodes.OperDot {
	if node.leftNode != nil || node.rightNode != nil {
		return &nodes.OperDot{} // not float
	}
	return nil
}

func DotCaseNum(node *OperNode) *nodes.ValExpr {
	if len(node.leftElems) > 1 || len(node.rightElems) > 1 {
		return nil // not-number case
	}

	sb := &strings.Builder{}
	// get int part
	if len(node.leftElems) == 1 {
		el := node.leftElems[0]
		if el.Type != Lt.Num {
			return nil // not number
		}
		sb.WriteString(el.Text)
	}
	sb.WriteString(".")
	// get decimal part
	if len(node.rightElems) == 1 {
		el := node.rightElems[0]
		if el.Type != Lt.Num {
			return nil // not number
		}
		sb.WriteString(el.Text)
	}
	nstr := sb.String()
	nval, err := strconv.ParseFloat(nstr, 64)
	if err != nil {
		return nil
	}
	return &nodes.ValExpr{Val: nval}
}

// func ProcOperTree(rNode *OperNode) (base.OperExpr, bool) {
func ProcOperTree(rNode *OperNode) (base.Expression, bool) {
	var expr base.OperExpr
	oper := rNode.oper
	// println("$$PROP#0:", oper)
	// PrintONode(rNode, 0)
	switch oper {
	case "=":
		expr = &nodes.OperAssign{Oper: oper}
	case ".":
		fnum := DotCaseNum(rNode)
		if fnum != nil {
			return fnum, true
		}
		expr = &nodes.OperDot{} // lambda
	case "*", "/", "+", "-", "**", "^/", "<<", ">>", "%", "|", "&", "^":
		expr = &nodes.OperBin{Oper: OperByStr(oper)}
	case "==", "!=", "<", "<=", ">", ">=", "&&", "||":
		expr = &nodes.OperBin{Oper: OperByStr(oper)}
	case "+=", "-=", "*=", "/=", "%=":
		expr = &nodes.OperBinAssign{Oper: OperByStr(oper)}
	case "=~", "?~", "/~":
		expr = &nodes.OperBin{Oper: OperByStr(oper)} // Rx opers
	case "..":
		expr = nodes.NewDots2(nil, nil)
	case ":":
		expr = &nodes.OperColon{Oper: OperByStr(oper)}
	case "->":
		expr = &nodes.LambdaExpr{} // lambda
	case "<-":
		expr = &nodes.LeftArrow{} // L-arrow
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
	// fmt.Printf("POT#2 %v (%T: %v) \n", oper, expr, expr)
	lArg, lok := OperSub(rNode.leftNode, rNode.leftElems)
	if lok {
		expr.SetLeft(lArg)
	} else {
		if slices.Contains(unary, oper) {
			expr = &nodes.UnaryLeft{Oper: OperByStr(oper)}
		}
	}
	rArg, rok := OperSub(rNode.rightNode, rNode.rightElems)
	// fmt.Printf("POT#4 %v (%T: %v) \n", oper, rArg, rArg)
	if rok {
		expr.SetRight(rArg)
	}
	// fmt.Println("POT#10:", "L:", nodes.OperArgsInfo(lArg), lok, "|| R:", nodes.OperArgsInfo(rArg), rok)

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

// list[index], list[:], dict[key], func(args), StructType{args}
func BracketsWithLeft(oper string, lexp base.Expression, subs base.Expression) (base.Expression, bool) {
	// fmt.Printf(" ??Coll([>> `%s` (%T, %v) \n", oper, subs, subs)
	seq, okc := subs.(*nodes.SequenceComma)
	switch oper {
	case "(":
		// func call
		fEx := lexp
		var fArgs []base.Expression
		if subs != nil {
			// has args
			if okc {
				// has more 1 arg
				fArgs = seq.Subs
			} else {
				fArgs = []base.Expression{subs}
			}
		}
		return nodes.NewFuncCall(fEx, fArgs), true

	case "[":
		// fmt.Printf("BWL#3 %s  %T %v\n", oper, lexp, lexp)
		// if pref, ok := lexp.(*nodes.ValExpr); ok {
		// 	// pp := pref.Val
		// 	fmt.Println(">> Pref", pref.Val, pref)
		// }

		switch subex := subs.(type) {
		case *nodes.OperColon:
			// Slice nn[a : b]
			// fmt.Printf(" Brackets [OperColon] (%T, %v) \n", subex, subex)
			return nodes.NewSlice(lexp, subex), true
		default:
			// Elem of collection: var[index|key]
			// fmt.Printf(" Brack [??] (%T, %v) \n", subex, subex)
			return &nodes.ColElemExpr{Col: lexp, Key: subs}, true
		}
	case "{":
		// fmt.Printf("BWL#4 struct %s  %T %v\n", oper, lexp, lok)
		var args []*nodes.OperColon
		var name string
		nexp, ok := lexp.(*nodes.VarExpr)
		if !ok {
			return nil, false
		}
		name = nexp.GetName()

		switch subex := subs.(type) {
		case *nodes.OperColon:
			// struct: Type{name : val}
			// fmt.Printf(" Brackets {OperColon} (%T, %v) \n", subex, subex)
			args = []*nodes.OperColon{subex}
			// return nil, true
		case *nodes.SequenceComma:
			// struct: Type{name : val}
			// fmt.Printf(" Brackets {,,} (%T, %v) \n", subex, subex)
			args = make([]*nodes.OperColon, len(subex.Subs))
			for i, sb := range subex.Subs {
				if colx, ok := sb.(*nodes.OperColon); ok {
					args[i] = colx
				} else {
					return nil, false
				}
			}
			// return nil, true
		case *nodes.EmptyExpr:
			// Elem of collection: var[index|key]
			// fmt.Printf(" Brack {} (%T, %v) \n", subex, subex)
			args = []*nodes.OperColon{}
		default:
			// Bad case
			// fmt.Printf(" Brack default (%T, %v) \n", subex, subex)
			args = []*nodes.OperColon{}
			// return nil, false
		}
		return nodes.NewStructConstr(name, args), true
	}
	return nil, false
}

func CaseBytes(rNode *OperNode, subs base.Expression) (base.Expression, bool) {
	var pref string
	if len(rNode.leftElems) > 0 {
		pref = rNode.leftElems[0].Text
	} else {
		pref = "0x"
	}
	// fmt.Printf("BytesExpr ##1 subs (%T, %v) \n", subs, subs)
	res := nodes.NewBytesExpr(pref, subs)
	err := res.Parse()
	if err != nil {
		panic("Bytes case has bad syntax: " + err.Error())
	}
	return res, true
}

func BracketsExpr(rNode *OperNode) (base.Expression, bool) {
	oper := rNode.oper
	bt := nodes.GetBrType(rNode.oper)
	subs, rok := OperSub(rNode.rightNode, rNode.rightElems)
	lexp, lok := OperSub(rNode.leftNode, rNode.leftElems)
	seq, okc := subs.(*nodes.SequenceComma)
	// fmt.Printf("BracketsExpr >>> \n")
	// PrintONode(rNode, 0)
	if lok {
		if _, ok := lexp.(*nodes.ValExpr); ok {
			if len(rNode.leftElems) > 0 && rNode.leftElems[0].Type == Lt.Num {
				// Bytes case
				// fmt.Printf("##00 %s [%T] \n", rNode.leftElems[0].Text, subs)
				switch subs.(type) {
				case *nodes.ValExpr, *nodes.VarExpr:
					// subs = &nodes.ByteLine{Src: rNode.rightElems}
					subs = nodes.NewByteLine(rNode.rightElems)
				case *nodes.NumField:
					// do nothing
				default:
					// fmt.Printf("##2 %s [%T] \n", rNode.leftElems[0].Text, subs)
					// subs = &nodes.ByteLine{Src: rNode.rightElems[0].Text}
					if len(rNode.leftElems) == 0 {
						subs = &nodes.ByteLine{Src: ""}
					} else {
						subs = nodes.NewByteLine(rNode.rightElems)
					}
				}
				// fmt.Printf("##1 %s [%T] \n", rNode.leftElems[0].Text, subs)
				return CaseBytes(rNode, subs)
			}
		}
		if !rok {
			subs = nil
		}
		return BracketsWithLeft(oper, lexp, subs)
	}
	// fmt.Printf("BEx#1 < %s >  L(%T %v), R(%T %v) comm: %v\n", oper, lexp, lok, subs, rok, okc)

	if !okc {
		// non comma-separated cases: (a + b), [1 .. 5]
		switch oper {
		case "(":
			return &nodes.Brackets{Sub: subs, Type: bt}, true
		case "[":
			// fmt.Printf("BEx#11 < %s >  L(%T %v), R(%T %v) comm: %v\n", oper, lexp, lok, subs, rok, okc)
			switch subex := subs.(type) {
			case *nodes.Dots2Expr:
				return subex.GetNumSeq(), true
			case *nodes.NumField:
				// println("NumField case#1")
				return CaseBytes(rNode, subex)
			}
		}

		// other non-comma-separated cases
		// possible empty, 1-elem in sequence, num-seq [a..b]
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
	// fmt.Printf("PET#0 %v\n", rNode.oper)
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
	last := rNode
	for node != nil && node.oper == rNode.oper {
		last = node
		sub, ok := OperSub(node.rightNode, node.rightElems)
		if ok {
			elems = append(elems, sub)
		}
		node = node.leftNode
	}
	first, ok := OperSub(last.leftNode, last.leftElems)
	// fmt.Printf("SeqSub#3 %T %v %T", last, ok, first)
	if ok {
		elems = append(elems, first)
	}
	slices.Reverse(elems)
	// fmt.Println("Seq#1:", elems)
	// sub, ok := OperSub(node.leftNode, node.leftElems)
	return elems, true
}

func OperSub(node *OperNode, elems []*lang.Elem) (base.Expression, bool) {
	// fmt.Println("OperSubs#1", node == nil, len(elems), FPrintElems(elems))
	// PrintONode(node, 0)
	if node != nil {
		return ProcExprTree(node)
	} else if len(elems) > 0 {
		return ProcSubElems(elems)
	}
	return nil, false
}

func ProcSubElems(elems []*lang.Elem) (base.Expression, bool) {
	xpr, ok := CaseVal(elems)
	// fmt.Println("ProcSubElems", ok, xpr)
	if ok {
		return xpr, true
	}
	return CaseVar(elems)
}

// if line have unclosed brackets || binary opers without right arg
type UnclosedExpr struct {
	Prev []*lang.Elem
	Tree *LineTree
	// Parents []*OperNode
}

func (q UnclosedExpr) Do(base.Context) error { return nil }
func (q UnclosedExpr) Get() *base.Val        { return nil }

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

func IsLKWord(elems []*lang.Elem, prevTree *LineTree) bool {
	if prevTree != nil && len(prevTree.Tree.leftElems) > 0 {
		_, ok := kwMap[prevTree.Tree.leftElems[0].Text]
		return ok
	}
	e1 := elems[0]
	if e1.Type != Lt.Word {
		return false
	}
	_, ok := kwMap[e1.Text]
	return ok
}

func CaseLKeyword(elems []*lang.Elem) (*CaseRes, bool) {
	if elems[0].Type != Lt.Word || !IsLKWord(elems, nil) {

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
