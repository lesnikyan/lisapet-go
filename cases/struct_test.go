package cases

import "testing"

func TestSimpleMethods(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		struct T1 a: int, b: bool
		#
		func t:T1 add(x:int)
			t.a = x
		#
		r = T1{a: 100}
		r.add(1003)
		#
		`, "r", Stf("T1", dk{"a": 1003, "b": false})},
		{`
		struct T1 a: int, b: bool
		#
		func t:T1 add(x:int)
			t.a = x + t.a
		#
		r = T1{a: 100}
		r.add(3)
		#
		`, "r", Stf("T1", dk{"a": 103, "b": false})},
		{`
		struct T1 a: int, b: bool
		#
		func t:T1 add(x:int)
			t.a += x
		#
		r = T1{a: 100}
		r.add(4)
		#
		`, "r", Stf("T1", dk{"a": 104, "b": false})},
		{`
		struct T1 a: int, b: bool
		#
		func t:T1 foo()
			t.a
		#
		s = T1{a: 106}
		r = s.foo()
		#
		`, "r", int64(106)},
		{`
		struct T1 a: int, b: bool
		#
		func t:T1 foo(x:int)
			t.a + x
		#
		s = T1{a: 100}
		r = s.foo(5)
		#
		`, "r", int64(105)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestNestedStructConstr(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		struct T1 a: int, b: bool
		struct T2 t: T1, c: string
		#
		r = T2{}
			t: T1{}
				a: 11
				b: true
			c: "Hey!"
		#
		`, "r", Stf("T2", dk{"t": Stf("T1", dk{"a": 11, "b": true}), "c": "Hey!"})},
		{`
		# empty sub elements
		struct T1 a: list, b: tuple
		struct T2 t: T1, d: dict
		#
		r = T2{}
			t: T1{}
				a: []
					#
				b: (,)
					#
			d: {}
				#
		#
		`, "r", Stf("T2", dk{"t": Stf("T1", dk{"a": Anis(), "b": Tanis()}), "d": dk{}})},
		{`
		struct T1 a: list, b: tuple
		struct T2 t: T1, d: dict
		#
		r = T2{}
			t: T1{}
				a: []
					11
					22
				b: (,)
					33
					44
			d: {}
				55: 5005
				'Mg': "Bombarda!"
				'TT': T1{}
					a: [101]
					b: (102,)
		#
		`, "r", Stf("T2", dk{
			"t": Stf("T1", dk{"a": Anis(11, 22), "b": Tanis(33, 44)}),
			"d": dk{55: 5005, "Mg": "Bombarda!", "TT": Stf("T1", dk{"a": Anis(101), "b": Tanis(102)})}})},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestStructBlockInst(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		struct T1 a: int, b: float
		r = T1{}
			a: 22
		`, "r", Stf("T1", dk{"a": 22, "b": 0.0})},
		{`
		struct T1 a: int, aa: string, b: float, c: bool
		r = T1{c: true}
			a: 22
			b: 10.
			aa: 'Hello block'
		`, "r", Stf("T1", dk{"a": 22, "aa": `Hello block`, "b": 10.0, "c": true})},
		{`
		struct T1
			a: int
			b: int
			c: int
			d: int
			e: int
			f: int
			g: int
			h: int
			i: int
			j: int
			aa: string
		x = 5
		r = T1{}
			a: 11
			b: 22
			c: 33
			d: 44
			e: 55
			f: 66
			g: 77
			h: 88
			i: 99
			j: 100 + x
			aa: 'Hello block'
		`, "r", Stf("T1", dk{"a": 11, "aa": "Hello block", "b": 22, "c": 33, "d": 44, "e": 55, "f": 66, "g": 77, "h": 88, "i": 99, "j": 105})},
		// {``, "r", Stf("T", dk{"a": 0, "b": false})},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestStructBlockDef(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		struct T1
			a: int
			b: float
		r = T1{a: 22}
		`, "r", Stf("T1", dk{"a": 22, "b": 0.0})},
		{`
		struct T1 a: int, aa: int
			b: float
			c: bool
		r = T1{a: 22}
		`, "r", Stf("T1", dk{"a": 22, "aa": 0, "b": 0.0, "c": false})},
		{`
		struct T1 a: int, aa: string
			b: float
			c: bool
		r = T1{b: 2.2, c: true}
		`, "r", Stf("T1", dk{"a": 0, "aa": "", "b": 2.2, "c": true})},
		// {``, "r", Stf("T", dk{"a": 0, "b": false})},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestStructTyping(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		struct T a: int
		r = T{a: true}
		# r.a = true
		`, "r", Stf("T", dk{"a": 1})},
		{`
		struct T a: float
		r = T{a: 25}
		`, "r", Stf("T", dk{"a": 25.0})},
		{`
		struct T a: int, b: float, s: string, nn: list
		r = T{}
		r.a = true
		r.b = 5
		r.s = "Nya"
		r.nn = [2,3,44]
		`, "r", Stf("T", dk{"a": 1, "b": 5.0, "s": "Nya", "nn": Anis(2, 3, 44)})},
		{`
		struct A a: int
		struct T obj: A
		r = T{obj: null}
		`, "r", Stf("T", dk{"obj": Tnull()})},
		{`
		struct A a: int
		struct T obj: A
		r = T{}
		`, "r", Stf("T", dk{"obj": Tnull()})},
		{`
		struct A a: int
		struct T obj: A
		r = T{}
		r.obj = null
		`, "r", Stf("T", dk{"obj": Tnull()})},
		// {``, "r", Stf("T", dk{"a": 0, "b": false})},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestStructFieldSet(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		struct T a: int
		r = T{a: 21}
		r.a = 23
		`, "r", Stf("T", dk{"a": 23})},
		{`
		struct T a: int, b: bool, s: string, nn: list
		r = T{}
		r.a = 21
		r.b = true
		r.s = "Nya"
		r.nn = [2,3,44]
		`, "r", Stf("T", dk{"a": 21, "b": true, "s": "Nya", "nn": Anis(2, 3, 44)})},
		{`
		struct T1 a: int, b: bool
		#
		r = T1{a: 100}
		r.a = 2 + r.a
		#
		`, "r", Stf("T1", dk{"a": 102, "b": false})},
		{`
		struct T1 a: int, b: bool
		#
		r = T1{a: 100}
		r.a += 3
		#
		`, "r", Stf("T1", dk{"a": 103, "b": false})},
		// {``, "r", Stf("T", dk{"a": 0, "b": false})},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestStructFieldGet(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		struct T1 a:int
		t = T1{a: 5}
		r = t.a
		`, "r", int64(5)},
		{`
		struct T1 a:bool
		t = T1{a: true}
		r = t.a
		`, "r", true},
		{`
		struct T1 a:string
		t = T1{a: 'Bambarbia'}
		r = t.a
		`, "r", "Bambarbia"},
		{`
		struct T1 a:list
		t = T1{a: [1,2,3]}
		r = t.a
		`, "r", Anis(1, 2, 3)},
		{`
		struct T1 a:tuple
		t = T1{a: (1,2,31)}
		r = t.a
		`, "r", Tanis(1, 2, 31)},
		{`
		struct T1 a:dict
		t = T1{a: {'b':222}}
		r = t.a
		`, "r", adk(dk{"b": 222})},
		{`
		struct T1 a:int
		struct T2 b: T1
		
		t1 = T1{a: 24}
		t = T2{b: t1}
		r = t.b
		`, "r", Stf("T1", dk{"a": 24})},
		{`
		struct T1 a:int
		struct T2 b: T1
		
		t1 = T1{a: 25}
		t = T2{b: t1}
		r = t.b.a
		`, "r", int64(25)},
		// {``, "r", Stf("T", dk{"a": 0, "b": false})},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestStructDef(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		struct FirstType a: int, b: bool
		r = FirstType{}
		`, "r", Stf("FirstType", dk{"a": 0, "b": false})},
		{`
		struct Type1 a: int, b: bool
		r = Type1{a:111}
		`, "r", Stf("Type1", dk{"a": 111, "b": false})},
		{`
		struct Type1 a: int, b: bool, s: string
		r = Type1{a:112, b: true, s: 'Hello structs!'}
		`, "r", Stf("Type1", dk{"a": 112, "b": true, "s": "Hello structs!"})},

		{`
		struct T1 a: int
		struct T2 a: int
		struct T3 a: int
		struct T4 a: int
		struct T5 a: int
		struct T6 a: int
		struct T7 a: int
		struct T8 a: int
		
		t1 = T1{a:1}
		t2 = T2{a:2}
		t3 = T3{a:3}
		t4 = T4{a:4}
		t5 = T5{a:5}
		t6 = T6{a:6}
		t7 = T7{a:7}
		t8 = T8{a:8}
		r = [t1, t2, t3, t4, t5, t6, t7, t8]
		`, "r", Anis(
			Stf("T1", dk{"a": 1}),
			Stf("T2", dk{"a": 2}),
			Stf("T3", dk{"a": 3}),
			Stf("T4", dk{"a": 4}),
			Stf("T5", dk{"a": 5}),
			Stf("T6", dk{"a": 6}),
			Stf("T7", dk{"a": 7}),
			Stf("T8", dk{"a": 8}),
		)},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
