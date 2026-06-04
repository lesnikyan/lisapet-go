package objects

import dt "github.com/lesnikyan/lisapet-go/lang/datatype"

type CVal struct {
	Type dt.DType
	Item int // index in list
}

func NewCVal(t dt.DType, ind int) *CVal {
	return &CVal{Type: t, Item: ind}
}

type CVar struct {
	Id   int // index in Context.vars
	Type dt.DType
	v    *CVal
}

type Context struct {
	// vars       map[string]Var // or map[int64]Val
	parent *Context

	vars []*CVar
	varm map[string]int // var map

	// vals table:
	valsInt    []int64
	valsFloat  []float64
	valsBool   []bool
	valsString []string
}

func NewContext(parent *Context) *Context {
	ctx := Context{parent: parent}
	ctx.vars = []*CVar{}
	ctx.varm = map[string]int{}
	// vals

	ctx.valsInt = []int64{}
	ctx.valsFloat = []float64{}
	ctx.valsString = []string{}
	return &ctx
}

func (ctx *Context) SetVar(name string, vr *CVar) {
	ctx.vars = append(ctx.vars, vr)
	id := len(ctx.vars) - 1
	vr.Id = id
	ctx.varm[name] = id
}

func (ctx *Context) GetVar(name string) (*CVar, *Context) {
	var nctx *Context = ctx
	for nctx != nil {
		vr, ok := nctx.varm[name]
		if ok {
			return nctx.vars[vr], nctx
		}
		nctx = nctx.parent
	}
	return nil, nil
}

func (ctx *Context) AddVal(tp dt.DType, value interface{}) (*CVal, bool) {
	cval := CVal{Type: tp}
	var id int = 0
	switch tp {
	case dt.Int:
		ival, ok := value.(int64)
		if !ok {
			return nil, false
		}
		id = len(ctx.valsInt)
		ctx.valsInt = append(ctx.valsInt, ival)
	case dt.Float:
		fval, ok := value.(float64)
		if !ok {
			return nil, false
		}
		id = len(ctx.valsFloat)
		ctx.valsFloat = append(ctx.valsFloat, fval)
	case dt.String:
		sval, ok := value.(string)
		if !ok {
			return nil, false
		}
		id = len(ctx.valsString)
		ctx.valsString = append(ctx.valsString, sval)
	}
	cval.Item = id
	return &cval, true
}

func (c *Context) GetInt(v *CVal) (int64, bool) {
	if v.Item < len(c.valsInt) {
		return 0, false
	}
	return c.valsInt[v.Item], true
}

func (ctx *Context) AddInt(value int64) (*CVal, bool) {
	id := len(ctx.valsInt)
	ctx.valsInt = append(ctx.valsInt, value)
	return &CVal{Type: dt.Int, Item: id}, true
}

type VType interface {
	int64 | float64 | string | byte | bool | []byte
}

type XVal[T VType] interface {
	GetVal() T
	// GetType() dt.DType
}

type U[T VType] struct {
	v T
}

func (n *U[T]) GetVal() T {
	return n.v
}

type Int struct {
	// v int64
	U[int64]
}

// func (n *Int) GetVal() int64 {
// 	return n.v
// }

// func GetType() dt.DType {
// 	return dt.Int
// }

type Float struct {
	// v float64
	U[int64]
}

func put[T VType](val T, target []T) (int, []T) {
	id := len(target)
	target = append(target, val)
	return id, target
}

func (ctx *Context) Add1(value TVal) *CVal {
	val := value.GetVal()
	var id int
	var tp dt.DType
	switch val := val.(type) {
	case int64:
		i, tg := put(val, ctx.valsInt)
		ctx.valsInt = tg
		tp, id = dt.Int, i
	case float64:
		i, tg := put(val, ctx.valsFloat)
		ctx.valsFloat = tg
		tp, id = dt.Int, i
	}
	return &CVal{Type: tp, Item: id}
}
