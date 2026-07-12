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
	OpPercent
	OpPow
	OpRoot
	OpEqual
	OpNotEqual
	OpMore
	OpMoreEqual
	OpLess
	OpLessEqual
	OpPlusAssign
	OpMinusAssign
	OpMultAssign
	OpDivAssign
	OpAnd
	OpOr
	OpNot
	OpBitAnd
	OpBitOr
	OpBitNot
	OpXor
	OpBitLShift
	OpBitRShift
	OpDot
	OpComma
	OpColon
	OpSemicolon
	OpLArrow
	OpRArrow
	OpDuoColon
	OpAt
	OpDollar
	OpQMark
	OpQmAndColon
	OpTildLArrow
	OpSlashColon
	OpColonQm
)

type Oper struct {
	Sign string
	Id   Opid
}

func OperIndex(s string) Opid {
	id, ok := _operMap[s]
	if !ok {
		return -1
	}
	return id
}
func NewOper(sg string, id Opid) *Oper {
	return &Oper{Sign: sg, Id: id}
}

var _operStrMap = map[Opid]string{
	OpAssign:      "=",
	OpPlus:        "+",
	OpMinus:       "-",
	OpMult:        "*",
	OpDiv:         "/",
	OpPow:         "**",
	OpRoot:        "^/",
	OpEqual:       "==",
	OpNotEqual:    "!=",
	OpLess:        "<",
	OpMore:        ">",
	OpMoreEqual:   ">=",
	OpLessEqual:   "<=",
	OpPlusAssign:  "+=",
	OpMinusAssign: "-=",
	OpMultAssign:  "*=",
	OpDivAssign:   "/=",
	OpAnd:         "&&",
	OpOr:          "||",
	OpNot:         "!",
	OpBitAnd:      "&",
	OpBitOr:       "|",
	OpBitNot:      "~",
	OpXor:         "^",
	OpBitLShift:   "<<",
	OpBitRShift:   ">>",
	OpDot:         ".",
	OpComma:       ",",
	OpColon:       ":",
	OpSemicolon:   ";",
	OpLArrow:      "<-",
	OpRArrow:      "->",
	OpDuoColon:    "::",
	OpAt:          "@",
	OpDollar:      "$",
	OpQMark:       "?",
	OpQmAndColon:  "?:",
	OpTildLArrow:  "~>",
	OpSlashColon:  "/:",
	OpColonQm:     ":?",
}

var _operMap = func() map[string]Opid {
	res := map[string]Opid{}
	for id, s := range _operStrMap {
		res[s] = id
	}
	return res
}()
