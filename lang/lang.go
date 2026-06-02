package lang

import Lt "github.com/lesnikyan/lisapet-go/lang/lt"

type Elem struct {
	Text string
	Type Lt.Lt
}

type CLine struct {
	Elems []Elem
	Src   string
}
