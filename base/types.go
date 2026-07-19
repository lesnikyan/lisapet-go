package base

type Type struct {
	Id   int
	Name string

	IsUserDef bool // mostly for users struct
	Def       any  // pointer to type definition
}
