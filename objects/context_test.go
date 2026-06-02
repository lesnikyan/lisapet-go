package objects

import (
	"fmt"
	"testing"

	Lt "github.com/lesnikyan/lisapet-go/lang/lt"
	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	tdata := []struct {
		c    rune
		prev Lt.Lt
		exp  Lt.Lt
	}{}

	for _, tt := range tdata {
		t.Run(fmt.Sprintf("e-type '%s', <%s> ", string(tt.c), Lt.TName(tt.exp)), func(t2 *testing.T) {
			res := elemType(tt.c, tt.prev)
			assert.Equal(t2, tt.exp, res, "")
		})
	}
}

func TestNewContext(t *testing.T) {

	tvals := struct {
		item int32
		dtype dt.DType
		val interface{}
	}{
		{0, dt.Int, 123},
		{, dt.Float, 11.5}, 
		{2, dt.Int, 25},
		{3, dt.String, "hello1"},
		{4, dt.String, "yello2"},
	}
	tdata := []struct {
		name string
		val CVal
		}{}
	parent := NewContext(nil)
	for 

	for _, tt := range tdata {
		t.Run(fmt.Sprintf("e-type '%s', <%s> ", string(tt.c), Lt.TName(tt.exp)), func(t2 *testing.T) {
			res := elemType(tt.c, tt.prev)
			assert.Equal(t2, tt.exp, res, "")
		})
	}
}
