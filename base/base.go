package base

type Var struct {
	Val  any
	Name string
	Type int
}

type Val struct {
	V any
}

func NewVal(v any) *Val {
	return &Val{V: v}
}

// L : R
type Pair [2]any

type Context interface {
	AddVar(vr *Var)
	GetVar(name string) *Var
	AddFunc(fn FuncVal)
	GetElem(name string) *ContextElem
	SubContext() Context
}

type ContextElem struct {
	V any
}

type Expression interface {
	Do(Context) error
	/*
		Get() in cases:
		- value
		- definition object
		- target object (for change)
	*/
	Get() *Val
}

type Block interface {
	Do(Context) error
	Get() *Val
	Add(sub Expression)
	IsParent() bool
}

type FuncVal interface {
	Do(Context) error
	Get() *Val
	GetName() string
	SetArgVals(vals []any)
	// IsParent() bool
}

// super-expression, expression that can have sub-expression
type SupExpr interface {
	Add(sub Expression)
}

type OperExpr interface {
	Do(Context) error
	Get() *Val
	SetLeft(Expression)
	SetRight(Expression)
}

type SequenceExpr interface {
	Get() *Val
	Do(ctx Context) error
	Add(elem Expression)
	SetSubs([]Expression)
}
