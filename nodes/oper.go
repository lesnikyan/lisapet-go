package nodes

// ob "github.com/lesnikyan/lisapet-go/objects"

type Opid int

const (
	nn Opid = iota
	OpAssign
	OpPlus
	OpMinus
	OpMult
	OpDiv
	OpPow
	OpRoot
	OpEqual
	OpNotEqual
	OpMoreEqual
	OpLessEqual
	OpPlusAssign
	OpMinusAssign
	OpMultAssign
	OpDivAssign
	OpAnd
	OpOr
	OpNot
	OpBiAnd
	OpBinOr
	OpBinNot
	OpXor
	OpDot
	OpComma
	OpColon
	OpSemicolon
)

type Oper struct {
	Sign string
	Id   Opid
}

func NewOper(sg string, id Opid) *Oper {
	return &Oper{Sign: sg, Id: id}
}
