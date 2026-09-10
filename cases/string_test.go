package cases

import (
	"strings"
	"testing"

	obb "github.com/lesnikyan/lisapet-go/objects"
)

/*
TODO:
ok 0x[], 0b[] ,0o[]
ok 0b[11110000 01010101]; 0x[ff00 1234 0000 0001 1000]
ok prefixless byteset [aa bb ff]
bytes actions:
nbytes[elem]
nbytes[sli : ce]
nbytes <- add
for bt <- nbytes

string(int|glif|[glif...]|bool|float)
string.glifs()
glif(int|string|0x[])
bytes([]int, string, glif, []glif)
*/

// Escaping bytes, multibytes (\377 \xff, \u00ff): thinking

// Escape sequences: \n \t \' \" \` \\

// """ .. \n .. \n """.lines()

func TestStringEscSeq(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = '1\n2'
		`, "r", "1\n2"},
		{`
		r = "2\n3\t4\'\"\\"
		`, "r", "2\n3\t4'\"\\"},
		{`
		r = '3 \' \" \$$ \\ \/ /'
		`, "r", "3 ' \" ` \\ / /"},
		{`
		r = $$\ \n \'$$
		`, "r", `\ \n \'`},
		{`
		r = """<node>
		Hello 2-quote \"\"\"
		AAA 1
		</node>"""
		`, "r", "<node>\nHello 2-quote \"\"\"\nAAA 1\n</node>"},
		{`
		r = """<node>
		Hello 2-quote \"""
		AAA 2
		</node>"""
		`, "r", "<node>\nHello 2-quote \"\"\"\nAAA 2\n</node>"},
		{`
		r = """<node>
		Hello 2-quote "\""
		</node> 3"""
		`, "r", "<node>\nHello 2-quote \"\"\"\n</node> 3"},
		{`
		r = """<node>
		Hello 2-quote ""\"
		</node> 4"""
		`, "r", "<node>\nHello 2-quote \"\"\"\n</node> 4"},
		{`
		r = '\''
		`, "r", "'"},
		{`
		r = ''' $$ $$$$ $%$ ' '' " "" """ '''
		`, "r", " ` `` ``` ' '' \" \"\" \"\"\" "},
		{`
		r = "\\ \\\\ \\\""
		`, "r", "\\ \\\\ \\\""},
		{`
		r = $%$ 11 $%$
		`, "r", " 11 "},
		{`
		r = $%$ "11" $%$
		`, "r", " \"11\" "},
		{`
		r = $%$ """\ $%$
		`, "r", " \"\"\"\\ "},
		{`
		r = $%$ $$ \$$ ''' """ \' \" \n \t \s $%$
		`, "r", " ` \\` ''' \"\"\" \\' \\\" \\n \\t \\s "},
		// {``, "r",  ""},
		// {``, "r",  ""},
		// {``, "r",  ""},
		// {``, "r",  ""},
		// {``, "r",  ""},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		tt.src = strings.ReplaceAll(tt.src, "$%$", "```")
		tt.src = strings.ReplaceAll(tt.src, "$$", "`")
		RunTCodeVarExp(t, i, tt)
	}
}

// Multiline strings
func TestMultistring(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r= '''
		One
		X'''
		`, "r", "\nOne\nX"},
		{`
		s = """
		Two
		XX
		YYY
		132
		"""
		r = s
		`, "r", "\nTwo\nXX\nYYY\n132\n"},
		{`
		s = """
			Three
				SSS
			"""
		r = s
		`, "r", "\n\tThree\n\t\tSSS\n\t"},
		{`
		r = [
		"""
		AAA BBB
		CCC DDD
		""",
		'''
		EEE FFF
		GGG HHH 1
		''']
		`, "r", Anis("\nAAA BBB\nCCC DDD\n", "\nEEE FFF\nGGG HHH 1\n")},
		{`
		r = []
			"""
		AAA BBB
		CCC DDD
		"""
			'''
		EEE FFF
		GGG HHH 2
		'''
		`, "r", Anis("\nAAA BBB\nCCC DDD\n", "\nEEE FFF\nGGG HHH 2\n")},
		{`
		r = []
			"""
			AAA BBB
			CCC DDD
			"""
			'''
			EEE FFF
			GGG HHH 3
			'''
		`, "r", Anis("\n\tAAA BBB\n\tCCC DDD\n\t", "\n\tEEE FFF\n\tGGG HHH 3\n\t")},
		{`
		r = []
			"""
		  AAA BBB
		    CCC DDD 4
			"""
		`, "r", Anis("\n  AAA BBB\n    CCC DDD 4\n\t")},
		{
			// it looks like perversion, yep
			strings.ReplaceAll(`
		r = []
			"""
		  AAA BBB
		    CCC DDD 5
			"""
		`, `"""`, "```"), "r", Anis("\n  AAA BBB\n    CCC DDD 5\n\t")},
		{strings.ReplaceAll(`
		r = [""" One """, ''' Two ''', !!! N-3 !!!]
		`, `!!!`, "```"), "r", Anis(" One ", " Two ", " N-3 ")},
		{`
		r = """ He110 """.replace('1', 'L').replace('0','o')
		`, "r", " HeLLo "},
		{`
		s = '''
		123
		456
		'''
		r = """ 
		ABC
		DDD
		07 """ + s
		`, "r", " \nABC\nDDD\n07 \n123\n456\n"},
		// {``, "r",  ""},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestStrigMethods(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// split
		{`
		r = []
		s = "aa bb cc"
		r <- s.split(' ')
		r <- "q-w-e-r-t-y".split(g'-')
		r <- 'lorem,iprum,dolor'.split(glif(','))
		`, "r", Anis(Anis("aa", "bb", "cc"),
			Anis("q", "w", "e", "r", "t", "y"), Anis("lorem", "iprum", "dolor"))},
		// join
		{`
		ss = ['uuu','roo','roo']
		r = '_'.join(ss)
		`, "r", "uuu_roo_roo"},
		{`
		ss = ('uuu','roo','boom')
		r = '_'.join(ss)
		`, "r", "uuu_roo_boom"},
		{`
		gg = [g'A', g'r', g'a', g'm', g'-', g'C', g'h', g'u', ]
		r = '/'.join(gg)
		`, "r", "A/r/a/m/-/C/h/u"},
		// bytes
		{`
		r  = "a b c def.ABCD,xz".bytes()
		`, "r", obb.Bytes{0x61, 0x20, 0x62, 0x20, 0x63, 0x20, 0x64, 0x65, 0x66, 0x2e, 0x41, 0x42, 0x43, 0x44, 0x2c, 0x78, 0x7a}},
		//glifs
		{`
		r = "ABC abc αβγ 零 一 二 九".glifs()
		`, "r", Anis('A', 'B', 'C', ' ', 'a', 'b', 'c', ' ', 'α', 'β', 'γ', ' ', '零', ' ', '一', ' ', '二', ' ', '九')},
		// replace
		{`
		s = 'a b c d'
		r = s.replace(' ', '@')
		`, "r", "a@b@c@d"},
		{`
		s = 'Some man woke up. Man has breakfast. Man went out and met another man.'
		r = s.replace('man', 'cat').replace('Man','Cat')
		`, "r", "Some cat woke up. Cat has breakfast. Cat went out and met another cat."},
		{`
		s = "This man shines over the city."
		r = s.replace(g'a', 'oo')
		`, "r", "This moon shines over the city."},
		{`
		s = "capy flappy saan"
		r = s.replace(g'a', g'o')
		`, "r", "copy floppy soon"},
		// has
		{`
		s = 'Hello Holland!'
		r = []
		r <- s.has('llo')
		r <- s.has('oll')
		r <- s.has('u')
		`, "r", Anis(true, true, false)},
		{`
		s = 'Logos'
		r = []
		r <- s.has(g'L')
		r <- s.has(g'i')
		r <- s.has(g's')
		`, "r", Anis(true, false, true)},
		//trim
		{`
		s = "     abc 132    .    "
		r = s.trim(' ')
		`, "r", "abc 132    ."},
		{`
		s = "     abc 134    .    "
		r = s.trim(g' ')
		`, "r", "abc 134    ."},
		{`
		s = "  =>   abc 134:,    .  ?,  "
		r = s.trim(' ,.?=><')
		`, "r", "abc 134:"},
		// TODO: lines for multiline strings
		// TODO: escape sequences in line
		// {`
		// s = "Lorem \n ipsum \n belor \n"
		// r = s.lines()
		// `, "r", ""},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestBytesActions(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		bb = [11 22 ff]
		r = bb[1]
		`, "r", byte(0x22)},
		{`
		bb = [11 22 ff]
		r = bb[-1]
		`, "r", byte(0xff)},
		{`
		bb = [0 1 2 3 4 5 6 7]
		r = bb[2:5]
		`, "r", obb.Bytes{2, 3, 4}},
		{`
		bb = [0 1 2 3 4 5 6 7]
		r = bb[2:-2]
		`, "r", obb.Bytes{2, 3, 4, 5}},
		{`
		bb = [0 1 2 3 4 5 6 7]
		r = bb[-5:-2]
		`, "r", obb.Bytes{3, 4, 5}},
		{`
		bb = [0 1 2 3 4 5 6 7]
		r = bb[:3]
		`, "r", obb.Bytes{0, 1, 2}},
		{`
		bb = [0 1 2 3 4 5 6 7]
		r = bb[5:]
		`, "r", obb.Bytes{5, 6, 7}},
		{`
		b1 = [0 1 2 3]
		b2 = [14 15]
		r = b1 + b2
		`, "r", obb.Bytes{0, 1, 2, 3, 0x14, 0x15}},
		{`
		r  = [0 1 2 a]
		b2 = [14 15]
		r += b2
		`, "r", obb.Bytes{0, 1, 2, 0xa, 0x14, 0x15}},
		{`
		bb = 0x[0 1]
		bb <- 5
		`, "bb", obb.Bytes{0, 1, 5}},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestBytesNoPref(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# bytes
		r = [11 22 ff]
		`, "r", obb.Bytes{0x11, 0x22, 0xff}},
		{`
		# bytes
		r = [11 22 33 44 55 66 77 88 99 aa bb cc dd ee ff]
		`, "r", obb.Bytes{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}},
		{`
		# bytes
		r = [12345 678910 abcdef a0b0c]
		`, "r", obb.Bytes{0x1, 0x23, 0x45, 0x67, 0x89, 0x10, 0xab, 0xcd, 0xef, 0xa, 0xb, 0xc}},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestBytes(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# bytes
		r = 0x[]
		`, "r", obb.Bytes{}},
		{`
		r = 0b[]
		`, "r", obb.Bytes{}},
		{`
		r = 0o[]
		`, "r", obb.Bytes{}},
		{`
		r = 0d[]
		`, "r", obb.Bytes{}},
		{`
		r = 0x[1 2 f]
		`, "r", obb.Bytes{1, 2, 0xf}},
		{`
		r = 0x[F 2 f]
		`, "r", obb.Bytes{0xf, 2, 0xf}},
		{`
		r = 0x[1 2 3 4 5 6 7 8 9 a b c d e f 0 10 11 12 13 14 15 16 f0 f5 ff]
		`, "r", obb.Bytes{1, 2, 3, 4, 5, 6, 7, 8, 9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf,
			0x0, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0xf0, 0xf5, 0xff}},
		{`
		r:bytes = 0x[f b]
		`, "r", obb.Bytes{0xf, 0xb}},
		{`
		r = 0x[f]
		`, "r", obb.Bytes{0xf}},
		{`
		r = 0x[1]
		`, "r", obb.Bytes{1}},
		{`
		# bytes
		r = 0b[1111]
		`, "r", obb.Bytes{0xf}},
		{`
		# bytes
		r = 0x[1111]
		`, "r", obb.Bytes{0x11, 0x11}},
		{`
		# bytes
		r = 0x[1111f]
		`, "r", obb.Bytes{1, 0x11, 0x1f}},
		{`
		# bytes
		r = 0x[f1122334455667788990000af]
		`, "r", obb.Bytes{0xf, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x0, 0x0, 0xaf}},
		{`
		# bytes
		r = 0x[0a1122334455667788990000af]
		`, "r", obb.Bytes{0xa, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x0, 0x0, 0xaf}},
		{`
		# bytes
		r = 0b[101111000011110000111100001111000011110000]
		`, "r", obb.Bytes{0x2, 0xf0, 0xf0, 0xf0, 0xf0, 0xf0}},
		{`
		# bytes
		r = 0b[10 11110000 11110000 11110000 11110000 11110000]
		`, "r", obb.Bytes{0x2, 0xf0, 0xf0, 0xf0, 0xf0, 0xf0}},

		{`
		# bytes
		r = 0o[7]
		`, "r", obb.Bytes{0x7}},
		{`
		# bytes
		r = 0o[1 2 11 22 44 55 66 77 100 200 300 377]
		`, "r", obb.Bytes{0x1, 0x2, 0x9, 0x12, 0x24, 0x2d, 0x36, 0x3f, 0x40, 0x80, 0xc0, 0xff}},
		{`
		# bytes
		r = 0d[100]
		`, "r", obb.Bytes{100}},
		{`
		# bytes
		r = 0d[1 2 9 10 50 100 199 200 255]
		`, "r", obb.Bytes{0x1, 0x2, 0x9, 0xa, 0x32, 0x64, 0xc7, 0xc8, 0xff}},
		{`
		# bytes
		r = 0x[1111 2222 fab]
		`, "r", obb.Bytes{0x11, 0x11, 0x22, 0x22, 0xf, 0xab}},
		{`
		# bytes
		r = 0x[1111 2222 ff0 f1a 0fff]
		`, "r", obb.Bytes{0x11, 0x11, 0x22, 0x22, 0xf, 0xf0, 0xf, 0x1a, 0xf, 0xff}},
		{`
		# bytes
		r = 0b[11110000 01010101]
		`, "r", obb.Bytes{0xf0, 0x55}},
		{`
		# bytes
		r = 0b[11110000 1 10 1000]
		`, "r", obb.Bytes{0xf0, 0x1, 0x2, 0x8}},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestGlifs(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# glif
		s = g"G"
		`, "s", Gf("G")},
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
		`, "s", Gf("G")},
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
