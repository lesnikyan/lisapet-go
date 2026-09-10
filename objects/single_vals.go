package objects

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/lesnikyan/lisapet-go/base"
)

// null value
type Null struct {
}

// internal value, can be produced from EmptyExpr
type EmptyVal struct {
}

// type StructVal struct {
// 	Fields []*base.Var
// }

type Regexp struct {
	Pattern *regexp.Regexp
	Src     string
	// Flags   []rexs.Flags
}

func (rx *Regexp) Init() error {
	// TODO: flags
	ptt, err := regexp.Compile(rx.Src)
	if err != nil {
		return err
	}
	rx.Pattern = ptt
	return nil
}

func (rx *Regexp) Match(s string) bool {
	return rx.Pattern.MatchString(s)
}

func (rx *Regexp) Find(s string) []string {
	ss := rx.Pattern.FindAllString(s, -1)
	// fmt.Printf("rx.Find#2: L=%d %v \n", len(ss), ss)
	return ss
}

func (rx *Regexp) FindSubs(s string) [][]string {
	ss := rx.Pattern.FindAllStringSubmatch(s, -1)
	return ss
}

func (rx *Regexp) Split(s string) []string {
	ss := rx.Pattern.Split(s, -1)
	return ss
}

func (rx *Regexp) Replace(s string, repl string) string {
	res := rx.Pattern.ReplaceAllString(s, repl)
	return res
}

var validFlags = []rune("aimsU")

func MakeRegexp(src string, flags string) *Regexp {
	pref := ""
	if flags != "" {
		for _, s := range []rune(flags) {
			if !slices.Contains(validFlags, s) {
				panic(fmt.Sprintf("regexp: incorrect flag of regexp: `%s`", string(s)))
			}
		}
		pref = fmt.Sprintf("(?%s)", flags)
	}
	ptt := fmt.Sprintf("%s%s", pref, src)
	rx := &Regexp{Src: ptt}
	err := rx.Init()
	if err != nil {
		panic(fmt.Sprintf("regexp: incorrect final pattern of regexp: `%s`", ptt))
	}
	// fmt.Printf("regexp ptt: \"%s\" \n", ptt)
	return rx
}

//====

type Glif = rune

//====

type Bytes []byte

func (bb Bytes) GetElem(i int64) (*base.Val, error) {
	index := int(i)
	if index < 0 {
		index = len(bb) + index
	}
	if index < 0 || index >= len(bb) {
		return nil, fmt.Errorf("Bytes sequence with len= %d doesn't have element %d", len(bb), index)
	}
	v := bb[int(index)]
	return base.NewVal(v), nil
}

func (bb Bytes) Add(n any) (Bytes, error) {
	var b byte
	switch v := n.(type) {
	case byte:
		b = v
	case int64:
		b = byte(v)
	default:
		return nil, fmt.Errorf("bytes.Add: incorrect argument type %T", v)
	}
	bb = append(bb, b)
	// fmt.Println("bytes: ad", bb)
	return bb, nil
}
