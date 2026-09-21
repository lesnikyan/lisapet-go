package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/lesnikyan/lisapet-go/cases"
	"github.com/lesnikyan/lisapet-go/objects"
	"github.com/lesnikyan/lisapet-go/parser"
	"github.com/stretchr/testify/assert"
)

type TTst = cases.TTst

func TestSpeedLoop1Million01(t *testing.T) {
	n := 1000_000
	tt := cases.TTst2{
		fmt.Sprintf(`
		r = len([0; n <- [1 .. %d]])
		`, n), "r", int64(1)}
	t1 := time.Now()
	clines := parser.SplitCode(tt.Src[1:])
	block, err := cases.TreeBlock(clines)
	assert.Nil(t, err)
	// t.Log("--- --- --- Do ...")
	ctx := objects.NewContext(nil)
	cases.PreloadContext(ctx)
	t2 := time.Now()
	terr := block.Do(ctx)
	t3 := time.Now()
	if terr != nil {
		assert.Fail(t, terr.Error())
		return
	}
	dt1 := t2.Sub(t1).Seconds()
	dt2 := t3.Sub(t2).Seconds()
	fmt.Printf("Speed0 LP-code       parse and load: %.06f sec \n", dt1)
	fmt.Printf("Speed1 by %d-iters loop,LP run: %.06f sec \n", n, dt2)

	// Go test by append
	t11 := time.Now()
	r := 0
	aa := []int{}
	for i := range n {
		aa = append(aa, i)
	}
	r = len(aa)
	t12 := time.Now()

	// dt11 := t12.Sub(t1).Seconds()
	dt12 := t12.Sub(t11).Seconds()
	fmt.Printf("Speed2 by %d-iters go-loop run: %.06f sec \n", r, dt12)

	// Go test by index
	t13 := time.Now()
	r = 0
	aa = make([]int, n)
	for i := range n {
		aa[i] = 0
	}
	r = len(aa)
	t14 := time.Now()

	// dt11 := t12.Sub(t1).Seconds()
	dt14 := t14.Sub(t13).Seconds()
	fmt.Printf("Speed3 by %d-iters go-make run: %.06f sec \n", r, dt14)

	// results 2026-09-19
	/*
		python console:
		>>> t1 = time.time(); x=[0 for i in range(1000000)]; t2 = time.time(); print("%.06f" % (t2 - t1))
		-----------------------------------  0.063960
		this test:
		Speed0 LP-code       parse and load: 0.000000 sec
		Speed1 by 1000000-iters loop,LP run: 0.185662 sec
		Speed2 by 1000000-iters go-loop run: 0.006323 sec
		Speed3 by 1000000-iters go-make run: 0.001123 sec

		python: 0.06  : x3 of LP
		LP    : 0.2   : we here
		go-app: 0.006 : x30 of LP
		go-ind: 0.001 : x200 %-)
	*/
}
