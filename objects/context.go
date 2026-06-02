package objects

import dt "github.com/lesnikyan/lisapet-go/lang/datatype"

type CVal struct {
	Type dt.DType
	Item int32 // index in list
}

func NewCVal(t dt.DType, ind int32) *CVal {
	return &CVal{Type: t, Item: ind}
}

type CVar struct {
	Id int32 // index in Context.vars
	v  *CVal
}

type Context struct {
	// vars       map[string]Var // or map[int64]Val
	parent *Context

	vars []*CVar
	varm map[string]int32 // var map

	// vals table:
	valsInt    []int64
	valsFloat  []float64
	valsBool   []bool
	valsString []string
}

func NewContext(parent *Context) *Context {
	ctx := Context{parent: parent}
	ctx.vars = []*CVar{}
	ctx.varm = map[string]int32{}
	// vals

	ctx.valsInt = []int64{}
	ctx.valsFloat = []float64{}
	ctx.valsString = []string{}
	return &ctx
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

func (ctx *Context) SetVar(name string, vr *CVar) {
	ctx.vars = append(ctx.vars, vr)
	id := int32(len(ctx.vars) - 1)
	vr.Id = id
	ctx.varm[name] = id
}

func (ctx *Context) AddVal(tp dt.DType, value interface{}) (*CVal, bool) {
	switch tp {
	case dt.Int:
		ival, ok := value.(int64)
		if !ok {
			return nil, false
		}
		ctx.valsInt = append(ctx.valsInt, ival)
	case dt.Float:
		fval, ok := value.(float64)
		if !ok {
			return nil, false
		}
		ctx.valsFloat = append(ctx.valsFloat, fval)
	case dt.String:
		sval, ok := value.(string)
		if !ok {
			return nil, false
		}
		ctx.valsString = append(ctx.valsString, sval)
	}
	return nil, false
}

func (c *Context) TakeVal(v *CVal) {

}
