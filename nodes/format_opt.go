package nodes

import (
	"fmt"
	re "regexp"
	"strconv"
	"strings"

	"github.com/lesnikyan/lisapet-go/objects"
)

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
var _validType = "dobxfesc%"

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
		if !strings.ContainsRune(_validType, tp) {
			return nil, fmt.Errorf("format string: bad type of data symbol: %c in `%s`", tp, src)
		}
		opt.Type = tp
	}
	// println(" -- tp:", string(tp))
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
	if strings.ContainsRune("feg%", tp) {
		// look for `.`
		dind := strings.Index(src, ".")
		// println("fe dind:", dind)
		if dind > -1 {
			// get decimap precision
			// i1:dind >> nums before dot, dind:i2 >> after dot
			// fmt.Printf(" -- #1 from `%s`[%d:%d] \n", string(runes), dind+1, i2)
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
	if strings.ContainsRune("dfegxob%", tp) {
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
		if i2 > 0 && runes[i2-1] == ',' {
			// comma-block for numbers
			// TODO: postponed
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
