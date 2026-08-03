package nodes

import "github.com/lesnikyan/lisapet-go/base"

//=== BREAK

type BreakExp struct{}

func (br *BreakExp) Do(cx base.Context) error {
	// nothing to do
	return nil
}

func (br *BreakExp) Get() *base.Val {
	return nil
}

// ===== CONTINUE

type ContinueExp struct{}

func (br *ContinueExp) Do(cx base.Context) error {
	// nothing to do
	return nil
}

func (br *ContinueExp) Get() *base.Val {
	return nil
}

// ==== RETURN

type ReturnExp struct {
	Sub base.Expression
	res *base.Val
}

func (rr *ReturnExp) Do(cx base.Context) error {
	rr.res = nil
	if rr.Sub == nil {
		return nil
	}
	err := rr.Sub.Do(cx)
	if err != nil {
		return err
	}
	rr.res = rr.Sub.Get()
	return nil
}

func (rr *ReturnExp) Get() *base.Val {
	return rr.res
}

func NewReturn(sub base.Expression) *ReturnExp {
	return &ReturnExp{Sub: sub}
}
