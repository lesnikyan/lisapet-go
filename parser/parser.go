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
var c_esc_map = map[rune]rune{'n': '\n', 't': '\t', 'r': '\r', '\\': '\\', '/': '/', '\'': '\'', '"': '"', '`': '`'}
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

var c_regex = `| / % `

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
	// case Lt.Comm:
	// 	return Lt.Comm
	// case Lt.Mtcomm:
	// 	return Lt.Comm
	case Lt.Esc:
		return Lt.Text
	case Lt.Num:
		// TODO: try use map of chars (rune:bool) instead of slice
		if strings.ContainsRune(extNum, c) {
			return Lt.Num
		}
		if c == '.' {
			return prev
		}
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
	case Lt.Word:
		return curType != Lt.Word
	case Lt.Space:
		return curType != prevType
	case Lt.Text:
		if cur[0] == next {
			return true // end of string
		}
	case Lt.Mttext:
		if curType != Lt.Quot {
			return false
		}
		// TODO: need 1-st line of mttext to detect valid close sequence
	case Lt.Num:
		if next == '.' {
			return slices.Contains(cur, '.') // possibly decimal point
		}
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

func SplitLine(runes []rune, ctx SplitContext) []*lang.Elem {

	res := []*lang.Elem{}
	cur := []rune{}
	ctype := ctx.Ltype // cur elem type
	for _, c := range runes {
		// cur = append(cur, c)

		switch ctype {
		case Lt.Mtcomm:
			// if not close
			cur = append(cur, c)
			continue
		case Lt.Mttext:
			// if not close
			cur = append(cur, c)
			continue
		}

		xtype := elemType(c, ctype) // next type

		switch xtype {
		case Lt.Quot:
			if ctype != Lt.Text {
				xtype = Lt.Text
			}
		}

		fin := finCond(cur, c, xtype, ctype)
		// log.Printf("SL1: %s  <%s>, %s : '%s'  ?%v", string(c), Lt.TName(ctype), Lt.TName(xtype), string(cur), fin)
		if fin {
			ntype := nextType(xtype, c)
			if xtype == Lt.Quot { // close string
				cur = append(cur, c)
			}
			text := nPart(cur)
			elem := &lang.Elem{Text: text, Type: ctype}
			res = append(res, elem)
			cur = []rune{}
			if ctype != Lt.Text {
				cur = append(cur, c)
			}
			ctype = ntype
			continue
		}
		ctype = xtype
		cur = append(cur, c)
	}
	if len(cur) > 0 {
		res = append(res, &lang.Elem{Text: nPart(cur), Type: ctype})
	}
	return res
}

var BaseIndent = 0

var rrr = regexp.MustCompile(`\s`)
var _spaces = " \t"

func cutIndent(rline []rune) ([]rune, int) {
	for i, r := range rline {
		if !strings.ContainsRune(_spaces, r) {
			return rline[i:], i
		}
	}
	return rline, 0
}

func SplitCode(code string) []*lang.CLine {
	ctx := SplitContext{}
	lines := Lines(code)
	indSize := 0
	res := make([]*lang.CLine, len(lines))
	if BaseIndent > 0 {

	}
	ln0 := Runes(lines[0])
	if strings.ContainsRune(_spaces, ln0[0]) {
		for i, r := range ln0 {
			if !strings.ContainsRune(_spaces, r) {
				BaseIndent = i - 1
			}
		}
		// cutIndent(line)
	}
	for i, line := range lines {
		rline := Runes(line)
		if BaseIndent > 0 {
			rline = rline[BaseIndent:]
		}
		curInd := 0
		// if strings.ContainsRune(_spaces, rline[0]) {}
		rline, size := cutIndent(rline)
		if size > 0 {
			// indent detected
			if indSize == 0 {
				// set size of indent
				indSize = size
			}
			curInd = size / indSize
		}
		lems := SplitLine(rline, ctx)
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
