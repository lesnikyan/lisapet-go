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

func (bb Bytes) Add(n any) (Bytes, error) {
	var b byte
	switch v := n.(type) {
	case byte:
		b = v
	case int64:
		b = byte(v)
	default:
		return nil, fmt.Errorf("bytes.Add: incorrect argument type %T", v)
	}
	bb = append(bb, b)
	// fmt.Println("bytes: ad", bb)
	return bb, nil
}
