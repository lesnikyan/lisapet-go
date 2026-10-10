package lang

import Lt "github.com/lesnikyan/lisapet-go/lang/lt"

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
	Src     string
	Indent  int
	Special Spec
}
