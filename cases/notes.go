package cases

/*
// single
-- var
name, _a, c22
-- val
12, 2.5, 'quot1', "quot2", `backticks`, '''tripl1''', """triple2""", ```triple-back```
-- prexif/postfix vals
g'D', re`[a-z]+`im
-- unary oper
-12, +12, !val, ~"{n}"
-- unary-right oper
args... , expr /:
-- binary oper
r = 1, expr <+> expr
-- keyword
func name(args)
struct Name f:type
if aexpr
for expr
while expr
match expr
break
continue
return |expr|
group Name
enum Name
-- special / service
@debug, @exit, @! (@del?)
-- sequences
1,2,3 // expr ; expr ; expr // T|int|string
-- brackets
(expr), [expr], {expr:expr, expr:expr}
[expr; v <- expr; expr], {expr; v,k ,- expr;..}, (= expr; v <- expr, expr)
-- brackets after word
foo(args)
func foo(args)
func v:T foo(args)
SType{f:val}
coll[key]
coll[i1:i2]
-- pref brackets
0x[12 13]
0b[11 00 11]


-- common cases:
single lexem
(brackets)
word expr
expr expr // as usual, definition
expr(brackets)
oper
	expr <oper> expr
	<oper>expr
	expr<oper>
sequence , ;
prefix: pref'quotes', pref[brackets]

-- resolving:

kword expr - definition | control

expr(n) - function: def, call
expr[n] - contaiter elem, expr[:] - slice
0n[] - bytes
pref"" - string cases: regexp, glif

n * n - bin oper
* n -- pref unary
n * -- post urary

(expr) - grouping
(,,) - collection
[;;] [ .. ] -- generator/comprehension

*/

/*
kword:
sub-nodes

func
	? var:type
	name
	(args)

if
	? pre-assign ;
	condition

for
	iter | ;sequence

struct|grup|enum
	name
	? elems

while|match
	expr

return
	?expr


---
-- match case
mcase
	? :? expr
	? /: ?expr


*/
