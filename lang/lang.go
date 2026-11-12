package lang

import (
	"fmt"
	"strings"

	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
)

type Elem struct {
	Text string
	Type Lt.Lt
}

type Spec int

const (
	SpecMatch = 1 // mark of matching patter string
)

type CLine struct {
	Elems   []*Elem
	Subs    map[int][]*Elem // preparsed strings, etc. key = index in Elems
	Src     string
	Indent  int
	Special Spec
}

func FPrintElems(elems []*Elem) string {
	ss := make([]string, len(elems))
	// println("lenEl", len(elems))
	for i := 0; i < len(elems); i++ {
		// println("?=", elems[i].Text)
		ss[i] = elems[i].Text
	}
	stt := strings.Join(ss, ",")
	return fmt.Sprintf("`%s`", stt)
}
