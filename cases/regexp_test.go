package cases

import (
	"strings"
	"testing"
)

/*
1. regexp object
2. =~ match operator
3. ?~ find operator
4. split
5. replace
*/

func TestRXOper(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// =~
		{`
		rx = re@@^[a-z\-\.]+$@@
		ss = ['Abc', 'abc', 'abc ', 'ab12', '132', 'qwerty', 'my-dom.com']
		r = []
		for s <- ss
			t = rx =~ s
			r <- s
			r <- t
		`, "r", Anis("Abc", false, "abc", true, "abc ", false, "ab12", false, "132", false, "qwerty", true, "my-dom.com", true)},
		{`
		rx = re@@^[a-z\-\.]+$@@
		ss = ['Abc', 'abc', 'abc ', "qwerty", 'my-dom-bom.com', 'v123.net']
		r = []
		for i, s <- ss
			if rx =~ s
				r <- s
		`, "r", Anis("abc", "qwerty", "my-dom-bom.com")},
		{`
		rx = re@@^[a-z\-\.]+$@@
		ss = ['Abc', 'abc', 'abc ', "qwerty", 'my-dom-bom.com', 'v123.net']
		r = []
		for i, s <- ss
			r <- i	
			r <- rx =~ s
		`, "r", Anis(0, false, 1, true, 2, false, 3, true, 4, true, 5, false)},
		// {``, "r",  ""},
		// ?~
		{`
		rx = re@@\d+[a-z]+@@
		s = "1 aaa bbb 2cc dd 3eee"
		ss = rx ?~ s
		r = []
		for n <- ss
			r <- n[0]
		`, "r", Anis("2cc", "3eee")},
		{`
		rx = re@@\d+([a-z]+)@@
		s = "1aaa bbb 2cc dd 3eee"
		ss = rx ?~ s
		r = []
		for n <- ss
			r <- n[1]
		`, "r", Anis("aaa", "cc", "eee")},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		tt.src = strings.ReplaceAll(tt.src, "$%$", "```")
		tt.src = strings.ReplaceAll(tt.src, "@@", "`")
		RunTCodeVarExp(t, i, tt)
	}
}

func TestRXBuiltinMethods(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		// match
		{`
		rx = re~^\d+[a-v]+$~i
		r = ['match1']
		ss = ['123', '12qqe', '0rty', '0rta', 'asdf', ' 123 zxc']
		for s <- ss
			m = rx.match(s)
			r <- s
			r <- m
		`, "r", Anis("match1", "123", false, "12qqe", true, "0rty", false, "0rta", true, "asdf", false, " 123 zxc", false)},
		{`
		rx = re~^(\d+)[\-\:\/]\d{2}[\-\:\/]\d{2}$~i
		r = ['match2']
		ss = ['1-1-1', '123-44-55', '1:22:33', 'abc-22-33', '1/23/45', 'qwerty', '- - -']
		for s <- ss
			m = rx.match(s)
			r <- s
			r <- m
		`, "r", Anis("match2", "1-1-1", false, "123-44-55", true, "1:22:33", true,
			"abc-22-33", false, "1/23/45", true, "qwerty", false, "- - -", false)},
		{`
		rx = re~^[a-z]+ \d{2,5}$~i
		r = ['match3']
		ss = ['Abc 123', '123 123', 'Unit 0', 'Qs 1234567', 'Amber 17',  '- - ?']
		for s <- ss
			m = rx.match(s)
			r <- s
			r <- m
		`, "r", Anis("match3", "Abc 123", true, "123 123", false, "Unit 0", false, "Qs 1234567", false, "Amber 17", true, "- - ?", false)},
		// find
		{`
		rx = re~\d+~i
		r = ['find1']
		ss = ['abc', 'as 12 rr 34 ty567', '222 333', '3344 asd', 'cc 55 gg 66', '\\d \\d', '(67)-55-44']
		for i, s <- ss
			parts = rx.find(s)
			r <- '%d) %s {{%s}}' << (i , s, parts.join('_#_'))
		`, "r", Anis("find1", "0) abc {{}}", "1) as 12 rr 34 ty567 {{12_#_34_#_567}}", "2) 222 333 {{222_#_333}}",
			"3) 3344 asd {{3344}}", "4) cc 55 gg 66 {{55_#_66}}", "5) \\d \\d {{}}", "6) (67)-55-44 {{67_#_55_#_44}}")},
		// findSubs
		{`
		rx = re~pref([0-9a-f]+)-(\d+)~i
		r = ['findSubs1']
		ss = ['pref0F-123', 'prefA-22 prefB-33', 'Hddff! 7123/prefBF-44 543/prefCE-55']
		for i, s <- ss
			subs = rx.findSubs(s)
			r <- subs
		`, "r", Anis("findSubs1",
			Anis(Anis("pref0F-123", "0F", "123")),
			Anis(Anis("prefA-22", "A", "22"), Anis("prefB-33", "B", "33")),
			Anis(Anis("prefBF-44", "BF", "44"), Anis("prefCE-55", "CE", "55")),
		)},
		// replace
		{`
		rx = re~A[a-z]+\d~
		r = ['replace1']
		ss = ['qwerty1', 'banana-Abanana2', 'Hello Apple-Apple3 Avo-Avocado4', 'n51Aa51 n52Ab52 n53Ac53 n54Ae54', 'XsAe0 Aeo']
		for i, s <- ss
			r <- rx.replace(s, '<%d>' << i)
		`, "r", Anis("replace1", "qwerty1", "banana-<1>", "Hello Apple-<2> Avo-<2>", "n51<3>1 n52<3>2 n53<3>3 n54<3>4", "Xs<4> Aeo")},
		{`
		# using backreference
		rx = re~c\.([a-z]+),\s*st\.([a-z]+)~i
		r = ['replace2']
		ss = ['Story 1 c.First, st.Red ', 'More c.Second, st.Green', 'After c.Third, st.White']
		for i, s <- ss
			r <- rx.replace(s, 'Location: <$1>.<$2>')
		`, "r", Anis("replace2", "Story 1 Location: <First>.<Red> ", "More Location: <Second>.<Green>", "After Location: <Third>.<White>")},
		// {``, "r",  ""},
		// split
		{`
		rx = re~[\s\-\=\/\:\?]+~
		#rx = re~[\-]+~
		r = ['split1']
		ss = ['aaa aaa aa-aaa/aaaaa', 'bbbbb=bbbbb1:bbbbb2', 'cccc ccc1 cccc ccc', 'dddd-dddd-ddddd', 'eeeee?ee?eeee']
		for i, s <- ss
			pp = rx.split(s)
			#r <- pp
			r <- pp.join(', ')
		`, "r", Anis("split1", "aaa, aaa, aa, aaa, aaaaa", "bbbbb, bbbbb1, bbbbb2", "cccc, ccc1, cccc, ccc", "dddd, dddd, ddddd", "eeeee, ee, eeee")},
		// {``, "r",  ""},
		// {``, "r",  Anis()},
	}
	for i, tt := range tdata {
		tt.src = strings.ReplaceAll(tt.src, "$%$", "```")
		tt.src = strings.ReplaceAll(tt.src, "~", "`")
		RunTCodeVarExp(t, i, tt)
	}
}

func TestRXSimple(t *testing.T) {
	tdata := []struct {
		src   string
		vname string
		res   any
	}{
		{`
		rx: regexp = re~abc~
		`, "rx", Rx("abc")},
		{`
		r = re'\'\"\~\\s\\d\\n\\b'
		`, "r", Rx("'\"`\\s\\d\\n\\b")},
		{`
		r = re~' " \ / \s \d \b \n ~
		`, "r", Rx("' \" \\ / \\s \\d \\b \\n ")},
		{`
		r = re~\d+~is
		`, "r", Rx("(?is)\\d+")},
		{`
		r = re~\d+~Ui
		`, "r", Rx("(?Ui)\\d+")},
		{`
		r = re~\d+~ms
		`, "r", Rx("(?ms)\\d+")},
		// {``, "r",  ""},
		// {``, "r",  ""},
		// {``, "r",  ""},
		// {`
		// rx = re~abc~
		// #r = []
		// #r <- ('qwerty', rx =~ 'qwerty')
		// #r <- ('abc', rx =~ 'abc')
		// `, "rx", Rx("abc")},
		// {``, "r",  ""},
	}
	for i, tt := range tdata {
		tt.src = strings.ReplaceAll(tt.src, "$%$", "```")
		tt.src = strings.ReplaceAll(tt.src, "~", "`")
		RunTCodeVarExp(t, i, tt)
	}
}
