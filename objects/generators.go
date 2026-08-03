package objects

import (
	"errors"

	"github.com/lesnikyan/lisapet-go/base"
)

// simple number generator [a .. b]
type NumSeqGen struct {
	Start  int64
	Max    int64
	Step   int64
	direct int64

	cur   int64
	index int64
}

func (ng *NumSeqGen) Finished() bool {
	// fmt.Printf("NGen#Fin? cu: %v, max: %v, ?: %v \n", ng.cur, ng.Max, ng.cur <= ng.Max)
	if ng.direct < 0 {
		return ng.cur < ng.Max
	}
	return ng.cur > ng.Max
}

func (ng *NumSeqGen) Init() {
	ng.cur = ng.Start
	ng.direct = 1
	if ng.Step < 0 {
		ng.direct = -1
	}
	ng.index = 0
}

func (ng *NumSeqGen) Next() (base.Pair, error) {
	// fmt.Printf("NGen#09, index: %v, cur: %v, step:%v\n", ng.index, ng.cur, ng.Step)
	if ng.Finished() {
		return base.Pair{}, errors.New("trying to Next of Finished iterator")
	}
	res := ng.cur
	i := ng.index
	ng.cur += ng.Step
	ng.index += 1
	return base.Pair{i, res}, nil
}

func (ng *NumSeqGen) GetList() *ListVal {
	ng.Init()
	diff := (ng.Max - ng.Start) * ng.direct
	count := int((diff)/ng.Step) + 1
	// fmt.Printf("NGen#10, start: %v, max: %v, step:%v, count: %v\n", ng.Start, ng.Max, ng.Step, count)
	rnums := make([]any, count)
	// id := 0
	v := ng.Start
	// for v := ng.Start; v <= ng.Max; v += ng.Step {
	for id := 0; id < count; id += 1 {
		// fmt.Printf("NGen#11, id: %v, v: %v\n", id, v)
		rnums[id] = v
		v += ng.Step
		// id += 1
	}
	res := NewListVal(rnums)
	return res
}

func NewNumSeqGen(start int64, end int64, step int64) *NumSeqGen {
	return &NumSeqGen{Start: start, Max: end, Step: step}
}
