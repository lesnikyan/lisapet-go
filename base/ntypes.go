package base

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
	Res    *Val
}

func NewPopUp(nt NodeType, res *Val) *PopUp {
	return &PopUp{Parent: nt, Res: res}
}
