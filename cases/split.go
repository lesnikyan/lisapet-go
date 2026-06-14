package cases

import (
	"errors"
	"log"
	"slices"
	"strings"

	"github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
)

// import	"github.com/lesnikyan/lisapet-go/nodes"

var obrs = "([{"
var cbrs = ")]}"

// map of valid close : open brs
var brmap = map[string]string{
	")": "(",
	"]": "[",
	"}": "{",
}

var _operPriorStr = `( ) [ ] { } 1 . 1 ~> 1 ... 1 -x ! ~ 1 ** ^/ 1 * / % 1 + - 1` +
	`<< >> 1 =~ ?~ /~1 < <= > >= !> ?> !?> 1 == != 1 & 1 ^ 1 | 1 :: 1 && 1 || 1 \\ 1 ->` +
	` 1 @ 1 $ 1 ?: 1 : 1 ? 1 , 1 .. 1 <- 1 @! 1 = += -= *= /= %= 1 ; 1 !: :? => 1 /: `

var unary = strings.Split("- ! ~ +", " ")
var seqSeprs = strings.Split(", ;", " ")

var priors = func() [][]string {
	ss := strings.Split(_operPriorStr, "1")
	var res [][]string
	for _, s := range ss {
		s = strings.TrimSpace(s)
		nn := strings.Split(s, " ")
		res = append(res, nn)
	}
	return res
}()

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

var badBracketsErr = errors.New("bad set of brackets")

type OperNode struct {
	oper  string       // if oper
	elems []*lang.Elem // if other expr
	left  *OperNode
	right *OperNode
}

/*
a * b - (c + d) / (2 / 5)
res:
Lowest: 3
Others: [[1, prec*], [9, prec/]]
Brackets: [[4, 8], [10, 14]]
*/
func OperSplit(elems []*lang.Elem) (*SplittedRes, error) {
	opris := priors[5:] // except solid opers
	brC := 0            // brackets depth
	brs := []string{}
	brpos := []ints2{}
	// brN := -1 // index of last opened bracked
	// var curin int
	var lowest ints2 = [2]int{-1, -1} // lowest precedence of found prior
	var others []ints2 = []ints2{}
	var cur *lang.Elem
	var prev *lang.Elem
	var closeBr = false
	for i, el := range elems {
		// curin = i
		tx := el.Text
		etp := el.Type
		if etp == Lt.Space {
			continue
		}
		// log.Println("cb0:", closeBr)
		prevCloseBr := closeBr
		closeBr = false
		prev = cur
		cur = el
		// log.Println("$1", tx, brC, brs, "preCl:", prevCloseBr)
		if etp != Lt.Oper {
			continue
		}
		// closing bracket
		if id := strings.Index(cbrs, tx); id > -1 {
			lastbr := brs[len(brs)-1]
			expBr, ok := brmap[tx]
			closeBr = true
			log.Println("$11", tx, brC, brs, lastbr, expBr, "clM", closeBr)
			if !ok || lastbr != expBr {
				// incorrect closing bracket
				return nil, badBracketsErr
			}
			brs = brs[:len(brs)-1]
			brC--
			if brC == 0 {
				brpos[len(brpos)-1][1] = i // closed br
			}
			continue
		}

		// open brackets
		if id := strings.Index(obrs, tx); id > -1 {
			if brC == 0 {
				brpos = append(brpos, ints2{i, -1}) // opened br
			}

			brs = append(brs, tx)
			brC++
			continue
		}
		if brC > 0 {
			continue // skip sequence in brackets (thinking about 1-pass logic)
		}

		log.Println("$102", tx, prevCloseBr)
		if prev.Type == Lt.Oper && !prevCloseBr {
			if len(tx) == 1 && slices.Contains(unary, tx) {
				// unary oper after another oper
				continue
			}
		}
		curpri := getPrior(opris, tx)
		log.Println("$2 pri", tx, curpri, lowest[1], curpri <= lowest[1])
		others = append(others, [2]int{i, curpri})
		if curpri >= lowest[1] {
			lowest[0] = i
			lowest[1] = curpri
		}

	}
	// -1 = solid expr.
	return &SplittedRes{Lowest: lowest[0]}, nil
}
