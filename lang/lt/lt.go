package lt

type Lt int

const (
	None    Lt = 0
	Space   Lt = 1 //# white space
	Word    Lt = 2
	Lang    Lt = 3 //# special words
	Num     Lt = 4
	Oper    Lt = 5
	Comm    Lt = 6
	Text    Lt = 7
	Quot    Lt = 8
	Esc     Lt = 9
	Block   Lt = 11 //# open block. start of block if more than shift of prev line
	Close   Lt = 12 //# close block
	Endline Lt = 13
	Indent  Lt = 14 //# set of spaces, 4-th by default,
	Mtcomm  Lt = 15 //# multiline comment
	Mttext  Lt = 16 //# multiline string
)

var typeNames = map[Lt]string{
	None:    "none",
	Space:   "space",
	Word:    "word",
	Lang:    "lang",
	Num:     "num",
	Oper:    "oper",
	Comm:    "comm",
	Text:    "text",
	Quot:    "quot",
	Esc:     "esc",
	Block:   "block", //# open block. start of block if more than shift of prev line
	Close:   "close", //# close block
	Endline: "endline",
	Indent:  "indent", //# set of spaces, 4-th by default,
	Mtcomm:  "mtcomm", //# multiline comment
	Mttext:  "mttext", //# multiline string
}

func TName(n Lt) string {
	if v, ok := typeNames[n]; ok {
		return v
	}
	return "-/-"
}

// const (
// 	LtNone    Lt = 0
// 	LtSpace   Lt = 1 //#white space
// 	LtOrd     Lt = 2
// 	LtLang    Lt = 3 //# special words
// 	LtNum     Lt = 4
// 	LtOper    Lt = 5
// 	LtComm    Lt = 6
// 	LtText    Lt = 7
// 	LtQuot    Lt = 8
// 	LtEsc     Lt = 9
// 	LtBlock   Lt = 11 //# open block. start of block if more than shift of prev line
// 	LtClose   Lt = 12 //# close block
// 	LtEndline Lt = 13
// 	LtIndent  Lt = 14 //# set of spaces, 4-th by default,
// 	LtMtcomm  Lt = 15 //# multiline comment
// 	LtMttext  Lt = 16 //# multiline string
// )

// type Lt int

// const (
// 	None    Lt = 0
// 	space   Lt = 1 //#white space
// 	word    Lt = 2
// 	lang    Lt = 3 //# special words
// 	num     Lt = 4
// 	oper    Lt = 5
// 	comm    Lt = 6
// 	text    Lt = 7
// 	quot    Lt = 8
// 	esc     Lt = 9
// 	block   Lt = 11 //# open block. start of block if more than shift of prev line
// 	close   Lt = 12 //# close block
// 	endline Lt = 13
// 	indent  Lt = 14 //# set of spaces, 4-th by default,
// 	mtcomm  Lt = 15 //# multiline comment
// 	mttext  Lt = 16 //# multiline string
// )

// type LangLTypes struct {
// 	none    int
// 	space   int //white space
// 	word    int
// 	lang    int // special words
// 	num     int
// 	oper    int
// 	comm    int
// 	text    int
// 	quot    int
// 	esc     int
// 	block   int // open block. start of block if more than shift of prev line
// 	close   int // close block
// 	endline int
// 	indent  int // set of spaces, 4-th by default,
// 	mtcomm  int // multiline comment
// 	mttext  int // multiline string
// 	mclose  int // close multiline cases
// }

// var Lt = LangLTypes{
// 	none:    0,
// 	space:   1,
// 	word:    2,
// 	lang:    3, // special words
// 	num:     4,
// 	oper:    5,
// 	comm:    6,
// 	text:    7,
// 	quot:    8,
// 	esc:     9,
// 	block:   11, // open block. start of block if more than shift of prev line
// 	close:   12, // close block
// 	endline: 13,
// 	indent:  14, // set of spaces, 4-th by default,
// 	mtcomm:  15, // multiline comment
// 	mttext:  16, // multiline string
// 	mclose:  17, // close multiline cases
// }
