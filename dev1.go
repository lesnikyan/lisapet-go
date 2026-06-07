package main

import "reflect"

type Null struct {
}

type Var struct {
	Val  any
	Name string
	Type int
}

type Context struct {
	vars map[string]*Var
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

func f1() {
	n := 1
	reflect.TypeOf(n)
}

type Expression interface {
	Do(*Context)
	Get() any
}

type Block struct {
	subs []Expression
	res  any
}

func (bk *Block) Do(cx *Context) {
	var res any = nil
	for _, exp := range bk.subs {
		exp.Do(cx)
	}
	bk.res = res
}
func (bk *Block) Get() any {
	return bk.res
}

// ValExpr
type ValExpr struct {
	Val any
}

func (vex *ValExpr) Do(cx *Context) {
	// do nothing
}
func (vex *ValExpr) Get() any {
	return vex.Val
}

// VarExpr
type VarExpr struct {
	name string
	vr   *Var
}

func (ex *VarExpr) Do(cx *Context) {
	vr := cx.GetVar(ex.name)
	if vr == nil {
		vr := &Var{Name: ex.name}
		cx.AddVar(vr)
	}
	ex.vr = vr
}
func (ex *VarExpr) Get() any {
	return ex.vr
}

type OperAssign struct {
	oper  string
	left  Expression
	right Expression
	res   any
}

func Any2Val(obj any) any {
	switch src := obj.(type) {
	case *Var:
		return src.Val
	}
	return nil
}

func (op *OperAssign) Do(cx *Context) {
	op.left.Do(cx)
	op.right.Do(cx)
	lop := op.left.Get()
	rval := op.right.Get()
	switch target := lop.(type) {
	case *Var:
		target.Val = Any2Val(rval)
	}
}

func (op *OperAssign) Get() any {
	return op.res // make sense for last expression in the Block
}
