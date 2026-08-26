package objects

import (
	"regexp"
	rexs "regexp/syntax"
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
	Src     string
	Flags   []rexs.Flags
}

type Glif = rune

type Bytes = []byte
