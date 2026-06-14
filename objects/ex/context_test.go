package objects

// func Test1(t *testing.T) {
// 	tdata := []struct {
// 		c    rune
// 		prev Lt.Lt
// 		exp  Lt.Lt
// 	}{}

// 	for _, tt := range tdata {
// 		t.Run(fmt.Sprintf("e-type '%s', <%s> ", string(tt.c), Lt.TName(tt.exp)), func(t2 *testing.T) {
// 			res := elemType(tt.c, tt.prev)
// 			assert.Equal(t2, tt.exp, res, "")
// 		})
// 	}
// }

// func TestNewContext(t *testing.T) {

// 	tvals := []struct {
// 		item  int32
// 		dtype dt.DType
// 		val   interface{}
// 	}{
// 		{0, dt.Int, int64(123)},
// 		{1, dt.Float, float64(11.5)},
// 		{2, dt.Int, int64(25)},
// 		{3, dt.String, "hello1"},
// 		{4, dt.String, "yello2"},
// 	}
// 	tdata := []struct {
// 		name string
// 		tp   dt.DType
// 		vid  int32
// 	}{
// 		{"x", dt.Int, 0},
// 		{"size", dt.Float, 1},
// 		{"y", dt.Int, 2},
// 		{"s1", dt.String, 3},
// 		{"s2", dt.String, 4},
// 	}
// 	parent := NewContext(nil)
// 	tctx := NewContext(parent)
// 	vvs := make([]interface{}, len(tvals))
// 	vals := make([]*CVal, len(tvals))
// 	vars := make([]*CVar, len(tvals))
// 	for i, tv := range tvals {
// 		t.Run(fmt.Sprintf("ctx-val '%d', ", (tv.item)), func(t2 *testing.T) {
// 			// cv, ok := tctx.AddVal(tv.dtype, tv.val)
// 			v := tVal(tv.val)
// 			// t.Logf("tt1>> %v > %v ", tv.val, v)
// 			cv := tctx.PutVal(v)
// 			// t.Logf("tt2>>%v > %v", v, cv)
// 			vals[i] = cv
// 			vr := &CVar{Type: tdata[i].tp, v: cv}
// 			vvs[i] = tv.val
// 			vars[i] = vr
// 			// assert.True(t, ok)
// 			tctx.SetVar(tdata[i].name, vr)
// 		})
// 	}
// 	for _, tv := range tdata {
// 		t.Run(fmt.Sprintf("ctx-var '%s', ", (tv.name)), func(t2 *testing.T) {
// 			vr, cx := tctx.GetVar(tv.name)
// 			vl := vr.v
// 			assert.Equal(t2, vr.Type, vl.Type)
// 			switch vl.Type {
// 			case dt.Int:
// 				val, ok := cx.GetInt(vl)
// 				t2.Logf("tt4 > %v, %v", val, ok)
// 				// assert.True(t2, ok)
// 				// assert.Equal(t2, vvs[tv.vid], val)
// 			}
// 		})
// 	}
// }
