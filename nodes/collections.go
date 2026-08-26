package nodes

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/lang"
	"github.com/lesnikyan/lisapet-go/objects"
)

type ListExpr struct {
	Seq  *SequenceComma
	Subs []base.Expression
	res  *objects.ListVal
}

func (cs *ListExpr) Add(sub base.Expression) {
	if cs.Subs == nil {
		cs.Subs = []base.Expression{}
	}
	// fmt.Printf("[] List Add: (%T, %v)  \n", sub, sub)
	cs.Subs = append(cs.Subs, sub)
}

func (cs *ListExpr) Get() *base.Val {
	return base.NewVal(cs.res)
}

// func (op *ListExpr) IsParent() bool {
// 	// don't used as Block by default
// 	return false
// }

func (cs *ListExpr) Do(ctx base.Context) error {
	src := make([]base.Expression, len(cs.Seq.Subs))
	copy(src, cs.Seq.Subs)
	src = append(src, cs.Subs...)
	res := make([]any, len(src))
	// fmt.Println("[] List Do:", len(src))
	for i, ex := range src {
		err := ex.Do(ctx)
		if err != nil {
			return err
		}
		res[i] = GetExprVal(ex, ctx)
	}
	cs.res = objects.NewListVal(res)
	return nil
}

// ===========

type ColElem struct {
	Src  any // ListVal, DictVal, etc
	KVal any // int64, string, etc
}

func (cc *ColElem) Get() *base.Val {
	switch src := cc.Src.(type) {
	case *objects.ListVal:
		index, ok := cc.KVal.(int64)
		if !ok {
			// return errors.New("incorrect type of index in collection-elem expr")
			panic("incorrect type of index in collection-elem expr")
		}
		res, err := src.GetElem(index)
		if err != nil {
			return nil
		}
		return res

	case *objects.TupleVal:
		index, ok := cc.KVal.(int64)
		if !ok {
			panic("incorrect type of index in collection-elem expr")
		}
		res, err := src.GetElem(index)
		if err != nil {
			return nil
		}
		return res

	case string:
		index, ok := cc.KVal.(int64)
		if !ok {
			panic("incorrect type of index in collection-elem expr")
		}
		res := src[int(index)]
		return base.NewVal(res)

	case *objects.DictVal:
		res, err := src.GetElem(cc.KVal)
		if err != nil {
			return nil
		}
		return res
	}
	return nil
}

func (cc *ColElem) Set(val any) error {
	switch src := cc.Src.(type) {
	case *objects.ListVal:
		index, ok := cc.KVal.(int64)
		if !ok {
			return errors.New("incorrect type of index in collection-elem expr")
		}
		err := src.Set(index, val)
		if err != nil {
			return nil
		}
	case *objects.DictVal:
		src.Set(cc.KVal, val)
	}
	return nil
}

type ColElemExpr struct {
	Col    base.Expression // list, dict, tuple, string object
	Key    base.Expression // index or key
	ColRes *ColElem
}

func (cc *ColElemExpr) Get() *base.Val {
	return base.NewVal(cc.ColRes)
}

func (cc *ColElemExpr) Do(ctx base.Context) error {
	// cc.KVal = nil
	// cc.Src = nil
	err := cc.Col.Do(ctx)
	if err != nil {
		return err
	}
	err = cc.Key.Do(ctx)
	if err != nil {
		return err
	}
	kval := GetExprVal(cc.Key, ctx)
	// cc.KVal = kval
	sval := GetExprVal(cc.Col, ctx)
	// cc.Src = sval
	cc.ColRes = &ColElem{Src: sval, KVal: kval}
	return nil
}

// ----  Slice:  collection[ start : end]
type ColSlice struct {
	Col  base.Expression // list, dict, tuple, string object
	Inds *ColonPair      // index or key

	res any // ListVal, TupleVal, string; TODO: xBytes
}

func (cs *ColSlice) Get() *base.Val {
	if cs.res == nil {
		return nil
	}
	return base.NewVal(cs.res)
}

func (cs *ColSlice) Do(cx base.Context) error {
	cs.res = nil
	err := cs.Col.Do(cx)
	if err != nil {
		return errors.New("slice: bad collection expression")
	}
	// colV := cs.Col.Get()
	colv := GetExprVal(cs.Col, nil)
	if colv == nil {
		return errors.New("slice: bad collection expression")
	}
	err = cs.Inds.Do(cx)
	if err != nil {
		fmt.Printf("ColSlice#inds err: %v\n", err)
		return err
	}
	inds := cs.Inds.GetPair()
	var start int64
	var end int64
	i0 := inds[0]
	i1 := inds[1]
	switch i0 := i0.(type) {
	case int64:
		start = i0
	case *objects.EmptyVal:
		start = 0
	default:
		return errors.New("slice: bad index start")
	}
	switch i1 := i1.(type) {
	case int64:
		end = i1
	case *objects.EmptyVal:
		switch col := colv.(type) {
		case *objects.ListVal:
			end = int64(len(col.Elems))
		case *objects.TupleVal:
			end = int64(len(col.Elems))
		case string:
			end = int64(len(col))
		}
	default:
		return errors.New("slice: bad index end")
	}

	// fmt.Printf("ColSlice#src: %T, %v\n", colv, colv)
	switch col := colv.(type) {
	case *objects.ListVal:
		vals := col.Elems[int(start):int(end)]
		cs.res = objects.NewListVal(vals)
	case *objects.TupleVal:
		vals := col.Elems[int(start):int(end)]
		cs.res = objects.NewTupleVal(vals)
	case string:
		val := col[int(start):int(end)]
		cs.res = val
	}
	return nil
}

func NewSlice(col base.Expression, sub *OperColon) *ColSlice {
	inds := sub.GetPair()
	return &ColSlice{Col: col, Inds: inds}
}

//===== Tuple

type TupleExpr struct {
	Seq  *SequenceComma
	Subs []base.Expression
	res  *objects.TupleVal
}

// func (op *TupleExpr) IsParent() bool {
// 	// don't used as Block by default
// 	return false
// }

func (cs *TupleExpr) Add(sub base.Expression) {
	if cs.Subs == nil {
		cs.Subs = []base.Expression{}
	}
	cs.Subs = append(cs.Subs, sub)
}

func (cs *TupleExpr) Get() *base.Val {
	return base.NewVal(cs.res)
}

func (cs *TupleExpr) Do(ctx base.Context) error {
	src := make([]base.Expression, len(cs.Seq.Subs))
	copy(src, cs.Seq.Subs)
	src = append(src, cs.Subs...)
	res := make([]any, len(src))
	// fmt.Println("(,) Tuple Do:", len(src))
	for i, ex := range src {
		err := ex.Do(ctx)
		if err != nil {
			return err
		}
		res[i] = GetExprVal(ex, ctx)
	}
	cs.res = objects.NewTupleVal(res)
	return nil
}

//==== ColonPair

type ColonPair struct {
	Left  base.Expression
	Right base.Expression
	res   base.Pair
}

func (cp *ColonPair) Get() *base.Val {
	return base.NewVal(cp.res)
}

func (cp *ColonPair) GetPair() base.Pair {
	return cp.res
}
func (cp *ColonPair) Do(ctx base.Context) error {
	err1 := cp.Left.Do(ctx)
	if err1 != nil {
		return err1
	}
	err2 := cp.Right.Do(ctx)
	if err1 != nil {
		return err2
	}
	cp.res = base.Pair{GetExprVal(cp.Left, ctx), GetExprVal(cp.Right, ctx)}
	return nil
}

//==== {dict}

type DictExpr struct {
	Seq  *SequenceComma
	Subs []base.Expression
	res  *objects.DictVal
}

func (op *DictExpr) IsParent() bool {
	// don't used as Block by default
	return false
}

func (cs *DictExpr) Add(sub base.Expression) {
	if cs.Subs == nil {
		cs.Subs = []base.Expression{}
	}
	cs.Subs = append(cs.Subs, sub)
}

func (cs *DictExpr) Get() *base.Val {
	return base.NewVal(cs.res)
}

func (cs *DictExpr) Do(ctx base.Context) error {
	// TODO: think about constraints of key type
	src := make([]base.Expression, len(cs.Seq.Subs))
	copy(src, cs.Seq.Subs)
	src = append(src, cs.Subs...)
	res := make(map[any]any, len(src))
	for _, ex := range src {
		opcol, ok := ex.(*OperColon)
		if !ok {
			return errors.New("Incorest subelement of dict expression")
		}
		cpair := opcol.GetPair()
		// fmt.Printf("DictDo1 (%T, %v) \n", cpair, cpair)
		err := cpair.Do(ctx)
		if err != nil {
			return err
		}
		pairVal := cpair.GetPair()
		res[pairVal[0]] = pairVal[1]
	}
	cs.res = objects.NewDictVal(res)
	return nil
}

type BytesExpr struct {
	Pref string
	Subs base.Expression

	res objects.Bytes
}

func (bx *BytesExpr) Get() *base.Val {
	return base.NewVal(bx.res)
}

// solid line: hex: f0f0f0f0f0f0f0 | bin: 101010101010101010101
// src should be longer than 2
func ParseByteLine(src string, base int) ([]byte, error) {
	// 256 - 2 ** 63
	// t := n
	// for
	// rm := t % 0xff
	// t := 5 >> 1
	st := 8
	switch base {
	case 16:
		st = 2
	}
	rrs := []rune(src)
	slen := len(rrs)
	if slen <= st {
		n, err := strconv.ParseInt(src, base, 64)
		if err != nil {
			return nil, err
		}
		return []byte{byte(n)}, nil
	}
	count := slen / st
	rm := slen % st
	start := 0
	rsize := count
	if rm > 0 {
		rsize++
	}

	res := make([]byte, rsize)
	if rm > 0 {
		// count++
		start = 1
		n0, err := strconv.ParseInt(string(rrs[0:rm]), base, 64)
		// fmt.Println("ppBL1:", rrs[0:rm], string(rrs[0:rm]), "=>", n0)
		if err != nil {
			return nil, err
		}
		res[0] = byte(n0)
	}
	// fmt.Printf("BytesX.ParseBL src:%s k:%d count: %v \n", src, st, count)
	for i := 0; i < count; i += 1 {
		nb := rm + i*st // n part begin
		// fmt.Println("ppBL2:", i, rrs[nb:nb+st], string(rrs[nb:nb+st]))
		n, err := strconv.ParseInt(string(rrs[nb:nb+st]), base, 64)
		if err != nil {
			return nil, err
		}
		res[i+start] = byte(n)
	}
	return res, nil
	// return nil, nil
}

// mane case: 0x[01 ae ff] ; 0b[1111 0000, 1001] ; 0d[1 30 255]
func ParseBytes(src []string, base int) ([]byte, error) {
	res := make([]byte, len(src))
	// res := []byte{}
	for i, s := range src {
		n, err := strconv.ParseInt(s, base, 64)
		if err != nil {
			return nil, err
		}
		if n > 255 {
			return nil, fmt.Errorf("Out of byte range: %s => %x", s, n)
		}
		res[i] = byte(n)
	}
	return res, nil
}

func (bx *BytesExpr) Parse() error {
	// 0 1 2 7
	// 0 1 ae ff
	// 00 11 00 11
	// 12 11 03 144 255
	k := 2 // 0b
	switch bx.Pref[1] {
	case 'x':
		k = 16
	case 'o':
		k = 8
	case 'd':
		k = 10
	}
	switch sx := bx.Subs.(type) {
	// case *NumField:
	// 	fmt.Printf("BytesX.Parse NumField (%T, %v) %v \n", sx, sx, sx.V)
	// 	nums, err := ParseBytes(sx.V, k)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	bx.res = nums
	case *NumField:
		// fmt.Printf("BytesX.Parse NumField (%T, %v) %v \n", sx, sx, sx.V)
		// nums, err := ParseBytes(sx.V, k)
		// if err != nil {
		// 	return err
		// }
		// bx.res = nums
		nums := []byte{}
		for _, ss := range sx.V {
			bb, err := ParseByteLine(ss, k)
			if err != nil {
				return err
			}
			nums = append(nums, bb...)
		}
		bx.res = nums
	case *ByteLine:
		// fmt.Printf("BytesX.Parse ByteLine (%T, %v) %v \n", sx, sx, sx.Src)
		if len(sx.Src) == 0 {
			bx.res = objects.Bytes{}
			return nil
		}
		bb, err := ParseByteLine(sx.Src, k)
		if err != nil {
			return err
		}
		bx.res = bb
		// case *ValExpr:
		// fmt.Printf("BytesX.Parse Val (%T, %v) %v \n", sx, sx, sx.Val)
		// ParseByteLine(sx.Src, k)
		// case *VarExpr:
		// fmt.Printf("BytesX.Parse Var (%T, %v) %v \n", sx, sx, sx.GetName())
	}
	return nil
}
func (bx *BytesExpr) Do(cx base.Context) error {
	// bx.Subs.Do(cx)
	return nil
}

func NewBytesExpr(pref string, subs base.Expression) *BytesExpr {
	return &BytesExpr{Pref: pref, Subs: subs}
}

type ByteLine struct {
	Src string
}

func (t *ByteLine) Do(cx base.Context) error {
	return nil
}
func (t *ByteLine) Get() *base.Val {
	return nil
}

func NewByteLine(elems []*lang.Elem) *ByteLine {
	// nn:= make([]string)
	var val string
	if len(elems) == 1 {
		val = elems[0].Text
	}
	return &ByteLine{Src: val}
}
