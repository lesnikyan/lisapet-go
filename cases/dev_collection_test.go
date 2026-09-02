package cases

import (
	"testing"

	obb "github.com/lesnikyan/lisapet-go/objects"
)

/*
TODO:
ok 1. slice: nn[a : b]
ok 5. Block in block: list:dict,tuple; dict: dict,list...
constructors:
ok 2. list()
ok 3. dict()
ok 4. tuple()
Methods:
list
tuple
dict
bytes: min, max,

*/

func TestTupleMethods(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// join
		{`
		tt = ('ab','CD', '45')
		r = tt.join('-')
		`, "r", "ab-CD-45"},
		{`
		tt = tuple("11 23 456".split(' '))
		r = tt.join(g'-')
		`, "r", "11-23-456"},
		// map
		{`
		func foo(x)
			x * 11
		#
		tt = (1,2,3,4,5)
		r = tt.map(foo)
		`, "r", Tanis(11, 22, 33, 44, 55)},
		{`
		func foo(x)
			glif(65 + x)
		#
		tt = (1,2,3,4,37)
		r = tt.map(foo)
		`, "r", Tanis('B', 'C', 'D', 'E', 'f')},
		// {``, "r",  Anis( )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestListMethods(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// join
		{`
		nn = ['Hello','honey','hill']
		r = nn.join(' | ')
		`, "r", "Hello | honey | hill"},
		{`
		nn = ['Hello','honey','hill']
		r = nn.join(g'|')
		`, "r", "Hello|honey|hill"},
		// map
		{`
		func x10(x)
			x * 10
		#
		ss = [1,2,3,4,5]
		r = ss.map(x10)
		`, "r", Anis(10, 20, 30, 40, 50)},
		{`
		func xcase(x:glif)
			n = int(x)
			if n < 65
				return x
			if n < 91
				return glif(n + 32)
			if n < 97
				return x
			if n < 123
				return glif(n - 32)
			x
		#
		ss = "Abs Hello (SQL)."
		r = ss.glifs().map(xcase).join('')
		`, "r", "aBS hELLO (sql)."},
		{`
		func rangeN(a, b)
			r = []
			for n <- [a .. b]
				r <- n
			r
		#
		r = rangeN(65, 75).map(glif).join('')
		`, "r", "ABCDEFGHIJK"},
		{`
		# change outer var from map(func)
		r1 = []
		func foo(x)
			r1 <- x
			(x, 100+x)
		#
		nn = [1,2,3,4]
		r2 = nn.map(foo)
		r = [r1, r2]
		# >> [[1,2,3,4], [(1,101),(2,102),(3,103),(4,104)]]
		`, "r", Anis(Anis(1, 2, 3, 4), Anis(Tanis(1, 101), Tanis(2, 102), Tanis(3, 103), Tanis(4, 104)))},
		// fold
		{`
		func sum(s, x)
			s + x
		#
		nn = [1,2,3,4,5]
		r = nn.fold(0, sum)
		`, "r", int64(15)},
		{`
		func foo(s:string, n:glif)
			p = "<%s>" << string(n)
			s + p
		#
		s = "abc".glifs()
		r = s.fold('', foo)
		`, "r", "<a><b><c>"},
		{`
		func max(s, n)
			if n > s
				return n
			s
		nn = [-20000,3,5,-9,12,6,8,3,0]
		r = nn.fold(-1000000, max)
		`, "r", int64(12)},
		// flat
		{`
		nn = [[1,2], [3,4,5], [6]]
		r = nn.flat()
		`, "r", Anis(1, 2, 3, 4, 5, 6)},
		{`
		nn = [0, [1,[11, 22],[111, 222, 333]], [[44, 45],[51, 52, 53],[61, 62]], [[1001]], 5005]
		r = nn.flat().flat()
		`, "r", Anis(0, 1, 11, 22, 111, 222, 333, 44, 45, 51, 52, 53, 61, 62, 1001, 5005)},
		{`
		func glifs(s: string)
			s.glifs()
		#
		s = "Hello Letters 1234"
		r = s.split(' ').map(glifs).flat()
		`, "r", Anynn([]obb.Glif("HelloLetters1234"))},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestCollBlockNested(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		nn = []
			[1,2,3]
			[4,5,6]
		`, "nn", Anis(Anis(1, 2, 3), Anis(4, 5, 6))},
		{`
		nn = []
			[]
				1
				2
				3
			[]
				4
				5
				6
		`, "nn", Anis(Anis(1, 2, 3), Anis(4, 5, 6))},
		{`
		nn = [77]
			[88]
				1
				2
				3
			[99]
				4
				5
				6
		`, "nn", Anis(77, Anis(88, 1, 2, 3), Anis(99, 4, 5, 6))},

		{`
		nn = (111,222)
			(1,2,3)
			(4,5,6)
		`, "nn", Tanis(111, 222, Tanis(1, 2, 3), Tanis(4, 5, 6))},

		{`
		# deep case
		nn = (,)
			[]
				(,)
					[]
						(,)
							[]
								4
								5
		`, "nn", Tanis(Anis(Tanis(Anis(Tanis(Anis(4, 5))))))},

		{`
		nn = (,)
			(,)
				1
				2
				3
			(,)
				4
				5
				6
		`, "nn", Tanis(Tanis(1, 2, 3), Tanis(4, 5, 6))},

		{`
		nn = {}
			'k1': {}
				'AA': "Hello Alladin!"
				'BB': "Hello Barbara!"
			'k2': {}
				'CC': "Hello Centaur!"
				'DD': "Hello Dambldor!"
		#
		`, "nn", adk(dk{
			"k1": dk{"AA": "Hello Alladin!", "BB": "Hello Barbara!"},
			"k2": dk{"CC": "Hello Centaur!", "DD": "Hello Dambldor!"}})},

		{`
		nn = {'00':'Zero point', '11':'Eleven elefants'}
			'H1':{}
				'AA': "Hello Alladin!"
			'H2':{}
				'BB': "Hello Barbara!"
		`, "nn", adk(dk{
			"00": "Zero point", "11": "Eleven elefants",
			"H1": adk(dk{"AA": "Hello Alladin!"}),
			"H2": adk(dk{"BB": "Hello Barbara!"})})},
		{`
		# deep case
		nn = {}
			'k1': {}
				'k2': {}
					'k3': {}
						'k4': {}
							'k5': {}
								'AA': "Hello Alladin!"
								'BB': "Hello Barbara!"
		#
		`, "nn", adk(dk{
			"k1": adk(dk{
				"k2": adk(dk{
					"k3": adk(dk{
						"k4": adk(dk{
							"k5": dk{"AA": "Hello Alladin!", "BB": "Hello Barbara!"},
						})})})})})},

		{`
		# combi
		nn = {'o1':[-1,-2]}
			'L1': []
				(1,2,3)
			'D2':{}
			 	'T3':(4,5,6)
			'T4': (,)
				7
				888
		`, "nn", adk(dk{
			"o1": Anis(-1, -2),
			"L1": Anis(Tanis(1, 2, 3)),
			"D2": adk(dk{"T3": Tanis(4, 5, 6)}),
			"T4": Tanis(7, 888),
		})},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestCollBlockConstr(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		nn = []
			"aaaaaaaaaa"
			"bbbbbbbbbb"
			"cccccccccc"
			"dddddddddd"
			"xxxxxxxxxx"
			"zzzzzzzzzz"
		`, "nn", Anis("aaaaaaaaaa", "bbbbbbbbbb", "cccccccccc", "dddddddddd", "xxxxxxxxxx", "zzzzzzzzzz")},
		{`
		nn = []
			132
			2 * 100 + 7 * 5 - 1
			300+40+5
		`, "nn", Anis(132, 234, 345)},
		{`
		nn = [123]
			200 + 34
			300+40+5
		#
		`, "nn", Anis(123, 234, 345)},
		{`
		nn = (,)
			132
			2 * 100 + 7 * 5 - 1
			300+40+5
		`, "nn", Tanis(132, 234, 345)},
		{`
		nn = (11,22)
			132
			2 * 100 + 36
			300+40+5
		`, "nn", Tanis(11, 22, 132, 236, 345)},
		{`
		nn = {}
			'AA': "Hello Alladin!"
			'BB': "Hello Barbara!"
			'CC': "Hello Centaur!"
			'DD': "Hello Dambldor!"
		`, "nn", adk(dk{"AA": "Hello Alladin!", "BB": "Hello Barbara!", "CC": "Hello Centaur!", "DD": "Hello Dambldor!"})},
		{`
		nn = {'00':'Zero point', '11':'Eleven elefants'}
			'AA': "Hello Alladin!"
			'BB': "Hello Barbara!"
		`, "nn", adk(dk{"00": "Zero point", "11": "Eleven elefants", "AA": "Hello Alladin!", "BB": "Hello Barbara!"})},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

// list[a : b], tuple[a, b], string[a, b]
// if arg skipped nn[2:], nn[:5], nn[:] (means start=0, end=length)
func TestCollSlice(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		nn = [1,2,3,4,5]
		r = nn[1:3]
		`, "r", Anis(2, 3)},
		{`
		nn = split('a b c d e f g h i j', ' ')
		r = nn[1:5]
		`, "r", Anis("b", "c", "d", "e")},
		{`
		tt = (1,2,3,4,5,6,7,8,9)
		r = tt[3:8]
		`, "r", Tanis(4, 5, 6, 7, 8)},
		{`
		s = "Hello test string!"
		r = s[3: 10]
		`, "r", "lo test"},
		{`
		nn = [1,2,3,4,5]
		r = nn[2:]
		`, "r", Anis(3, 4, 5)},
		{`
		nn = [1,2,3,4,5]
		r = nn[:3]
		`, "r", Anis(1, 2, 3)},
		{`
		tt = (1,2,3,4,5,6,7)
		r = tt[3:]
		`, "r", Tanis(4, 5, 6, 7)},
		{`
		tt = (1,2,3,4,5,6,7)
		r = tt[:4]
		`, "r", Tanis(1, 2, 3, 4)},
		{`
		s = 'ABCDEFGHIJKLMN'
		r = s[5:]
		`, "r", "FGHIJKLMN"},
		{`
		s = 'ABCDEFGHIJKLMN'
		r = s[:7]
		`, "r", "ABCDEFG"},
		{`
		nn = [1,2,3,4,5]
		r = [1,2,3,4,5][1:3]
		`, "r", Anis(2, 3)},
		{`
		r = "Hello test string!"[8:14]
		`, "r", "st str"},
		{`
		nn = [1,2,3,4,5]
		r = nn[:]
		`, "r", Anis(1, 2, 3, 4, 5)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestCollDelElemOper(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		nn = [1,2,3,4,5]
		nn - [2]
		`, "nn", Anis(1, 2, 4, 5)},
		{`
		nn = [1,2,303,4,5]
		r = nn - [2]
		`, "r", int64(303)},
		{`
		dd = {'a':11, 'b':22, 'c':33}
		dd - ['b']
		`, "dd", adk(dk{"a": 11, "c": 33})},
		{`
		dd = {'a':11, 'b':22, 'c':33}
		r = dd - ['b']
		`, "r", int64(22)},
		// {``, "r",  Anis(11, )},
		// {``, "r", adk(dk{"a": 11, "c": 33})},
		// {``, "r",  int64(205)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestCollectionsAddOper(t *testing.T) {
	// tuple + tuple
	// list += list
	// dict += dict
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		nn = [1,2,3]
		x = 1
		x += 5
		nn += [15]
		nn += [16]
		`, "nn", Anis(1, 2, 3, 15, 16)},
		{`
		a = [1,2,3]
		b = [44, 55]
		r = a + b
		`, "r", Anis(1, 2, 3, 44, 55)},
		{`
		a = (1,2)
		b = (33,44)
		r = a + b
		`, "r", Tanis(1, 2, 33, 44)},
		{`
		r = (1,2)
		b = (33, 55)
		r += b
		`, "r", Tanis(1, 2, 33, 55)},
		{`
		a = {'a': 11}
		b = {'b':22, 'c':33}
		r = a + b
		`, "r", adk(dk{"a": 11, "b": 22, "c": 33})},
		{`
		r = {'a':11}
		b = {'b':222, 'c':333}
		r += b
		`, "r", adk(dk{"a": 11, "b": 222, "c": 333})},
		// {``, "r",  Anis(11, )},
		// {``, "r",  int64(205)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestListsCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# empty list
		nn = []
		`, "nn", Anynn([]any{})},
		{`
		# non-empty list
		nn = [1,2,3,4,5]
		`, "nn", Anynn([]int64{1, 2, 3, 4, 5})},
		{`
		# write to list
		nn = [1,2,3]
		nn[1] = 222
		`, "nn", Anynn([]int64{1, 222, 3})},
		{`
		# write to list in loop
		a = [0,0,0,0,0]
		for i=0; i < 5; i += 1
			a[i] = 10 + i
		`, "a", Anynn([]int64{10, 11, 12, 13, 14})},
		{`
		# read from list
		nn = [1,2,3,4,100]
		a = 0
		for i = 0; i < 5; i += 1
			a += nn[i]
		`, "a", int64(110)},
		{`
		# from list to list
		nn = [11,12,13,14,15]
		a = [0,0,0,0,0]
		for i = 0; i < 5; i += 1
			a[i] = nn[i]
		`, "a", Anynn([]int64{11, 12, 13, 14, 15})},
		{`
		# list[i] += v
		a = [0,0,0,0,100]
		for i=0; i < 5; i += 1
			a[i] += 10 + i
		`, "a", Anynn([]int64{10, 11, 12, 13, 114})},
		{`
		# list[i] += list[i]
		nn = [11,12,13,14,15]
		a = [0,0,0,0,100]
		for i=0; i < 5; i += 1
			a[i] += nn[i]
		`, "a", Anynn([]int64{11, 12, 13, 14, 115})},
		// {`` "nn", int64(1)},
		// {``, "nn", int64(1)},
	}

	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestTupleCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# empty tuple
		nn = (,)
		`, "nn", tanynn([]any{})},
		{`
		# empty tuple
		nn = (1,)
		`, "nn", tanynn([]any{int64(1)})},
		{`
		# empty tuple
		nn = ('abc',)
		`, "nn", tanynn([]any{"abc"})},
		{`
		# empty tuple
		nn = (1,2)
		`, "nn", tanynn([]any{1, 2})},
		{`
		# non-empty tuple
		nn = (1,2,3,4,5)
		`, "nn", tanynn([]int64{1, 2, 3, 4, 5})},
		{`
		# tuple, different type elems
		nn = (1,2,true, false, "5", "hello")
		`, "nn", tanynn([]any{1, 2, true, false, "5", "hello"})},
		{`
		# read from tuple
		nn = (1,2,3,4,100])
		a = 0
		for i = 0; i < 5; i += 1
			a += nn[i]
		`, "a", int64(110)},
		{`
		# from tuple to list
		tt = (11, 12, 13, 14, 115)
		a = [0,0,0,0,0]
		for i = 0; i < 5; i += 1
			a[i] = tt[i]
		`, "a", Anynn([]any{11, 12, 13, 14, 115})},
		{`
		# list[i] += tuple[i]
		tt = (11, 12, 13, 14, 115)
		a = [10,20,30,40,50]
		for i = 0; i < 5; i += 1
			a[i] += tt[i]
		`, "a", Anynn([]any{21, 32, 43, 54, 165})},
		// {``, "nn", int64(1)},
	}

	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestDictCase(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# empty dict
		dd = {}
		`, "dd", dk{}},
		{`
		# dict 1 elem
		dd = {'a': 11, }
		`, "dd", adk(dk{"a": 11})},
		{`
		# dict 1 elem
		dd = {'a': 12}
		`, "dd", adk(dk{"a": 12})},
		{`
		# dict
		dd = {'a': 11, 'b':22, 'c':33}
		`, "dd", adk(dk{"a": 11, "b": 22, "c": 33})},
		{`
		# dict str:str
		dd = {'a': 'abc'}
		`, "dd", dk{"a": "abc"}},
		{`
		# dict different type
		dd = {'a': 'abc', 2:22, 'c':[1,2,3]}
		`, "dd", adk(dk{"a": "abc", 2: 22, "c": Anynn([]int64{1, 2, 3})})},
		{`
		# dict from tuple
		tt = (1,2)
		dd = {tt[0] : tt[1]}
		`, "dd", adk(dk{1: 2})},
		{`
		# get elem
		dd = {'a':1, 'b': 20}
		a = dd['a']
		a += dd['b']
		`, "a", int64(21)},
		{`
		# set elem
		dd = {}
		dd['a'] = 11
		dd['b'] = 20
		dd['a'] = 111
		`, "dd", adk(dk{"a": 111, "b": 20})},
		{`
		# dd[k] += dd[k]
		dd = {'a':3, 'b': 2}
		dd['a'] += dd['b']
		`, "dd", adk(dk{"a": 5, "b": 2})},
		{`
		# read from dict
		nn = ['a','b','c']
		dd = {'a': 2, 'b': 3, 'c': 20}
		a = 0
		for i = 0; i < 3; i += 1
			k = nn[i]
			a += dd[k]
		`, "a", int64(25)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
