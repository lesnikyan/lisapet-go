package objects

type Context struct {
	parent *Context
	vars   map[string]*Var
}

func NewContext(parent *Context) *Context {
	c := Context{parent: parent, vars: map[string]*Var{}}
	return &c
}

func (cx *Context) AddVar(vr *Var) {
	cx.vars[vr.Name] = vr
}

func (cx *Context) GetVar(name string) *Var {
	vr, ok := cx.vars[name]
	if ok {
		return vr
	}
	return nil
}
