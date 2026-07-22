package nodes

import "github.com/lesnikyan/lisapet-go/base"

// popup

type NodeType int

const (
	NodeBreak = iota + 101
	NodeContinue
	NodeReturn
	NodeFuncDef
	NodeStructDef
	NodeIf
	NodeElse
	NodeMatch
)

type PopUp struct {
	Parent NodeType
	Res    *base.Val
}
