package cases

import "testing"

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
