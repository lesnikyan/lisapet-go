package cases

import (
	"errors"
	"slices"
	"strings"

	"github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
)

// import	"github.com/lesnikyan/lisapet-go/nodes"

func List2Keys[T comparable](data []T) map[T]bool {
	r := make(map[T]bool, len(data))
	for _, k := range data {
		r[k] = true
	}
	return r
}

var keywords = strings.Split("func|if|for|while|match|enum|grup|struct|else|break|continue|import|return|const|run", "|")
var maxKWlen = 7
var kwMap = List2Keys(keywords)

func IsKeywd(tx string) bool {
	ltx := len(tx)
	if ltx < 2 || ltx > maxKWlen {
		return false
	}
	_, ok := kwMap[tx]
	return ok
}

var obrs = "([{"
var cbrs = ")]}"
var OBRS = "([{" //

// map of valid close : open brs
var brmap = map[string]string{
	")": "(",
	"]": "[",
	"}": "{",
}

var _operPriorStr = `( ) [ ] { } 1 . 1 ~> 1 ... 1 -x ! ~ 1 ** ^/ 1 * / % 1 + - 1` +
	`<< >> 1 =~ ?~ /~1 < <= > >= !> ?> !?> 1 == != 1 & 1 ^ 1 | 1 :: 1 && 1 || 1 \\ 1 ->` +
	` 1 @ 1 $ 1 ?: 1 : 1 ? 1 , 1 .. 1 <- 1 @! 1 = += -= *= /= %= 1 ; 1 !: :? => 1 /: `

var _operPriorStr2 = `1 . 1 ~> 1 ... 1 ** ^/ 1 * / % 1 + - 1` +
	`<< >> 1 =~ ?~ /~1 < <= > >= !> ?> !?> 1 == != 1 & 1 ^ 1 | 1 :: 1 && 1 || 1 \\ 1 ->` +
	` 1 @ 1 $ 1 ?: 1 : 1 ? 1 func: 1 fun= 1 , 1 .. 1 <- 1 @! 1 = += -= *= /= %= 1 ; 1 !: :? => 1 /: `

// spec: fun=

var _operPriorFBr = `1 . 1 ~> 1 ... 1 ** ^/ 1 * / % 1 + - 1` +
	`<< >> 1 =~ ?~ /~1 < <= > >= !> ?> !?> 1 == != 1 & 1 ^ 1 | 1 :: 1 && 1 || 1 \\ 1 ->` +
	` 1 @ 1 $ 1 ?: 1 : 1 ? 1 = += -= *= /= %= 1 , 1 .. 1 <- 1 @! 1 ; 1 !: :? => 1 /: `

var unary = strings.Split("- ! ~ +", " ")
var unaryR = strings.Split("... ~>", " ")
var seqSeprs = strings.Split(", ;", " ")
var kwPrior = 100000

func priorSet(opers string) [][]string {
	ss := strings.Split(_operPriorStr, "1")
	var res [][]string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		nn := strings.Split(s, " ")
		res = append(res, nn)
	}
	return res
}

var priors2 = priorSet(_operPriorStr2)
var priorsFuBr = priorSet(_operPriorFBr)
var priors = priorSet(_operPriorStr)

// func() [][]string {
// 	ss := strings.Split(_operPriorStr, "1")
// 	var res [][]string
// 	for _, s := range ss {
// 		s = strings.TrimSpace(s)
// 		nn := strings.Split(s, " ")
// 		res = append(res, nn)
// 	}
// 	return res
// }()

// var priors2 = func() [][]string {
// 	ss := strings.Split(_operPriorStr2, "1")
// 	var res [][]string
// 	for _, s := range ss {
// 		s = strings.TrimSpace(s)
// 		nn := strings.Split(s, " ")
// 		res = append(res, nn)
// 	}
// 	return res
// }()

type ints2 [2]int

type SplittedRes struct {
	Lowest   int     // index of lower operator
	Others   []ints2 // other operators [index, precedence by `priors`]
	Brackets []ints2
	Inners   []*SplittedRes // TODO: think about 1-pass work of split
}

func getPrior(priorNN [][]string, oper string) int {
	for i, nn := range priorNN {
		// log.Println("prr>>", oper, i, nn)
		if slices.Contains(nn, oper) {
			return i
		}
	}
	return -1
}

var unaryPrior = getPrior(priors2, "...") + 1
var solidExPrior = getPrior(priors2, "...")

var badBracketsErr = errors.New("bad set of brackets")

// type DDD interface {
// 	*OperNode | []*lang.Elem
// }

type OperNode struct {
	oper         string // if oper
	kword        string
	prior        int
	IsOper       bool
	IsBrackets   bool
	BracketOpen  string
	BracketClose string
	bracketType  string // "(" | "[" | "{"
	leftNode     *OperNode
	rightNode    *OperNode
	subNodes     []*OperNode // if sequence (,;) or keyword subs
	// elems      []*lang.Elem // if other expr
	// elems       []*lang.Elem // should be a left part in transfer to operator
	leftElems  []*lang.Elem // left operand
	rightElems []*lang.Elem //
}

/*
operands:
123, "string", pref"string"suf
varname
cont[key]
funcCall(args)
StructInst{a:1, b:2}
0x[1 2 f]
[list, ], (tuple, ), {dict, k:v}
\x->lambda
-1, !true, ~"{n}", nn...
*/

func (nd *OperNode) AddRElem(el *lang.Elem) {
	if nd.rightElems == nil {
		nd.rightElems = []*lang.Elem{}
	}
	nd.rightElems = append(nd.rightElems, el)
}

func (nd *OperNode) AddLElem(el *lang.Elem) {
	if nd.leftElems == nil {
		nd.leftElems = []*lang.Elem{}
	}
	nd.leftElems = append(nd.leftElems, el)
}

// func AddLeft[T *OperNode | *lang.Elem](targ *OperNode, elem T) {
// 	switch val := elem.(type) {
// 	case *OperNode:
// 		print(val)
// 	}
// }

// in function brackets foo(a,b,c=9)

/*
a * b - (c + d) / (2 / 5)
res:
Lowest: 3
Others: [[1, prec*], [9, prec/]]
Brackets: [[4, 8], [10, 14]]
*/

/*
1. oper precedence: 1 * 2 + foo(3) * obj.nn ** 2 - 2^/25 / stt.nn[2]
2. brackets: [] , [{}, {}, []] + [{}]
3. ,;-sequence
4. other cases: keyword, space-separated case, n|n|n-sequence, L-unary, R-unary
*/

// x, y = (foo(a=(1,2), b=5), 12), 13

type PriorKey int

const (
	_             PriorKey = iota
	PriorKeyBase           // basic priors
	PriorKeyFunBr          // func definition or call
	// PriorKeyMatch // case of match
)

var priorMap = map[PriorKey][][]string{
	PriorKeyBase:  priors2,
	PriorKeyFunBr: priorsFuBr,
}

func notEmptyLeft(node *OperNode) bool {
	return node.leftNode != nil || len(node.leftElems) > 0
}

func Line2tree(elems []*lang.Elem, prevTree *LineTree) (*LineTree, error) {
	opris := priors2 // here all opers, except unary
	brC := 0         // brackets depth
	var others []ints2 = []ints2{}
	var cur *lang.Elem
	var prev *lang.Elem
	// var rprev *lang.Elem // real prev
	var closeBr = false
	var rNode *OperNode
	var parents []*OperNode
	var cNode *OperNode // curent node
	opcx := 1
	if prevTree != nil {
		rNode = prevTree.Tree
		parents = prevTree.Parents
		cNode = parents[len(parents)-1]
		brC = prevTree.BracketsCount
		if prevTree.InFuncBr {
			opcx = 3
		}
	} else {
		rNode = &OperNode{prior: 10000, oper: "ЫХ"} // root node
		parents = []*OperNode{rNode}
		cNode = rNode
	}
	curPart := []int{} // indexes of elements
	var slashLambda bool = false

	var solidEnd = strings.Split(") ] } ... ~>", " ")
	for i, el := range elems {
		// curin = i
		tx := el.Text
		etp := el.Type
		// rprev = cur
		if etp == Lt.Space {
			continue
		}
		// log.Println("cb0:", closeBr)
		// log.Println("split:", i, tx, "operT:", Lt.Oper, "curT:", etp)
		if etp == Lt.Comm {
			// comment up to end of line
			break
		}
		if etp == Lt.Mtcomm {
			// just ignore
			continue
		}
		prevCloseBr := closeBr
		closeBr = false
		prev = cur
		cur = el

		curPart = append(curPart, i)
		if etp != Lt.Oper {
			cNode.AddRElem(el)
			continue
		}

		// Closing bracket
		// TODO: \ arg -> lambda, `->` as a closing oper of `\\`-opened arg-subnode
		if slashLambda && tx == `->` {
			// TODO: move to bracket-closing section
		}
		if id := strings.Index(cbrs, tx); id > -1 {
			// lastbr := brs[len(brs)-1]
			// expBr, ok := brmap[tx]
			closeBr = true
			brC -= 1
			// log.Println("$11", tx, brC, brs, ok, expBr, "clM", closeBr)
			// finding last brackets in parent
			tInd := 0
			for k := len(parents) - 1; k >= 0; k-- {
				if parents[k].IsBrackets {
					tInd = k
					break
				}
			}
			if tInd == 0 {
				// cur brackets was a root, not sure if it correct case
				// upd1: valid case for brackets-only, but need empty root:
				//  chain call (foo)()()
				//  collectios or gen cases [f1, f2][1]

			} else {
				// new cur is a parent of closed brackets node
				fParent := parents[tInd-1]
				cNode = fParent
			}
			parents = parents[:tInd]

			continue
		}

		// 1. brackets after operator - different node
		// a * (b - c)
		// r = [a, b, c], [a+b, c-2, d+3]
		// r = {1:11, 2:20+2, 3:3*11}
		// a=1; b=2; a + b
		// 2. brackets after word - the same node
		// varlist[key]
		// foo(arg1, arg2)
		// obj.member[key].foo(1,2,3).mem2
		//
		// Open brackets,
		// TODO: \ arg -> lambda, `\\` as a open sub-node with child commas
		if id := strings.Index(obrs, tx); id > -1 {
			// if brC == 0 {
			// 	brpos = append(brpos, ints2{i, -1}) // opened br
			// }
			brC += 1
			curpri := solidExPrior
			tNode := &OperNode{BracketOpen: tx, IsBrackets: true, oper: tx, prior: curpri}
			// TODO: func call shoud have special case: word, (,...,)
			// if prevNode.Type = Lt.Word, Lt.Oper(closeBr), ~>,
			// if call / getElem, etc expr
			/*
				foo(123)
				names[ind]
				Struct{a:1}
				[1,2][0]
				['f1'](123)
				f~>(1)(2)
				(f)(1)
				obj.mem[1][2](3)(4)
			*/
			fParent := cNode
			if prev != nil && (prev.Type == Lt.Word || prev.Type == Lt.Text ||
				(prev.Type == Lt.Oper && slices.Contains(solidEnd, prev.Text)) || // ) ] } ~>
				(prev.Type == Lt.Num && prev.Text[0] == '0')) { // 0x[]
				// log.Println("$_if_brop1", tx)
				// func call(), collect[elem]
				// find solid expression in left:
				// name | (brackets) | expr.expr |
				// tNode := &OperNode{oper: tx, prior: unaryPrior}
				tInd := 0
				for k := len(parents) - 1; k >= 0; k-- {
					// log.Println("$o(pk:", k, "", len(parents), parents[k])
					if curpri < parents[k].prior || parents[k].IsBrackets || (parents[k].IsOper && slices.Contains(unaryR, parents[k].oper)) {
						// log.Println("$_brackets-open", "tPar:", k, parents[k].prior, parents[k].oper)
						tInd = k
						break
					}
				}
				fParent = parents[tInd]
				if fParent.rightNode != nil {
					tNode.leftNode = fParent.rightNode
				} else {
					tNode.leftElems = fParent.rightElems
				}
				parents = parents[:tInd+1] // cut right branch of parent
			} else {
				// cNode.rightNode = tNode
			}
			fParent.rightNode = tNode
			// if arg of oper
			cNode = tNode
			parents = append(parents, tNode)
			continue
		}

		// log.Println("$102", tx, "isBr:", cNode.IsBrackets, fmt.Sprintf("cNode:(%T, %v)", cNode, cNode.oper))
		// opris = priors2
		isFuBr := cNode.IsBrackets && cNode.oper == "(" && notEmptyLeft(cNode)
		txex := tx // changing oper var for special context
		if cNode.IsBrackets {
			opcx = 2
			// debug
			// fmt.Printf(" -- nodeOper: %s, is fuBr: %v \n", cNode.oper, isFuBr)
			if isFuBr {
				// cur parent should be a function def or call
				// log.Println("Change oper priors")
				// opris = priorsFuBr
				opcx = 3
			} else {
				// opris = priors2
				opcx = 4
			}
		}
		if opcx == 3 {
			// log.Println("Change oper to spec context")
			switch tx {
			case "=":
				txex = "func="
			case ":":
				txex = "func:"

			}
		}
		curpri := getPrior(opris, txex)
		// p1 := getPrior(opris, txex)
		// p2 := getPrior(priorsFuBr, txex)
		// fmt.Printf(" -+ isFuBr %v (ocx:%d); oper: <%s> ; prior: %v \n", isFuBr, opcx, tx, curpri)
		// fmt.Printf(" - oper: <%s> ; prior1: %v ; prior2: %v \n", txex, p1, p2)

		// L-unary section
		if prev != nil && (prev.Type == Lt.Oper && !prevCloseBr) {
			if len(tx) == 1 && slices.Contains(unary, tx) {
				// unary oper after another oper
				tNode := &OperNode{oper: tx, prior: unaryPrior}
				cNode.rightNode = tNode
				cNode = tNode
				parents = append(parents, tNode)
				// cNode.AddRElem(el)
				continue
			}
		}

		// R-unary section
		if slices.Contains(unaryR, tx) {

			curpri = getPrior(opris, txex)
			tNode := &OperNode{oper: tx, prior: curpri}
			tInd := 0
			for k := len(parents) - 1; k >= 0; k-- {
				if curpri < parents[k].prior || parents[k].IsBrackets {
					// log.Println("$_r-unary", "tPar:", parents[k].prior)
					tInd = k
					break
				}
			}
			fParent := parents[tInd]
			if fParent.rightNode != nil {
				tNode.leftNode = fParent.rightNode
			} else {
				tNode.leftElems = fParent.rightElems
			}

			fParent.rightNode = tNode
			parents = parents[:tInd+1]
			parents = append(parents, tNode)
			continue
		}

		// Bin Oper section
		// log.Println("$2 pri", tx, curpri, lowest[1], curpri <= lowest[1])
		// log.Println("$2 pri", tx, curpri, cNode.oper, cNode.prior, "less?:", cNode.prior > curpri)
		others = append(others, [2]int{i, curpri})
		// hwere we reach the binary operator
		// prev sequence should be a left operand
		tNode := &OperNode{oper: tx, prior: curpri}
		//
		// 1 * 2 - 3/4 * [5] + 6
		// (((1 * 2) - ((3/4) * 5)) + 6) - (7 * (2 ** 8))
		// *1,2
		// -(*1,2),3
		// -(*1,2),(/3,4)
		// ind(*) < ind(+)
		// cur have lesser precedence (prior index >)

		// operNNs := func(pp []*OperNode) string {
		// 	ss := []string{}
		// 	for _, op := range pp {
		// 		ss = append(ss, op.oper)
		// 	}
		// 	return strings.Join(ss, ",")
		// }
		if cNode.prior > curpri {
			// log.Println("$1_lesser", "rN", cNode.rightNode != nil, "rE", len(cNode.rightElems), "lE", len(cNode.leftElems))
			// next is sub of prev
			if cNode.rightNode != nil {
				// if oper was right of parent
				tNode.leftNode = cNode.rightNode
			} else {
				tNode.leftElems = cNode.rightElems
				cNode.rightElems = nil
			}
			cNode.rightNode = tNode
			parents = append(parents, tNode)
			cNode = tNode
			// log.Println("$_less_res:", operNNs(parents), "cNode:", cNode)
			continue
		}
		// else, cur have more o equal precedence (index <=), ** after *
		// 1. find first node from right with precedence lesser than current, + for cur *
		lInd := -1 // no index
		// for k := len(ndStack) - 1; k >= 0; k-- {
		// 	if curpri > ndStack[k].prior {
		for k := len(parents) - 1; k >= 0; k-- {
			// log.Println("split: find parent> ", curpri, tx, " < ", parents[k].prior, parents[k].oper, curpri < parents[k].prior)
			if curpri < parents[k].prior || parents[k].IsBrackets {
				lInd = k
				break
			}
		}
		var fParent *OperNode = rNode // found parent
		var rBranch *OperNode         // right branch of nodes of prev parent to be a first arg of cur oper
		var rElems []*lang.Elem
		if lInd == -1 {
			// correct parent not found, current will be top parent
			// log.Println("parent Not found")
			rBranch = rNode
			// fParent = tNode
		} else {
			fParent = parents[lInd]
			if fParent.rightNode != nil {
				rBranch = fParent.rightNode
				tNode.leftNode = rBranch // if operNode
			} else if fParent.rightElems != nil {
				rElems = fParent.rightElems
				tNode.leftElems = rElems
			}
		}
		// if fParent == nil {
		// 	rNode = tNode
		// } else {
		// 	fParent.leftNode = tNode
		// }
		// log.Println("$_parent", lInd, fParent)
		fParent.rightNode = tNode
		parents = parents[:lInd+1]
		parents = append(parents, tNode)
		cNode = tNode

		// lNode := ndStack[lInd] // TODO: need take  upper-lvl oper-chain a - |b * c| - d ..
		// tNode.leftNode = lNode
		// prevNode := rNode // upper level node

		// prevNode.leftNode = tNode
		// cNode = tNode
		// ndStack = append(ndStack, tNode)
		// if curpri >= lowest[1] {
		// 	lowest[0] = i
		// 	lowest[1] = curpri
		// }
	}
	finished := true // TODO
	if brC > 0 {
		finished = false
	}
	// log.Println("$_split_res:", rNode, rNode.rightNode, "// finised:", finished)
	// return rNode, nil
	return &LineTree{Tree: rNode, Finished: finished, Parents: parents, BracketsCount: brC, InFuncBr: opcx == 3}, nil
}
