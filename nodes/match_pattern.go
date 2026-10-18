package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

type MatchPattern interface {
	// should return bool
	// set local vars
	Do(base.Context) error
	Get() *base.Val
	// PutArg(arg any)
	Match(cx base.Context, val any) (bool, error)
}

// ====

type MCaseVal struct {
	Val any
}

func (mc *MCaseVal) Do(base.Context) error {
	return nil
}
func (mc *MCaseVal) Get() *base.Val {
	return nil
}

func (mc *MCaseVal) Match(cx base.Context, arg any) (bool, error) {
	// fmt.Printf("MCaseVal (%v). Match: %T: %v \n", mc.Val, arg, arg)
	if !SameType(mc.Val, arg) {
		return false, nil
	}
	return EqualInternTypes(mc.Val, arg), nil
}

// ====

type MCaseUnder struct {
}

func (mc *MCaseUnder) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseUnder) Get() *base.Val {
	return nil
}

func (mc *MCaseUnder) Match(cx base.Context, arg any) (bool, error) {
	// fmt.Printf("MCaseUnder  Match: %T: %v \n", arg, arg)
	return true, nil
}

// ====

type MCaseType struct {
	TypeExp base.Expression
	Var     *MCaseVar
	TypeVal *base.Type
}

func (mc *MCaseType) Do(cx base.Context) error {
	vt, err := ExprGetType(mc.TypeExp, cx)
	if err != nil {
		return err
	}
	mc.TypeVal = vt
	return nil
}
func (mc *MCaseType) Get() *base.Val {
	return nil
}

func (mc *MCaseType) Match(cx base.Context, arg any) (bool, error) {
	// fmt.Printf("MCaseType  Match: %T: %v :: type: %v \n", arg, arg, mc.TypeExp)
	err := mc.Do(cx)
	if err != nil {
		return false, err
	}
	// fmt.Printf("MCaseType  Match: tval= %T, typeId: %v \n", mc.TypeVal, mc.TypeVal.Id)
	ok, err := CheckTypeEqual(arg, mc.TypeVal)
	// fmt.Printf("MCaseType  Match: ok= %v :: err= %v \n", ok, err)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	// assign var if presents
	if mc.Var != nil {
		ok, err = mc.Var.Match(cx, arg)
		if err != nil {
			return false, err
		}

	}
	return ok, nil
}

// ====

type MCaseVar struct {
	Expr *VarExpr
}

func (mc *MCaseVar) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseVar) Get() *base.Val {
	return nil
}

func (mc *MCaseVar) Match(cx base.Context, arg any) (bool, error) {
	// fmt.Printf("MCaseVar  Match: %T: %v \n", arg, arg)
	// mc.Do(cx)
	mc.Expr.NewVar(cx)
	vr := mc.Expr.GetVar()
	if vr == nil {
		return false, fmt.Errorf("var patter couldn't add var %s", mc.Expr.GetName())
	}
	vr.Val = arg
	// fmt.Printf("MCaseVar  Match: arg= %v :: vr = %v \n", arg, vr)
	return true, nil
}

// ====

// *
type MCaseStar struct {
}

func (mc *MCaseStar) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseStar) Get() *base.Val {
	return nil
}

func (mc *MCaseStar) Match(cx base.Context, arg any) (bool, error) {
	return true, nil
}

// ====

// ?
type MCaseQMark struct {
}

func (mc *MCaseQMark) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseQMark) Get() *base.Val {
	return nil
}

func (mc *MCaseQMark) Match(cx base.Context, arg any) (bool, error) {
	return true, nil
}

// ====

// [1, 2, 3, 4]
// [1, 2, 4]
// [1, 2, ?, 4]

func mpMatchList(cx base.Context, vals []any, subs []MatchPattern) (bool, error) {
	slen := len(subs)
	vlen := len(vals)
	// fmt.Printf("MCaseList  Match: %T: %v len: %d \n", arg, arg, slen)
	vid := 0 // index of val element
	varp := 0
	opts := 0 // counter of optional by question mark
	sed := -1 // sid was iterated
	ved := -1 // vid was iterated
	vmax := vlen - 1

	// doVal := func(sid int, mp MatchPattern) (bool, error) {
	// 	vback := opts // vals for search valid value
	// 	if varp > 0 {
	// 		vback -= varp
	// 	}
	// 	curVal := -1
	// 	fmt.Printf("    >> : opts: %d, varp: %d , vid:%d, sid: %d \n", opts, varp, vid, sid)
	// 	finId := vid + vback + 1
	// 	if finId > vlen {
	// 		finId = vlen
	// 	}
	// 	// mp == nil if ending call
	// 	if mp != nil {
	// 		for ti, tv := range vals[vid:finId] {
	// 			ok, err := mp.Match(cx, tv)
	// 			// fmt.Printf("    #4: ti: %d (%d), tv:%v, mp: %v, ok:%v, err: %v \n", ti, ti+vid, tv, mp.Val, ok, err)
	// 			if err != nil {
	// 				return false, err
	// 			}
	// 			if ok {
	// 				// found 1st valid value
	// 				curVal = ti + vid
	// 			}
	// 		}
	// 		if curVal == -1 {
	// 			// valid value not found
	// 			return false, nil
	// 		}
	// 	} else {
	// 		if varp > vlen {
	// 			return false, nil
	// 		}
	// 		curVal = vlen - varp
	// 	}
	// 	vid = curVal
	// 	ved = vid
	// 	if varp > 0 {
	// 		// resolve accumulated vars, underscores
	// 		if varp > 0 {
	// 			tvid := curVal - varp
	// 			fmt.Printf("    >> : curVal: %d, tvid: %d , vid:%d, sid: %d \n", curVal, tvid, vid, sid)
	// 			for ti, tv := range vals[tvid:curVal] {
	// 				tsp := subs[sid-varp+ti]
	// 				tsp.Match(cx, tv)
	// 				// tsid++
	// 			}
	// 			varp = 0
	// 		}
	// 	}
	// 	opts = 0
	// 	return true, nil
	// }

	for sid, sp := range subs {
		// fmt.Printf(" --- LC: vid=%d, sid=%d, ptt: %T\n", vid, sid, sp)
		sed = sid
		switch sp.(type) {
		case *MCaseQMark:
			// TODO: optimization: list-pattern without ?, * can be ckecked by: slen == vlen
			opts += 1
			continue
		case *MCaseStar:
			// noMax = true
			opts += vlen - ved
			continue
		default:
			// vals has ended
			if vid > vmax {
				return false, nil
			}
		}
		ved = vid
		v := vals[vid]
		// fmt.Printf(" --- Lc: vid=%d, v=%v, sid=%d, ptt: %T\n", vid, v, sid, sp)
		// fmt.Printf("    #2: opts: %d, varp: %d \n", opts, varp)
		switch mp := sp.(type) {
		// TODO: implement any sub cases with strikt condition like MCaseVal
		// ::type, var@sub, re'gexp', [], (), {}, N{}, _{}, some(),
		case *MCaseVal:
			if opts > 0 {
				vback := opts // vals for search valid value
				if varp > 0 {
					vback -= varp
				}
				curVal := -1
				// fmt.Printf("    == : opts: %d, varp: %d \n", opts, varp)
				finId := vid + vback + 1
				if finId > vlen {
					finId = vlen
				}
				for ti, tv := range vals[vid:finId] {
					ok, err := mp.Match(cx, tv)
					// fmt.Printf("    . #4: ti: %d (%d), tv:%v, mp: %v, ok:%v, err: %v \n", ti, ti+vid, tv, mp.Val, ok, err)
					if err != nil {
						return false, err
					}
					if ok {
						// found 1st valid value
						curVal = ti + vid
						break
					}
				}
				if curVal == -1 {
					// valid value not found
					return false, nil
				}
				vid = curVal
				ved = vid
				// resolve accumulated vars, underscores
				if varp > 0 {
					tvid := curVal - varp
					for ti, tv := range vals[tvid:curVal] {
						tsp := subs[sid-varp+ti]
						tsp.Match(cx, tv)
						// tsid++
					}
					varp = 0
				}
				opts = 0
				vid++
				continue
			}
			// regular check
			ok, err := mp.Match(cx, v)
			if err != nil {
				return false, err
			}
			// fmt.Printf("    #3: ok=%v, err: %v \n", ok, err)
			if !ok {
				return false, nil
			}
			vid++
		case *MCaseVar:
			// TODO: to group non-matching cases: var, under,  (@-with-non-matched ... thinking)
			// do we need case: a @ _ ?
			// examp: [1, a@_, 2] the same [1, a, 2] - no need
			// case a @ ? or a @ * - doesn't make sense
			// conclusion: @-case - for matching right-part only, so @ is a matching case
			if opts > 0 {
				varp++
				continue
			}
			ok, err := mp.Match(cx, v)
			if err != nil {
				return false, err
			}
			if !ok {
				// can't happen in normal situation
				return false, nil
			}
			vid++
		case *MCaseUnder:
			// fmt.Printf("case _ \n")
			if opts > 0 {
				varp++
				continue
			}
			vid++
		default:
			// noMin = true
		}
	}
	// fmt.Printf("  MCaseList  fin:  ved: %d, vlen-1:%d, opts:%d, sed:%d, varp: %d  \n", ved, vlen-1, opts, sed, varp)
	// post proicess:
	// if vars in the end after ?
	if varp > 0 {
		if varp+ved+opts < vlen {
			return false, nil
		}
		tvid := vlen - varp
		tss := slen - varp
		// fmt.Printf("    >> : vlen: %d, tvid: %d , vid:%d, sid: %d \n", vlen, tvid, vid, sed)
		for ti, tv := range vals[tvid:vlen] {
			tsp := subs[tss+ti]
			// fmt.Printf("    -== : ti: %d, tv: %v , tsp: %T\n", ti, tv, tsp)
			tsp.Match(cx, tv)
			ved = tvid + ti
			sed = tss + ti
			// tsid++
		}
		varp = 0
		opts = 0
	}
	if opts > 0 {
		// resolve ending vars
		// ok, err := doVal(sed, nil)
		// fmt.Printf("  MCaseList  fin Do:  ok: %v, err: %v  \n", ok, err)
	}

	// if star in the end

	// if ? in the end
	// fmt.Printf("    opts? : vlen: %d, ved: %d , vmax:%d, opts: %d \n", vlen, ved, vmax, opts)
	if vmax-ved > opts {
		return false, nil
	}
	// fmt.Printf("  MCaseList  Match.End: sed = %v , slen-1 = %v, fail: %v \n", sed, slen-1, sed < slen-1)
	if sed < slen-1 {
		return false, nil
	}

	// fmt.Printf("  MCaseList  Match.End: slen = %v , vlen = %v \n", slen, vlen)
	return true, nil
}

// ====

type MCaseList struct {
	Subs []MatchPattern

	// optimization fields
	minLen int
	maxLen int
}

func (mc *MCaseList) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseList) Get() *base.Val {
	return nil
}

func (mc *MCaseList) Match(cx base.Context, arg any) (bool, error) {
	// min := 0
	// max := 0
	// noMin := false
	// noMax := false
	arglist, ok := arg.(*objects.ListVal)
	if !ok {
		return false, nil
	}
	vals := arglist.Elems
	// vid := 0 // vals index
	subs := mc.Subs
	return mpMatchList(cx, vals, subs)

}

func NewMCaseList(subs []MatchPattern) *MCaseList {
	// min := 0
	// max := 0
	// noMin := false
	// noMax := false
	// for i, sp := range subs {
	// 	switch mp := sp.(type) {
	// 	case *MCaseVal:
	// 	case *MCaseVar:
	// 	case *MCaseUnder:
	// 	case *MCaseStar:
	// 		noMax = true
	// 	case *MCaseQMark:

	// 	default:
	// 		noMin = true
	// 	}
	// }

	return &MCaseList{Subs: subs}
}
