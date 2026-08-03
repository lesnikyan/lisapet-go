package nodes

import "github.com/lesnikyan/lisapet-go/base"

type BrType int

const (
	UnknownBr BrType = 100
	RoundBr   BrType = 101
	SquareBr  BrType = 102
	CurlyBr   BrType = 103
)

var brTypeMap = map[string]BrType{
	"(": RoundBr,
	"[": SquareBr,
	"{": CurlyBr,
}

func GetBrType(s string) BrType {
	t, ok := brTypeMap[s]
	if ok {
		return t
	}
	return UnknownBr
}

type Brackets struct {
	Type BrType
	Sub  base.Expression
}

func (br *Brackets) Get() *base.Val {
	return br.Sub.Get()
}

func (br *Brackets) Do(ctx base.Context) error {
	return br.Sub.Do(ctx)
}
