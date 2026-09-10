package parser

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	lang "github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
)

var spaceSet = lang.Kmap([]rune{' ', '\t', '\n', '\r'})
var c_esc = "'\"ntr\\/`"
var escMap = map[rune]rune{'n': '\n', 't': '\t', 'r': '\r', '\\': '\\', '/': '/', '\'': '\'', '"': '"', '`': '`'}
var c_nums = `1234567890`
var c_oper = `+~-*/=%^&!?<>()[]:.;,|${}@\\`

// # single-line comment, to and of line
var c_comm = '#'

// # blok-comment, multiline or inline
var c_opcomm = `#@`
var c_ndcomm = `@#`
var c_mlines = strings.Split("''' \"\"\" ```", ` `)
var rxChar = regexp.MustCompile(`[a-zA-Z_\$]`)

var charSet = lang.Kmap([]rune("qwertyuiopasdfghjklzxcvbnm" +
	"QWERTYUIOPASDFGHJKLZXCVBNM_$"))

var quotSet = lang.Kmap([]rune("'\"`"))

// var c_regex = `| / % `

// var ext_in = []rune{'j', 'x', 'b', 'o', 'a', 'b', 'c', 'd', 'e', 'f'}

var extNum = "boxabcdef1234567890"

var opers = strings.Split(
	`; , ... $$ .. ** ++ -- += -= *= /= %=  && || == != <= >= << >>`+
		` => ?> !?> -> <- !- := ?: /: !: :? :> :: =~ ?~ /~ ~> ^/ @!`+
		` < > = + - * / | \\ { } [ ] . , : ? ~ ! % ^ & * $ ( ) @`,
	` `)

// # if i > 0 and el.text in ['-', '+', '!', '~'] and elems[i-1].type == Lt.oper and elems[i-1].text != ')'
var unarOpers = `- + ! ~`

//

func validOper(cur []rune, c rune) bool {
	// simplest finding, need optimization
	// define if new sequence is out of correct operators
	oper := make([]rune, len(cur)+1)
	copy(oper, cur)
	oper[len(cur)] = c
	return slices.Contains(opers, string(oper))
}

func has[K comparable, V any](m map[K]V, v K) bool {
	_, ok := m[v]
	return ok
}

var textTypes = []Lt.Lt{Lt.Text, Lt.Mttext}
var absorbTypes = []Lt.Lt{Lt.Text, Lt.Mttext, Lt.Comm, Lt.Mtcomm}

func elemType(c rune, prev Lt.Lt) Lt.Lt {
	if slices.Contains(textTypes, prev) {
		if c == '\\' {
			return Lt.Esc
		}
	}
	if has(quotSet, c) {
		return Lt.Quot
	}
	if slices.Contains(absorbTypes, prev) {
		return prev
	}

	switch prev {
	// case Lt.Text:
	// 	return Lt.Text
	// case Lt.Mttext:
	// 	return Lt.Text
	case Lt.Comm:
		return Lt.Comm
	case Lt.Mtcomm:
		return Lt.Comm
	case Lt.Esc:
		return Lt.Text
	case Lt.Num:
		// TODO: try use map of chars (rune:bool) instead of slice
		if strings.ContainsRune(extNum, c) {
			return Lt.Num
		}
		// if c == '.' {
		// 	return prev
		// }
	case Lt.Word:
		if strings.ContainsRune(c_nums, c) || has(charSet, c) {
			return Lt.Word
		}
	}
	if has(charSet, c) {
		return Lt.Word
	}
	if has(spaceSet, c) {
		return Lt.Space
	}
	if c == c_comm {
		return Lt.Comm
	}
	if strings.ContainsRune(c_oper, c) {
		return Lt.Oper
	}
	if strings.ContainsRune(c_nums, c) {
		return Lt.Num
	}

	return Lt.None
}

// func f() {
// 	// lt.LtClose
// 	fmt.Print(Lt.Block)
// }

// is string, comment, etc
type SplitContext struct {
	Ltype  Lt.Lt
	strval string
	elems  []*lang.Elem
	prev   []rune
	indent int
}

var spaceInds = []rune(" \t")

func Lines(src string) []string {
	src = strings.ReplaceAll(src, "\r\n", "\n")
	var lines []string = strings.Split(src, "\n")
	if len(lines) == 0 || len(lines[0]) == 0 {
		return lines
	}
	if slices.Contains(spaceInds, rune(lines[0][0])) {
		res := make([]string, len(lines))
		ind := 0
		sp := rune(lines[0][0])
		for _, c := range lines[0] {
			if c != sp {
				break
			}
			ind++
		}
		for i, ln := range lines {
			res[i] = ln[ind:]
		}
		lines = res
	}
	return lines
}

func Runes(src string) []rune {
	return []rune(src)
}

/*
TODO: check unicode cases
https://pkg.go.dev/unicode/utf8#DecodeRuneInString
*/
func nPart(cc []rune) string {
	return string(cc)
}

// if time to finalize cur set and start new
// like whitespace after number or word, or end of "string"
func finCond(cur []rune, next rune, curType Lt.Lt, prevType Lt.Lt) bool {
	if len(cur) == 0 {
		return false
	}

	switch prevType {
	case Lt.Comm:
		return false
	case Lt.Word:
		return curType != Lt.Word
	case Lt.Space:
		return curType != prevType
	case Lt.Text:
		if curType != Lt.Quot {
			return false
		}
		if cur[0] == next {
			return true // end of string
		}
	case Lt.Mttext:
		if curType != Lt.Quot {
			return false
		}
		// TODO: need 1-st line of mttext to detect valid close sequence
		// clen := len(cur)
		// if clen < 5 {
		// 	return false
		// }
		// if next == cur[clen-1] && next == cur[clen-2] {
		// 	return true
		// }
		return closeMult(cur, closeMstr, next)

	case Lt.Num:

		return curType != Lt.Num
	case Lt.Oper:
		if curType != Lt.Oper {
			return true
		}
		// check if `next` breaks valid oper
		return !validOper(cur, next)
	}
	// log.Printf("f1--'%s'", string(next))
	switch curType {
	case Lt.Space:
		return prevType != Lt.Space // space breaks any other
	}
	return false
}

func nextType(ntype Lt.Lt, c rune) Lt.Lt {
	if ntype == Lt.Quot {
		return Lt.None
	}
	return ntype
}

type Er rune // escaped rune

// var closeMstrEnd = []int{0, -3, -2}
var closeMstr = []int{0, -2, -1}

func closeMult(prev []rune, cc []int, last rune) bool {
	clen := len(prev)
	if clen < 5 {
		return false
	}
	// last := prev[clen-1]
	for _, c := range cc {
		ci := c
		if c < 0 {
			ci = clen + c
		}
		// fmt.Printf("CM?1: %d, <%s> [%d>>%d] \n", clen, string(last), c, ci)
		// fmt.Printf("CM?2: <%s> [%d] == %s \n", string(last), ci, string(prev[ci]))
		if prev[ci] != last {
			return false
		}
	}
	return true
	// return cmp.Compare(cur[0:3], cur[clen-3:clen-1]) == 0
}

var endLine = '\n'

func SplitLine(runes []rune, ctx *SplitContext) []*lang.Elem {
	res := []*lang.Elem{}
	cur := []rune{}
	ctype := Lt.None
	if len(ctx.prev) > 0 {
		cur = ctx.prev
		ctype = ctx.Ltype // cur elem type
		res = ctx.elems
	}
	ctx.Ltype = Lt.None
	ctx.prev = nil
	ctx.elems = nil
	esc := false
	lastEsc := -1

	// log.Printf("SpLine#1:  , c<%s: %d> esc:%v  \"%s\" ", Lt.TName(ctype), ctype, esc, string(cur))
	// log.Printf("SpLine#2:  \\\\%s//", string(runes))
	// muq := '*'         // multiline string quote char
	for _, c := range runes {
		// cur = append(cur, c)

		switch ctype {
		case Lt.Comm:
			cur = append(cur, c)
			continue
			// case Lt.Mtcomm:
			// 	// if not close
			// 	cur = append(cur, c)
			// 	continue
			// case Lt.Mttext:
			// 	// if not close
			// 	cur = append(cur, c)
			// 	continue
		}

		curStr := ctype == Lt.Text || ctype == Lt.Mttext // prev is text
		xtype := elemType(c, ctype)                      // next type
		if curStr {
			if !esc && xtype == Lt.Esc {
				// escMap
				// log.Printf("esc#1:  , x<%s>  _%s_ opnch: <%s>", Lt.TName(xtype), string(c), string(cur[0]))
				if cur[0] != '`' {
					// backtics string doesn't support escapes
					esc = true
					lastEsc = len(cur)
					continue
				} else {
					// log.Printf("#no esc by  _%s_ ", string(c))
					xtype = ctype
				}
			}
			if esc {
				// log.Printf("esc#2:  , c<%s>  _%s_ ", Lt.TName(xtype), string(c))
				esc = false
				if rep, ok := escMap[c]; ok {
					// c = rep
					// cur = cur[0 : len(cur)-2]
					cur = append(cur, rep)
					continue
				} else {
					panic(fmt.Sprintf("Incorrect escape sequence in string: `%s`", string(c)))
				}
			}
			if c != cur[0] {
				xtype = ctype
			}

		}
		switch xtype {
		case Lt.Quot:
			if !curStr {
				xtype = Lt.Text
				lastEsc = -1
			}
			if len(res) > 0 && len(cur) == 0 {
				// if prev == [quot, quot]
				rlast := res[len(res)-1]
				if rlast.Type == Lt.Text && len(rlast.Text) == 2 {
					preRr := []rune(rlast.Text)
					if c == preRr[0] {
						// s = ''' '''
						// start of multiline string
						// muq = c
						res = res[0 : len(res)-1]
						cur = preRr
						xtype = Lt.Mttext
					}
				}
			}
		}

		fin := finCond(cur, c, xtype, ctype)
		if fin {
			switch ctype {
			case Lt.Mttext:
				if lastEsc > len(cur)-3 {
					fin = false
				}
			case Lt.Text:
				if lastEsc > len(cur)-1 {
					fin = false
				}
			}
		}
		// log.Printf("SL1:  , c<%s> : cur=\"%s\", _%s_  ?%v", Lt.TName(ctype), string(cur), string(c), fin)
		// log.Printf("SL2: %s  c<%s>, x<%s> : s=\"%s\"  ?%v", string(c), Lt.TName(ctype), Lt.TName(xtype), string(cur), fin)
		if fin {
			// log.Printf("SL2: %s  c<%s>, x<%s> : s=\"%s\"  ?%v", string(c), Lt.TName(ctype), Lt.TName(xtype), string(cur), fin)
			// fmt.Printf("-- fin |%s|, |%s| \n", string(cur), string(c))
			ntype := nextType(xtype, c)
			switch xtype {
			case Lt.Quot:
				// last qoute in string
				cur = append(cur, c)
			}
			text := nPart(cur)
			elem := &lang.Elem{Text: text, Type: ctype}
			res = append(res, elem)
			cur = []rune{}
			// if ctype != Lt.Text && ctype != Lt.Mttext{
			if !curStr {
				cur = append(cur, c)
			}
			ctype = ntype
			continue
		} else {
			switch xtype {
			case Lt.Quot:
				switch ctype {
				case Lt.Mttext:
					xtype = Lt.Mttext
				}
			}
		}
		if ctype == Lt.Mttext {

		}
		ctype = xtype
		cur = append(cur, c)
	}
	if len(cur) > 0 {
		if ctype == Lt.Mttext {
			ctx.prev = cur
			ctx.Ltype = ctype
		} else {
			res = append(res, &lang.Elem{Text: nPart(cur), Type: ctype})
		}
	}
	return res
}

var BaseIndent = 0

var rrr = regexp.MustCompile(`\s`)
var _spaces = " \t"

func cutIndent(rline []rune, fixCut int) ([]rune, int) {
	// cutSize := fixCut + 1
	for i, r := range rline {
		if i < fixCut && !strings.ContainsRune(_spaces, r) {
			return rline[i:], i
		}
	}
	return rline, 0
}

func SplitCode(code string) []*lang.CLine {
	ctx := &SplitContext{}
	lines := Lines(code)
	indSize := 0
	res := make([]*lang.CLine, len(lines))
	if BaseIndent > 0 {

	}
	// prepare lines
	for _, ln := range lines {
		if len(ln) == 0 {
			continue
		}
		lnx := Runes(ln)
		if strings.ContainsRune(_spaces, lnx[0]) {
			for i, r := range lnx {
				if !strings.ContainsRune(_spaces, r) {
					BaseIndent = i - 1
				}
			}
			// cutIndent(line)
		}
		break
	}
	// parse lines
	for i, line := range lines {
		rline := Runes(line)
		if BaseIndent > 0 {
			rline = rline[BaseIndent:]
		}
		curInd := 0
		// if strings.ContainsRune(_spaces, rline[0]) {}
		cutSize := len(rline)
		if ctx.Ltype == Lt.Mttext {
			// continuation of multiline text
			curInd = ctx.indent
			cutSize = ctx.indent
		}
		rline, size := cutIndent(rline, cutSize)
		if size > 0 {
			// indent detected
			if indSize == 0 {
				// set size of indent
				indSize = size
			}
			curInd = size / indSize
		}
		lems := SplitLine(rline, ctx)
		// lastEl := lems[len(lems)-1]
		// fmt.Printf("pars.Line last: type: %s, tx: %s\n", Lt.TName(ctx.Ltype), string(ctx.prev))
		switch ctx.Ltype {
		case Lt.Mttext:
			// ctx.Ltype = lastEl.Type
			ctx.elems = append(ctx.elems, lems...)
			ctx.prev = append(ctx.prev, endLine)
			ctx.indent = curInd
			continue
		default:
			ctx = &SplitContext{}
		}
		if len(lems) == 0 {
			continue
		}

		res[i] = &lang.CLine{Elems: lems, Src: line, Indent: curInd}
	}
	return res
}

func Crop(code string) string {
	code = strings.ReplaceAll(code, "\r\n", "\n")
	if code[0] == '\n' {
		code, _ = strings.CutPrefix(code, "\n")
	}
	return code
}

func Foo() {
	fmt.Println("parser 1")
}
