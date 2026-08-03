package cases

import (
	"fmt"
	"testing"
)

func TestStringBuiltins(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# replace
		s = 'marapita'
		r = replace(s, 'a', 'e')
		`, "r", "merepite"},
		{`
		# replace
		s = 'a-b-c- - qwerty'
		r = replace(s, '-', '#%')
		`, "r", "a#%b#%c#% #% qwerty"},
		{`
		# join
		ss = ['Hello', 'there', 'all']
		r = join(ss, ' ')
		`, "r", "Hello there all"},
		{`
		# split
		s = "lorem ipsum dolor"
		r = split(s, ' ')
		`, "r", Anis("lorem", "ipsum", "dolor")},
		// {``, "r",  Anis(11, )},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
	var aa []any = []any{1, 2, 3}
	fmt.Println(aa)
}
