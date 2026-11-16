package nodes

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

func IsLess(val any, arg float64) bool {
	switch v := any(val).(type) {
	case int:
		return float64(v) < arg
	case int64:
		return float64(v) < arg
	case objects.Glif:
		return float64(v) < arg
	case byte:
		return float64(v) < arg
	case float64:
		return (v) < float64(arg)
	}
	panic(fmt.Sprintf("IsLess: Bad type of value: %T", val))
}

func floatVal(fv float64, opt *SFormOptions) string {
	decsize := 6
	if opt.DecimSize > -1 {
		fds := math.Pow10(opt.DecimSize)
		fv = math.Round(fv*fds) / fds
		decsize = opt.DecimSize
	}
	ptr := fmt.Sprintf("%%.%df", decsize)
	tstr := fmt.Sprintf(ptr, fv)
	// fmt.Printf(" #4: ptr: %s tstr: %s \n", ptr, tstr)
	if opt != nil {
		if opt.DecimSize > -1 {
			frun := []rune(tstr)
			dind := strings.IndexRune(tstr, '.')
			tv := frun[:dind+1]
			decv := frun[dind+1:]
			if len(decv) > opt.DecimSize {
				decv = decv[0:opt.DecimSize]
				tv = append(tv, decv...)
			}
			tstr = fmt.Sprintf("%s.%s", string(frun[0:dind]), string(decv))
		}
	}
	return tstr
}

func val2Str(val any, opt *SFormOptions) string {
	var tstr string
	negative := false
	xtype := 's'
	if opt != nil {
		xtype = opt.Type
	} else {
		xtype = DefaultFormatType(val)
		opt = &SFormOptions{Type: xtype, DecimSize: -1}
	}

	// val prepare
	var fv float64 // float value
	switch opt.Type {
	case 'f', 'e', '%':
		switch f := val.(type) {
		case int64:
			fv = float64(f)
		case byte:
			fv = float64(f)
		case float64:
			fv = f
		case objects.Glif:
			fv = float64(f)
		default:
			panic(fmt.Sprintf("s-format cant represend %v as `e` type.", val))
		}
	}

	switch xtype {
	case 's':
		if opt.Align == AlignDefault {
			opt.Align = AlignDefStr
		}
		switch s := val.(type) {
		case string:
			tstr = s
		default:
			panic("template `%s` wait a string value")
		}
	case 'f':
		tstr = floatVal(fv, opt)
	case 'e':
		opt.Zero = false
		dsz := ""
		dpoint := ""
		// esig := ""
		if opt.DecimSize > -1 {
			// println("== opt.DSize", opt.DecimSize)
			dpoint = "."
			dsz = fmt.Sprintf("%d", opt.DecimSize)
		}
		ePatr := fmt.Sprintf("%%%s%se", dpoint, dsz)
		tstr = fmt.Sprintf(ePatr, fv)
		// println(" #2.0 tstr:", tstr)
	case '%':
		tstr = floatVal(fv*100.0, opt)
		tstr += "%"
	case 'v':
		switch vv := val.(type) {
		case *objects.Maybe:
			if vv.IsNone() {
				return "none"
			}
			return fmt.Sprintf("some(%v)", vv.Val)
		}
		tstr = fmt.Sprintf("%v", val)
		// no additional formatting for extra types
		return tstr
	case 'c':
		// glifs
		switch c := val.(type) {
		case objects.Glif:
			tstr = string(c)
		case int64:
			tstr = string(c)
		case byte:
			tstr = string(c)
		}
	case 'd':
		switch d := val.(type) {
		case int64, objects.Glif, byte:
			tstr = fmt.Sprintf("%d", d)
		default:
			panic(fmt.Sprintf("s-format cant represend %v as `d` type.", val))
		}

	case 'x':
		switch d := val.(type) {
		case int64, objects.Glif, byte:
			tstr = fmt.Sprintf("%x", d)
		default:
			panic(fmt.Sprintf("s-format cant represend %v as `x` type.", val))
		}

	case 'b':
		switch d := val.(type) {
		case int64, objects.Glif, byte:
			tstr = fmt.Sprintf("%b", d)
		default:
			panic(fmt.Sprintf("s-format cant represend %v as `b` type.", val))
		}
	case 'o':
		switch d := val.(type) {
		case int64, objects.Glif, byte:
			tstr = fmt.Sprintf("%o", d)
		default:
			panic(fmt.Sprintf("s-format cant represend %v as `o` type.", val))
		}
	}

	// TODO: maybe not actual
	if opt == nil {
		return tstr
	}

	// sign
	// sign := ""
	switch xtype {
	case 'd', 'f', 'g', 'e', 'x', 'b', 'o':
		if IsLess(val, 0) {
			negative = true
		}
		if opt.Align == AlignDefault {
			opt.Align = AlignDefNum
		}
	}
	// indent and align
	tlen := utf8.RuneCountInString(tstr)
	if opt.Sign == 0 && (opt.Size == 0 || tlen >= opt.Size) {
		// no size and align
		return tstr
	}
	tsig := '~'
	sglen := 0
	// fmt.Printf(" #5-align: %d , negat: %v, Sign: %c, val: %v type: %c \n", opt.Align, negative, opt.Sign, val, xtype)
	// fmt.Printf(" #6: Less:: %v  \n", IsLess(val, 0))
	if negative || opt.Sign == '+' {
		tsig = '+'
		if negative {
			tsig = '-'
			tstr = string([]rune(tstr)[1:])
			tlen--
		}
		sglen = 1
	}

	fill := " "
	sig0 := ""
	sigN := ""

	// zero fill
	if opt.Zero {
		switch xtype {
		case 'd', 'f', 'x', 'b', 'o':
			fill = "0"
		}
	} else if opt.Fill != 0 {
		fill = string(opt.Fill)
	}
	nLeft := 0
	nRight := 0

	// println(" #7.0 tstr:", tstr)
	// fmt.Printf(" #7: tlen: %v  tsig: %c \n", tlen, tsig)
	switch opt.Align {
	case AlignLeft, AlignDefStr:
		nRight = opt.Size - tlen - sglen
		if tsig != '~' {
			sig0 = string(tsig)
		}
	case AlignRight:
		nLeft = opt.Size - tlen - sglen
		if tsig != '~' {
			sigN = string(tsig)
		}
	case AlignDefNum:
		nLeft = opt.Size - tlen - sglen
		if tsig != '~' {
			if opt.Zero {
				sig0 = string(tsig)
			} else {
				sigN = string(tsig)
			}

		}
	case AlignCenter:
		if tsig != '~' {
			sigN = string(tsig)
		}
		diff := opt.Size - tlen - sglen
		nLeft = diff / 2
		nRight = diff - nLeft
	default:
		nRight = opt.Size - tlen
	}
	pref := ""
	suff := ""
	if opt.Size > 0 {
		if nLeft > 0 {
			pref = strings.Repeat(fill, nLeft)
		}
		if nRight > 0 {
			suff = strings.Repeat(fill, nRight)
		}
	}
	return sig0 + pref + sigN + tstr + suff
}

// ==

func FPartsConcate(parts []*FormatElem, cx base.Context) (string, error) {
	if len(parts) == 0 {
		return "", nil
	}
	sparts := make([]string, len(parts))
	forms := make([]*SFormOptions, len(parts))
	var pval any
	for i, part := range parts {
		expr := part.Expr
		form := part.Format
		forms[i] = form
		var spart string
		// fmt.Printf(" FS-opt: %T, %v | %T \n", form, form, expr)
		switch next := expr.(type) {
		case *TextNode:
			spart = next.Text
		default:
			err := next.Do(cx)
			if err != nil {
				return "", err
			}
			pval = GetExprVal(next, nil)
			// fmt.Printf(" FS-pval: %T, %v \n", pval, pval)
			if form != nil {
				spart = val2Str(pval, form)
			} else {
				spart = val2Str(pval, form)
			}
		}
		sparts[i] = spart
	}
	res := strings.Join(sparts, "")
	return res, nil
}
