package nodes

import (
	"fmt"

	"github.com/lesnikyan/lisapet-go/base"
	"github.com/lesnikyan/lisapet-go/objects"
)

type GenLoop struct {
	IterExp *LeftArrow
	Loop    *ForSourceNode
	Subs    []base.Expression
	Guard   base.Expression
	SubLoop *GenLoop

	rex base.Expression
}

func (gn *GenLoop) Get() *base.Val {
	return nil
}

func (gn *GenLoop) DoSub(cx base.Context, post base.Expression) error {
	snd := gn.Loop
	err := snd.Init(cx)
	if err != nil {
		return err
	}
	for {
		if snd.assign.Finished() {
			break // correct finich
		}
		err = snd.assign.Next()
		if err != nil {
			return err
		}
		err = snd.Block.Do(cx)
		// fmt.Printf("GLoop.DoSub#err2: %v\n", err)
		if err != nil {
			return err
		}
		if gn.Guard != nil {
			gn.Guard.Do(cx)
			gres := GetExprVal(gn.Guard, nil)
			gv, ok := gres.(bool)
			if !ok {
				return fmt.Errorf("generators guard should return bool but received : %T", gres)
			}
			if !gv {
				// skip result iteration by condition
				continue
			}
		}
		if gn.SubLoop == nil {
			err = post.Do(cx)
			if err != nil {
				return err
			}
		} else {
			err = gn.SubLoop.DoSub(cx, post)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gn *GenLoop) Do(cx base.Context) error {
	post := gn.rex
	err := gn.DoSub(cx, post)
	return err
}

func NewGenLoop(arrEx *LeftArrow, subs []base.Expression, guard base.Expression) *GenLoop {
	gn := &GenLoop{IterExp: arrEx, Guard: guard}
	gn.Loop = NewForSource(gn.IterExp)
	for _, sx := range subs {
		gn.Loop.Add(sx)
	}
	return gn
}

// ===
type GenResExpr struct {
	Exp base.Expression
	res []any
}

func (gg *GenResExpr) Get() *base.Val {
	return base.NewVal(gg.res)
}

func (gg *GenResExpr) Do(cx base.Context) error {
	if gg.res == nil {
		gg.res = []any{}
	}
	// fmt.Printf("?GenResMap.Do r = (%T)\n", gg.Exp)
	err := gg.Exp.Do(cx)
	if err != nil {
		// fmt.Printf("!GenResMap.Do err  ( %v)\n", err)
		return err
	}
	rval := GetExprVal(gg.Exp, nil)
	// fmt.Printf("-==GenResMap.Do r = (%T), (%T, %v)\n", gg.Exp, rval, rval)
	gg.res = append(gg.res, rval)
	return nil
}

//===

type ListComprh struct {
	ResEx    base.Expression
	LoopNode *GenLoop

	res *objects.ListVal
}

func (cm *ListComprh) Get() *base.Val {
	if cm.res == nil {
		return nil
	}
	return base.NewVal(cm.res)
}

func (cm *ListComprh) Do(cx base.Context) error {
	// fmt.Printf("ListCompr.Do\n")
	inCx := cx.SubContext()
	rr := []any{}
	post := &GenResExpr{Exp: cm.ResEx, res: rr}
	cm.LoopNode.rex = post
	err := cm.LoopNode.Do(inCx)
	if err != nil {
		// fmt.Printf("!ListComprh.Do err  ( %v)\n", err)
		return err
	}
	// Run ResEx in the each iter of last sub loop
	rr = post.res
	cm.res = objects.NewListVal(rr)
	return nil
}

func NewListComprh(res base.Expression, loop *GenLoop) *ListComprh {
	return &ListComprh{ResEx: res, LoopNode: loop}
}

// ===
type GenResMap struct {
	Exp base.Expression
	res map[any]any
}

func (gg *GenResMap) Get() *base.Val {
	return base.NewVal(gg.res)
}

func (gg *GenResMap) Do(cx base.Context) error {
	if gg.res == nil {
		gg.res = map[any]any{}
	}
	err := gg.Exp.Do(cx)
	if err != nil {
		return err
	}
	rval := GetExprVal(gg.Exp, nil)
	// fmt.Printf("GenResMap.Do r = (%T), %T\n", gg.Exp, rval)
	tt, ok := rval.(*ColonPair)
	if !ok {
		return fmt.Errorf("dict generator expects colon-separated pair, but %T returned", rval)
	}
	err = tt.Do(cx)
	if err != nil {
		return err
	}
	// fmt.Println("TT:", tt.GetPair())
	ttpair := tt.GetPair()
	key := ttpair[0]
	val := ttpair[1]
	gg.res[key] = val
	return nil
}

//===

type DictComprh struct {
	ResEx    base.Expression
	LoopNode *GenLoop

	res *objects.DictVal
}

func (cm *DictComprh) Get() *base.Val {
	if cm.res == nil {
		return nil
	}
	return base.NewVal(cm.res)
}

func (cm *DictComprh) Do(cx base.Context) error {
	// fmt.Printf("DictCompr.Do\n")
	inCx := cx.SubContext()
	rr := map[any]any{}
	post := &GenResMap{Exp: cm.ResEx, res: rr}
	cm.LoopNode.rex = post
	cm.LoopNode.Do(inCx)
	// Run ResEx in the each iter of last sub loop
	rr = post.res
	cm.res = objects.NewDictVal(rr)
	return nil
}

func NewDictComprh(res base.Expression, loop *GenLoop) *DictComprh {
	return &DictComprh{ResEx: res, LoopNode: loop}
}

// ===
var genFinish = NewGenFinish()

// ===

type YieldLoop struct {
	Loop    *ForSourceNode
	SubLoop *YieldLoop
	Subs    []base.Expression
	Guard   base.Expression
	Started bool
	id      int

	// rex base.Expression // result expression
}

func (sg *YieldLoop) Init(cx base.Context) error {
	// fmt.Printf("# YieldLoop.Init [%d] Cx: (%T, %v) loop: (%T, %v)\n", sg.id, cx, cx, sg.Loop.assign, sg.Loop.assign)
	// fmt.Printf("# --- YieldLoop.Init [%d] \n", sg.id)
	err := sg.Loop.Init(cx)
	if err != nil {
		return err
	}
	if sg.SubLoop != nil {
		return sg.SubLoop.Init(cx)
	}
	return nil
}

func (sg *YieldLoop) DoGuard(cx base.Context) (bool, error) {
	if sg.Guard == nil {
		return true, nil
	}
	err := sg.Guard.Do(cx) // err, res, bool
	if err != nil {
		return false, err
	}
	ganval := GetExprVal(sg.Guard, nil)
	gval, ok := ganval.(bool)
	if !ok {
		return false, fmt.Errorf("generator guard maust returm boo, %t received", ganval)
	}
	return gval, nil
}

func (sg *YieldLoop) Finished() bool {
	cur := sg.Loop.assign.Finished()
	if sg.SubLoop != nil {
		return cur && sg.SubLoop.Finished()
	}
	return cur
}

func (sg *YieldLoop) DoSubs(cx base.Context) error {
	for _, ex := range sg.Subs {
		err := ex.Do(cx)
		if err != nil {
			return err
		}
	}
	return nil
}

// `Next` method for last loop in chain
func (sg *YieldLoop) NextHere(cx base.Context) error {
	// fmt.Printf("YieldLoop.NextHere1[%d]\n", sg.id)
	for {
		if sg.Finished() {
			// fmt.Printf("YieldLoop.NextHere2[%d] Fin\n", sg.id)
			return genFinish
			// if sg.SubLoop == nil {
			// 	return genFinish
			// }
		}
		var err error
		err = sg.Loop.assign.Next()
		if err != nil {
			return err
		}
		err = sg.DoSubs(cx)
		if err != nil {
			return err
		}
		gval, err := sg.DoGuard(cx)
		if err != nil {
			return err
		}
		if gval {
			return nil
		}
	}
}

func (sg *YieldLoop) SubInit(cx base.Context) error {
	// fmt.Printf("# --- YieldLoop.SubInit [%d] \n", sg.id)
	err := sg.Loop.Init(cx)
	if err != nil {
		return err
	}
	// err = sg.DoSubs(cx)
	err = sg.NextHere(cx)
	if err != nil {
		return err
	}
	if sg.SubLoop != nil {
		// sg.SubLoop.Init(cx)
		err = sg.SubLoop.SubInit(cx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sg *YieldLoop) NextSub(cx base.Context) error {
	if !sg.Started {
		err := sg.NextHere(cx)
		if err != nil {
			return err
		}
		sg.Started = true
	}
	i := -1
	for {
		i++
		// fmt.Printf("YieldLoop.NextSub1 %d , %d\n", sg.id, i)
		if sg.SubLoop.Finished() {
			err := sg.NextHere(cx)
			if err != nil {
				return err
			}
			// fmt.Printf("YieldLoop.NextSub2 %d SubLoop.Init, %d\n", sg.id, i)
			err = sg.SubLoop.SubInit(cx)
			if err != nil {
				switch err.(type) {
				case *GenFinish:
					continue
				default:
					return err
				}
			}
			return nil
			// err = sg.SubLoop.DoSubs(cx)
			// if err != nil {
			// 	return err
			// }
		} else {
			err := sg.SubLoop.Next(cx)
			if err == nil {
				// success iter
				return nil
			}
			switch err.(type) {
			case *GenFinish:
				continue
			default:
				return err
			}
		}
		// fmt.Printf("YieldLoop.NextSub3 %d , %d\n", sg.id, i)
	}
}

func (sg *YieldLoop) Next(cx base.Context) error {
	if sg.SubLoop == nil {
		return sg.NextHere(cx)
	}
	return sg.NextSub(cx)
}

var yid = 0

func NewYieldLoop(arrEx *LeftArrow, subs []base.Expression, guard base.Expression) *YieldLoop {
	loop := NewForSource(arrEx)
	// for _, sx := range subs {
	// 	loop.Add(sx)
	// }
	// fmt.Printf("# NewYieldLoop  LArr: (%T, %v) loop: (%T, %v) guard: (%T, %v)\n", arrEx, arrEx, loop, loop, guard, guard)
	yid++
	return &YieldLoop{Loop: loop, Subs: subs, Guard: guard, id: yid}
	// gn := &GenLoop{IterExp: arrEx, Guard: guard}
	// 	return gn
}

// ===

// Sequence generator (: x ; x <- nn)
type SeqGenerator struct {
	MainLoop *YieldLoop // ?
	Rex      base.Expression
	Ctx      base.Context

	nRes    any
	counter int
}

func (sg *SeqGenerator) SetContext(cx base.Context) {
	sg.Ctx = cx.SubContext()
}

func (sg *SeqGenerator) Init() {
	sg.MainLoop.Init(sg.Ctx)
	sg.counter = 0
}

func (sg *SeqGenerator) Finished() bool {
	return sg.MainLoop.Finished()
}

func (sg *SeqGenerator) Next() ([]any, error) {
	cx := sg.Ctx
	err := sg.MainLoop.Next(cx)
	if err != nil {
		switch err.(type) {
		case *GenFinish:
			// TODO: handle unexpected finish of main loop
			// panic("TODO: handle unexpected finish in gen")
			return nil, err
		default:
			return nil, err
		}
	}
	err = sg.Rex.Do(cx)
	if err != nil {
		return nil, err
	}
	rval := GetExprVal(sg.Rex, nil)
	id := sg.counter
	sg.counter++
	// fmt.Printf("Rval: %v \n", rval)
	return []any{id, rval}, nil
}

// ===

type SeqGenExpr struct {
	Loop   *YieldLoop
	ResExp base.Expression

	res *SeqGenerator
}

func (gx *SeqGenExpr) Get() *base.Val {
	if gx.res == nil {
		return nil
	}
	return base.NewVal(gx.res)
}

func (gx *SeqGenExpr) Do(cx base.Context) error {
	// Make SeqGenerator

	gen := &SeqGenerator{Rex: gx.ResExp, MainLoop: gx.Loop}
	gen.SetContext(cx)
	// gen.Init()
	gx.res = gen
	return nil
}

func NewSeqGenExpr(rex base.Expression, loop *YieldLoop) *SeqGenExpr {
	return &SeqGenExpr{ResExp: rex, Loop: loop}
}
