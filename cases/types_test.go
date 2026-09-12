package cases

import (
	"testing"

	oob "github.com/lesnikyan/lisapet-go/objects"
)

func TestConstructDict(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r:dict = dict()
		`, "r", adk(dk{})},
		{`
		r = dict({11:111, '22':222})
		`, "r", adk(dk{11: 111, "22": 222})},
		{`
		r = dict([(3,33), ('4', 44), (7,777)])
		`, "r", adk(dk{3: 33, "4": 44, 7: 777})},
		{`
		r = dict((55, '66', '8'), ('5', 6, 888))
		`, "r", adk(dk{55: "5", "66": 6, "8": 888})},
		{`
		r = dict([1,2,3,4,5], split('6,7,8,9,0,10,12,13', ","))
		`, "r", adk(dk{1: "6", 2: "7", 3: "8", 4: "9", 5: "0"})},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestConstructBytes(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r:bytes = bytes(0)
		`, "r", oob.Bytes{}},
		{`
		r = bytes()
		`, "r", oob.Bytes{}},
		{`
		r = bytes(12)
		`, "r", oob.Bytes{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{`
		r = bytes(byte(4))
		`, "r", oob.Bytes{0, 0, 0, 0}},
		{`
		r = bytes(00x5)
		`, "r", oob.Bytes{0, 0, 0, 0, 0}},
		{`
		r = bytes([])
		`, "r", oob.Bytes{}},
		{`
		r = bytes([00xe5, 00x01])
		`, "r", oob.Bytes{0xe5, 1}},
		{`
		r = bytes((11,22, 255))
		`, "r", oob.Bytes{11, 22, 255}},
		{`
		r = bytes(0x[0 1 2 3 f 1f a0 ff])
		`, "r", oob.Bytes{0, 1, 2, 3, 0xf, 0x1f, 0xa0, 0xff}},
		{`
		r = bytes('Hello 123!')
		`, "r", oob.Bytes{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x31, 0x32, 0x33, 0x21}},
		//  string([]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x31, 0x32, 0x33, 0x21}
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestConstructList(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r:list = list(0)
		`, "r", Anis()},
		{`
		r = list()
		`, "r", Anis()},
		{`
		r = list(10)
		`, "r", Anis(0, 0, 0, 0, 0, 0, 0, 0, 0, 0)},
		{`
		r = list([])
		`, "r", Anis()},
		{`
		r = list([1, 2, 'Q'])
		`, "r", Anis(1, 2, "Q")},
		{`
		r = list((11,22))
		`, "r", Anis(11, 22)},
		{`
		r = list(byte(4))
		`, "r", Anis(0, 0, 0, 0)},
		{`
		r = list(0x[0 1 2 3 f 1f a0 ff])
		`, "r", Anis(byte(0), byte(1), byte(2), byte(3), byte(0xf), byte(0x1f), byte(0xa0), byte(0xff))},
		{`
		r = list('Hello 123!')
		`, "r", Anis('H', 'e', 'l', 'l', 'o', ' ', '1', '2', '3', '!')},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestConstructTuple(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r:tuple = tuple(0)
		`, "r", Tanis()},
		{`
		r = tuple()
		`, "r", Tanis()},
		{`
		r = tuple(10)
		`, "r", Tanis(0, 0, 0, 0, 0, 0, 0, 0, 0, 0)},
		{`
		r = tuple([])
		`, "r", Tanis()},
		{`
		r = tuple([1, 2, 'Q'])
		`, "r", Tanis(1, 2, "Q")},
		{`
		r = tuple((11,22))
		`, "r", Tanis(11, 22)},
		{`
		r = tuple(byte(4))
		`, "r", Tanis(0, 0, 0, 0)},
		{`
		r = tuple(0x[0 1 2 3 f 1f b0 ff])
		`, "r", Tanis(byte(0), byte(1), byte(2), byte(3), byte(0xf), byte(0x1f), byte(0xb0), byte(0xff))},
		{`
		r = tuple('Hello 123!')
		`, "r", Tanis('H', 'e', 'l', 'l', 'o', ' ', '1', '2', '3', '!')},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestConstructIBFG(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		# int
		r:int = int(0)
		`, "r", int64(0)},
		{`
		# int
		r = int(1)
		`, "r", int64(1)},
		{`
		# int
		r = int('100500')
		`, "r", int64(100500)},
		{`
		# int
		r = int(true)
		`, "r", int64(1)},
		{`
		# int
		bb = 0x[1 2 3]
		r = int(bb[2])
		`, "r", int64(3)},
		{`
		# int
		r = int(g'A')
		`, "r", int64(65)},
		{`
		# int
		r = int(2.1)
		`, "r", int64(2)},
		{`
		# int
		r = int(null)
		`, "r", int64(0)},
		{`
		r:bool = bool(0)
		`, "r", false},
		{`
		r = bool(2)
		`, "r", true},
		{`
		r = bool(1.111)
		`, "r", true},
		{`
		r = bool('true')
		`, "r", true},
		{`
		r = bool('false')
		`, "r", false},
		{`
		r = bool(null)
		`, "r", false},
		{`
		r:float = float(0)
		`, "r", float64(0)},
		{`
		r = float(1)
		`, "r", float64(1)},
		{`
		r = float(1.25)
		`, "r", float64(1.25)},
		{`
		r = float(true)
		`, "r", float64(1)},
		{`
		r = float(null)
		`, "r", float64(0)},
		{`
		r = float(false)
		`, "r", float64(0)},
		{`
		r = float('1.22')
		`, "r", float64(1.22)},
		{`
		r:glif = glif(100)
		`, "r", Gf("d")},
		{`
		bb = [0 61]
		r = glif(bb[1])
		`, "r", Gf("a")},
		{`
		r = glif('百')
		`, "r", '百'},
		{`
		bb = [e7 99 be]
		r = glif(bb)
		`, "r", '百'},
		{`
		b = byte(101)
		r = glif(b)
		`, "r", 'e'},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestConstructMaybe(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = some(1)
		`, "r", Tmay(1)},
		{`
		r:maybe = some(2)
		`, "r", Tmay(2)},
		{`
		r = some('Hello')
		`, "r", Tmay("Hello")},
		{`
		r = some([1,2,3])
		`, "r", Tmay(Anis(1, 2, 3))},
		{`
		struct Abc a: int
		a1 = Abc{a:5}
		r = some(a1)
		`, "r", Tmay(Stf("Abc", dk{"a": 5}))},

		{`
		r:maybe = none
		`, "r", oob.None()},
		{`
		a1 = some('in-some')
		r = some(a1)
		`, "r", Tmay(Tmay("in-some"))},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
func TestConstructByte(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r:byte = byte(1)
		`, "r", byte(1)},
		{`
		r = byte(100)
		`, "r", byte(100)},
		{`
		r = byte(255)
		`, "r", byte(255)},
		{`
		r = byte(0x105)
		`, "r", byte(5)},
		{`
		r = byte(0xf)
		`, "r", byte(0xf)},
		{`
		r = byte(true)
		`, "r", byte(1)},
		{`
		r = byte(false)
		`, "r", byte(0)},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestConstructString(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r = string('Hello')
		`, "r", "Hello"},
		{`
		r = string(65)
		`, "r", "65"},
		{`
		bb = 0x[fa]
		r = string(bb[0])
		`, "r", "250"},
		{`
		r = string([65, 66, 67, 97, 115, 116])
		`, "r", "ABCast"},
		{`
		bb = 0x[48 65 6c 6c 6f 20 62 79 74 65 73 21]
		r = string(bb)
		`, "r", "Hello bytes!"},
		{`
		r = string([c6 90 20 c6 8d 20 c6 80 20 c6 8b 20 c6 95 20 c6 a9 20 c6 b1 20 c6 b3 20 c6 9b])
		`, "r", "Ɛ ƍ ƀ Ƌ ƕ Ʃ Ʊ Ƴ ƛ"},
		{`
		r = string([g'G', g'L', g'i', g'P', g'h', g'S'])
		`, "r", "GLiPhS"},
		{`
		r = string(null)
		`, "r", "null"},
		{`
		r = string(true)
		`, "r", "true"},
		{`
		r = string(false)
		`, "r", "false"},
		{`
		r = string('')
		`, "r", ""},
		{`
		r = string()
		`, "r", ""},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

/*
float = int
int = bool
float = bool
float : int + float
int : bool + int
*/
func TestTypesCast(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		r : int = false
		`, "r", int64(0)},
		{`
		r: int = null
		`, "r", int64(0)},
		{`
		r: float = 12
		`, "r", float64(12)},
		{`
		r: float = true
		`, "r", float64(1)},
		{`
		r: float = null
		`, "r", float64(0)},
		{`
		r: list = null
		`, "r", Tnull()},
		// {``, "r",  float64(205)},
		// {``, "r",  float64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestTypesVarsSameType(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		a:int = 123
		`, "a", int64(123)},
		{`
		r: string = 'karamba!'
		`, "r", "karamba!"},
		{`
		r: float = 1.25
		`, "r", float64(1.25)},
		{`
		r: bool = true
		`, "r", true},
		{`
		r: bool = false
		`, "r", false},
		{`
		r: list = [1,2,3]
		`, "r", Anis(1, 2, 3)},
		{`
		r: tuple = (1,2,'S')
		`, "r", Tanis(1, 2, "S")},
		{`
		r : dict = {'a': 15}
		`, "r", adk(dk{"a": 15})},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
