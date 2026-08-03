package lang

/*
*
list of runes to map keys
*/
func Kmap[K comparable](keys []K) map[K]bool {
	r := make(map[K]bool)
	for _, n := range keys {
		r[n] = true
	}
	return r
}
