package nodes

import (
	"fmt"
	"strings"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

func parr(vv []any) string {
	ss := make([]string, len(vv))
	for i, v := range vv {
		ss[i] = fmt.Sprintf("%v", v)
	}
	return strings.Join(ss, ", ")
}

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

type MCaseRegexp struct {
	Regexp *objects.Regexp
}

func (mc *MCaseRegexp) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseRegexp) Get() *base.Val {
	return nil
}

func (mc *MCaseRegexp) Match(cx base.Context, arg any) (bool, error) {
	s, ok := arg.(string)
	if !ok {
		return false, nil
	}
	// fmt.Printf("MCaseRegexp.Match#0 src: `%s`, arg: '%s' \n", mc.Regexp.Src, s)
	return mc.Regexp.Match(s), nil
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
	return &MCaseList{Subs: subs}
}

// ====

type MCaseTuple struct {
	Subs []MatchPattern

	// optimization fields
	minLen int
	maxLen int
}

func (mc *MCaseTuple) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseTuple) Get() *base.Val {
	return nil
}

func (mc *MCaseTuple) Match(cx base.Context, arg any) (bool, error) {
	arglist, ok := arg.(*objects.TupleVal)
	if !ok {
		return false, nil
	}
	vals := arglist.Elems
	// vid := 0 // vals index
	subs := mc.Subs
	return mpMatchList(cx, vals, subs)

}

func NewMCaseTuple(subs []MatchPattern) *MCaseTuple {
	return &MCaseTuple{Subs: subs}
}

// ====

// key : val
type MCaseKeyVal struct {
	Key MatchPattern
	Val MatchPattern
}

func (mc *MCaseKeyVal) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseKeyVal) Get() *base.Val {
	return nil
}

func (mc *MCaseKeyVal) Match(cx base.Context, arg any) (bool, error) {
	// argDc, ok := arg.(*objects.DictVal)
	// if !ok {
	// 	return false, nil
	// }
	// vals := argDc.Vmap
	// var found any
	// // check key pattern
	// switch kp := mc.Key.(type) {
	// case *MCaseVal:
	// 	k := kp.Val
	// 	v, ok := vals[k]
	// 	if !ok {
	// 		return false, nil
	// 	}
	// 	found = v
	// case *MCaseVar:
	// 	// var nextKey any
	// 	// vals[nextKey]
	// 	// kp.Match(cx, )
	// default:

	// }

	// //check val
	// switch vp := mc.Key.(type) {
	// case *MCaseVal:
	// 	return vp.Match(cx, found)
	// case *MCaseVar:
	// 	return vp.Match(cx, found)
	// default:
	// 	return vp.Match(cx, found)
	// }

	return false, nil
}

func NewMCaseKeyVal(key MatchPattern, val MatchPattern) *MCaseKeyVal {
	return &MCaseKeyVal{Key: key, Val: val}
}

// ====

// {sub:Elems,...}
type MCaseDict struct {
	Subs []MatchPattern // *MCaseDictSub or optional elems: *, ?
}

func (mc *MCaseDict) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseDict) Get() *base.Val {
	return nil
}

func (mc *MCaseDict) Match(cx base.Context, arg any) (bool, error) {
	argDc, ok := arg.(*objects.DictVal)
	if !ok {
		return false, nil
	}
	vals := argDc.Vmap
	// vid := 0 // vals index
	// fmt.Printf(" lenV=%d lenS=%d \n", len(vals), len(mc.Subs))
	if len(mc.Subs) == 0 {
		if len(vals) == 0 {
			return true, nil
		}
		return false, nil
	}
	subs := mc.Subs
	skip := map[any]any{}
	opt := 0
	star := false
	// need := 0
	for _, sb := range subs {
		// fmt.Printf("  MDP sub %T\n", sb)
		switch pair := sb.(type) {
		case *MCaseKeyVal:
			// fmt.Printf(" D:Pair K:%T, V:%T \n", pair.Key, pair.Val)
			if len(skip) >= len(vals) {
				// too small count of pairs of arg
				return false, nil
			}
			switch kk := pair.Key.(type) {
			case *MCaseVal:
				// if constant key-pattern we can check key in arg dict
				v, ok := vals[kk.Val]
				// fmt.Printf(" DKc-CaseVal k:%v, ok:%v \n", kk.Val, ok)
				if !ok {
					return false, nil
				}
				vok, err := pair.Val.Match(cx, v)
				if err != nil {
					return false, err
				}
				// fmt.Printf(" DV@1 v:%v, err:%v \n", vok, err)
				if !vok {
					// here key has matched, but val has not
					return false, nil
				}
				skip[kk.Val] = true
			default:
				need := 0
				got := 0
				for k, v := range vals {
					// think about optimization
					// fmt.Printf(" DKa-CaseVal k:%v, ok:%v \n", k, v)
					if _, ok := skip[k]; ok {
						// skip previously found key
						continue
					}
					// first check more strict pattern
					// order: val|regexp|container, type, var|underscore
					// TODO: @-pattern will need some fix
					// mp1 := pair.Key
					// mp2 := pair.Val
					kfirst := true
					switch pair.Key.(type) {
					case *MCaseVar, *MCaseUnder:
						switch pair.Val.(type) {
						case *MCaseType:
							kfirst = false
							need++
						case *MCaseVal, *MCaseRegexp:
							kfirst = false
							need++
						}
					case *MCaseType:
						switch pair.Val.(type) {
						case *MCaseVal, *MCaseRegexp:
							kfirst = false
							need++
						}
					default:

					}

					// fmt.Printf(" DKa-CaseVal, k-1st:%v \n", kfirst)
					// TODO: think about code optimization
					if kfirst {
						// first check key-pattern
						kok, err := pair.Key.Match(cx, k)
						if err != nil {
							return false, err
						}
						if !kok {
							continue
						}
						// check value-pattern
						vok, err := pair.Val.Match(cx, v)
						if err != nil {
							return false, err
						}
						if vok {
							// found pair
							skip[k] = true
							got++
							break
						}
					} else {
						// first check value-pattern
						vok, err := pair.Val.Match(cx, v)
						// fmt.Printf(" -- DKv1st#1:,ok:%v, err:%v \n", vok, err)
						if err != nil {
							return false, err
						}
						if !vok {
							continue
						}
						kok, err := pair.Key.Match(cx, k)
						// fmt.Printf(" -- DKv1st#2:,ok:%v, err:%v \n", kok, err)
						if err != nil {
							return false, err
						}
						if kok {
							skip[k] = true
							got++
							break
						}
					}
				}
				if need > 0 {
					if got < 1 {
						return false, nil
					}
				}
			}
		case *MCaseQMark:
			opt++
		case *MCaseStar:
			star = true
		default:
			return false, fmt.Errorf("dict mp-case: incorrect subcase of dict: %T", sb)
		}
	}
	// after all subcases we should hav `skip` filled by keys of found pairs
	// fmt.Printf(" == Dc-fin lskip:%v, lvals:%v  \n", len(skip), len(vals))
	if star {
		return true, nil
	}
	return len(skip)+opt >= len(vals), nil
	// return false, nil
}

// sized slice
type SzSl struct {
	Vals []any
	size int
	Id   int
}

func (sz *SzSl) Add(v any) {
	if sz.Id >= sz.size {
		panic("sized slice is out of range")
	}
	sz.Vals[sz.Id] = v
	sz.Id++
}

func NewSz(size int) *SzSl {
	vv := make([]any, size)
	return &SzSl{Vals: vv, size: size, Id: 0}
}

const (
	KVal   = 0
	KType  = 1
	VVal   = 2
	VType  = 3
	KVar   = 4
	KUnd   = 5
	MTQmk  = 6
	MTStar = 7
)

func NewMCaseDict(subs []MatchPattern) *MCaseDict {
	slen := len(subs)
	mps := make([]MatchPattern, slen) // sub patterns: pair | * | ?
	/*
		// kc + tk + _v + tv + vrr + _k + qm + stars
			0 key: MCval, regexp
			1 - type
			2 - K var
			3 val: MCVal, regexp
			4 - type
			5 - K _
			6 ?
			7 *
			// - V var
			// - V _
	*/
	order := make([]int, slen*8)
	nums := make([]int, 8)
	mid := -1
	for i, pair := range subs {
		mid = -1
		switch ppt := pair.(type) {
		case *MCaseKeyVal:
			mk := ppt.Key
			mv := ppt.Val
			// fmt.Printf(" NewMCD i:%d k:%T, v:%T \n", i, mk, mv)
			// Key can be: num, bool, string, null
			switch mk.(type) {
			case *MCaseVal, *MCaseRegexp:
				mid = KVal
			case *MCaseType:
				mid = KType
			case *MCaseVar:
				switch mv.(type) {
				case *MCaseVal, *MCaseRegexp:
					mid = VVal
				case *MCaseType:
					mid = VType
				default:
					mid = KVar
				}
			case *MCaseUnder:
				switch mv.(type) {
				case *MCaseVal, *MCaseRegexp:
					mid = VVal
				case *MCaseType:
					mid = VType
				default:
					mid = KUnd
				}
			}
		case *MCaseQMark:
			mid = MTQmk
		case *MCaseStar:
			mid = MTStar
		}
		// println(" // mid:", mid)
		nx := nums[mid]
		order[slen*mid+nx] = i
		nums[mid]++
	}
	id := 0
	// println("-- nums:", parr(vals2anis(nums)))
	for i := range len(nums) {
		n := nums[i]
		// println("-- -- :", i, n)
		for j := range n {
			mps[id] = subs[order[slen*i+j]]
			id++
		}
	}
	return &MCaseDict{Subs: mps}
}

// ====

// A{}, _{}
type MCaseStruct struct {
	Type MatchPattern // | MatchPattern ?
	TDef *objects.StructDef
	Subs []MatchPattern
}

func (mc *MCaseStruct) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseStruct) Get() *base.Val {
	return nil
}

func (mc *MCaseStruct) Match(cx base.Context, arg any) (bool, error) {
	inst, ok := arg.(*objects.StructInst)
	if !ok {
		// not struct instance in arg
		return false, nil
	}

	switch typeVar := mc.Type.(type) {
	case *MCaseUnder:
		// universal | noname type
		// universal | noname type
	case *MCaseVar:
		// concrete type
		typeName := typeVar.Expr.GetName()
		tpv := cx.GetType(typeName) // type val
		if tpv == nil {
			// type not found
			return false, fmt.Errorf("match struct: type `%s` not found", typeName)
		}
		ok, err := CheckTypeEqual(arg, tpv)
		if err != nil {
			return false, err
		}
		tdef := tpv.Def
		sdef := tdef.(*objects.StructDef)
		if !ok {
			// correct negative result
			return false, nil
		}
		subs := mc.Subs
		if len(subs) == 0 {
			// no fields defined in pattern
			return true, nil
		}
		for _, sub := range subs {
			// fmt.Printf("Str.fld ptr: %T \n", sub)
			switch fld := sub.(type) {
			case *MCaseKeyVal:
				// main case: field : value
				// Key should be VarExpr
				fnameMP, ok := fld.Key.(*MCaseVar)
				if !ok {
					return false, fmt.Errorf("bad pattern of struct field %T", fld.Key)
				}
				fname := fnameMP.Expr.GetName()
				if _, ok := sdef.Fields[fname]; !ok {
					// no such field
					return false, fmt.Errorf("struct `%s` doesn't have field `%s`", typeName, fname)
				}
				fval := inst.Get(fname)
				if fval == nil {
					return false, fmt.Errorf("struct doesnt have any value if field `%s`", fname)
				}
				res, err := fld.Val.Match(cx, fval.V)
				// fmt.Printf("Str.fld Val: %T fval:%v ;; >> %v, %v \n", fld.Val, fval, res, err)
				if err != nil {
					return false, err
				}
				if !res {
					return false, nil
				}
			}
		}

	default:
		return false, fmt.Errorf("match struct: incorrect empty type name")
	}
	return true, nil
}

func NewMCaseStruct(stype MatchPattern, subs []MatchPattern) *MCaseStruct {
	return &MCaseStruct{Type: stype, Subs: subs}
}

// ====

// @ TODO
type MCaseAt struct {
}

func (mc *MCaseAt) Do(cx base.Context) error {
	return nil
}
func (mc *MCaseAt) Get() *base.Val {
	return nil
}

func (mc *MCaseAt) Match(cx base.Context, arg any) (bool, error) {
	return true, nil
}
