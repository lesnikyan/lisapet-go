package cases

import (
	"testing"
)

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
		// {`
		// r : int = false
		// `, "r", int64(0)},
		// {``, "r",  int64(205)},
		// {``, "r",  float64(205)},
		// {``, "r",  float64(205)},
		// {``, "r",  float64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}

func TestTypesVars(t *testing.T) {
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
		{`
		r: list = null
		`, "r", Tnull()},
		// {``, "r",  int64(205)},
		// {``, "r",  Anis(11, )},
	}
	for i, tt := range tdata {
		RunTCodeVarExp(t, i, tt)
	}
}
