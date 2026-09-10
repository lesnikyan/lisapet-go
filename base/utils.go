package base

func CropStr(s string, pref int, suff int) string {
	rr := []rune(s)
	rlen := len(rr)
	if pref+suff > rlen {
		return "" // fast hack
	}
	return string(rr[pref : rlen-suff])
}
