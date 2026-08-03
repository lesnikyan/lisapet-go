package base

type Type struct {
	Id   int
	Name string

	IsUserDef bool // mostly for users struct
	Def       any  // pointer to type definition
}

var typeId = 1001

func DefineType(name string, usdef bool) *Type {
	return &Type{Name: name, IsUserDef: usdef}
}

func CompareType(a *Type, b *Type) bool {
	return a.Id == b.Id
}
