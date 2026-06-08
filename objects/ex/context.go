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

type UVal[T VType] struct {
	v T
}

func (n *UVal[T]) GetVal() any {
	return n.v
}

type Int struct {
	UVal[int64]
}

func IntVal(v int64) T2Val {
	return &Int{UVal[int64]{v}}
}

// func (n *Int) GetVal() int64 {
// 	return n.v
// }

// func GetType() dt.DType {
// 	return dt.Int
// }

type Float struct {
	// v float64
	UVal[float64]
}

func FloatVal(v float64) T2Val {
	return &Float{UVal[float64]{v}}
}

type String struct {
	UVal[string]
}

func StringVal(v string) T2Val {
	return &String{UVal[string]{v}}
}

type TT interface {
}

type T2Val interface {
	GetVal() any
	// GetType() dt.DType
}

func TestIn() {
	v := &Int{UVal: UVal[int64]{123}}
	ctx := NewContext(nil)
	ctx.PutVal(v)

}

func put[T VType](val T, target []T) (int, []T) {
	id := len(target)
	target = append(target, val)
	return id, target
}

func tVal(val any) T2Val {
	switch v := val.(type) {
	case int64:
		return IntVal(v)
	case float64:
		return FloatVal(v)
	case string:
		return StringVal(v)
	}
	return nil
}

func (ctx *Context) PutVal(value T2Val) *CVal {
	val := value.GetVal()
	var id int
	var tp dt.DType
	switch vT := val.(type) {
	case int64:
		i, tg := put(vT, ctx.valsInt)
		ctx.valsInt = tg
		tp, id = dt.Int, i
	case float64:
		i, tg := put(vT, ctx.valsFloat)
		ctx.valsFloat = tg
		tp, id = dt.Int, i
	case string:
		i, tg := put(vT, ctx.valsString)
		ctx.valsString = tg
		tp, id = dt.Int, i
	}

	return &CVal{Type: tp, Item: id}
}

func (ctx *Context) Add2(value TVal) *CVal {
	val := value.GetVal()
	var id int
	var tp dt.DType
	switch vT := val.(type) {
	case int64:
		i, tg := put(vT, ctx.valsInt)
		ctx.valsInt = tg
		tp, id = dt.Int, i
	case float64:
		i, tg := put(vT, ctx.valsFloat)
		ctx.valsFloat = tg
		tp, id = dt.Int, i
	case string:
		i, tg := put(vT, ctx.valsString)
		ctx.valsString = tg
		tp, id = dt.Int, i
	}

	return &CVal{Type: tp, Item: id}
}

func (ctx *Context) Update(value T2Val, target *CVal) bool {
	val := value.GetVal()
	id := target.Item
	if id < 0 {
		return false
	}
	var tp dt.DType = target.Type
	switch vT := val.(type) {
	case int64:
		if tp != dt.Int || id >= len(ctx.valsInt) {
			return false
		}
		ctx.valsInt[target.Item] = vT
		return true
	}
	return false
}

func (ctx *Context) GetVal(vv *CVal) T2Val {
	return IntVal(111)
}

type T2Var struct {
	val  T2Val
	Type dt.DType
}
