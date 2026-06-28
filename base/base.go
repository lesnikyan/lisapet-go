package base

type Var struct {
	Val  any
	Name string
	Type int
}

type Context interface {
	AddVar(vr *Var)
	GetVar(name string) *Var
}

type Expression interface {
	Do(Context) error
	Get() any
}
type Block interface {
	Do(Context) error
	Get() any
	Add(sub Expression)
}

// super-expression, expression that can have sub-expression
type SupExpr interface {
	Add(sub Expression)
}

type BinOperExpr interface {
	Do(Context) error
	Get() any
	SetLeft(Expression)
	SetRight(Expression)
}

type SequenceExpr interface {
	Get() any
	Do(ctx Context) error
	Add(elem Expression)
	SetSubs([]Expression)
}
