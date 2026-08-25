package objects

import (
	"regexp"
)

// null value
type Null struct {
}

// internal value, can be produced from EmptyExpr
type EmptyVal struct {
}

// type StructVal struct {
// 	Fields []*base.Var
// }

type Regexp struct {
	Pattern *regexp.Regexp
}

type Glif = rune

// func Str2Glif(s string) (Glif, bool) {
// 	rrs := []rune(s)
// 	if len(rrs) != 1 {
// 		return 0, false
// 	}
// 	println("len(rrs)", len(rrs), 'Ы', 'Ф', ">>", s, ":", rrs[0])
// 	return rrs[0], true
// }
