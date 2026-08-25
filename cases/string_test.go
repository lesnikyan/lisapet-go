package cases

import (
	"testing"
)

/*
TODO:
0x[], 0b[] ,0o[]
string(int|glif|[glif...]|bool|float)
string.glifs()
glif(int|string|0x[])
*/

func TestGlifs(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# glif
		s = g"G"
		`, "s", "G"},
		{`
		# glifs in list
		r = [g'A', g'Z', g'@', g'Ы', g'ф', g'百']
		`, "r", Anis('A', 'Z', '@', Gf("Ы"), 'ф', '百')},
		{`
		# glifs in dict
		r = {g'@': g'Ы', g'ф': g'百'}
		`, "r", adk(dk{'@': 'Ы', 'ф': '百'})},
		{`
		# glifs in list, block-syntax
		s = g"G"
		`, "s", "G"},
		{`
		r = []
			g'A'
			g'Z'
			g'@'
			g'Ы'
			g'Ф'
			g'百'
		`, "r", Anis('A', 'Z', '@', Gf("Ы"), 'Ф', '百')},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

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
	// var aa []any = []any{1, 2, 3}
	// fmt.Println(aa)
}
