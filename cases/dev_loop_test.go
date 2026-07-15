package cases

import (
	"fmt"
	"testing"

	obb "github.com/lesnikyan/lisapet-go/objects"
	par "github.com/lesnikyan/lisapet-go/parser"
	"github.com/stretchr/testify/assert"
)

/*
TODO:
ok 3. for i=0; i<5; i +=1
ok 1. create list: nn = [1,2,3]
ok 2.1 read list v = nn[0]
ok 2.2 set elem: nn[0] = 1
5. create tuple: (1,2,3)
5.1 tuple elem t[i]
6. create dict: d1 = {'a':1, 'b': 2}
6.1 dict elem: d1[k]
4. LeftArrow:
4.1 loop by list: for n <- nn
4.2 append to list: nn <- v
6.2 loop by dict for k, v <- d1
6.3 append to dict: d1 <- (k, v)
*/
func TestForCountCase(t *testing.T) {
	tdata := []struct {
		src string
		res any
	}{
		{`
		a = 0
		for i = 1; i < 6; i = i + 1
			a = a + i
		`, int64(15)},
		{`
		a = 0
		for i=0; i < 10; i += 1
			a += i
		`, int64(45)},
		{`
		a = 0
		for i=0; i < 10; i+= 1
			if i % 2 != 0
				a += i
		`, int64(25)},
		{`
		a = 0
		for i=1; i < 10; i += 1
			for j=1; j < 10; j += 1
				for k=1; k < 10; k += 1
					a += i + j + k
		`, int64(10935)},
		{`
		a = 0
		for i = 0; i < 10 ; i += 1
			for j=1; j < 10; j += 1
				if i % 2 > 0
					if j % 2 >0
						a += i + j
		`, int64(250)},
		// {``, int64(1)},
		// {``, int64(1)},
		// {``, int64(1)},
		// {``, int64(1)},
	}
	for _, tt := range tdata {
		t.Run(fmt.Sprintf("ExprDo, %s >>", tt.src), func(t2 *testing.T) {
			clines := par.SplitCode(tt.src[1:])
			block, err := TreeBlock(clines)
			assert.Nil(t, err)
			// t.Log("--- --- --- Do ...")
			ctx := obb.NewContext(nil)
			block.Do(ctx)
			vr := ctx.GetVar("a")
			// t.Log("tt#vr", vr)

			switch rr := vr.Val.(type) {
			case *obb.ListVal:
				// fmt.Printf("tt#ListVal#1  (%T, %v)  (%T, %v) len: %d \n", vr, vr, col, col, len(col.Elems))
				assert.Equal(t2, tt.res, rr.Elems)
			default:
				assert.Equal(t2, tt.res, vr.Val)
				// fmt.Println("tt3>", vr, vr.Name, vr.Val)
			}
		})
	}

	// a := 0
	// for i := 1; i < 10; i++ {
	// 	for j := 1; j < 10; j++ {
	// 		// for k := 1; k < 10; k++ {
	// 		// a += i + j + k
	// 		if i%2 > 0 && j%2 > 0 {
	// 			a += i + j
	// 		}
	// 	}
	// }
	// fmt.Println("i+j+k", a)
}
