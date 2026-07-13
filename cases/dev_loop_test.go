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
3. for i=0; i<5; i +=1
1. create list: nn = [1,2,3]
2.1 read list v = nn[0]
2.2 set elem: nn[0] = 1
4. LeftArrow:
4.1 loop by list: for n <- nn
4.2 append to list: nn <- v
5. create tuple: (1,2,3)
6. create dict: dict1 = {'a':1, 'b': 2}
6.1 dict elem: dict1[k]
6.2 loop by dict for k, v <- dict1
6.3 append to dict: dict1 <- (k, v)
*/
func TestForCase(t *testing.T) {
	tdata := []struct {
		src string
		res any
	}{
		{`
		a = 0
		for i = 1; i < 6; i = i + 1
			a = a + i
		`, int64(15)},
		// {``, int64(1)},
		// {``, int64(1)},
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
			assert.Equal(t2, tt.res, vr.Val)
			// fmt.Println("tt3>", vr, vr.Name, vr.Val)
		})
	}
}
