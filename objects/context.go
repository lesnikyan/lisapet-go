package objects

import (
	"github.com/lesnikyan/lisapet-go/base"
)

type Context struct {
	parent *Context
	vars   map[string]*base.Var
	funcs  map[string]*Function

	typeNames map[string]*base.Type
	typeIds   map[int]*base.Type
}

func (cx *Context) AddVar(vr *base.Var) {
	// fmt.Printf("o.Ctx.AddVar1 %T, %v \n", vr, vr.Name)
	cx.vars[vr.Name] = vr
	// fmt.Printf("o.Ctx.AddVar2 %T, %v \n", cx.vars, cx.vars)
}

func (cx *Context) GetVar(name string) *base.Var {
	curCx := cx
	for curCx != nil {
		vr, ok := curCx.vars[name]
		if ok {
			return vr
		}
		curCx = curCx.parent
	}
	// fmt.Printf("o.Ctx.GetVar no such var: %v \n", name)
	return nil
}

func NewContext(parent base.Context) *Context {
	octx, ok := parent.(*Context)
	if !ok {
		// return nil
		octx = nil
	}
	c := Context{parent: octx, vars: map[string]*base.Var{}}
	return &c
}
