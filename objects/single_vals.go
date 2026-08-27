package objects

import (
	"fmt"
	"regexp"
	rexs "regexp/syntax"

	"github.com/lesnikyan/lisapet-go/base"
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

type Bytes []byte

func (bb Bytes) GetElem(i int64) (*base.Val, error) {
	index := int(i)
	if index < 0 {
		index = len(bb) + index
	}
	if index < 0 || index >= len(bb) {
		return nil, fmt.Errorf("Bytes sequence with len= %d doesn't have element %d", len(bb), index)
	}
	v := bb[int(index)]
	return base.NewVal(v), nil
}
