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
ok list
ok tuple
ok dict
ok bytes

*/

type Bytes = obb.Bytes

func TestBytesMethods(t *testing.T) {
	// src := []byte{0, 0, 0xff, 0x1, 0x2, 0x3, 0xf, 0xf0}
	// src := []byte{0b10000000, 0, 0, 0, 0, 0, 0, 0}
	// src := []byte{0b00000000, 0, 0, 0, 0, 0, 1, 0}
	// tn := binary.BigEndian.Uint64(src[0:8])
	// fmt.Printf(" tt: %v\n - - %v \n", tn, int64(tn))
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// bits
		{`
		bb = 0x[01]
		r = bb.bits()
		`, "r", Anis(0, 0, 0, 0, 0, 0, 0, 1)},
		{`
		bb = 0x[f3]
		r = bb.bits()
		`, "r", Anis(1, 1, 1, 1, 0, 0, 1, 1)},
		{`
		bb = 0x[f0 07]
		r = bb.bits()
		`, "r", Anis(1, 1, 1, 1, 0, 0, 0, 0, 0x0, 0, 0, 0, 0, 1, 1, 1)},
		{`
		bb = 0b[11110000 10100000]
		r = bb.bits()
		`, "r", Anis(1, 1, 1, 1, 0, 0, 0, 0, 0x1, 0, 1, 0, 0, 0, 0, 0)},
		{`
		bb = 0x[01 87 0f f9]
		r = bb.bits()
		`, "r", Anis(
			0, 0, 0, 0, 0, 0, 0, 1, 0x1, 0, 0, 0, 0, 1, 1, 1,
			0, 0, 0, 0, 1, 1, 1, 1, 0x1, 1, 1, 1, 1, 0, 0, 1)},
		// blocks
		{`
		bb = 0x[01 02 0e f1]
		r = bb.blocks(2)
		`, "r", Anis(Bytes{1, 2}, Bytes{0x0e, 0xf1})},
		// // blocks
		{`
		bb = 0x[00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f]
		r = bb.blocks(4)
		`, "r", Anis(Bytes{0x0, 0x1, 0x2, 0x3}, Bytes{0x4, 0x5, 0x6, 0x7}, Bytes{0x8, 0x9, 0xa, 0xb}, Bytes{0xc, 0xd, 0xe, 0xf})},
		{`
		bb = 0x[00 00 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d]
		r = bb.blocks(4)
		`, "r", Anis(Bytes{0x0, 0x0, 0x0, 0x1}, Bytes{0x2, 0x3, 0x4, 0x5}, Bytes{0x6, 0x7, 0x8, 0x9}, Bytes{0xa, 0xb, 0xc, 0xd})},
		{`
		bb = 0x[01 02 03 04 05 06 07 08 09 0a 0b 0c 0d]
		r = bb.blocks(4)
		`, "r", Anis(Bytes{0x0, 0x0, 0x0, 0x1}, Bytes{0x2, 0x3, 0x4, 0x5}, Bytes{0x6, 0x7, 0x8, 0x9}, Bytes{0xa, 0xb, 0xc, 0xd})},
		{`
		bb = 0x[01 02 03 04 05 06 07 08 09 0a 0b 0c 0d]
		r = bb.blocks(10)
		`, "r", Anis(Bytes{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3}, Bytes{0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd})},
		// map
		{`
		func b2i(x: bytes)
			nn = []
			for n <- x
				nn <- int(n)
			nn
		#
		bb = 0x[01 02 03 04 05 06 07 08 09 0a 0b 0c 0d]
		bs = bb.blocks(4)
		r = bs.map(b2i)
		`, "r", Anis(Anis(0, 0, 0, 1), Anis(2, 3, 4, 5), Anis(6, 7, 8, 9), Anis(0xa, 0xb, 0xc, 0xd))},
		{`
		func foo(b: byte)
			00xff - b
		#
		bb = 0x[00 01 02 03 f1 f2]
		r = bb.map(foo)
		`, "r", Bytes{0xff, 0xfe, 0xfd, 0xfc, 0xe, 0xd}},
		{`
		func foo(b: byte)
			00xf0 + b
		#
		bb = 0x[00 01 02 03 0f 0e 0d]
		r = bb.map(foo)
		`, "r", Bytes{0xf0, 0xf1, 0xf2, 0xf3, 0xff, 0xfe, 0xfd}},
		// reverse
		{`
		bb = 0x[01 02 f1 f5]
		r = bb.reverse()
		`, "r", Bytes{0xf5, 0xf1, 0x2, 0x1}},
		{`
		r = 0x[1 2 3 4 5 6 7 8 9 10 11 12 13 14 15].reverse()
		`, "r", Bytes{0x15, 0x14, 0x13, 0x12, 0x11, 0x10, 0x9, 0x8, 0x7, 0x6, 0x5, 0x4, 0x3, 0x2, 0x1}},
		{`
		func foo(b: byte)
			00xff - b
		bb = 0x[01 02 f1 f5]
		r = bb.reverse().map(foo)
		`, "r", Bytes{0xa, 0xe, 0xfd, 0xfe}},
		// fold
		{`
		func foo (s:int, n:byte)
			s + int(n)
		#
		bb = [1 2 3 4 5 f0]
		r = bb.fold(0, foo)
		 `, "r", int64(255)},
		{`
		func foo (s:list, n)
			s <- n
			s
		#
		bb = [1 2 3 4 5 f0]
		r = bb.fold([], foo)
		 `, "r", Anynn([]byte{1, 2, 3, 4, 5, 0xf0})},
		{`
		func foo (s:list, n)
			s <- int(n) * 5
		#
		bb = [1 2 3 4 5 f0]
		r = bb.fold([], foo)
		 `, "r", Anis(5, 10, 15, 20, 25, 1200)},
		{`
		func foo(s, n)
			s + string(glif(n))
		bb = [48 65 6c 6c 6f 2c 20 46 6f 6c 64 21]
		r= bb.fold('result ', foo)
		`, "r", "result Hello, Fold!"},
		// nums
		{`
		bb = 0d[0 1 0 2 0 3 0 4 0 5]
		r = bb.nums(2)
		`, "r", Anis(1, 2, 3, 4, 5)},
		{`
		bb = 0x[0 1 2 3 5 9 a b f]
		r = bb.nums(1)
		`, "r", Anis(0, 1, 2, 3, 5, 9, 0xa, 0xb, 0xf)},
		{`
		bb = 0d[100 1 100 2 100 3 100 4 100 5]
		r = bb.nums(2)
		`, "r", Anis(25601, 25602, 25603, 25604, 25605)},
		{`
		bb = 0d[10 1 10 2 10 3 10 4 10 5]
		r = bb.nums(2)
		`, "r", Anis(2561, 2562, 2563, 2564, 2565)},
		{`
		bb = 0x[ff 00 80 01 00 01 7f ff]
		r = bb.nums(2)
		`, "r", Anis(65280, 32769, 1, 32767)},
		{`
		bb = 0x[ff 00 00 00 80 00 00 01 00 00 00 01 7f ff ff ff]
		r = bb.nums(4)
		`, "r", Anis(4278190080, 2147483649, 1, 2147483647)},
		{`
		bb = 0x[00 00 00 00 00 00 00 01 00 00 00 00 00 00 00 00 80 00 00 00 00 00 00 01 ff ff ff ff ff ff ff fb]
		r = bb.nums(8)
		`, "r", Anis(1, 0, -9223372036854775807, -5)},
		// {``, "r",  Anis(1)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestListSortRev(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// sort
		{`
		nn = [5,4,6,3,7,2,-8,1,0,-1000000,6,4,2]
		r = nn.sort()
		`, "r", Anis(-1000000, -8, 0, 1, 2, 2, 3, 4, 4, 5, 6, 6, 7)},
		{`
		nn = [1,2,3,4,5,999999999]
		r = nn.sort()
		`, "r", Anis(1, 2, 3, 4, 5, 999999999)},
		{`
		nn = [9999999,  '5', g'5',  'A', 'S', 'Aaaa', g'a', g'S', true, g'A', -100,-10000000]
		r = nn.sort()
		`, "r", Anis(true, -10000000, -100, 9999999, '5', 'A', 'S', 'a', "5", "A", "Aaaa", "S")},
		{`
		nn = [9999999, 100, 5, '5', g'5', byte(5), 3, 1 == 3, 1, byte(122), false, 'A', 'S', 'Aaaa', g'a', g'S', byte(255), true, g'A', 4==4, -100,-10000000]
		r = nn.sort()
		`, "r", Anis(false, false, true, true, -10000000, -100, 1, 3, byte(0x5), 5, 100, byte(0x7a), byte(0xff), 9999999, '5', 'A', 'S', 'a', "5", "A", "Aaaa", "S")},
		{`
		nn = ['a','aa','A', 'S', g'a', g'A', g'S', g'z']
		r = nn.sort()
		`, "r", Anis('A', 'S', 'a', 'z', "A", "S", "a", "aa")},
		{`
		nn = [300, 500, -1, g'a', g'A', g'S', g'z', 700, 99999]
		r = nn.sort()
		`, "r", Anis(-1, 300, 500, 700, 99999, 'A', 'S', 'a', 'z')},
		{`
		nn = ['a', 'aa', 'aA', 'az', g'a']
		r = nn.sort()
		`, "r", Anis('a', "a", "aA", "aa", "az")},
		{`
		nn = [00x1, 00x10, 00xfa, 00x01, 00xf0, 00x0f]
		r = nn.sort()
		`, "r", Anynn([]byte{0x1, 0x1, 0xf, 0x10, 0xf0, 0xfa})},
		{`
		nn = [1.2, 100.002, 2, 3, byte(250), 'Q', g'Q', true]
		r = nn.sort()
		`, "r", Anis(true, 1.2, 2, 3, 100.002, byte(0xfa), 'Q', "Q")},
		{`
		nn = [1, byte(1), byte(100), 100, 'A', 1==1, g'A']
		r = nn.sort()
		`, "r", Anis(true, 1, byte(0x1), byte(0x64), 100, 'A', "A")},
		// reverse
		{`
		nn = [1,2,3,4,5]
		r = nn.reverse()
		`, "r", Anis(5, 4, 3, 2, 1)},
		{`
		nn = [1,2,3,4,5, '','A','b','c',true, false, byte(0xee), g'Q', 1.2]
		r = nn.reverse()
		`, "r", Anis(1.2, 'Q', byte(0xee), false, true, "c", "b", "A", "", 5, 4, 3, 2, 1)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestDictMethods(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// map
		{`
		func foo(k, v)
			k, v + 100
		#
		dd = {'1':11, 'a':22}
		r = dd.map(foo)
		`, "r", adk(dk{"1": 111, "a": 122})},
		{`
		func foo(k, v)
			'<%s>' << k, '(%d)' << v
		#
		dd = {'1':11, 'a':22}
		r = dd.map(foo)
		`, "r", adk(dk{"<1>": "(11)", "<a>": "(22)"})},
		// k-map
		{`
		func foo(k)
			'_%s' << k
		#
		dd = {'1':11, 'b':22}
		r = dd.kmap(foo)
		`, "r", adk(dk{"_1": 11, "_b": 22})},
		// v-map
		{`
		func foo(v)
			v + 1000
		#
		dd = {'3':33, 'c':44}
		r = dd.vmap(foo)
		`, "r", adk(dk{"3": 1033, "c": 1044})},
		// keys
		{`
		dd = {1:11, 2:22, 'aa':'33', 'bb':44}
		kk = dd.keys()
		r = {}
		for k <- kk
			r <- (k, 0)
		`, "r", adk(dk{"aa": 0, "bb": 0, 1: 0, 2: 0})},
		{`
		func foo(s)
			'~%s' << s
		#
		dd = {'111':11, '2':22, 'aa':'33', 'bb':44}
		k2 = dd.keys().map(foo)
		r = {}
		i = 1
		for k <- k2
			r[k] = k[1:]
			i += 1
		`, "r", adk(dk{"~111": "111", "~2": "2", "~aa": "aa", "~bb": "bb"})},
		// vals
		{`
		dd = {'111':11, '2':22, 'aa':'33', 'bb':44}
		vv = dd.vals()
		r = {}
		for v <- vv
			r <- (v, 0)
		`, "r", adk(dk{"33": 0, 11: 0, 22: 0, 44: 0})},
		// {``, "r",  Anis( )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

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
