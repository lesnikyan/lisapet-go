package cases

// Single element

import (
	"slices"
	"strconv"

	re "regexp"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/lesnikyan/lisapet-go/nodes"
	obb "github.com/lesnikyan/lisapet-go/objects"
)

var _valLexms = []Lt.Lt{Lt.Num, Lt.Text, Lt.Word}

// var _constLex = []string{`null`, `true`, `false`}
var _contsVals = map[string]any{
	`null`:  &obb.Null{},
	`true`:  true,
	`false`: false,
}
var rxInt = re.MustCompile(`^^[0-9]+$`)
var rxInt16 = re.MustCompile(`^0x[0-9a-fA-F]+$`)
var rxInt8 = re.MustCompile(`^0o[0-7]+$`)
var rxInt2 = re.MustCompile((`^0b[01]+$`))
var rxFloat = re.MustCompile(`^[0-9]+\.[0-9]*$`)
var rxVar = re.MustCompile(`^[_a-zA-Z][_a-zA-Z0-9]*$`)

func valex(v any) *nodes.ValExpr {
	return &nodes.ValExpr{Val: v}
}

func MakeNumField(ee []*lang.Elem) *nodes.NumField {
	nums := make([]string, len(ee))
	for i, nx := range ee {
		nums[i] = nx.Text
	}
	return &nodes.NumField{V: nums}
}

func CaseVal(ee []*lang.Elem) (base.Expression, bool) {

	elen := len(ee)
	// if len(ee) != 1 {
	// 	return nil, false
	// }
	etype := ee[0].Type
	etext := ee[0].Text
	if !slices.Contains(_valLexms, ee[0].Type) {
		return nil, false
	}
	if elen > 1 {
		switch etype {
		case Lt.Num:
			// t0 := ee[0].Text[0]
			// println(" >>>>>>>>>> #!!!", ee[0].Text, ">>", t0)
			// nums := make([]string, len(ee))
			// for i, nx := range ee {
			// 	nums[i] = nx.Text
			// }
			// res := &nodes.NumField{V: nums}
			res := MakeNumField(ee)
			return res, true
		case Lt.Word:
			// try get bytes sequence
			t0 := ee[0].Text[0]
			// println(" >>>>>>>>>> #!!!", t0)
			if (t0 > 64 && t0 < 71) || (t0 > 96 && t0 < 103) {
				// println("make NumField")
				res := MakeNumField(ee)
				return res, true
			}
			// string-prefixes
			switch elen {
			case 2:
				// glif, etc
				if ee[1].Type == Lt.Text {
					// fmt.Printf("CaseVal prefix[2] %s, %s \n", ee[0].Text, ee[1].Text)
					switch ee[0].Text {
					case "g":
						// glif
						// fmt.Printf("CaseVal Glif %s, %s \n", ee[0].Text, ee[1].Text)
						rrs := []rune(ee[1].Text)
						if len(rrs) != 3 {
							// fmt.Printf("Error bad Glif %s \n", ee[1].Text)
							return nil, false
						}
						return valex(rrs[1]), true
					case "re":
						//regexp
					}
				}
			case 3:
				// regexp
				if ee[1].Type == Lt.Text {
					// fmt.Printf("CaseVal prefix[3] %s, %s, %s \n", ee[0].Text, ee[1].Text, ee[2].Text)
				}
			}
		}
	}

	var res base.Expression = nil
	// var ok = false
	// var val any
	switch etype {
	case Lt.Num:
		switch elen {
		case 1:
			var err error = nil
			var val any
			if rxFloat.MatchString(etext) {
				val, err = strconv.ParseFloat(etext, 64)
			} else if rxInt.MatchString(etext) {
				val, err = strconv.ParseInt(etext, 10, 64)
			} else if rxInt16.MatchString(etext) {
				val, err = strconv.ParseInt(etext[2:], 16, 64)
			} else if rxInt8.MatchString(etext) {
				val, err = strconv.ParseInt(etext[2:], 8, 64)
			} else if rxInt2.MatchString(etext) {
				val, err = strconv.ParseInt(etext[2:], 2, 64)
			}
			// log.Println("", etext, val, err)
			if err == nil {
				return valex(val), true
			}
		}
	case Lt.Text:
		return valex(etext), true
	case Lt.Word:
		switch elen {
		case 1:
			cv, ok := _contsVals[etext]
			if ok {
				return valex(cv), true
			}
		}
	}
	return res, res != nil
}

// should call after keyword and val cases
func CaseVar(ee []*lang.Elem) (base.Expression, bool) {
	if len(ee) != 1 {
		return nil, false
	}
	etype := ee[0].Type
	if etype != Lt.Word {
		return nil, false
	}
	_, ok := _contsVals[ee[0].Text]
	if ok {
		return nil, false
	}
	etext := ee[0].Text
	if !rxVar.MatchString(etext) {
		return nil, false
	}
	res := nodes.NewVarExpr(etext)
	return res, true
}
