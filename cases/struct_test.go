package cases

import "testing"

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
		// {``, "r", Stf("T", dk{"a": 0, "b": false})},
		// {``, "r", Stf("T", dk{"a": 0, "b": false})},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
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
