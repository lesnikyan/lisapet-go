package objects

import "github.com/lesnikyan/lisapet-go/base"

type Context struct {
	parent *Context
	vars   map[string]*base.Var
}

func NewContext(parent *Context) *Context {
	c := Context{parent: parent, vars: map[string]*base.Var{}}
	return &c
}

func (cx *Context) AddVar(vr *base.Var) {
	cx.vars[vr.Name] = vr
}

func (cx *Context) GetVar(name string) *base.Var {
	vr, ok := cx.vars[name]
	if ok {
		return vr
	}
	return nil
}
