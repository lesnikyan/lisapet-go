package objects

import (
	"github.com/lesnikyan/lisapet-go/base"
)

func NewCxElem(v any) *base.ContextElem {
	return &base.ContextElem{V: v}
}

type ArgInfo struct {
	Name string // if named arg in call
	Type *base.Type
}

// func domain is interface of function or method domain
// that contains all obverloaded function by one name
type FuncDom interface {
	GetFunc(args []ArgInfo) base.FuncVal
	// Funcs []FuncVal
}

// type MethMap map[string]base.FuncVal

type Context struct {
	parent *Context
	vars   map[string]*base.Var
	funcs  map[string]base.FuncVal
	types  map[string]*base.Type

	typeNames map[string]*base.Type
	typeIds   map[base.TypeId]*base.Type
}

func (cx *Context) SubContext() base.Context {
	return NewContext(cx)
}

func (cx *Context) AddType(tp *base.Type) {
	cx.types[tp.Name] = tp
}

func (cx *Context) GetType(name string) *base.Type {
	curCx := cx
	for curCx != nil {
		tp, ok := curCx.types[name]
		if ok {
			return tp
		}
		curCx = curCx.parent
	}
	// fmt.Printf("o.Ctx.GetVar no such var: %v \n", name)
	return nil
}

func (cx *Context) AddVar(vr *base.Var) {
	// fmt.Printf("o.Ctx.AddVar1 %T, %v \n", vr, vr.Name)
	if _, ok := cx.vars[vr.Name]; ok {
		panic("ctx: attempt to add an existing var " + vr.Name)
	}
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

func (cx *Context) AddFunc(fn base.FuncVal) {
	// fmt.Printf("cx.AddFunc %s \n", fn.GetName())
	cx.funcs[fn.GetName()] = fn
}

func (cx *Context) GetFunc(name string) base.FuncVal {
	curCx := cx
	for curCx != nil {
		vr, ok := curCx.funcs[name]
		if ok {
			return vr
		}
		curCx = curCx.parent
	}
	return nil
}

func (cx *Context) GetElem(name string) *base.ContextElem {
	curCx := cx
	for curCx != nil {
		// fmt.Printf("cx.GetEl#1 vars(%T, %v ) funcs(%T, %v )\n", cx.vars, len(cx.vars), cx.funcs, len(cx.funcs))
		// aa, ak := curCx.vars[name]
		// tt, tk := curCx.types[name]
		// bb, bk := curCx.funcs[name]

		// fmt.Printf("cx.GetEl#2 <%s> vars(%v, %v ) typs(%v, %v )  funcs(%v, %v )\n", name, aa, ak, tt, tk, bb, bk)
		var ok bool
		vr, ok := curCx.vars[name]
		if ok {
			return NewCxElem(vr)
		}
		tp, ok := curCx.types[name]
		if ok {
			return NewCxElem(tp)
		}
		fn, ok := curCx.funcs[name]
		if ok {
			return NewCxElem(fn)
		}
		curCx = curCx.parent
	}
	return nil
}

func NewContext(parent base.Context) *Context {
	octx, ok := parent.(*Context)
	if !ok {
		// return nil
		octx = nil
	}
	varMap := map[string]*base.Var{}
	funMap := map[string]base.FuncVal{}
	typeMap := map[string]*base.Type{}
	c := Context{parent: octx, vars: varMap, funcs: funMap, types: typeMap}
	return &c
}
