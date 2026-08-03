package main

import (
	"fmt"

	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/lesnikyan/lisapet-go/parser"
)

func main() {
	nn := Lt.Block
	fmt.Println("lisapet main", nn)

	parser.Foo()
}
