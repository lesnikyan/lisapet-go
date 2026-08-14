package base

type Type struct {
	Id   TypeId
	Name string

	IsUserDef bool // mostly for users struct
	Def       any  // pointer to type definition
}

type TypeInf interface {
	GetId() TypeId
	GetName() string
	IsCustom() bool // if user defined
}

type UserType interface {
	GetInf() TypeInf
	GetDefinition() any // definition of type (struct def, etc)
}

type TypeId = int

const (
	TypeAny TypeId = iota + 1001
	TypeNull
	TypeBool
	TypeInt
	TypeFloat
	TypeString
	TypeList
	TypeTuple
	TypeDict
	TypeBytes
	TypeGlif
	TypeFunction
	TypeEnum
	TypeGrup
	TypeMaybe
	// Last defined type:
	TypeUndefined TypeId = 2001
)

var typeId = TypeUndefined

// TODO: resolve possible issues with multi-task access
func NextTypeId() TypeId {
	typeId += 1
	return typeId
}

func BaseType(name string, id TypeId) *Type {
	// tid := NextTypeId()
	return &Type{Id: id, Name: name}
}

func DefineUserType(name string, def any) *Type {
	tid := NextTypeId()
	return &Type{Id: tid, Name: name, IsUserDef: true, Def: def}
}

// func CompareType(a *Type, b *Type) bool {
// 	return a.Id == b.Id
// }

func DefaultVal(ti TypeId) any {
	switch ti {
	case TypeBool:
		return false
	case TypeInt:
		return int64(0)
	case TypeFloat:
		return float64(0)
	case TypeString:
		return ""
	case TypeList:
		return nil // maybe []
	case TypeTuple:
		return nil
	case TypeDict:
		return nil
	case TypeBytes:
		return nil
	case TypeGlif:
		return 0
	case TypeFunction:
		return nil
	case TypeEnum:
		return nil
	case TypeGrup:
		return nil
	}
	return nil
}
