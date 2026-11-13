package nodes

import (
	"fmt"
	"math"
	re "regexp"
	"strconv"
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

func DefaultFormatType(val any) rune {
	switch v := val.(type) {
	case string:
		return 's'
	case int64:
		return 'd'
	case byte:
		return 'x'
	case objects.Glif:
		return 'c'
	case float64:
		return 'f'
	case bool:
		return 'v'
	case *objects.Maybe:
		return 'v'
	// Need some representation of extra types, like !r in py f-string
	// case bool:
	// 	return ??
	// case *objects.Null:
	// null
	// case *objects.Maybe:
	// if v.IsNone(): none
	default:
		panic(fmt.Sprintf("formatting: val %v of type %T can't be displayed by default format", v, v))
	}
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
		switch val.(type) {
		case int64, byte, float64, objects.Glif:
			var fv float64
			switch f := val.(type) {
			case int64:
				fv = float64(f)
			case byte:
				fv = float64(f)
			case float64:
				fv = float64(f)
			case objects.Glif:
				fv = float64(f)
			}
			if opt.DecimSize > -1 {
				fds := math.Pow10(opt.DecimSize)
				fv = math.Round(fv*fds) / fds
			}
			tstr = fmt.Sprintf("%f", fv)
			// fmt.Printf(" #4: tstr: %s \n", tstr)
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
		}
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
		}
	case 'd':
		switch d := val.(type) {
		case int64, objects.Glif, byte:
			tstr = fmt.Sprintf("%d", d)
		default:
			panic(fmt.Sprintf("s-format cant represend %v as `d` type.", val))
		}
		// switch d := val.(type) {
		// case int64, objects.Glif, byte:
		// 	if IsLess(d, 0) {
		// 		negative = true
		// 	}
		// }
	case 'x':
		switch d := val.(type) {
		case int64, objects.Glif, byte:
			tstr = fmt.Sprintf("%x", d)
		default:
			panic(fmt.Sprintf("s-format cant represend %v as `x` type.", val))
		}
		// switch xx := val.(type) {
		// case int64:
		// 	tstr = fmt.Sprintf("%x", xx)
		// case byte:
		// 	tstr = fmt.Sprintf("%x", xx)
		// case objects.Glif:
		// 	tstr = fmt.Sprintf("%x", xx)
		// }
	case 'b':
		switch d := val.(type) {
		case int64, objects.Glif, byte:
			tstr = fmt.Sprintf("%b", d)
		default:
			panic(fmt.Sprintf("s-format cant represend %v as `b` type.", val))
		}
		// switch xx := val.(type) {
		// case int64:
		// 	tstr = fmt.Sprintf("%x", xx)
		// case byte:
		// 	tstr = fmt.Sprintf("%x", xx)
		// case objects.Glif:
		// 	tstr = fmt.Sprintf("%x", xx)
		// }
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
		// if opt.Sign != 0 {
		// if negative {
		// sign = "-"
		// } else {
		// if opt.Sign == '+' {
		// 	sign = "+"
		// }
		// }
		// }
		if opt.Align == AlignDefault {
			opt.Align = AlignDefNum
		}
	}

	// if !negative && sign == "+" {
	// 	tstr = sign + tstr
	// }

	// indent and align
	tlen := utf8.RuneCountInString(tstr)
	if opt.Sign == 0 && (opt.Size == 0 || tlen >= opt.Size) {
		// no size and align
		return tstr
	}

	// >>> f"{-12:>04d}"
	// '0-12'
	// >>> f"{-12:04d}"
	// '-012'
	// tstr = tstr

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
		// switch opt.Align {
		// case AlignDefNum:
		// 	// sign first
		// case AlignRight:
		// 	// sign just before num

		// }

	}

	fill := " "
	sig0 := ""
	sigN := ""

	// zero fill
	if opt.Zero {
		switch xtype {
		case 'd', 'f', 'g', 'e', 'x', 'b', 'o':
			fill = "0"
		}
	} else if opt.Fill != 0 {
		fill = string(opt.Fill)
	}
	nLeft := 0
	nRight := 0

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
		// case *ValExpr:
		// 	pval = next.Val
		// spart = val2Str(pval, form)
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

type Align int

const (
	AlignDefault Align = 0
	AlignLeft    Align = 1
	AlignRight   Align = 2
	AlignCenter  Align = 3
	AlignDefNum  Align = 4
	AlignDefStr  Align = 5
)

var _alignMap = map[rune]Align{
	'<': AlignLeft,
	'>': AlignRight,
	'^': AlignCenter,
}

type SFormOptions struct {
	Indent int
	Zero   bool
	Size   int
	// DecimZero bool
	DecimSize int
	Sign      rune // '-', '+'
	Align     Align
	Type      rune
	Fill      rune
}

var digRx = re.MustCompile("^[0-9]$")
var digitsRx = re.MustCompile("^[0-9]+$")
var arrRx = re.MustCompile(`[<>\^]`)

func FormOption(src string) (*SFormOptions, error) {
	// error if prefix is incorrect
	runes := []rune(src)
	opt := &SFormOptions{Align: AlignDefault, DecimSize: -1}
	rlen := len(runes)
	if rlen == 0 {
		// no format, TODO default
		return nil, nil
	}
	var tp rune = '?'
	if rlen > 0 {
		// at least data type char
		tp = runes[rlen-1]
		opt.Type = tp
	}
	ind := arrRx.FindStringIndex(src)
	i1 := 0
	i2 := rlen - 1
	if ind != nil {
		// have align
		arr1, arr2 := ind[0], ind[1]
		// fmt.Printf(" - .. fs#1 %d, %d >> %s \n", arr1, arr2, string(runes[arr1]))
		align, ok := _alignMap[runes[arr1]]
		if !ok {
			return nil, fmt.Errorf("format string: bad align key: %c ", runes[arr1])
		}
		opt.Align = align
		if arr1 > 0 {
			// has fill char before align
			opt.Fill = runes[arr1-1]
		}
		i1 = arr2
	}
	// pref := runes[0 : rlen-1]
	// val : =^+05,.4f
	if rlen == 1 {
		// just 1 char format {val : d}
		return opt, nil
	}
	if tp == 'f' {
		// look for `.`
		dind := strings.Index(src, ".")
		if dind > -1 {
			// get decimap precision
			// i1:dind >> nums before dot, dind:i2 >> after dot
			if i2-dind < 1 {
				// no precision digits
				return nil, fmt.Errorf("format string: no precision digits after dot")
			}
			if i2 > dind {
				decim := string(runes[dind+1 : i2])
				// fmt.Printf(" -- decim part: `%s` from `%s`[%d:%d] \n", decim, string(runes), dind+1, i2)
				if digitsRx.MatchString(decim) {
					// has decimal part
					decval, err := strconv.ParseInt(decim, 10, 64)
					if err != nil {
						return nil, err
					}
					opt.DecimSize = int(decval)
					i2 = dind
				}
			}
		}
	}
	if strings.ContainsRune("dfexob", tp) {
		// numeric value
		if strings.ContainsRune("+-", runes[i1]) {
			// signed
			opt.Sign = runes[i1] // diff between +, - ?
			i1++
		}
		if runes[i1] == '0' {
			// zero lead nums
			opt.Zero = true
		}
		if runes[i2-1] == ',' {
			// comma-block for numbers
			// TODO
			i2--
		}
	}
	if i1 < i2 {
		// no size
		zerosize := string(runes[i1:i2])
		// fmt.Printf(" -- zerosize part: `%s` i1: %d, i2: %d  \n", zerosize, i1, i2)
		if !digitsRx.MatchString(zerosize) {
			return nil, fmt.Errorf("format string: bad size format: `%s`", string(zerosize))
		}
		size, err := strconv.ParseInt(zerosize, 10, 64)
		if err != nil {
			return nil, err
		}
		opt.Size = int(size)
	}
	if strings.ContainsRune("sc", tp) {
		// ?
	}

	return opt, nil
}

// func FormatVal(val any, pref string) {
// string: left indent

// num: sign, indent, zero lead,

// float: sign, indent, decimal part

// hex: indent,  zero lead,

// binary: zero lead

// octa: zero lead

// float E
// }
